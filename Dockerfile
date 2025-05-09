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
# Copy the built web assets into the expected directory
RUN mkdir -p /app/web && cp -r /app/web/build /app/web/
RUN go build -buildvcs=false -o seanime .

# --- Stage 3: Production image ---
FROM debian:bookworm-slim
WORKDIR /app
# Install runtime dependencies (adjust as needed)
RUN apt-get update && apt-get install -y ca-certificates ffmpeg && rm -rf /var/lib/apt/lists/*
COPY --from=go-build /app/seanime /app/seanime
COPY --from=go-build /app/data /app/data
COPY --from=go-build /app/web /app/web
COPY --from=go-build /app/seanime.db /app/seanime.db
EXPOSE 43211
ENV TZ=UTC
CMD ["/app/seanime", "--datadir", "/app/data"]
