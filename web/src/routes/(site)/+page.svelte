<script lang="ts">
	import { onMount } from 'svelte';

	let outerSize = '35rem';

	let ws: WebSocket;
	let stubMsg = { type: 'stub_msg', data: { id: 69, nice: true } };
	let payload = JSON.stringify(stubMsg);

	onMount(() => {
		ws = new WebSocket('wss://juicer-dev.xyz/ws');

		ws.onopen = () => {
			console.log('opened');
		};
		ws.onclose = () => {
			console.log('closed');
		};
		ws.onerror = e => {
			console.log('errored: ', e);
		};
		ws.onmessage = e => {
			console.log(`recv: ${e.data}`);
		};

		return () => {
			ws.close();
		};
	});
</script>

<button
	class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded mb-2"
	on:click={() => {
		ws.send(payload);
	}}>SEND STUB MSG</button
>

<div style="width: {outerSize}; height: {outerSize};">
	<juicer-board fen="start" coords="inside" files="start" ranks="start" interactive show-ghost></juicer-board>
</div>
