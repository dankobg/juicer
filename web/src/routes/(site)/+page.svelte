<script lang="ts">
	import { onMount } from 'svelte';
	import { JuicerWS } from '$lib/ws/ws';
	import { Spinner } from 'flowbite-svelte';
	import GameChat from '$lib/gamechat/GameChat.svelte';
	import {
		AbortGame,
		AcceptDraw,
		CancelSeekGame,
		Chat,
		Echo,
		Message,
		OfferDraw,
		PlayMoveUCI,
		SeekGame,
	} from '$lib/gen/juicer_pb';
	import { chatMessages } from '$lib/gamechat/messages';
	import type { JuicerBoard, MoveCancelEvent, MoveStartEvent, MoveFinishEvent } from '@dankop/juicer-board';
	import { BLACK, WHITE } from 'chess.js';

	let outerSize = '35rem';

	let ws: JuicerWS = new JuicerWS();
	let wsErr: string = '';

	let board: JuicerBoard;

	let opponentDisconnectIntervalId: NodeJS.Timeout | undefined;
	let opponentDisconnectTimer = 0;
	let showpOpponentDisconnectTimer = false;

	let gameClockIntervalId: NodeJS.Timeout | undefined;
	let whiteRemainingTime: number = 300;
	let blackRemainingTime: number = 300;
	let fen: string = 'start';
	let uci: string = '';
	let san: string = '';
	let lan: string = '';
	let ply: number = 0;
	let legalMoves: string[] = [];

	let lobbyCount = 0;
	let roomsCount = 0;
	let seekingCount = 0;
	let playingCount = 0;

	let roomId: string = '';
	let gameId: string = '';
	let side = 'w';

	let state: 'idle' | 'seeking' | 'playing' = 'idle';
	let gameResult: string = '';
	let gameStatus: string = '';

	let drawOffered: boolean = false;

	function onChatMessage(event: CustomEvent<{ text: string }>) {
		const msg = event.detail.text;
		if (!msg) {
			return;
		}
		ws.send(new Message({ event: { case: 'chat', value: new Chat({ message: msg }) } }));
	}

	function startOpponentDisconnectTimer() {
		stopOpponentDisconnectTimer();

		showpOpponentDisconnectTimer = true;
		opponentDisconnectTimer = 10;

		opponentDisconnectIntervalId = setInterval(() => {
			if (opponentDisconnectTimer > 0) {
				opponentDisconnectTimer--;
			} else {
				clearInterval(opponentDisconnectIntervalId);
				opponentDisconnectIntervalId = undefined;
			}
		}, 1000);
	}

	function startPlayerClocks() {
		gameClockIntervalId = setInterval(() => {
			if (ply >= 2) {
				if (ply % 2 === 0) {
					whiteRemainingTime -= 1;
				} else {
					blackRemainingTime -= 1;
				}
			}
		}, 1000);
	}

	function updatePlayerClocks(whiteTime?: number, blackTime?: number) {
		if (whiteTime) {
			whiteRemainingTime = whiteTime;
		}
		if (blackTime) {
			blackRemainingTime = blackTime;
		}
	}

	function stopOpponentDisconnectTimer() {
		showpOpponentDisconnectTimer = false;

		if (opponentDisconnectIntervalId) {
			clearTimeout(opponentDisconnectIntervalId);
		}
		opponentDisconnectTimer = 0;
	}

	function onBoardMoveStart(event: MoveStartEvent) {}

	function onBoardMoveCancel(event: MoveCancelEvent) {}

	function onBoardMoveFinish(event: MoveFinishEvent) {
		if ((side === WHITE && ply % 2 !== 0) || (side === BLACK && ply % 2 === 0)) {
			console.log('not your turn');
			event.preventDefault();
			return;
		}

		const uci = event.data.srcCoord + event.data.destCoord;
		console.log('uci:', uci);
		ws.send(new Message({ event: { case: 'playMoveUci', value: new PlayMoveUCI({ move: uci }) } }));
	}

	function addBoardEventListeners() {
		if (board) {
			board.addEventListener('movestart', onBoardMoveStart);
			board.addEventListener('movecancel', onBoardMoveCancel);
			board.addEventListener('movefinish', onBoardMoveFinish);
		}
	}

	function removeBoardEventListeners() {
		if (board) {
			board.removeEventListener('movestart', onBoardMoveStart);
			board.removeEventListener('movecancel', onBoardMoveCancel);
			board.removeEventListener('movefinish', onBoardMoveFinish);
		}
	}

	onMount(() => {
		if (board) {
			addBoardEventListeners();
		}

		ws.connect();

		ws.onOpen = (event: Event) => {
			console.debug(`ws open: ${event}`);
		};

		ws.onClose = (event: CloseEvent) => {
			console.debug(`ws closed: code:${event.code}, reason: ${event.reason}, wasClean: ${event.wasClean}`);
		};

		ws.onError = (event: Event) => {
			console.debug(`ws error: ${event}`);
		};

		ws.onMessage = (event: MessageEvent) => {
			try {
				const msg = Message.fromJsonString(event.data);
				console.debug(`ws recv: ${msg.event.case} - ${msg.event.value}`);

				switch (msg.event.case) {
					case 'error':
						wsErr = msg.event.value.message;
						break;
					case 'echo':
						break;
					case 'chat':
						const chatMsg = msg.event.value.message;
						chatMessages.update(msgs => [...msgs, { own: false, text: chatMsg }]);
						break;
					case 'clientConnected':
						stopOpponentDisconnectTimer();
						break;
					case 'clientDisconnected':
						startOpponentDisconnectTimer();
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
						side = msg.event.value.color;
						drawOffered = false;
						startPlayerClocks();
						break;
					case 'gameFinished':
						state = 'idle';
						roomId = '';
						gameId = '';
						gameResult = msg.event.value.result;
						gameStatus = msg.event.value.status;
						drawOffered = false;
						showpOpponentDisconnectTimer = false;
						break;
					case 'offerDraw':
						drawOffered = true;
						break;
					case 'move':
						const clocks = msg.event.value.clocks;
						updatePlayerClocks(clocks?.white, clocks?.black);
						uci = msg.event.value.uci;
						san = msg.event.value.san;
						lan = msg.event.value.lan;
						fen = msg.event.value.fen;
						ply = msg.event.value.ply;
						legalMoves = msg.event.value.legalMoves;
						const src = uci.slice(0, 2) as any;
						const dest = uci.slice(2, 4) as any;
						board.movePiece(src, dest);
						console.log({ w: clocks?.white, b: clocks?.black, uci, san, lan, fen, legalMoves });
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
			removeBoardEventListeners();
			clearInterval(opponentDisconnectIntervalId);
			clearInterval(gameClockIntervalId);
		};
	});

	$: board && addBoardEventListeners();
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

{#if showpOpponentDisconnectTimer}
	<p>Opponent time to rejoin the game: {opponentDisconnectTimer}</p>
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
			drawOffered = false;
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
				drawOffered = false;
			}}
		>
			Accept draw</button
		>
	{/if}

	<div style="display:flex;flex-wrap:wrap;gap:1rem;">
		<div style="width: {outerSize}; height: {outerSize};">
			<juicer-board
				bind:this={board}
				orientation={side}
				{fen}
				coords="inside"
				files="start"
				ranks="start"
				interactive
				show-ghost
			></juicer-board>
		</div>

		<GameChat on:message={onChatMessage} />
	</div>
{/if}

{#if state === 'playing'}
	Side: <h1>{side}</h1>
	white clock:
	<h3>{whiteRemainingTime.toFixed(0)}</h3>
	black clock:
	<h3>{blackRemainingTime.toFixed(0)}</h3>
{/if}
