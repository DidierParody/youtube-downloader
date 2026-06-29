FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ffmpeg yt-dlp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build all worker binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/downloader ./cmd/downloader/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/metadata ./cmd/metadata/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/ocr ./cmd/ocr/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/embedding ./cmd/embedding/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/optimization ./cmd/optimization/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/thumbnail ./cmd/thumbnail/main.go

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates ffmpeg yt-dlp

WORKDIR /root/

COPY --from=builder /app/downloader /app/metadata /app/ocr /app/embedding /app/optimization /app/thumbnail ./

# The specific worker will be specified at runtime via CMD
CMD ["./downloader"]
