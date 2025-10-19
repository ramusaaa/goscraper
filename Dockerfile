FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goscraper-server cmd/server/main.go

FROM alpine:3.18

RUN apk add --no-cache chromium ca-certificates curl dumb-init

RUN addgroup -g 1001 -S goscraper && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G goscraper -g goscraper goscraper

WORKDIR /app

COPY --from=builder /app/goscraper-server /app/goscraper-server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN chown -R goscraper:goscraper /app && chmod +x /app/goscraper-server

USER goscraper

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

EXPOSE 8080 9090

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/goscraper-server"]