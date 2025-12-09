FROM golang:1.25-alpine AS builder

RUN apk update && apk add --no-cache \
	ca-certificates \
	git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" \
	-trimpath \
	-o ytdlp-api ./cmd/server


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
	ca-certificates \
	wget \
	unzip && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.cache

RUN curl -fsSL https://bun.sh/install | bash && \
	mv /root/.bun/bin/bun /usr/local/bin/bun

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux \
	-o /usr/local/bin/yt-dlp && chmod +x /usr/local/bin/yt-dlp

RUN groupadd -r appgroup && \
	useradd -r -g appgroup -u 1000 -d /app -s /sbin/nologin appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/ytdlp-api /app/ytdlp-api

RUN mkdir -p /app/cookies /app/cookies-cache /app/logs && \
	chown -R appuser:appgroup /app/cookies /app/cookies-cache /app/logs && \
	chmod 755 /app/cookies /app/cookies-cache /app/logs

VOLUME ["/app/cookies", "/app/cookies-cache", "/app/logs"]

USER appuser

EXPOSE 5000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
	CMD wget --no-verbose --tries=1 --spider http://localhost:5000/health || exit 1

CMD ["./ytdlp-api"]
