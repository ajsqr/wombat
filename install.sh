#!/usr/bin/env bash

set -euo pipefail

###############################################################################
# Choose component
###############################################################################

if [[ $# -ge 1 ]]; then
    case "$1" in
        agent)
            binary="wombat-agent"
            ;;
        server)
            binary="wombat-server"
            ;;
        *)
            echo "Usage: $0 [agent|server]"
            exit 1
            ;;
    esac
else
    echo "Which component would you like to install?"
    echo "  1) Agent"
    echo "  2) Server"

    read -rp "Choice [1-2]: " choice </dev/tty

    case "$choice" in
        1) binary="wombat-agent" ;;
        2) binary="wombat-server" ;;
        *)
            echo "Invalid choice."
            exit 1
            ;;
    esac
fi

###############################################################################
# Detect OS
###############################################################################

case "$(uname -s)" in
    Darwin)
        os="darwin"
        config_dir="$HOME/Library/Application Support/wombat"
        ;;
    Linux)
        os="linux"
        config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/wombat"
        ;;
    *)
        echo "Unsupported operating system."
        exit 1
        ;;
esac

###############################################################################
# Detect architecture
###############################################################################

case "$(uname -m)" in
    x86_64)
        arch="amd64"
        ;;
    arm64|aarch64)
        arch="arm64"
        ;;
    *)
        echo "Unsupported architecture: $(uname -m)"
        exit 1
        ;;
esac

###############################################################################
# Create config directory
###############################################################################

mkdir -p "$config_dir"

###############################################################################
# Download
###############################################################################

url="https://github.com/ajsqr/wombat/releases/latest/download/${binary}_${os}_${arch}.tar.gz"

echo "Downloading ${binary}..."

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

archive="$tmpdir/${binary}.tar.gz"

curl -fsSL "$url" -o "$archive"

###############################################################################
# Extract
###############################################################################

echo "Extracting..."

tar -xzf "$archive" -C "$tmpdir"

###############################################################################
# Install
###############################################################################

if [[ ! -w /usr/local/bin ]]; then
    echo "Installing to /usr/local/bin (you may be prompted for your password)..."
fi

sudo install -m 755 "$tmpdir/$binary" "/usr/local/bin/$binary"

###############################################################################
# Done
###############################################################################

echo
echo "✅ Successfully installed ${binary}"
echo
echo "Binary:"
echo "  /usr/local/bin/${binary}"
echo
echo "Config directory:"
echo "  ${config_dir}"
echo
echo "Run:"
echo "  ${binary} run"