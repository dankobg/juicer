<script lang="ts">
	import { onMount } from 'svelte';

	let outerSize = '35rem';

	let ws: WebSocket;
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
		};

		return () => {
			ws.close();
		};
	});
</script>

<p>clients in lobby: {lobbyClients}</p>
<p>total rooms: {activeRooms}</p>

<button
	class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded mb-2"
	on:click={() => {
		const msg = JSON.stringify({ t: 'stub_msg', d: { id: 69, nice: true } });
		ws.send(msg);
	}}>SEND STUB MSG</button
>

<div style="width: {outerSize}; height: {outerSize};">
	<juicer-board fen="start" coords="inside" files="start" ranks="start" interactive show-ghost></juicer-board>
</div>
