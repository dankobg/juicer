import { timestampNow, timestampDate, type Timestamp } from '@bufbuild/protobuf/wkt';
import { create, fromJsonString } from '@bufbuild/protobuf';
import type { JuicerBoard, Coord, PieceFenSymbol } from '@dankop/juicer-board';
import {
	Color,
	GameResult,
	GameResultStatus,
	GameState,
	GameTimeControlSchema,
	MessageSchema,
	type ClientConnected,
	type ClientDisconnected,
	type GameAbort,
	type GameAcceptDraw,
	type GameChatReceive,
	type GameChatRetrieve,
	type GameDeclineDraw,
	type GameFinished,
	type GameOfferDraw,
	type GameResign,
	type GameTimeControl,
	type HistoryMoveInfo,
	type HubInfo,
	type MatchFound,
	type OpponentInfo,
	type Problem,
	type ReceiveMove
} from '$lib/gen/juicer_pb';
import { goto } from '$app/navigation';
import { Clock } from '$lib/state/clock.svelte';
import { Timer } from '$lib/state/timer.svelte';
import { ws } from './ws-state.svelte';
import { PersistedState } from 'runed';
import { soundManager } from './sound-manager.svelte';

export type ChatMsg = { text: string; received?: boolean };

const RECONNECT_TIMEOUT_MS = 15_000;
const FIRST_MOVE_TIMEOUT_MS = 10_000;

