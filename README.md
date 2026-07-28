# 🦡 Wombat

[![Release](https://img.shields.io/github/v/release/ajsqr/wombat)](https://github.com/ajsqr/wombat/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ajsqr/wombat)](https://go.dev/)
[![License](https://img.shields.io/github/license/ajsqr/wombat)](LICENSE)
[![CI](https://github.com/ajsqr/wombat/actions/workflows/release.yml/badge.svg)](https://github.com/ajsqr/wombat/actions/workflows/release.yml)

> **A lightweight reverse TCP tunneling service written in Go.**

Wombat securely exposes local TCP services to the internet through persistent, TLS-encrypted tunnels. Each configured service gets its own dedicated tunnel, while multiple client connections are multiplexed efficiently over that tunnel using Wombat's custom binary framing protocol.

## Learn more

- 📖 **[Getting Started Guide](https://ajsqr.dev/blog/wombat-reverse-tunnel-tutorial/)** – Set up Wombat in minutes.
- 📝 **[Building Wombat: A Reverse TCP Tunnel](https://ajsqr.dev/blog/building-wombat-a-reverse-tcp-tunnel/)** – Learn how Wombat's architecture and custom framing protocol work.

---

## Features

- ✅ TLS-encrypted persistent tunnels
- ✅ One persistent tunnel per configured service
- ✅ Multiplex multiple TCP connections over dedicated tunnels
- ✅ Token-based tunnel authentication
- ✅ Custom binary framing protocol
- ✅ Automatic tunnel reconnection
- ✅ Cross-platform (Linux, macOS, Windows)
- ✅ Forward any TCP service (HTTP, HTTPS, SSH, databases, MQTT, etc.)

---

## Installation

### Unix-like systems

```bash
curl -fsSL https://raw.githubusercontent.com/ajsqr/wombat/main/install.sh | bash
```

For Windows and other platforms, download the appropriate binary from the latest GitHub release.

### Verify the installation

```bash
wombat-server version
wombat-agent version
```

Example output:

```text
wombat-server v0.2.0 build abc1234
wombat-agent v0.2.0 build abc1234
```

---

## Quick Start


Run the guided setup to generate certificates, create the required configuration files, and configure your first Wombat tunnel.

```bash
curl -fsSL https://raw.githubusercontent.com/ajsqr/wombat/quickstart-cli/scripts/quickstart/quickstart.sh | bash -s -- --force-certs
```

To reuse your existing Wombat certificates, omit the --force-certs flag:

```bash
curl -fsSL https://raw.githubusercontent.com/ajsqr/wombat/quickstart-cli/scripts/quickstart/quickstart.sh | bash
```


## Manual Setup

Wombat uses TLS to secure the tunnel between the server and the agent. The following steps walk you through generating the required certificates, configuring authentication, and starting your first tunnel.


### 1. Generate a private Certificate Authority

```bash
mkdir certs
cd certs

openssl genrsa -out ca.key 4096

openssl req \
    -x509 \
    -new \
    -nodes \
    -key ca.key \
    -sha256 \
    -days 3650 \
    -out ca.crt \
    -subj "/CN=Wombat CA"
```

### 2. Generate a server certificate

```bash
openssl genrsa -out server.key 4096
```

Create `server.cnf`

```ini
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext

[dn]
CN = localhost

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
```

Generate the CSR:

```bash
openssl req \
  -new \
  -key server.key \
  -out server.csr \
  -config server.cnf
```

Sign it:

```bash
openssl x509 \
  -req \
  -in server.csr \
  -CA ca.crt \
  -CAkey ca.key \
  -CAcreateserial \
  -out server.crt \
  -days 365 \
  -sha256 \
  -extensions req_ext \
  -extfile server.cnf
```

The server uses `server.crt` and `server.key`. Agents only require `ca.crt`.

---

### 3. Shared Token

Every tunnel is authenticated using a shared secret.

Each tunnel specifies a `tokenName`. This is the **name of the environment variable** that contains the shared secret.

Example configuration:

```json
{
  "tokenName": "echo"
}
```

Before starting both the server and the agent:

#### Linux / macOS

```bash
export echo="my-super-secret-token"
```

#### Windows PowerShell

```powershell
$env:echo="my-super-secret-token"
```

The server and the corresponding agent must use the same token value.

---

### 4. Configuration Directory

Wombat automatically loads its configuration from the operating system's user configuration directory.

| Platform | Directory |
|----------|-----------|
| Linux | `~/.config/wombat/` |
| macOS | `~/Library/Application Support/wombat/` |
| Windows | `%AppData%\wombat\` |

Place the following files there:

- `server-config.json`
- `agent-config.json`

TLS certificates may be stored anywhere and referenced by path.

---

### 5. Configure the Server

```json
{
  "certPath": "/path/to/server.crt",
  "keyPath": "/path/to/server.key",
  "tunnels": [
    {
      "name": "echo",
      "tunnel": "0.0.0.0:4001",
      "public": "0.0.0.0:8001",
      "tokenName": "echo"
    }
  ]
}
```

Run:

```bash
wombat-server run
```

Check version:

```bash
wombat-server version
```

---

### 6. Configure the Agent

```json
{
  "caCertPath": "/path/to/ca.crt",
  "serverName": "localhost",
  "tunnels": [
    {
      "name": "echo",
      "tunnel": "your-server:4001",
      "local": "127.0.0.1:8080",
      "tokenName": "echo"
    }
  ]
}
```

`serverName` must match one of the DNS names in the server certificate's Subject Alternative Name (SAN).

Run:

```bash
wombat-agent run
```

Check version:

```bash
wombat-agent version
```

Once connected, clients can connect to the configured public endpoint.

---

## Architecture

```text
Internet Client
       │
       ▼
Wombat Server (Public TCP Listener)
       │
       ▼
TLS-encrypted Persistent Tunnel
       │
       ▼
Wombat Agent
       │
       ▼
Local TCP Service
```

The public listener forwards raw TCP streams. Wombat encrypts only the persistent tunnel between the server and the agent, allowing protocols such as HTTPS, SSH and PostgreSQL to provide their own end-to-end encryption.

---

## Design Philosophy

The transport layer is intentionally protocol-agnostic.

- **Tunnel** transports frames.
- **Session** represents a single TCP connection.
- **Server** accepts public connections and mirrors sessions.
- **Agent** connects to local services and mirrors the server's session lifecycle.

---

## Protocol

Each logical TCP connection is multiplexed over a persistent tunnel.

Each frame contains:

- Frame Type
- Connection ID
- Payload Length
- Payload

| Frame | Purpose |
|-------|---------|
| OpenConnection | Create mirrored session |
| CloseConnection | Tear down session |
| DataFrame | Transport TCP payload |
| Ping | Reserved |
| Pong | Reserved |

---


## Roadmap

- Optional mutual TLS
- Heartbeats / keepalive
- Flow control
- Protocol versioning
- Compression
- Metrics & observability

---

## Why "Wombat"?

Wombats are exceptional tunnel builders.

Just like the animal, Wombat creates reliable tunnels that quietly connect distant places through a single network.
