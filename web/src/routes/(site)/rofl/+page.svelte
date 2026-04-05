<script lang="ts">
	import { ws } from '$lib/state/ws-state.svelte';
	import { create } from '@bufbuild/protobuf';
	import type { PageProps } from './$types';
	import { MessageSchema } from '$lib/gen/juicer_pb';
	import { page } from '$app/state';

	let { data }: PageProps = $props();

	$effect(() => {
		console.log('CALLED ROFL SHITTE');

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
			console.log('recv: ', event.data);
			// gameManager.handleWebsocketMessage(event);
		};
		return () => {
			ws.close();
			// gameManager.cancelSeekGame();
		};
	});
</script>

<h1>ROFL PAGE</h1>

<button
	class="w-fit rounded-lg bg-yellow-800 p-2"
	onclick={() => {
		const msg = create(MessageSchema, {
			event: {
				case: 'test',
				value: { message: 'picketine' }
			}
		});
		ws.send(msg);
	}}
>
	send ws
</button>

<a class="w-fit bg-purple-700" href="/">GOTO HOME</a>