class GameManager {
	chatDialogElm: HTMLDialogElement | null = null;
	chatDialogOpen: boolean = $state(false);
	chatLastId = new PersistedState<string>('chat-last-id', '');
	chatGameId = new PersistedState<string>('chat-game-id', '');
	chatMessages = new PersistedState<ChatMsg[]>('chat-messages', []);
	board!: JuicerBoard;
	wsError = $state('');
	uiState: 'idle' | 'seeking' | 'playing' = $state('idle');
	lobbyCount = $state(0);
	seekingCount = $state(0);
	playingCount = $state(0);
	gameId: string = $state('');
	clientId: string = $state('');
	fen: string = $state('start');
	color: Color = $state(Color.WHITE);
	orientation: 'w' | 'b' = $state('w');
	opponentInfo: OpponentInfo | null = $state(null);
	uci: string = $state('');
	lan: string = $state('');
	san: string = $state('');
	ply: number = $state(0);
	legalMoves: string[] = $state([]);
	historyMovesInfo: HistoryMoveInfo[] = $state([]);
	historyIndex: number = $state(0);
	startTime?: Timestamp = $state();
	reconnectTimeoutMs: number = $state(RECONNECT_TIMEOUT_MS);
	firstMoveTimeoutMs: number = $state(FIRST_MOVE_TIMEOUT_MS);
	clock?: Clock = new Clock({ whiteTimeMs: 0, blackTimeMs: 0, showPreciseFn: ms => ms < 10_000 });
	opponentReconnectTimer: Timer = new Timer(RECONNECT_TIMEOUT_MS);
	whiteFirstMoveTimer: Timer = new Timer(this.firstMoveTimeoutMs);
	blackFirstMoveTimer: Timer = new Timer(this.firstMoveTimeoutMs);
	showOpponentReconnectTimer: boolean = $state(false);
	showWhiteFirstMoveTimer: boolean = $state(false);
	showBlackFirstMoveTimer: boolean = $state(false);
	whiteFirstMoveTimerStartTime = new PersistedState<{ val: number | null }>('juicer-white-first-move-timer-start', {
		val: null
	});
	blackFirstMoveTimerStartTime = new PersistedState<{ val: number | null }>('juicer-black-first-move-timer-start', {
		val: null
	});
	poolLast = new PersistedState<{ val: ClockControl | null }>('juicer-pool-last', {
		val: null
	});
	gameTimeControl?: GameTimeControl = $state();
	gameResult: GameResult = $state(GameResult.UNSPECIFIED);
	gameResultStatus: GameResultStatus = $state(GameResultStatus.UNSPECIFIED);
	gameState: GameState = $state(GameState.UNSPECIFIED);
	promotionPieceSymbol = $state('');
	promotionSrcDest = $state('');
	check: boolean = $derived(this.san.includes('+'));
	checkmate: boolean = $derived(this.san.includes('#'));
	hasIncrement: boolean = $derived(
		this.gameTimeControl?.increment?.seconds !== 0n && this.gameTimeControl?.increment?.nanos !== 0
	);
	opponentColor: Color = $derived(this.color === Color.WHITE ? Color.BLACK : Color.WHITE);
	ownFirstMoveTimer = $derived(this.color === Color.WHITE ? this.whiteFirstMoveTimer : this.blackFirstMoveTimer);
	opponentFirstMoveTimer = $derived(this.color === Color.WHITE ? this.blackFirstMoveTimer : this.whiteFirstMoveTimer);
	showOwnFirstMoveTimer = $derived(
		(this.color === Color.WHITE && this.showWhiteFirstMoveTimer) ||
			(this.color === Color.BLACK && this.showBlackFirstMoveTimer)
	);
	showOpponentFirstMoveTimer = $derived(
		(this.color === Color.WHITE && this.showBlackFirstMoveTimer) ||
			(this.color === Color.BLACK && this.showWhiteFirstMoveTimer)
	);
	ownGameTimer = $derived(this.color === Color.WHITE ? this.clock?.white : this.clock?.black);
	opponentGameTimer = $derived(this.color === Color.WHITE ? this.clock?.black : this.clock?.white);
	isWhiteTurn = $derived(this.ply % 2 === 0);
	currentTurn = $derived(this.isWhiteTurn ? Color.WHITE : Color.BLACK);
	hasActiveTurn = $derived(
		(this.gameState === GameState.IN_PROGRESS && this.color === Color.WHITE && this.isWhiteTurn) ||
			(this.color === Color.BLACK && !this.isWhiteTurn)
	);
	opponentHasActiveTurn = $derived(
		(this.gameState === GameState.IN_PROGRESS && this.opponentColor === Color.WHITE && this.isWhiteTurn) ||
			(this.opponentColor === Color.BLACK && !this.isWhiteTurn)
	);
	opponenOfferedDraw: boolean = $state(false);
	showOfferDraw = $derived(this.gameState === GameState.IN_PROGRESS && gameManager.ply >= 2);
	showResign = $derived(this.gameState === GameState.IN_PROGRESS && gameManager.ply >= 2);
	showAbort = $derived(
		this.gameState === GameState.IN_PROGRESS &&
			((gameManager.color === Color.WHITE && gameManager.ply === 0) ||
				(gameManager.color === Color.BLACK && gameManager.ply <= 1))
	);
	showGameResultMessage = $derived(this.gameState === GameState.FINISHED || this.gameState === GameState.INTERRUPTED);
	gameResultMessage = $derived.by(() => {
		if (this.gameResult === GameResult.UNSPECIFIED) {
			return '';
		}
		let msg = '';
		switch (this.gameResult) {
			case GameResult.DRAW:
				msg = 'Draw';
				break;
			case GameResult.WHITE_WON:
				msg = 'White won';
				break;
			case GameResult.BLACK_WON:
				msg = 'Black won';
				break;
			case GameResult.INTERRUPTED:
				msg = 'Interrupted';
				break;
			default:
				break;
		}
		switch (this.gameResultStatus) {
			case GameResultStatus.CHECKMATE:
				msg += ' by checkmate';
				break;
			case GameResultStatus.INSUFFICIENT_MATERIAL:
				msg += ' by insufficient material';
				break;
			case GameResultStatus.THREEFOLD_REPETITION:
				msg += ' by threefold repetition';
				break;
			case GameResultStatus.FIVEFOLD_REPETITION:
				msg += ' by fivefold repetition';
				break;
			case GameResultStatus.FIFTY_MOVE_RULE:
				msg += ' by fifty move rule';
				break;
			case GameResultStatus.SEVENTYFIVE_MOVE_RULE:
				msg += ' by seventy five move rule';
				break;
			case GameResultStatus.STALEMATE:
				msg += ' by stalemate';
				break;
			case GameResultStatus.RESIGNATION:
				msg += ' by resignation';
				break;
			case GameResultStatus.DRAW_AGREED:
				msg += ', draw agreed';
				break;
			case GameResultStatus.FLAGGED:
				msg += ', opponent flagged';
				break;
			case GameResultStatus.ADJUDICATION:
				msg += ' by adjudication';
				break;
			case GameResultStatus.TIMED_OUT:
				msg += ', opponent timed out';
				break;
			case GameResultStatus.ABORTED:
				msg += ', opponent aborted the game';
				break;
			default:
				break;
		}
		return msg;
	});
	currentPreview = $derived(this.historyMovesInfo[this.historyIndex]);
	checkSquare: Coord | undefined = $derived.by(() => {
		if (!this.currentPreview?.check) {
			return;
		}
		const isWhiteTurn = this.historyIndex % 2 === 0;
		const king = isWhiteTurn ? 'K' : 'k';
		for (const [coord, pieceData] of gameManager.board.position) {
			if (pieceData.piece === king) {
				return coord;
			}
		}
	});
	totalMoves = $derived(this.historyMovesInfo.length);
	isLatestHistory = $derived(this.historyIndex === this.ply);
	moveDurationsMs = $derived.by(() => {
		if (!this.startTime || this.historyMovesInfo.length < 1) {
			return [];
		}
		const startDate = timestampDate(this.startTime);
		return this.historyMovesInfo.reduce<number[]>((acc, cur, i) => {
			if (cur.playedAt) {
				const prev = this.historyMovesInfo[i - 1];
				const curDate = timestampDate(cur.playedAt);
				const prevDate = prev?.playedAt ? timestampDate(prev.playedAt) : startDate;
				acc.push(curDate.getTime() - prevDate.getTime());
			} else {
				acc.push(0);
			}
			return acc;
		}, []);
	});

