#!/usr/bin/env bash
set -euo pipefail

###############################################################################
# Colors
###############################################################################

GREEN='\033[0;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
CYAN='\033[1;36m'
BOLD='\033[1m'
RESET='\033[0m'

prompt() {
    local var="$1"
    local message="$2"
    read -r -p "$message" "$var" </dev/tty
}

FORCE_CERTS=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --force-certs) FORCE_CERTS=true ;;
        *) echo "Unknown option: $1"; exit 1 ;;
    esac
    shift
done

case "$(uname -s)" in
  Darwin) CONFIG_DIR="$HOME/Library/Application Support/wombat" ;;
  Linux) CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/wombat" ;;
  *) echo "Unsupported operating system."; exit 1 ;;
esac

CERT_DIR="$CONFIG_DIR/certs"

command -v openssl >/dev/null || { echo "openssl is required."; exit 1; }

mkdir -p "$CONFIG_DIR"
mkdir -p "$CERT_DIR"

if [[ "$FORCE_CERTS" == "true" ]]; then
    echo -e "${YELLOW}Regenerating certificates...${RESET}"
    rm -f "$CERT_DIR"/ca.* "$CERT_DIR"/server.*
fi



echo
echo -e "${CYAN}${BOLD}Let's configure your first tunnel.${RESET}"
echo
echo -e "${YELLOW}Press Enter to accept the defaults.${RESET}"

echo
echo -e "${GREEN}${BOLD}1. Local service${RESET}"
echo "The address of the TCP service running on the same machine as wombat-agent which you would like to expose."
echo "Examples: 127.0.0.1:8080 (web app), localhost:5432 (PostgreSQL), localhost:22 (SSH)"
prompt LOCAL_ADDR "Local service [127.0.0.1:8080]: "
LOCAL_ADDR=${LOCAL_ADDR:-127.0.0.1:8080}

echo
echo -e "${GREEN}${BOLD}2. Public address${RESET}"
echo "This is where your users will connect."
echo "Example: demo.example.com:80 or 13.45.67.89:8001"
prompt PUBLIC_ADDR "Public address [0.0.0.0:8001]: "
PUBLIC_ADDR=${PUBLIC_ADDR:-0.0.0.0:8001}

echo
echo -e "${GREEN}${BOLD}3. Tunnel address${RESET}"
echo "The address wombat-agent uses to connect to wombat-server."
echo "This address is used internally by Wombat and should not be used by anyone else."
echo "Example: demo.example.com:4001"

prompt TUNNEL_ADDR "Tunnel address [localhost:4001]: "
TUNNEL_ADDR=${TUNNEL_ADDR:-localhost:4001}
echo

SERVER_NAME="${TUNNEL_ADDR%:*}"

TOKEN="$(openssl rand -hex 32)"

cat > "$CONFIG_DIR/env.sh" <<EOF
#!/usr/bin/env bash
export WOMBAT_TOKEN="$TOKEN"
EOF
chmod +x "$CONFIG_DIR/env.sh"

if [[ -f "$CERT_DIR/ca.crt" && -f "$CERT_DIR/server.crt" ]]; then
    echo -e "${GREEN}Using existing certificates.${RESET}"
else

echo -e "${CYAN}Generating Certificate Authority...${RESET}"
openssl genrsa -out "$CERT_DIR/ca.key" 4096 >/dev/null 2>&1
openssl req -x509 -new -nodes -key "$CERT_DIR/ca.key" -sha256 -days 3650 -subj "/CN=Wombat CA" -out "$CERT_DIR/ca.crt" >/dev/null 2>&1

echo -e "${CYAN}Generating server certificate...${RESET}"
openssl genrsa -out "$CERT_DIR/server.key" 4096 >/dev/null 2>&1

cat > "$CERT_DIR/server.cnf" <<EOF
[req]
prompt=no
distinguished_name=dn
req_extensions=req_ext

[dn]
CN=$SERVER_NAME

[req_ext]
subjectAltName=@alt_names

[alt_names]
DNS.1=$SERVER_NAME
EOF

openssl req -new -key "$CERT_DIR/server.key" -out "$CERT_DIR/server.csr" -config "$CERT_DIR/server.cnf" >/dev/null 2>&1
openssl x509 -req -in "$CERT_DIR/server.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial -out "$CERT_DIR/server.crt" -days 365 -sha256 -extensions req_ext -extfile "$CERT_DIR/server.cnf" >/dev/null 2>&1
fi

cat > "$CONFIG_DIR/server-config.json" <<EOF
{
  "certPath":"$CERT_DIR/server.crt",
  "keyPath":"$CERT_DIR/server.key",
  "tunnels":[
    {
      "name":"quickstart",
      "tunnel":"$TUNNEL_ADDR",
      "public":"$PUBLIC_ADDR",
      "tokenName":"WOMBAT_TOKEN"
    }
  ]
}
EOF

cat > "$CONFIG_DIR/agent-config.json" <<EOF
{
  "caCertPath":"$CERT_DIR/ca.crt",
  "serverName":"$SERVER_NAME",
  "tunnels":[
    {
      "name":"quickstart",
      "tunnel":"$TUNNEL_ADDR",
      "local":"$LOCAL_ADDR",
      "tokenName":"WOMBAT_TOKEN"
    }
  ]
}
EOF

echo
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${GREEN}${BOLD}✅ Wombat quickstart environment created${RESET}"
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"

echo
echo -e "${YELLOW}${BOLD}Before you begin${RESET}"
echo "Make sure your application is already running and listening on:"
echo
echo "  $LOCAL_ADDR"
echo
echo "If your application uses a different address,"
echo "run this script again and update the Local service."

echo
echo -e "${CYAN}${BOLD}Step 1 - Start the Wombat server${RESET}"
echo "Open a new terminal and run:"
echo
echo "  source \"$CONFIG_DIR/env.sh\""
echo "  wombat-server run"

echo
echo -e "${CYAN}${BOLD}Step 2 - Start the Wombat agent${RESET}"
echo "Open another terminal and run:"
echo
echo "  source \"$CONFIG_DIR/env.sh\""
echo "  wombat-agent run"

echo
echo -e "${CYAN}${BOLD}Step 3 - Connect to your service${RESET}"
echo "Once the agent connects to the server,"
echo
echo "  $LOCAL_ADDR"
echo
echo "will be available at"
echo
echo "  $PUBLIC_ADDR"

echo
echo -e "${GREEN}Happy tunneling! 🦡${RESET}"