<script lang="ts">
	export let connected: boolean = false;
	export let peerCount: number = 0;
</script>

<footer class="status-bar">
	<span class="indicator" class:online={connected} class:offline={!connected}>
		{connected ? '● Connected' : '○ Disconnected'}
	</span>
	{#if connected}
		<span class="sep">|</span>
		<span>{peerCount} peer{peerCount !== 1 ? 's' : ''} online</span>
	{/if}
	<span class="sep">|</span>
	<span class="version">Rimesh v0.1.0</span>
</footer>

<style>
	.status-bar {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		background: var(--surface);
		border-top: 1px solid #1e2a40;
		padding: 0.35rem 1.5rem;
		font-size: 0.78rem;
		color: var(--text-muted);
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.indicator {
		font-weight: 600;
	}

	.indicator.online {
		color: var(--success);
	}

	.indicator.offline {
		color: var(--text-muted);
	}

	.sep {
		opacity: 0.3;
	}
</style>
