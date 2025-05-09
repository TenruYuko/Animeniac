#!/bin/sh
set -e

# Force DNS to use Gluetun's internal gateway (VPN protected)
# Remove /etc/resolv.conf if it exists or is a symlink, then write correct DNS
touch /etc/resolv.conf
rm -f /etc/resolv.conf
printf 'nameserver 127.0.0.1\n' > /etc/resolv.conf

# Start qBittorrent-nox Web UI on 912 (username is always 'admin').
# On first launch, set the password via the Web UI at http://localhost:912, then use that password in Seanime.
qbittorrent-nox --webui-port=912 --profile=/app/qbittorrent-profile &

# Wait for qBittorrent-nox to be ready
for i in $(seq 1 15); do
  if nc -z localhost 912; then
    break
  fi
  sleep 1
done

# Start Seanime
exec /app/seanime --datadir /app/data