	loadPersistedWhiteFirstMoveTimer() {
		if (this.ply !== 0) {
			return;
		}
		if (!this.whiteFirstMoveTimerStartTime.current.val) {
			return;
		}
		if (Number.isNaN(this.whiteFirstMoveTimerStartTime.current.val)) {
			this.whiteFirstMoveTimerStartTime.current.val = null;
			throw new Error('invalid white first move timer start time');
		}
		const elapsedMs = Date.now() - this.whiteFirstMoveTimerStartTime.current.val;
		const remainingMs = Math.max(0, FIRST_MOVE_TIMEOUT_MS - elapsedMs);
		if (remainingMs > 0) {
			this.whiteFirstMoveTimer.synchronize(remainingMs);
		}
	}

	loadPersistedBlackFirstMoveTimer() {
		if (this.ply !== 1) {
			return;
		}
		if (!this.blackFirstMoveTimerStartTime.current.val) {
			return;
		}
		if (Number.isNaN(this.blackFirstMoveTimerStartTime.current.val)) {
			this.blackFirstMoveTimerStartTime.current.val = null;
			throw new Error('invalid white first move timer start time');
		}
		const elapsedMs = Date.now() - this.blackFirstMoveTimerStartTime.current.val;
		const remainingMs = Math.max(0, FIRST_MOVE_TIMEOUT_MS - elapsedMs);
		if (remainingMs > 0) {
			this.blackFirstMoveTimer.synchronize(remainingMs);
		}
	}

	updateTimersState(): void {
		if (this.ply === 0) {
			this.showWhiteFirstMoveTimer = true;
			this.whiteFirstMoveTimer.start();
			if (!this.whiteFirstMoveTimerStartTime.current.val) {
				this.whiteFirstMoveTimerStartTime.current.val = Date.now();
			}
		} else if (this.ply === 1) {
			this.showWhiteFirstMoveTimer = false;
			this.whiteFirstMoveTimer.reset();
			this.showBlackFirstMoveTimer = true;
			this.blackFirstMoveTimer.start();
			if (!this.blackFirstMoveTimerStartTime.current.val) {
				this.blackFirstMoveTimerStartTime.current.val = Date.now();
			}
			this.whiteFirstMoveTimerStartTime.current.val = null;
		} else if (this.ply === 2) {
			this.showBlackFirstMoveTimer = false;
			this.blackFirstMoveTimer.reset();
			this.clock?.toggle();
			this.blackFirstMoveTimerStartTime.current.val = null;
		} else if (this.ply > 2) {
			this.clock?.toggle();
		}
	}

