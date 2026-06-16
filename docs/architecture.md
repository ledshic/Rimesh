# Architecture

## Overview

```
┌─────────────────────────────────────────────┐
│                  Rimesh App                  │
│                                              │
│  ┌──────────────────────────────────────┐   │
│  │         Svelte/SvelteKit UI           │   │
│  │  (src/)                               │   │
│  └──────────────┬───────────────────────┘   │
│                 │ Tauri IPC (invoke)         │
│  ┌──────────────▼───────────────────────┐   │
│  │         Tauri Shell (Rust)            │   │
│  │  (src-tauri/)                         │   │
│  └──────────────┬───────────────────────┘   │
│                 │ stdin/stdout JSON-RPC       │
│  ┌──────────────▼───────────────────────┐   │
│  │      rimesh-service (Go sidecar)      │   │
│  │  (service/)                           │   │
│  └──────────────┬───────────────────────┘   │
└─────────────────┼───────────────────────────┘
                  │ UDP (LAN broadcast / ICE)
              ┌───▼───┐
              │  P2P  │
              └───────┘
```

## Components

### 1. Svelte/SvelteKit UI (`src/`)

The frontend is a standard SvelteKit application compiled to a static bundle (`adapter-static`). Tauri loads it from the `build/` output directory.

Key pages / routes:

| Route | Description |
|---|---|
| `/` | Main dashboard – connect, view peers |

Key components:

| Component | Purpose |
|---|---|
| `PeerList.svelte` | Displays connected peers |
| `StatusBar.svelte` | Footer showing connection state |

### 2. Tauri Host (`src-tauri/`)

Written in Rust, the Tauri shell:

- Provides the native window frame and OS integration.
- Exposes Tauri **commands** (`invoke`) to the frontend:
  - `get_status` – returns current connection state
  - `connect(alias)` – starts the Go sidecar and announces presence
  - `disconnect` – stops the sidecar
- Manages the Go sidecar process via `tauri-plugin-shell`.

### 3. Go Service (`service/`)

The Go binary `rimesh-service` is the networking engine:

| Package | Responsibility |
|---|---|
| `internal/network` | Virtual subnet management, UDP signalling, peer discovery |
| `internal/peer` | Peer state machine (disconnected → connecting → connected) |
| `cmd/rimesh-service` | Entry point; exposes JSON-RPC over stdin/stdout |

#### Virtual IP Allocation

Each node deterministically derives its `10.88.x.x` address from its UUID so collisions are rare and no central IPAM is needed.

#### Signalling Protocol

Nodes broadcast UDP messages on port **7788** within the local subnet:

| Message | Purpose |
|---|---|
| `join` | Announce presence; carries `PeerInfo` |
| `leave` | Graceful departure |
| `heartbeat` | Keep-alive every 10 s |
| `peer_list` | Full peer roster (future) |

## Build Matrix

| Target | GOARCH | Rust target |
|---|---|---|
| macOS Apple Silicon | `arm64` | `aarch64-apple-darwin` |
| Windows 11 (x64) | `amd64` | `x86_64-pc-windows-msvc` |

## Release Workflow

See [`.github/workflows/release.yml`](../.github/workflows/release.yml).
