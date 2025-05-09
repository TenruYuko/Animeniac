# Mullvad OpenVPN Integration (Sweden)

To use Mullvad OpenVPN with your docker-compose stack, you need:
- Your Mullvad account number
- A Mullvad OpenVPN configuration file for a Swedish server (e.g., se-stockholm-xyz.mullvad.net)

## Steps:
1. Download a Swedish OpenVPN config from https://mullvad.net/download/openvpn-config/ (choose Sweden).
2. Place the `.ovpn` file and any required certs in this directory (`vpn/`).
3. Set your Mullvad account number as an environment variable or in a file as needed.

The docker-compose file will use the `qmcgaw/gluetun` image for VPN routing.
