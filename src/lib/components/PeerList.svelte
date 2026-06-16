<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let peers: { id: string; alias: string; virtual_ip: string; public_ip: string }[] = [];

	const dispatch = createEventDispatcher();
</script>

<div class="peer-list">
	<div class="header-row">
		<h2>Peers ({peers.length})</h2>
		<button class="secondary small" on:click={() => dispatch('refresh')}>↻ Refresh</button>
	</div>

	{#if peers.length === 0}
		<p class="empty">No peers found. Share your Virtual IP with friends to get started.</p>
	{:else}
		<ul>
			{#each peers as peer (peer.id)}
				<li class="peer-item">
					<span class="dot online"></span>
					<div class="peer-info">
						<strong>{peer.alias}</strong>
						<span class="peer-ips">
							<code>{peer.virtual_ip}</code>
							{#if peer.public_ip}
								<span class="muted">({peer.public_ip})</span>
							{/if}
						</span>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.peer-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.header-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	h2 {
		font-size: 1.1rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.small {
		padding: 0.3rem 0.75rem;
		font-size: 0.8rem;
	}

	.empty {
		color: var(--text-muted);
		font-size: 0.9rem;
		text-align: center;
		padding: 1.5rem 0;
	}

	ul {
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.peer-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem;
		background: var(--surface-2);
		border-radius: var(--radius);
	}

	.peer-info {
		display: flex;
		flex-direction: column;
		gap: 0.15rem;
	}

	.peer-ips {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.85rem;
	}

	.dot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		flex-shrink: 0;
	}

	.dot.online {
		background: var(--success);
		box-shadow: 0 0 5px var(--success);
	}

	code {
		font-size: 0.85rem;
		background: rgba(0, 0, 0, 0.3);
		padding: 0.15rem 0.4rem;
		border-radius: 4px;
	}

	.muted {
		color: var(--text-muted);
	}
</style>
