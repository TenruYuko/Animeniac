# --- Stage 1: Build the web frontend ---
FROM node:20 AS web-build
WORKDIR /app/web
COPY seanime-web/ ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi \
    && npm run build

# --- Stage 2: Build the Go backend ---
FROM golang:1.23.8 AS go-build
WORKDIR /app
COPY . .
RUN go build -buildvcs=false -o seanime .

# --- Stage 3: Production image ---
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates ffmpeg qbittorrent-nox tini netcat-openbsd && rm -rf /var/lib/apt/lists/*
COPY --from=go-build /app/seanime /app/seanime
COPY --from=go-build /app/data /app/data
COPY --from=web-build /app/web/out /app/web
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 43211
EXPOSE 912
ENV TZ=UTC
ENTRYPOINT ["/entrypoint.sh"]
