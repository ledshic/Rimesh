<script lang="ts">
	import { invoke } from '@tauri-apps/api/core';
	import { onMount } from 'svelte';

	import PeerList from '$lib/components/PeerList.svelte';
	import StatusBar from '$lib/components/StatusBar.svelte';

	interface Peer {
		id: string;
		alias: string;
		virtual_ip: string;
		public_ip: string;
	}

	let peers: Peer[] = [];
	let virtualIP = '';
	let nodeID = '';
	let connected = false;
	let alias = '';
	let loading = false;
	let errorMsg = '';

	onMount(async () => {
		await refresh();
	});

	async function refresh() {
		try {
			const status = await invoke<{ id: string; virtual_ip: string; peers: Peer[] }>('get_status');
			nodeID = status.id;
			virtualIP = status.virtual_ip;
			peers = status.peers ?? [];
			connected = true;
		} catch {
			connected = false;
		}
	}

	async function handleConnect() {
		if (!alias.trim()) {
			errorMsg = 'Please enter a display name.';
			return;
		}
		loading = true;
		errorMsg = '';
		try {
			await invoke('connect', { alias: alias.trim() });
			await refresh();
		} catch (e) {
			errorMsg = String(e);
		} finally {
			loading = false;
		}
	}

	async function handleDisconnect() {
		loading = true;
		try {
			await invoke('disconnect');
			connected = false;
			peers = [];
			virtualIP = '';
			nodeID = '';
		} catch (e) {
			errorMsg = String(e);
		} finally {
			loading = false;
		}
	}
</script>

<main>
	<header>
		<h1>Rimesh</h1>
		<p class="tagline">P2P Virtual LAN for Rimworld multiplayer</p>
	</header>

	<section class="card connect-card">
		{#if !connected}
			<h2>Connect to Network</h2>
			<div class="form-row">
				<input
					bind:value={alias}
					placeholder="Your display name"
					disabled={loading}
					on:keydown={(e) => e.key === 'Enter' && handleConnect()}
				/>
				<button class="primary" on:click={handleConnect} disabled={loading}>
					{loading ? 'Connecting…' : 'Connect'}
				</button>
			</div>
			{#if errorMsg}
				<p class="error">{errorMsg}</p>
			{/if}
		{:else}
			<div class="status-connected">
				<div class="status-info">
					<span class="dot online"></span>
					<strong>Connected</strong>
				</div>
				<div class="node-details">
					<span class="label">Your Virtual IP</span>
					<code>{virtualIP}</code>
				</div>
				<div class="node-details">
					<span class="label">Node ID</span>
					<code class="muted">{nodeID}</code>
				</div>
				<button class="secondary" on:click={handleDisconnect} disabled={loading}>
					Disconnect
				</button>
			</div>
		{/if}
	</section>

	{#if connected}
		<section class="card">
			<PeerList {peers} on:refresh={refresh} />
		</section>
	{/if}

	<StatusBar {connected} peerCount={peers.length} />
</main>

<style>
	main {
		max-width: 720px;
		margin: 0 auto;
		padding: 2rem 1rem 4rem;
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
	}

	header {
		text-align: center;
		margin-bottom: 0.5rem;
	}

	h1 {
		font-size: 2.5rem;
		font-weight: 800;
		color: var(--accent);
		letter-spacing: -1px;
	}

	.tagline {
		color: var(--text-muted);
		font-size: 0.95rem;
		margin-top: 0.25rem;
	}

	.card {
		background: var(--surface);
		border-radius: var(--radius);
		padding: 1.5rem;
		border: 1px solid #1e2a40;
	}

	h2 {
		font-size: 1.1rem;
		margin-bottom: 1rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.form-row {
		display: flex;
		gap: 0.75rem;
	}

	.form-row input {
		flex: 1;
	}

	.error {
		color: var(--danger);
		font-size: 0.85rem;
		margin-top: 0.5rem;
	}

	.status-connected {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.status-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 1rem;
	}

	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		display: inline-block;
	}

	.dot.online {
		background: var(--success);
		box-shadow: 0 0 6px var(--success);
	}

	.node-details {
		display: flex;
		flex-direction: column;
		gap: 0.2rem;
	}

	.label {
		font-size: 0.75rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	code {
		font-size: 0.9rem;
		background: var(--surface-2);
		padding: 0.2rem 0.5rem;
		border-radius: 4px;
		display: inline-block;
	}

	.muted {
		color: var(--text-muted);
		font-size: 0.75rem;
	}
</style>
