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
		Error,
		ClientConnected,
		ClientDisconnected,
		HubInfo,
		SeekingCount,
		MatchFound,
		MatchRejoined,
		GameFinished,
		Move,
	} from '$lib/gen/juicer_pb';
	import { chatMessages } from '$lib/gamechat/messages';
	import type { JuicerBoard, MoveCancelEvent, MoveStartEvent, MoveFinishEvent } from '@dankop/juicer-board';
	import { BLACK, WHITE } from '$lib/shared/shared';

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

	function synchronizeClocks(whiteTime?: number, blackTime?: number) {
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
		if ((side === WHITE && !isWhiteTurn) || (side === BLACK && isWhiteTurn)) {
			console.log('not your turn');
			event.preventDefault();
			return;
		}

		const uci = event.data.srcCoord + event.data.destCoord;
		ws.send(new Message({ event: { case: 'playMoveUci', value: new PlayMoveUCI({ move: uci }) } }));
		ply++;
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

	function handleErr(errMsg: Error) {
		wsErr = errMsg.message;
	}

	function handleEcho(echoMsg: Echo) {}

	function handleChat(chatMsg: Chat) {
		chatMessages.update(msgs => [...msgs, { own: false, text: chatMsg.message }]);
	}

	function handleClientConnected(clientConnectedMsg: ClientConnected) {
		stopOpponentDisconnectTimer();
	}
	function handleClientDisconnected(clientDisconnectedMsg: ClientDisconnected) {
		startOpponentDisconnectTimer();
	}

	function handleHubInfo(hubInfoMsg: HubInfo) {
		lobbyCount = hubInfoMsg.lobby;
		roomsCount = hubInfoMsg.rooms;
		playingCount = hubInfoMsg.playing;
	}

	function handleSeekingCount(seekingCountMsg: SeekingCount) {
		seekingCount = seekingCountMsg.count;
	}

	function handleMatchFound(matchFoundMsg: MatchFound) {
		state = 'playing';
		roomId = matchFoundMsg.roomId;
		gameId = matchFoundMsg.gameId;
		side = matchFoundMsg.color;
		drawOffered = false;
	}

	function handleMatchRejoined(matchRejoinedMsg: MatchRejoined) {
		state = 'playing';
		const clocks = matchRejoinedMsg.clocks;
		synchronizeClocks(clocks?.white, clocks?.black);
		roomId = matchRejoinedMsg.roomId;
		gameId = matchRejoinedMsg.gameId;
		side = matchRejoinedMsg.color;
		drawOffered = false;
		fen = matchRejoinedMsg.fen;
		ply = matchRejoinedMsg.ply;
		legalMoves = matchRejoinedMsg.legalMoves;
	}

	function handleGameFinished(gameFinishedMsg: GameFinished) {
		state = 'idle';
		roomId = '';
		gameId = '';
		gameResult = gameFinishedMsg.result;
		gameStatus = gameFinishedMsg.status;
		drawOffered = false;
		showpOpponentDisconnectTimer = false;
	}

	function handleOfferDraw(offerDrawMsg: OfferDraw) {
		drawOffered = true;
	}

	function handleMove(moveMsg: Move) {
		const clocks = moveMsg.clocks;
		synchronizeClocks(clocks?.white, clocks?.black);
		uci = moveMsg.uci;
		san = moveMsg.san;
		lan = moveMsg.lan;
		ply = moveMsg.ply;
		legalMoves = moveMsg.legalMoves;
		const src = uci.slice(0, 2) as any;
		const dest = uci.slice(2, 4) as any;
		board.movePiece(src, dest);
	}

	onMount(() => {
		if (board) {
			addBoardEventListeners();
		}

		ws.connect();

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
			try {
				const msg = Message.fromJsonString(event.data);
				console.debug('ws recv:', msg.event.case, '-', msg.event.value);

				switch (msg.event.case) {
					case 'error':
						handleErr(msg.event.value);
						break;
					case 'echo':
						handleEcho(msg.event.value);
						break;
					case 'chat':
						handleChat(msg.event.value);
						break;
					case 'clientConnected':
						handleClientConnected(msg.event.value);
						break;
					case 'clientDisconnected':
						handleClientDisconnected(msg.event.value);
						break;
					case 'hubInfo':
						handleHubInfo(msg.event.value);
						break;
					case 'seekingCount':
						handleSeekingCount(msg.event.value);
						break;
					case 'matchFound':
						handleMatchFound(msg.event.value);
						break;
					case 'matchRejoined':
						handleMatchRejoined(msg.event.value);
						break;
					case 'gameFinished':
						handleGameFinished(msg.event.value);
						break;
					case 'offerDraw':
						handleOfferDraw(msg.event.value);
						break;
					case 'move':
						handleMove(msg.event.value);
						break;
					default:
						console.log('unknown message', msg.event.case, msg.event.value);
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

	$: isWhiteTurn = ply % 2 === 0;

	$: if (ply >= 2 && !gameClockIntervalId) {
		console.log(' i am called ');

		gameClockIntervalId = setInterval(() => {
			if (isWhiteTurn) {
				whiteRemainingTime -= 1;
			} else {
				blackRemainingTime -= 1;
			}
		}, 1000);
	}
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
	<div>
		Side: <strong>{side}</strong>
	</div>
	<div>
		Current turn: <strong>{isWhiteTurn ? 'w' : 'b'}</strong>
	</div>
	<div>
		white clock:<strong>{whiteRemainingTime.toFixed(0)}</strong>
	</div>
	<div>
		black clock: <strong>{blackRemainingTime.toFixed(0)}</strong>
	</div>
{/if}
