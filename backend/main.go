package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var dbPool *pgxpool.Pool

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initDB() *pgxpool.Pool {
	dsn := getEnv("DATABASE_URL", "postgres://ytd:ytd_password@localhost:5432/ytd_db?sslmode=disable")
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil
	}
	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("failed to ping database", "error", err)
		return nil
	}
	slog.Info("connected to database")
	return pool
}

func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getEnv("JWT_SECRET", "dev-secret")))
}

type registerInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type downloadInput struct {
	YoutubeURL string `json:"youtube_url"`
	Quality    string `json:"quality"`
	Format     string `json:"format"`
}

func handleRegister(c *fiber.Ctx) error {
	var input registerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	id := uuid.New()
	now := time.Now().UTC()

	_, err = dbPool.Exec(context.Background(),
		`INSERT INTO "Usuario" (id, email, username, password_hash, plan, status, storage_used_bytes, storage_quota_bytes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'free', 'active', 0, 0, $5, $5)`,
		id, input.Email, input.Username, string(hash), now,
	)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
	}

	token, err := generateToken(id.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"user_id":      id.String(),
		"access_token": token,
	})
}

func handleLogin(c *fiber.Ctx) error {
	var input loginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	var userID string
	var passwordHash string
	err := dbPool.QueryRow(context.Background(),
		`SELECT id, password_hash FROM "Usuario" WHERE email = $1`,
		input.Email,
	).Scan(&userID, &passwordHash)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := generateToken(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.JSON(fiber.Map{"access_token": token})
}

func handleCreateDownload(c *fiber.Ctx) error {
	var input downloadInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	userID := "user-id-placeholder" // TODO: extract from JWT token

	id := uuid.New()
	now := time.Now().UTC()

	_, err := dbPool.Exec(context.Background(),
		`INSERT INTO "Descarga" (id, user_id, status, requested_quality, requested_format, created_at, updated_at)
		VALUES ($1, $2, 'pending', $3, $4, $5, $5)`,
		id, userID, input.Quality, input.Format, now,
	)
	if err != nil {
		slog.Error("failed to create download", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to create download"})
	}

	return c.Status(201).JSON(fiber.Map{
		"download_id": id.String(),
		"status":      "pending",
	})
}

func handleGetDownload(c *fiber.Ctx) error {
	downloadID := c.Params("id")

	var status, quality, format string
	var createdAt time.Time

	err := dbPool.QueryRow(context.Background(),
		`SELECT status, requested_quality, requested_format, created_at FROM "Descarga" WHERE id = $1`,
		downloadID,
	).Scan(&status, &quality, &format, &createdAt)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "download not found"})
	}

	return c.JSON(fiber.Map{
		"id":               downloadID,
		"status":           status,
		"requested_quality": quality,
		"requested_format":  format,
		"created_at":       createdAt,
	})
}

func main() {
	dbPool = initDB()
	if dbPool == nil {
		panic("failed to connect to database")
	}
	defer dbPool.Close()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api/v1")
	api.Post("/auth/register", handleRegister)
	api.Post("/auth/login", handleLogin)
	api.Post("/downloads", handleCreateDownload)
	api.Get("/downloads/:id", handleGetDownload)

	port := getEnv("PORT", "3000")
	slog.Info("server starting", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