	handleWebsocketMessage(event: MessageEvent) {
		try {
			const msg = fromJsonString(MessageSchema, event.data);
			console.debug('ws recv:', msg.event.case, '-', msg.event.value);

			switch (msg.event.case) {
				case 'problem':
					this.onProblem(msg.event.value);
					break;
				case 'gameChatReceive':
					this.onGameChatReceive(msg.event.value);
					break;
				case 'gameChatRetrieve':
					this.onGameChatRetrieve(msg.event.value);
					break;
				case 'clientConnected':
					this.onClientConnected(msg.event.value);
					break;
				case 'clientDisconnected':
					this.onClientDisconnected(msg.event.value);
					break;
				case 'hubInfo':
					this.onHubInfo(msg.event.value);
					break;
				case 'matchFound':
					this.onMatchFound(msg.event.value);
					break;
				case 'gameAbort':
					this.onGameAbort(msg.event.value);
					break;
				case 'gameOfferDraw':
					this.onGameOfferDraw(msg.event.value);
					break;
				case 'gameResign':
					this.onGameResign(msg.event.value);
					break;
				case 'gameDeclineDraw':
					this.onGameDeclineDraw(msg.event.value);
					break;
				case 'gameAcceptDraw':
					this.onGameAcceptDraw(msg.event.value);
					break;
				case 'receiveMove':
					this.onReceiveMove(msg.event.value);
					break;
				case 'gameFinished':
					this.onGameFinished(msg.event.value);
					break;
				default:
					console.error('unknown message', msg.event.case, msg.event.value);
					break;
			}
		} catch (error) {
			console.error('json parse msg data', error);
		}
	}

	synchronizeClocks(whiteTimeMs?: number, blackTimeMs?: number) {
		if (whiteTimeMs && blackTimeMs) {
			this.clock?.synchronize(whiteTimeMs, blackTimeMs);
		}
	}

	seekGame(control: ClockControl): void {
		this.gameTimeControl = create(GameTimeControlSchema, {
			clock: { seconds: BigInt(control.clock) },
			increment: { seconds: BigInt(control.increment) }
		});

		const seekGameMsg = create(MessageSchema, {
			event: {
				case: 'seekGame',
				value: { timeControl: this.gameTimeControl }
			}
		});
		ws.send(seekGameMsg);
		this.uiState = 'seeking';
		this.poolLast.current.val = { clock: control.clock, increment: control.increment };
		this.chatMessages.current = [];
	}

	cancelSeekGame(): void {
		if (this.uiState !== 'seeking' || !this.gameTimeControl) {
			return;
		}
		const cancelSeekMsg = create(MessageSchema, {
			event: {
				case: 'cancelSeekGame',
				value: { timeControl: this.gameTimeControl }
			}
		});
		ws.send(cancelSeekMsg);
		this.uiState = 'idle';
		this.gameTimeControl = undefined;
		this.poolLast.current.val = null;
	}

	gameAbort() {
		const gameAbortMsg = create(MessageSchema, { event: { case: 'gameAbort', value: {} } });
		ws.send(gameAbortMsg);
	}

	gameOfferDraw() {
		const gameOfferDrawMsg = create(MessageSchema, { event: { case: 'gameOfferDraw', value: {} } });
		ws.send(gameOfferDrawMsg);
	}

	gameResign() {
		const gameResignMsg = create(MessageSchema, { event: { case: 'gameResign', value: {} } });
		ws.send(gameResignMsg);
	}

	gameDeclineDraw() {
		const gameDeclineDrawMsg = create(MessageSchema, { event: { case: 'gameDeclineDraw', value: {} } });
		ws.send(gameDeclineDrawMsg);
		this.opponenOfferedDraw = false;
	}

	gameAcceptDraw() {
		const gameAcceptDrawMsg = create(MessageSchema, { event: { case: 'gameAcceptDraw', value: {} } });
		ws.send(gameAcceptDrawMsg);
		this.opponenOfferedDraw = false;
	}

	playMoveUci(uci: string) {
		const playMoveUciMsg = create(MessageSchema, { event: { case: 'playMoveUci', value: { move: uci } } });
		ws.send(playMoveUciMsg);
		if (this.ply <= 1) {
			this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
		}
		this.ply++;
		this.updateTimersState();
	}

