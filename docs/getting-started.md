# Getting Started

## Requirements

| Platform | Minimum version |
|---|---|
| macOS | 13 Ventura (Apple Silicon) |
| Windows | 11 |

## Installation

### Pre-built Releases

Download the latest release from the [Releases page](https://github.com/ledshic/Rimesh/releases).

| Platform | Installer |
|---|---|
| macOS (Apple Silicon) | `Rimesh_x.x.x_aarch64.dmg` |
| Windows 11+ | `Rimesh_x.x.x_x64-setup.exe` |

### Building from Source

See [Architecture](architecture.md) for full build instructions.

## Usage

1. **Launch** the Rimesh app.
2. Enter your **display name** (shown to other players).
3. Click **Connect**.
4. Share your **Virtual IP** (shown after connecting) with your friends.
5. Friends launch Rimesh, connect, and use *your* Virtual IP as the server address in Rimworld.

## Troubleshooting

### Cannot find peers

- Ensure all players are running the same version of Rimesh.
- Check that your firewall allows UDP port **7788**.

### High latency

Rimesh uses direct UDP connections between peers. High latency usually indicates a suboptimal NAT traversal path. Try restarting Rimesh to re-negotiate the connection.
