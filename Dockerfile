# Stage 1: Build the static Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ ./src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /bin/eval-diff ./src

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -S eval-diff \
    && adduser -S -D -H -G eval-diff eval-diff

COPY --from=builder /bin/eval-diff /usr/local/bin/eval-diff
RUN chmod +x /usr/local/bin/eval-diff

USER eval-diff

ENTRYPOINT ["/usr/local/bin/eval-diff"]
