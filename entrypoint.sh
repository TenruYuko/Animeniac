#!/bin/sh
set -e

echo 'nameserver 127.0.0.1' > /etc/resolv.conf

# Start qBittorrent-nox Web UI on 912 (username is always 'admin').
qbittorrent-nox --webui-port=912 --profile=/app/qbittorrent-profile > /tmp/qbittorrent.log 2>&1 &
QBIT_PID=$!

# Wait for qBittorrent-nox to be ready
for i in $(seq 1 15); do
  if nc -z localhost 912; then
    break
  fi
  sleep 1
done

# Check if qBittorrent-nox is still running
if ! kill -0 $QBIT_PID 2>/dev/null; then
  echo "qBittorrent-nox failed to start. Log output:"
  cat /tmp/qbittorrent.log
  exit 1
fi

cat /tmp/qbittorrent.log

# Start Seanime
exec /app/seanime --datadir /app/data
