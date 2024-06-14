<script lang="ts">
	import { JuicerWS } from '$lib/ws/ws';
	import { onMount } from 'svelte';
	import { Spinner } from 'flowbite-svelte';
	import Chat from '$lib/chat/Chat.svelte';

	let outerSize = '35rem';

	let ws: JuicerWS = new JuicerWS();
	let wsErr = '';

	let lobbyCount = 0;
	let roomsCount = 0;
	let seekingCount = 0;
	let inGameCount = 0;

	let roomId = '';
	let gameId = '';

	let state: 'idle' | 'seeking' | 'playing' = 'idle';

	onMount(() => {
		ws.connect();

		ws.onmessage = (event: MessageEvent) => {
			try {
				const msg = JSON.parse(event.data);
				console.debug(msg.t, msg.d);

				switch (msg.t) {
					case 'error':
						wsErr = msg.d;
						break;

					case 'clients_count':
						lobbyCount = msg.d.lobby;
						roomsCount = msg.d.rooms;
						inGameCount = msg.d.in_game;
						break;

					case 'seeking_count':
						seekingCount = msg.d;
						break;

					case 'match_found':
						state = 'playing';
						roomId = msg.d.room_id;
						gameId = msg.d.game_id;
						break;

					default:
						break;
				}
			} catch (error) {
				console.error('json parse msg data', error);
			}
		};

		return () => {
			ws.close();
		};
	});
</script>

{#if wsErr}
	err: {wsErr}
{/if}
<p>In lobby: <strong>{lobbyCount}</strong></p>
<p>rooms: <strong>{roomsCount}</strong></p>
<p>In game: <strong>{inGameCount}</strong></p>
<p>seeking game: <strong>{seekingCount}</strong></p>

{#if gameId}
	<p>Game ID: <strong>{gameId}</strong></p>
{/if}
{#if roomId}
	<p>Room ID: <strong>{roomId}</strong></p>
{/if}

<button
	class="bg-purple-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		ws.send(JSON.stringify({ t: 'echo', d: 'hello bozo' }));
	}}
>
	SEND ECHO MSG</button
>

<button
	class="bg-green-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		ws.send(JSON.stringify({ t: 'seek_game', d: { game_mode: 'blitz' } }));
		state = 'seeking';
	}}
>
	SEEK GAME</button
>

<button
	class="bg-orange-500 text-white py-2 px-4 rounded mb-2"
	on:click={() => {
		ws.send(JSON.stringify({ t: 'cancel_seek_game' }));
		state = 'idle';
	}}
>
	CANCEL SEEK GAME</button
>

{#if state === 'seeking'}
	<div class="flex gap-3 m-4">
		<Spinner />
		<p class="">Searching for game...</p>
	</div>
{/if}

{#if state === 'playing'}
	<div style="display:flex;flex-wrap:wrap;gap:1rem;">
		<div style="width: {outerSize}; height: {outerSize};">
			<juicer-board fen="start" coords="inside" files="start" ranks="start" interactive show-ghost></juicer-board>
		</div>

		<Chat />
	</div>
{/if}
