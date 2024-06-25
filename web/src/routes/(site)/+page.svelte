<script lang="ts">
	import { JuicerWS } from '$lib/ws/ws';
	import { onMount } from 'svelte';
	import { Spinner } from 'flowbite-svelte';
	import GameChat from '$lib/gamechat/GameChat.svelte';
	import { AbortGame, AcceptDraw, CancelSeekGame, Chat, Echo, Message, OfferDraw, SeekGame } from '$lib/gen/juicer_pb';

	let outerSize = '35rem';

	let ws: JuicerWS = new JuicerWS();
	let wsErr = '';

	let lobbyCount = 0;
	let roomsCount = 0;
	let seekingCount = 0;
	let playingCount = 0;

	let roomId = '';
	let gameId = '';

	let state: 'idle' | 'seeking' | 'playing' = 'idle';
	let gameResult = '';
	let gameStatus = '';

	let drawOffered = false;

	function onChatMessage(event: CustomEvent<{ text: string }>) {
		const msg = event.detail.text;
		if (!msg) {
			return;
		}
		ws.send(new Message({ event: { case: 'chat', value: new Chat({ message: msg }) } }));
	}

	onMount(() => {
		ws.connect();

		ws.onmessage = (event: MessageEvent) => {
			try {
				const msg = Message.fromJsonString(event.data);

				switch (msg.event.case) {
					case 'error':
						wsErr = msg.event.value.message;
						break;
					case 'echo':
						console.log('echo resp:', msg.event.value);
						break;
					case 'clientConnected':
						console.log('client joined:', msg.event.value.id);
						break;
					case 'clientDisconnected':
						console.log('client left:', msg.event.value.id);
						break;
					case 'hubInfo':
						lobbyCount = msg.event.value.lobby;
						roomsCount = msg.event.value.rooms;
						playingCount = msg.event.value.playing;
						break;
					case 'seekingCount':
						seekingCount = msg.event.value.count;
						break;
					case 'matchFound':
						state = 'playing';
						roomId = msg.event.value.roomId;
						gameId = msg.event.value.gameId;
						break;
					case 'gameFinished':
						state = 'idle';
						roomId = '';
						gameId = '';
						gameResult = msg.event.value.result;
						gameStatus = msg.event.value.status;
						break;
					case 'offerDraw':
						drawOffered = true;
						break;
					default:
						console.log('unkown message', msg.event.case, msg.event.value);
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

{#if state !== 'playing'}
	<p>In lobby: <strong>{lobbyCount}</strong></p>
	<p>Rooms: <strong>{roomsCount}</strong></p>
	<p>Playing: <strong>{playingCount}</strong></p>
{/if}

{#if gameResult}
	<p>Game result: <strong>{gameResult}</strong></p>
{/if}
{#if gameStatus}
	<p>Game status: <strong>{gameStatus}</strong></p>
{/if}

{#if gameId}
	<p>Game ID: <strong>{gameId}</strong></p>
{/if}
{#if roomId}
	<p>Room ID: <strong>{roomId}</strong></p>
{/if}

{#if state === 'idle'}
	<button
		class="bg-purple-500 text-white py-2 px-4 rounded mb-2"
		on:click={() => {
			ws.send(new Message({ event: { case: 'echo', value: new Echo({ message: 'hello bozo' }) } }));
		}}
	>
		SEND ECHO MSG</button
	>

	<button
		class="bg-green-500 text-white py-2 px-4 rounded mb-2"
		on:click={() => {
			ws.send(new Message({ event: { case: 'seekGame', value: new SeekGame({ gameMode: 'blitz' }) } }));
			state = 'seeking';
			gameResult = '';
			gameStatus = '';
		}}
	>
		SEEK GAME</button
	>
{/if}

{#if state === 'seeking'}
	<button
		class="bg-orange-500 text-white py-2 px-4 rounded mb-2"
		on:click={() => {
			ws.send(new Message({ event: { case: 'cancelSeekGame', value: new CancelSeekGame() } }));
			state = 'idle';
			gameResult = '';
			gameStatus = '';
		}}
	>
		CANCEL SEEK GAME</button
	>

	<div class="flex flex-col gap-3 m-4 p-8 border-solid border-2 border-sky-500">
		<Spinner />
		<p>Players seeking: <strong>{seekingCount}</strong></p>
		<p class="">Searching for game...</p>
	</div>
{/if}

{#if state === 'playing'}
	<button
		class="bg-orange-500 text-white py-2 px-4 rounded mb-2"
		on:click={() => {
			ws.send(new Message({ event: { case: 'abortGame', value: new AbortGame() } }));
			state = 'idle';
			gameResult = '';
			gameStatus = '';
		}}
	>
		Abort game</button
	>

	<button
		class="bg-blue-500 text-white py-2 px-4 rounded mb-2"
		on:click={() => {
			ws.send(new Message({ event: { case: 'offerDraw', value: new OfferDraw() } }));
		}}
	>
		Offer draw</button
	>

	{#if drawOffered}
		<button
			class="bg-green-500 text-white py-2 px-4 rounded mb-2"
			on:click={() => {
				ws.send(new Message({ event: { case: 'acceptDraw', value: new AcceptDraw() } }));
			}}
		>
			Accept draw</button
		>
	{/if}

	<div style="display:flex;flex-wrap:wrap;gap:1rem;">
		<div style="width: {outerSize}; height: {outerSize};">
			<juicer-board fen="start" coords="inside" files="start" ranks="start" interactive show-ghost></juicer-board>
		</div>

		<GameChat on:message={onChatMessage} />
	</div>
{/if}
