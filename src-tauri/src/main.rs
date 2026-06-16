// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use rimesh_lib::{connect, disconnect, get_status, AppState, ServiceHandle};
use std::sync::Mutex;

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .manage(AppState {
            service: Mutex::new(ServiceHandle::default()),
        })
        .invoke_handler(tauri::generate_handler![get_status, connect, disconnect])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
