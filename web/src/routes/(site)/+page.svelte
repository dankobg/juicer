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
	import { BLACK, WHITE, type Color } from '$lib/shared/shared';
	import { goto } from '$app/navigation';

	let outerSize = '35rem';

	let ws: JuicerWS = new JuicerWS();
	let wsErr: string = '';

	let board: JuicerBoard;

	let opponentDisconnectIntervalId: NodeJS.Timeout | undefined;
	let opponentDisconnectTimer = 0;
	let showOpponentDisconnectTimer = false;

	let gameClockIntervalId: NodeJS.Timeout | undefined;
	let whiteRemainingTime: number = 0;
	let blackRemainingTime: number = 0;
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
	let side: Color = 'w';

	let state: 'idle' | 'seeking' | 'playing' = 'idle';
	let gameResult: string = '';
	let gameStatus: string = '';

	let drawOffered: boolean = false;

	let promotionPopoverElm: HTMLDivElement;
	let promotionPieceSymbol = '';
	let promotionSrcDest = '';

	const CASTLE_MOVES_PAIR = new Map<string, string>([
		['e1g1', 'h1f1'],
		['e1c1', 'a1d1'],
		['e8g8', 'h8f8'],
		['e8c8', 'a8d8'],
	]);

	const promos = ['q', 'r', 'n', 'b'];

	function getRookCastleMove(uci: string): string | undefined {
		return CASTLE_MOVES_PAIR.get(uci);
	}

	// move is src-dest
	function isPromotionMove(move: string, legalMoves: string[]): boolean {
		if (move.length !== 4) {
			return false;
		}
		return legalMoves.some(m => m.length === 5 && m.startsWith(move) && promos.includes(m[4]));
	}

	function isPromotionUciMove(uci: string): boolean {
		return uci.length === 5 && promos.includes(uci[4]);
	}

	function promotePiece(promotionPieceSymbol: string) {
		const src = promotionSrcDest.slice(0, 2) as any;
		const dest = promotionSrcDest.slice(2, 4) as any;
		board.removePiece(src, true);
		board.setPiece(dest, promotionPieceSymbol as any, true);
	}

	function sendMove(uci: string) {
		ws.send(new Message({ event: { case: 'playMoveUci', value: new PlayMoveUCI({ move: uci }) } }));
		ply++;
	}

	function onChatMessage(event: CustomEvent<{ text: string }>) {
		const msg = event.detail.text;
		if (!msg) {
			return;
		}
		ws.send(new Message({ event: { case: 'chat', value: new Chat({ message: msg }) } }));
	}

	function startOpponentDisconnectTimer() {
		stopOpponentDisconnectTimer();

		showOpponentDisconnectTimer = true;
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
		showOpponentDisconnectTimer = false;

		if (opponentDisconnectIntervalId) {
			clearTimeout(opponentDisconnectIntervalId);
		}
		opponentDisconnectTimer = 0;
	}

	function onBoardMoveStart(event: MoveStartEvent) {}

	function onBoardMoveCancel(event: MoveCancelEvent) {}

	function onBoardMoveFinish(event: MoveFinishEvent) {
		if ((side === WHITE && !isWhiteTurn) || (side === BLACK && isWhiteTurn)) {
			console.debug('not your turn');
			event.preventDefault();
			return;
		}

		const move = event.data.srcCoord + event.data.destCoord;

		const isPromo = isPromotionMove(move, legalMoves);
		if (isPromo) {
			promotionSrcDest = move;
			const dest = move.slice(2, 4);
			const promoSquareElm = board.shadowRoot?.querySelector(`juicer-square[coord='${dest}']`) ?? null;
			if (!promoSquareElm) {
				console.debug('no promotion square element');
				return;
			}
			const rect = promoSquareElm.getBoundingClientRect();
			promotionPopoverElm.style.left = `${rect.left}px`;
			promotionPopoverElm.style.top = `${rect.top}px`;
			promotionPopoverElm.showPopover();
			promotionPopoverElm.classList.remove('hidden');
			event.preventDefault();
			return;
		}

		if (!legalMoves.includes(move)) {
			console.debug('invalid move attempt:', move, legalMoves);
			event.preventDefault();
			return;
		}

		const rookMove = getRookCastleMove(move);
		if (rookMove) {
			const src = rookMove.slice(0, 2) as any;
			const dest = rookMove.slice(2, 4) as any;
			board.movePiece(src, dest, true);
		}

		ws.send(new Message({ event: { case: 'playMoveUci', value: new PlayMoveUCI({ move: move }) } }));
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
		const clocks = matchFoundMsg.clocks;
		synchronizeClocks(clocks?.white, clocks?.black);
		roomId = matchFoundMsg.roomId;
		gameId = matchFoundMsg.gameId;
		side = matchFoundMsg.color as Color;
		drawOffered = false;
		fen = matchFoundMsg.fen;
		ply = matchFoundMsg.ply;
		legalMoves = matchFoundMsg.legalMoves;
		// goto(`/game/${gameId}`);
	}

	function handleMatchRejoined(matchRejoinedMsg: MatchRejoined) {
		state = 'playing';
		const clocks = matchRejoinedMsg.clocks;
		synchronizeClocks(clocks?.white, clocks?.black);
		roomId = matchRejoinedMsg.roomId;
		gameId = matchRejoinedMsg.gameId;
		side = matchRejoinedMsg.color as Color;
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
		showOpponentDisconnectTimer = false;
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

		const isPromo = isPromotionUciMove(uci);
		if (isPromo) {
			promotionSrcDest = uci.slice(0, 4);
			promotionPieceSymbol = side === 'w' ? uci[4] : uci[4].toUpperCase();
			promotePiece(promotionPieceSymbol);
			return;
		}

		board.movePiece(src, dest);

		const rookMove = getRookCastleMove(uci);
		if (rookMove) {
			const src = rookMove.slice(0, 2) as any;
			const dest = rookMove.slice(2, 4) as any;
			board.movePiece(src, dest, true);
		}
	}

	function cleanupGame() {
		state = 'idle';
		gameResult = '';
		gameStatus = '';
		drawOffered = false;
		roomId = '';
		gameId = '';
		clearInterval(opponentDisconnectIntervalId);
		clearInterval(gameClockIntervalId);
	}

	function handlePromotionPiecePick(symbol: string) {
		promotionPieceSymbol = side === 'w' ? symbol.toUpperCase() : symbol;
		promotionPopoverElm.hidePopover();
		promotionPopoverElm.classList.add('hidden');
		promotePiece(promotionPieceSymbol);
		sendMove(promotionSrcDest + promotionPieceSymbol.toLowerCase());
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
						console.error('unknown message', msg.event.case, msg.event.value);
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

{#if showOpponentDisconnectTimer}
	<p>Opponent time to rejoin the game: {opponentDisconnectTimer.toFixed(2)}</p>
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
			cleanupGame();
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
				cleanupGame();
			}}
		>
			Accept draw</button
		>
	{/if}

	{#if state === 'playing'}
		<p>promotionPieceSymbol: {promotionPieceSymbol}</p>
		<p>Side: <strong>{side}</strong></p>
		<p>Current turn: <strong>{isWhiteTurn ? 'w' : 'b'}</strong></p>
		<p>white clock:<strong>{whiteRemainingTime.toFixed(2)}</strong></p>
		<p>black clock: <strong>{blackRemainingTime.toFixed(2)}</strong></p>
		<p>Legal moves: {legalMoves}</p>
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

<div
	bind:this={promotionPopoverElm}
	id="promotion"
	popover="manual"
	style="left:280px;right:324px;"
	class="absolute w-[70px] flex flex-col flex-wrap m-0 p-0 bg-orange-100 hidden"
>
	<button
		class="h-[70px] bg:blue-300 hover:bg-orange-300"
		popovertarget="promotion"
		popovertargetaction="hide"
		on:click={() => handlePromotionPiecePick('q')}>Q</button
	>
	<button
		class="h-[70px] bg:blue-300 hover:bg-orange-300"
		popovertarget="promotion"
		popovertargetaction="hide"
		on:click={() => handlePromotionPiecePick('r')}>R</button
	>
	<button
		class="h-[70px] bg:blue-300 hover:bg-orange-300"
		popovertarget="promotion"
		popovertargetaction="hide"
		on:click={() => handlePromotionPiecePick('n')}>N</button
	>
	<button
		class="h-[70px] bg:blue-300 hover:bg-orange-300"
		popovertarget="promotion"
		popovertargetaction="hide"
		on:click={() => handlePromotionPiecePick('b')}>B</button
	>
</div>