	promotePiece(promotionPieceSymbol: string) {
		const src = this.promotionSrcDest.slice(0, 2) as Coord;
		const dest = this.promotionSrcDest.slice(2, 4) as Coord;
		const pos = new Map(this.board.position);
		pos.delete(src);
		pos.set(dest, { id: crypto.randomUUID(), piece: promotionPieceSymbol as PieceFenSymbol });
		this.board.setPosition(pos);
	}

	handlePromotionPiecePick(promotionPopoverElm: HTMLElement, symbol: string) {
		this.promotionPieceSymbol = this.color === Color.WHITE ? symbol.toUpperCase() : symbol;
		promotionPopoverElm.hidePopover();
		this.promotePiece(this.promotionPieceSymbol);
		this.playMoveUci(this.promotionSrcDest + this.promotionPieceSymbol.toLowerCase());
	}

	movesSkipToStart() {
		this.historyIndex = 0;
	}

	movesStepBack() {
		if (this.historyIndex > 0) {
			this.historyIndex--;
		}
	}

	movesStepForward() {
		if (this.historyIndex < this.totalMoves - 1) {
			this.historyIndex++;
			soundManager.play(this.currentPreview?.move?.san.includes('+') ? 'Capture' : 'Move');
		}
	}

	movesSkipToEnd() {
		if (this.historyIndex < this.totalMoves - 1) {
			this.historyIndex = this.totalMoves - 1;
			soundManager.play(this.currentPreview?.move?.san.includes('+') ? 'Capture' : 'Move');
		}
	}

	movesJumpTo(num: number) {
		this.historyIndex = Math.min(Math.max(num, 0), this.totalMoves - 1);
	}

	closeChatDialog() {
		this.chatDialogOpen = false;
		this.chatDialogElm?.close();
	}

	openChatDialog() {
		this.chatDialogOpen = true;
		this.chatDialogElm?.showModal();
	}

	toggleChatDialog() {
		this.chatDialogOpen = !this.chatDialogOpen;
		if (this.chatDialogElm?.open) {
			this.closeChatDialog();
		} else {
			this.openChatDialog();
		}
	}

	onProblem(problem: Problem) {
		this.wsError = problem.message;
	}

	onClientConnected(clientConnectedMsg: ClientConnected) {
		if (this.uiState === 'playing') {
			this.showOpponentReconnectTimer = false;
			this.opponentReconnectTimer.reset();
		}
	}
	onClientDisconnected(clientDisconnectedMsg: ClientDisconnected) {
		if (this.uiState === 'playing') {
			this.showOpponentReconnectTimer = true;
			this.opponentReconnectTimer.start();
		}
	}

	onHubInfo(hubInfoMsg: HubInfo) {
		this.lobbyCount = hubInfoMsg.lobby;
		this.playingCount = hubInfoMsg.playing;
	}

	onMatchFound(matchFoundMsg: MatchFound) {
		if (matchFoundMsg.opponentInfo) {
			this.opponentInfo = matchFoundMsg.opponentInfo;
		}
		this.uiState = 'playing';
		this.gameId = matchFoundMsg.gameId;
		this.clientId = matchFoundMsg.clientId;
		this.gameState = matchFoundMsg.gameState;
		this.gameTimeControl = matchFoundMsg.timeControl;
		this.color = matchFoundMsg.color;
		this.orientation = this.color === Color.WHITE ? 'w' : 'b';
		this.fen = matchFoundMsg.fen;
		this.ply = matchFoundMsg.ply;
		this.legalMoves = matchFoundMsg.legalMoves;
		this.historyMovesInfo = matchFoundMsg.historyMoveInfos;
		if (this.totalMoves > 0) {
			this.historyIndex = this.totalMoves - 1;
		}
		this.startTime = matchFoundMsg.startTime;
		this.reconnectTimeoutMs = Number(matchFoundMsg.reconnectTimeoutMs);
		this.firstMoveTimeoutMs = Number(matchFoundMsg.firstMoveTimeoutMs);
		this.clock?.setIncrement(Number(this.gameTimeControl?.increment?.seconds ?? 0) * 1000);
		this.clock?.white.pause();
		this.clock?.black.pause();
		this.clock?.setCurrentTurn(this.isWhiteTurn ? 'w' : 'b');
		const whiteTimeMs: number =
			Number(matchFoundMsg.clocks?.white?.seconds ?? 0) * 1000 +
			Number(matchFoundMsg.clocks?.white?.nanos ?? 0) / 1_000_000;
		const blackTimeMs: number =
			Number(matchFoundMsg.clocks?.black?.seconds ?? 0) * 1000 +
			Number(matchFoundMsg.clocks?.black?.nanos ?? 0) / 1_000_000;
		this.synchronizeClocks(whiteTimeMs, blackTimeMs);
		this.loadPersistedWhiteFirstMoveTimer();
		this.loadPersistedBlackFirstMoveTimer();
		this.updateTimersState();
		goto(`/game/${this.gameId}`);
	}

