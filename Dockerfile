FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates git

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/backend .

EXPOSE 10000

CMD ["./backend"]