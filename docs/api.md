# API Reference

## Tauri Commands (Frontend ↔ Host)

All commands are called via `invoke(commandName, params)` from the Svelte frontend.

### `get_status`

Returns the current connection status.

**Parameters:** none

**Response:**

```ts
{
  id: string;          // local node UUID
  virtual_ip: string;  // assigned 10.88.x.x address
  peers: PeerInfo[];
}
```

**Errors:** `"Not connected"` when the service is not running.

---

### `connect`

Starts the Go service sidecar and announces the node on the virtual LAN.

**Parameters:**

```ts
{ alias: string }  // display name shown to peers
```

**Response:** `void`

**Errors:** `"Already connected"`, `"Alias must not be empty"`

---

### `disconnect`

Stops the Go service sidecar and clears local state.

**Parameters:** none

**Response:** `void`

---

## PeerInfo

```ts
interface PeerInfo {
  id: string;         // peer node UUID
  alias: string;      // peer display name
  virtual_ip: string; // peer 10.88.x.x address
  public_ip: string;  // peer external IP:port (for diagnostics)
}
```

---

## Go Service JSON-RPC (Host ↔ Service)

The Go sidecar communicates over **stdin / stdout** using newline-delimited JSON.

### Request

```json
{ "method": "peers", "id": 1 }
```

### Response

```json
{ "id": 1, "result": [ { "id": "…", "alias": "…", "virtual_ip": "…", "public_ip": "…" } ] }
```

### Methods

| Method | Description |
|---|---|
| `status` | Returns node ID and virtual IP |
| `peers` | Returns the current peer list |
| `announce` | Re-broadcasts a join message |