	onGameFinished(gameFinishedMsg: GameFinished) {
		this.uiState = 'idle';
		this.gameId = '';
		this.chatGameId.current = '';
		this.chatLastId.current = '';
		this.clientId = '';
		this.gameResult = gameFinishedMsg.gameResult;
		this.gameResultStatus = gameFinishedMsg.gameResultStatus;
		this.gameState = gameFinishedMsg.gameState;
		this.showOpponentReconnectTimer = false;
		this.showWhiteFirstMoveTimer = false;
		this.showBlackFirstMoveTimer = false;
		this.opponentReconnectTimer.reset();
		this.whiteFirstMoveTimer.reset();
		this.blackFirstMoveTimer.reset();
		this.clock?.pause();
		if (this.gameResultStatus === GameResultStatus.FLAGGED) {
			if (this.gameResult === GameResult.WHITE_WON) {
				this.clock?.white.synchronize(0);
			} else if (this.gameResult === GameResult.BLACK_WON) {
				this.clock?.white.synchronize(0);
			}
		}
		this.opponenOfferedDraw = false;
		if (this.clock) {
			this.clock.state = 'idle';
		}
		this.whiteFirstMoveTimerStartTime.current.val = null;
		this.blackFirstMoveTimerStartTime.current.val = null;
	}

	onGameAbort(msg: GameAbort) {}

	onGameOfferDraw(msg: GameOfferDraw) {
		this.opponenOfferedDraw = true;
	}

	onGameResign(msg: GameResign) {}

	onGameDeclineDraw(msg: GameDeclineDraw) {}

	onGameAcceptDraw(msg: GameAcceptDraw) {}

	onGameChatReceive(chat: GameChatReceive) {
		if (chat.clientId !== this.clientId) {
			if (this.chatGameId.current !== this.gameId) {
				this.chatMessages.current = [];
				this.chatGameId.current = this.gameId;
			}
			this.chatMessages.current.push({ text: chat.message, received: true });
		}
		this.chatLastId.current = chat.id;
	}

	onGameChatRetrieve(chat: GameChatRetrieve) {
		const last = chat.gameChat.at(-1);
		if (this.chatGameId.current !== this.gameId) {
			this.chatMessages.current = [];
			this.chatGameId.current = this.gameId;
		}
		if (last) {
			if (last.id === this.chatLastId.current) {
				return;
			}
			chat.gameChat.forEach(chat => {
				this.chatMessages.current.push({ received: chat.clientId !== this.clientId, text: chat.message });
			});
			this.chatLastId.current = last.id;
		}
	}

