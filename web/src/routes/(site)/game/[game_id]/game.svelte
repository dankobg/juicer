<script lang="ts">
	import { page } from '$app/state';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import { ws } from '$lib/state/ws-state.svelte';
	import type { PageProps } from './$types';

	let { data, params }: PageProps = $props();

	$effect(() => {
		if (ws.readyState !== WebSocket.OPEN) {
			const params = new URLSearchParams();
			params.set('path', page.url.pathname);
			ws.connect(params);
		}

		ws.onOpen = (event: Event) => {
			console.debug('ws open:', event);
		};

		ws.onClose = (event: CloseEvent) => {
			console.debug(`ws closed: code: ${event.code}, reason: ${event.reason}, wasClean: ${event.wasClean}`);
		};

		ws.onError = (event: Event) => {
			console.debug('ws error:', event);
		};

		ws.onMessage = (event: MessageEvent) => {
			gameManager.handleWebsocketMessage(event);
		};

		return () => {
			ws.close();
			gameManager.cancelSeekGame();
		};
	});
</script>

<h1>Game page {params.game_id}</h1>

<button class="mt-4 bg-red-800 p-2 text-white" onclick={gameManager.echo}>SEND ECHO</button>

<button class="mt-4 bg-cyan-700 p-2 text-white" onclick={gameManager.sendGameChat}>SEND GAME CHAT</button>
