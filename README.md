# Rimesh

> P2P Virtual LAN tool for Rimworld multiplayer experience.

[![Build](https://github.com/ledshic/Rimesh/actions/workflows/build.yml/badge.svg)](https://github.com/ledshic/Rimesh/actions/workflows/build.yml)

Rimesh creates a virtual local area network between players so everyone can connect to a Rimworld server as if they're on the same home network — no port forwarding, no central relay server.

## Tech Stack

| Layer | Technology |
|---|---|
| Desktop app shell | [Tauri v2](https://tauri.app) (Rust) |
| UI | [Svelte 5](https://svelte.dev) / [SvelteKit](https://kit.svelte.dev) |
| Networking service | [Go 1.22](https://go.dev) (sidecar binary) |

## Supported Platforms

- **macOS** – Apple Silicon (arm64), 13 Ventura+
- **Windows** – 11+ (x64)

## Quick Start

```sh
# Install Node dependencies
npm install

# Start the Tauri dev server (builds Go sidecar first)
npm run tauri dev
```

## Building a Release

Push a tag that matches `v*.*.*` to trigger the automated [release workflow](.github/workflows/release.yml):

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Documentation

- [Getting Started](docs/getting-started.md)
- [Architecture](docs/architecture.md)
- [API Reference](docs/api.md)

## License

[MIT](LICENSE)
