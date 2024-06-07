<script lang="ts">
	import { onMount } from 'svelte';

	let outerSize = '35rem';

	let ws: WebSocket;
	let wsErr = '';
	let lobbyClients = 0;
	let activeRooms = 0;

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
			const msg = JSON.parse(e.data);
			if (msg.t === 'clients_count') {
				lobbyClients = msg.d.lobby;
				activeRooms = msg.d.rooms;
			}
			if (msg.t === 'error') {
				wsErr = msg.d;
			}
		};

		return () => {
			ws.close();
		};
	});
</script>

{#if wsErr}error: {wsErr}{/if}
<p>clients in lobby: <strong>{lobbyClients}</strong></p>
<p>total rooms: <strong>{activeRooms}</strong></p>

<button
	class="bg-purple-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		const msg = JSON.stringify({ t: 'stub_msg', d: { id: 69, nice: true } });
		ws.send(msg);
	}}
>
	SEND STUB MSG</button
>

<button
	class="bg-green-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		const msg = JSON.stringify({ t: 'seek_game', d: { game_mode: 'blitz' } });
		ws.send(msg);
	}}
>
	SEEK GAME</button
>

<button
	class="bg-orange-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		const msg = JSON.stringify({ t: 'cancel_seek_game' });
		ws.send(msg);
	}}
>
	CANCEL SEEK GAME</button
>

<div style="width: {outerSize}; height: {outerSize};">
	<juicer-board fen="start" coords="inside" files="start" ranks="start" interactive show-ghost></juicer-board>
</div>
