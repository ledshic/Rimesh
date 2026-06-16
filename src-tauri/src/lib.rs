use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use tauri::State;

/// Shared application state managed across Tauri commands.
pub struct AppState {
    pub service: Mutex<ServiceHandle>,
}

/// A handle to the running Go sidecar process and its current status.
pub struct ServiceHandle {
    pub connected: bool,
    pub virtual_ip: String,
    pub node_id: String,
    pub peers: Vec<PeerInfo>,
}

impl Default for ServiceHandle {
    fn default() -> Self {
        ServiceHandle {
            connected: false,
            virtual_ip: String::new(),
            node_id: String::new(),
            peers: Vec::new(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PeerInfo {
    pub id: String,
    pub alias: String,
    pub virtual_ip: String,
    pub public_ip: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct StatusResponse {
    pub id: String,
    pub virtual_ip: String,
    pub peers: Vec<PeerInfo>,
}

/// Retrieve the current network status.
#[tauri::command]
pub fn get_status(state: State<AppState>) -> Result<StatusResponse, String> {
    let svc = state.service.lock().map_err(|e| e.to_string())?;
    if !svc.connected {
        return Err("Not connected".into());
    }
    Ok(StatusResponse {
        id: svc.node_id.clone(),
        virtual_ip: svc.virtual_ip.clone(),
        peers: svc.peers.clone(),
    })
}

/// Connect to the Rimesh virtual LAN with the given alias.
#[tauri::command]
pub async fn connect(alias: String, state: State<'_, AppState>) -> Result<(), String> {
    let mut svc = state.service.lock().map_err(|e| e.to_string())?;
    if svc.connected {
        return Err("Already connected".into());
    }
    if alias.trim().is_empty() {
        return Err("Alias must not be empty".into());
    }

    // In the real implementation the Tauri app would spawn the Go sidecar via
    // tauri-plugin-shell and communicate over stdin/stdout JSON-RPC.
    // For now we populate stub state so the UI can be exercised end-to-end.
    svc.connected = true;
    svc.node_id = format!("node-{}", uuid_v4_stub(&alias));
    svc.virtual_ip = "10.88.1.1".to_string();
    svc.peers = Vec::new();

    Ok(())
}

/// Disconnect from the virtual LAN and stop the service sidecar.
#[tauri::command]
pub fn disconnect(state: State<AppState>) -> Result<(), String> {
    let mut svc = state.service.lock().map_err(|e| e.to_string())?;
    svc.connected = false;
    svc.virtual_ip = String::new();
    svc.node_id = String::new();
    svc.peers = Vec::new();
    Ok(())
}

/// Naive deterministic "UUID-like" string derived from a seed string.
/// Replaced by a proper UUID crate in production.
fn uuid_v4_stub(seed: &str) -> String {
    let hash: u64 = seed
        .bytes()
        .enumerate()
        .fold(0u64, |acc, (i, b)| acc ^ ((b as u64).wrapping_shl((i as u32) % 64)));
    format!("{:016x}", hash)
}
