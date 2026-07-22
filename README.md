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
