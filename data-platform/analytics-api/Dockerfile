# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY data-platform/analytics-api/go.mod data-platform/analytics-api/go.sum ./
RUN go mod download

COPY data-platform/analytics-api/ .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o analytics-api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/analytics-api .

EXPOSE 3001

CMD ["./analytics-api"]
