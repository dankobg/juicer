<script lang="ts">
	import { page } from '$app/state';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import { ws } from '$lib/state/ws-state.svelte';
	import { create } from '@bufbuild/protobuf';
	import type { PageProps } from './$types';
	import { MessageSchema } from '$lib/gen/juicer_pb';

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

	let ack = $state(0);
	let mov = $state('');
</script>

<h1>Game page {params.game_id}</h1>

<button class="mt-4 bg-red-800 p-2 text-white" onclick={gameManager.echo}>SEND ECHO</button>

<button class="mt-4 bg-cyan-700 p-2 text-white" onclick={gameManager.sendGameChat}>SEND GAME CHAT</button>

<div>
	<input placeholder="UCI" bind:value={mov} />
	<button
		class="mt-4 bg-cyan-700 p-2 text-white"
		onclick={() => {
			let rofl = Number(window.location.pathname.split('/').at(-1));
			ack = ack + 1;

			ws.send(create(MessageSchema, { event: { case: 'playMoveUci', value: { gameId: rofl, uci: mov, ack: ack } } }));
		}}
	>
		SEND MOVE
	</button>
</div>
