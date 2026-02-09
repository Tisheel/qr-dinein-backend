FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

# Render provides PORT env var; GoFr expects HTTP_PORT
ENTRYPOINT sh -c "HTTP_PORT=\${PORT:-9000} exec ./server"
