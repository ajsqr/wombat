# 🦡 Wombat

> **A lightweight reverse TCP tunneling service written in Go.**

Wombat is a lightweight TCP tunneling service that lets users easily expose their local services to the internet, bypassing NAT restrictions.
Wombat creates isolated TCP tunnels for each service and uses it to multipex different client sessions. 

---

## Features

* Multiplex multiple TCP sessions over isolated persistent tunnels
* Custom binary framing protocol
* Symmetric client/server architecture
* Mirrored session lifecycle
* Cross-platform (Linux, macOS, Windows)

---

## Installation

For *nix systems wombat server / agent can be easily installed using the guided installer 

```.sh
curl -fsSL https://raw.githubusercontent.com/ajsqr/wombat/main/install.sh | bash
```

For other operating systems, a matching version can be downloaded from the latest release.

## Usage

Wombat consists of two components:

* **wombat-server** — Runs on a publicly accessible server and accepts incoming client connections.
* **wombat-agent** — Runs on the machine hosting your local services and connects to the server.

> **Note:** Sample configuration files are available in the repository's `example/` directory.

### Running `wombat-server`

Create a `server-config.json` file in your user configuration directory:

```json
{
  "tunnels": [
    {
      "name": "echo",
      "tunnel": "<tunnel-ip>:<tunnel-port>",
      "public": "<public-ip>:<public-port>"
    }
  ]
}
```

The `tunnel` address is where the agent connects, while the `public` address is where clients connect.

Start the server:

```bash
wombat-server
```

Or run it in the background using `nohup`:

```bash
nohup wombat-server > wombat-server.log 2>&1 &
```

### Running `wombat-agent`

Create an `agent-config.json` file in your user configuration directory:

```json
{
  "tunnels": [
    {
      "name": "echo",
      "tunnel": "<tunnel-ip>:<tunnel-port>",
      "local": "<local-ip>:<local-port>"
    }
  ]
}
```

The `tunnel` address must match the server's tunnel endpoint. The `local` address is the local service you want to expose.

Ensure `wombat-server` is already running, then start the agent:

```bash
wombat-agent
```

Or run it in the background:

```bash
nohup wombat-agent > wombat-agent.log 2>&1 &
```

Once the agent connects successfully, clients can access your local service through the server's configured `public` endpoint.

> **Deployment:** Wombat does not require a specific deployment method. You are free to run `wombat-server` and `wombat-agent` however best suits your environment, whether that's directly as binaries, using `systemd`, `nohup`, `Docker`, or any other process manager.

## Architecture

```text
                         Internet
                             │
                      Client Connection
                             │
                     Wombat Server (VPS)
                             │
                    Persistent TCP Tunnel(s)
                             │
                        Wombat Agent
                             │
                     Local TCP Service
                (HTTP, API, Database, etc.)
```

The tunnel itself has no knowledge of the application protocol being transported. It simply carries framed TCP streams between the server and the agent.

---

## Design Philosophy

Wombat is built around explicit ownership and separation of responsibilities.

### Tunnel

Responsible only for transporting protocol frames between two endpoints.

The tunnel has no knowledge of client sockets, local applications, or session creation. It simply delivers frames.

### Session

Represents a single TCP connection.

A session converts bytes read from a socket into protocol frames and reconstructs incoming frames back into a byte stream.

### Server

The server accepts incoming client connections, creates sessions, and mirrors them on the remote agent.

### Agent

The agent establishes local connections to services running on the host machine and mirrors the session lifecycle initiated by the server.

This separation allows the transport layer to remain completely independent of the applications using it.

---

## Protocol

Logical TCP connections are multiplexed over a single persistent connection using a custom binary framing protocol.

Each frame contains:

* Frame Type
* Connection ID
* Payload Length
* Payload

Supported frame types:

| Frame           | Purpose                   |
| --------------- | ------------------------- |
| OpenConnection  | Create a mirrored session |
| CloseConnection | Tear down a session       |
| DataFrame       | Transport TCP payload     |
| Ping            | TBD                       |
| Pong            | TBD                       |

Connection IDs uniquely identify individual streams, allowing many independent sessions to coexist over one tunnel.

---

## Current Capabilities

The current implementation supports:

* Reverse TCP tunneling
* HTTP forwarding
* Concurrent multiplexed sessions
* Large streaming responses
* Automatic session mirroring
* Graceful handling of peer disconnects

Current testing includes:

---

## Roadmap

Planned improvements include:

* TLS encryption
* Mutual authentication
* Automatic tunnel reconnection
* Heartbeats and keepalive
* Flow control
* Protocol versioning
* Compression
* Metrics and observability
* Benchmark suite
* Configuration file support

---

## Why "Wombat"?

Wombats are exceptional tunnel builders.

They construct robust underground tunnel systems that quietly connect distant places through a single network.

The name reflects the purpose of this project: building reliable tunnels that connect services without exposing them directly to the internet.
