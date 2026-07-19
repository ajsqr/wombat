# 🦡 Wombat

> **A reverse TCP tunneling service written in Go.**

Wombat securely exposes locally running TCP services through a persistent reverse tunnel. It multiplexes multiple independent client connections over a single TCP connection while maintaining mirrored session state at both ends of the tunnel.

Designed around a small set of composable components, Wombat separates transport, session management, and endpoint lifecycle, making each part independently testable and easy to reason about.

---

## Features

* Multiplex multiple TCP streams over a single persistent tunnel
* Custom binary framing protocol
* Symmetric client/server architecture
* Mirrored session lifecycle
* Transport-agnostic tunnel implementation
* Concurrent session handling
* Component-level unit tests
* Internet-tested deployment

---

## Architecture

```text
                         Internet
                             │
                      Client Connection
                             │
                     Wombat Server (VPS)
                             │
                    Persistent TCP Tunnel
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
| Ping            | Reserved                  |
| Pong            | Reserved                  |

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

* Concurrent HTTP requests
* Streaming downloads (100 MB+)
* Multi-session multiplexing over a single tunnel
* Deployment across an AWS EC2 instance with a locally running agent

---

## Motivation

Reverse tunneling systems hide a surprising amount of engineering behind a simple interface.

Wombat was built to explore those implementation details from first principles by designing and implementing each major component independently, including:

* Binary framing
* Session multiplexing
* Transport abstraction
* Concurrent stream management
* Session lifecycle coordination

The objective was to understand the architecture behind production-grade reverse tunneling systems by building one from scratch.

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
