FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY data-platform/export-service/go.mod data-platform/export-service/go.sum ./
RUN go mod download

COPY data-platform/export-service/ .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o export-service ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/export-service .

EXPOSE 3002

CMD ["./export-service"]