	onReceiveMove(moveMsg: ReceiveMove) {
		const ownMove = moveMsg.ply === this.ply;
		this.uci = moveMsg.uci;
		this.lan = moveMsg.lan;
		this.san = moveMsg.san;
		this.fen = moveMsg.fen;
		this.ply = moveMsg.ply;
		this.legalMoves = moveMsg.legalMoves;
		this.historyMovesInfo.push({
			$typeName: 'pb.HistoryMoveInfo',
			fen: this.fen,
			check: this.san.includes('+'),
			move: { $typeName: 'pb.HistoryMove', uci: this.uci, san: this.san },
			playedAt: timestampNow()
		});
		const tmpIndex = this.historyIndex + 1;
		if (tmpIndex === this.ply) {
			this.historyIndex++;
		}
		const whiteTimeMs: number =
			Number(moveMsg.clocks?.white?.seconds ?? 0) * 1000 + Number(moveMsg.clocks?.white?.nanos ?? 0) / 1_000_000;
		const blackTimeMs: number =
			Number(moveMsg.clocks?.black?.seconds ?? 0) * 1000 + Number(moveMsg.clocks?.black?.nanos ?? 0) / 1_000_000;
		this.synchronizeClocks(whiteTimeMs, blackTimeMs);
		if (!ownMove) {
			if (this.ply <= 2) {
				this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
			}
			this.updateTimersState();
			if (!gameManager.isLatestHistory) {
				return;
			}
			const src = this.uci.slice(0, 2) as Coord;
			const dest = this.uci.slice(2, 4) as Coord;
			const isPromo = isPromotionUciMove(this.uci);
			if (isPromo) {
				this.promotionSrcDest = this.uci.slice(0, 4);
				this.promotionPieceSymbol = this.color === Color.WHITE ? this.uci[4]! : this.uci[4]!.toUpperCase();
				this.promotePiece(this.promotionPieceSymbol);
				return;
			}
			const rookMove = getRookCastleMove(this.uci);
			if (rookMove) {
				const rookSrc = rookMove.slice(0, 2) as Coord;
				const rookDest = rookMove.slice(2, 4) as Coord;
				const pos = new Map(this.board.position);
				const pieceData = this.board.getPiece(src)!;
				const rookPieceData = this.board.getPiece(rookSrc)!;
				pos.delete(src);
				pos.set(dest, pieceData);
				pos.delete(rookSrc);
				pos.set(rookDest, rookPieceData);
				this.board.setPosition(pos);
				return;
			}
			const enpOppPieceCoordToDelete = isEnpassantLanMove(this.lan, this.opponentColor === Color.WHITE ? 'w' : 'b');
			if (enpOppPieceCoordToDelete) {
				const pos = new Map(this.board.position);
				const pieceData = this.board.getPiece(src)!;
				pos.delete(enpOppPieceCoordToDelete);
				pos.delete(src);
				pos.set(dest, pieceData);
				this.board.setPosition(pos);
				return;
			}
			this.board.movePiece(src, dest);
		}
	}
}

export const gameManager = new GameManager();

export type ClockControl = {
	clock: number;
	increment: number;
};

const CASTLE_MOVES_PAIR = new Map<string, string>([
	['e1g1', 'h1f1'],
	['e1c1', 'a1d1'],
	['e8g8', 'h8f8'],
	['e8c8', 'a8d8']
]);

export const PROMOS = ['q', 'r', 'n', 'b'];

export function getRookCastleMove(uci: string): string | undefined {
	return CASTLE_MOVES_PAIR.get(uci);
}

export function isPromotionMove(move: string, legalMoves: string[]): boolean {
	if (move.length !== 4) {
		return false;
	}
	return legalMoves.some(m => m.length === 5 && m.startsWith(move) && PROMOS.includes(m[4]!));
}

function isPromotionUciMove(uci: string): boolean {
	return uci.length === 5 && PROMOS.includes(uci[4]!);
}

export function isEnpassantLanMove(lan: string, currentTurn: 'w' | 'b'): Coord | undefined {
	if (!lan.includes('x')) {
		return;
	}
	const [src, dest] = lan.split('x');
	if (src?.length !== 2 || dest?.length !== 2) {
		return;
	}
	const [srcRank, destRank] = [src[1], dest[1]];
	if (currentTurn === 'w') {
		if (srcRank === '5' && destRank === '6') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	} else {
		if (srcRank === '4' && destRank === '3') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	}
}

export function playedEnpassantMove(move: string): Coord | undefined {
	if (move.length !== 4) {
		return;
	}
	const src = move.slice(0, 2) as Coord;
	const dest = move.slice(2, 4) as Coord;
	const pieceData = gameManager.board.getPiece(src);
	if (!(pieceData?.piece === 'p' || pieceData?.piece === 'P')) {
		return;
	}
	const [srcRank, destRank] = [src[1], dest[1]];
	if (gameManager.isWhiteTurn) {
		if (!(srcRank === '5' && destRank === '6')) {
			return;
		}
		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'p') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	} else {
		if (!(srcRank === '4' && destRank === '3')) {
			return;
		}
		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'P') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	}
}
