package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered account in the platform.
type User struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	Username          string     `json:"username"`
	PasswordHash      string     `json:"-"`
	DisplayName       string     `json:"display_name"`
	Plan              string     `json:"plan"`
	Status            string     `json:"status"`
	StorageUsedBytes  int64      `json:"storage_used_bytes"`
	StorageQuotaBytes int64      `json:"storage_quota_bytes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
}
