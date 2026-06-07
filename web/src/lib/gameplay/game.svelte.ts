import {
	Color,
	GameResult,
	GameResultStatus,
	GameState,
	GameTimeCategory,
	GameTimeKind,
	MessageSchema,
	type PlayerInfo,
	GameVariant,
	type GameTimeControl,
	type GameMove,
	type DrawOffer
} from '$lib/gen/juicer_pb';
import { soundManager } from '$lib/sound/sound-manager.svelte';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { create } from '@bufbuild/protobuf';
import { type Timestamp, type Duration, timestampDate } from '@bufbuild/protobuf/wkt';
import type { JuicerBoard, Coord, PieceFenSymbol } from '@dankop/juicer-board';
import { SvelteMap } from 'svelte/reactivity';

export type GameOptions = {
	gameId?: number;
	white?: PlayerInfo;
	black?: PlayerInfo;
	gameVariant?: GameVariant;
	gameTimeKind?: GameTimeKind;
	gameTimeCategory?: GameTimeCategory;
	gameTimeControl?: GameTimeControl;
	gameState?: GameState;
	gameResult?: GameResult;
	gameResultStatus?: GameResultStatus;
	reconnectTimeoutMs?: number;
	firstMoveTimeoutMs?: number;
	lastMove?: Timestamp;
	startTime?: Timestamp;
	endTime?: Timestamp;
	rated?: boolean;
	version?: number;
	ack?: number;
	repetitions?: number;
	ply?: number;
	gameMoves?: GameMove[];
	myColor?: Color;
	orientation?: Color;
	historyPointer?: number;
	legalMoves?: string[];
	whiteRemainingGameTime?: Duration;
	blackRemainingGameTime?: Duration;
	pendingDrawOffers?: Record<string, DrawOffer>;
};

export class Game {
	board!: JuicerBoard;
	gameId = $state<number | undefined>();
	white = $state<PlayerInfo | undefined>();
	black = $state<PlayerInfo | undefined>();
	gameVariant = $state<GameVariant>(GameVariant.UNSPECIFIED);
	gameTimeKind = $state<GameTimeKind>(GameTimeKind.UNSPECIFIED);
	gameTimeCategory = $state<GameTimeCategory>(GameTimeCategory.UNSPECIFIED);
	gameTimeControl = $state<GameTimeControl | undefined>();
	gameState = $state<GameState>(GameState.UNSPECIFIED);
	gameResult = $state<GameResult>(GameResult.UNSPECIFIED);
	gameResultStatus = $state<GameResultStatus>(GameResultStatus.UNSPECIFIED);
	reconnectTimeoutMs = $state<number | undefined>();
	firstMoveTimeoutMs = $state<number | undefined>();
	startTime = $state<Timestamp | undefined>();
	endTime = $state<Timestamp | undefined>();
	lastMove = $state<Timestamp | undefined>();
	rated = $state<boolean | undefined>();
	ply = $state<number>(0);
	version = $state<number>(0);
	ack = $state<number>(0);
	repetitions = $state<number>(0);
	gameMoves = $state<GameMove[]>([]);
	legalMoves = $state<string[]>([]);
	whiteRemainingGameTime = $state<Duration>();
	blackRemainingGameTime = $state<Duration>();
	pendingDrawOffers = $state<SvelteMap<string, DrawOffer>>(new SvelteMap());
	historyPointer = $state<number>(Math.max(0, this.gameMoves.length - 1));
	myColor = $state<Color>(Color.UNSPECIFIED);
	opponentColor = $derived(
		this.myColor === Color.UNSPECIFIED ? Color.UNSPECIFIED : this.myColor === Color.WHITE ? Color.BLACK : Color.WHITE
	);
	promotionSrcDest = $state<string>('');
	promotionPieceSymbol = $state<string>('');

	constructor(options?: Partial<GameOptions>) {
		this.configure(options);
	}

	configure(options?: Partial<GameOptions>): void {
		const defaultOptions: GameOptions = {
			gameVariant: GameVariant.UNSPECIFIED,
			gameTimeKind: GameTimeKind.UNSPECIFIED,
			gameTimeCategory: GameTimeCategory.UNSPECIFIED,
			gameState: GameState.UNSPECIFIED,
			gameResult: GameResult.UNSPECIFIED,
			gameResultStatus: GameResultStatus.UNSPECIFIED,
			version: 0,
			ack: 0,
			repetitions: 0,
			ply: 0,
			gameMoves: [],
			myColor: Color.UNSPECIFIED,
			orientation: Color.UNSPECIFIED,
			// historyPointer: 0,
			legalMoves: []
		};

		const opts = { ...defaultOptions, ...options };

		this.gameId = opts?.gameId;
		this.white = opts?.white;
		this.black = opts?.black;
		this.gameVariant = opts?.gameVariant ?? GameVariant.UNSPECIFIED;
		this.gameTimeKind = opts?.gameTimeKind ?? GameTimeKind.UNSPECIFIED;
		this.gameTimeCategory = opts?.gameTimeCategory ?? GameTimeCategory.UNSPECIFIED;
		this.gameTimeControl = opts?.gameTimeControl;
		this.gameState = opts?.gameState ?? GameState.UNSPECIFIED;
		this.gameResult = opts?.gameResult ?? GameResult.UNSPECIFIED;
		this.gameResultStatus = opts?.gameResultStatus ?? GameResultStatus.UNSPECIFIED;
		this.reconnectTimeoutMs = opts?.reconnectTimeoutMs;
		this.firstMoveTimeoutMs = opts?.firstMoveTimeoutMs;
		this.lastMove = opts?.lastMove;
		this.startTime = opts?.startTime;
		this.endTime = opts?.endTime;
		this.rated = opts?.rated;
		this.version = opts?.version ?? 0;
		this.ack = opts?.ack ?? 0;
		this.repetitions = opts?.repetitions ?? 0;
		this.ply = opts?.ply ?? 0;
		this.gameMoves = opts?.gameMoves ?? [];
		this.myColor = opts?.myColor ?? Color.UNSPECIFIED;
		this.orientation = opts?.orientation ?? Color.UNSPECIFIED;
		if (opts.historyPointer) {
			this.historyPointer = opts.historyPointer;
		}
		this.legalMoves = opts?.legalMoves ?? [];
		this.whiteRemainingGameTime = opts?.whiteRemainingGameTime;
		this.blackRemainingGameTime = opts?.blackRemainingGameTime;
		if (opts.pendingDrawOffers) {
			this.pendingDrawOffers = new SvelteMap();
			Object.entries(opts.pendingDrawOffers).forEach(([k, v]) => {
				this.pendingDrawOffers.set(k, v);
			});
		}
	}

	isPlaying = $derived(this.myColor !== Color.UNSPECIFIED);
	isSpectating = $derived(this.myColor === Color.UNSPECIFIED);

	getPlayerByColor(color: Color): PlayerInfo | undefined {
		if (color === Color.UNSPECIFIED) {
			return;
		}
		return color === Color.WHITE ? this.white : this.black;
	}
	mePlayer = $derived.by(() => {
		if (!this.isPlaying) {
			return;
		}
		return this.getPlayerByColor(this.myColor);
	});
	opponentPlayer = $derived.by(() => {
		if (!this.isPlaying) {
			return;
		}
		return this.getPlayerByColor(this.opponentColor);
	});

	animationFrameId = $state<number | null>(null);
	whiteRemaininGameTimeMs = $derived(this.whiteRemainingGameTime ? protoDurationToMs(this.whiteRemainingGameTime) : 0);
	blackRemaininGameTimeMs = $derived(this.blackRemainingGameTime ? protoDurationToMs(this.blackRemainingGameTime) : 0);
	whiteDisplayTimeMs = $state<number>(this.whiteRemaininGameTimeMs);
	blackDisplayTimeMs = $state<number>(this.blackRemaininGameTimeMs);

	getWhiteDisplayTimeMs(ts: number): number {
		if (!this.lastMove) {
			return this.whiteRemaininGameTimeMs;
		}
		const elapsed = Date.now() - timestampDate(this.lastMove).getTime();
		return this.isWhiteTurn ? Math.max(0, this.whiteRemaininGameTimeMs - elapsed) : this.whiteRemaininGameTimeMs;
	}

	getBlackDisplayTimeMs(ts: number): number {
		if (!this.lastMove) {
			return this.blackRemaininGameTimeMs;
		}
		const elapsed = Date.now() - timestampDate(this.lastMove).getTime();
		return this.isBlackTurn ? Math.max(0, this.blackRemaininGameTimeMs - elapsed) : this.blackRemaininGameTimeMs;
	}

	startClockTimerLoop() {
		if (this.animationFrameId !== null) {
			return;
		}

		const tick = (timestamp: number) => {
			this.whiteDisplayTimeMs = this.getWhiteDisplayTimeMs(timestamp);
			this.blackDisplayTimeMs = this.getBlackDisplayTimeMs(timestamp);
			this.animationFrameId = requestAnimationFrame(tick);
		};

		this.animationFrameId = requestAnimationFrame(tick);
	}

	stopClockTimerLoop() {
		if (this.animationFrameId) {
			cancelAnimationFrame(this.animationFrameId);
		}
	}

	checkSquare: Coord | undefined = $derived.by(() => {
		if (!this.currentPosition?.san || !sanMoveIsCheck(this.currentPosition?.san)) {
			return;
		}
		const isWhiteTurn = this.historyPointer % 2 === 0;
		const king = isWhiteTurn ? 'K' : 'k';
		for (const [coord, pieceData] of this.board.position) {
			if (pieceData.piece === king) {
				return coord;
			}
		}
	});

	movesCount = $derived<number>(this.gameMoves.length);
	playedMovesCount = $derived<number>(Math.max(0, this.gameMoves.length - 1));

	currentPosition = $derived(this.gameMoves.at(this.historyPointer));
	isViewingLatestPosition = $derived(this.historyPointer === this.ply);
	isCheck = $derived(Boolean(this?.currentPosition?.san && sanMoveIsCheck(this?.currentPosition?.san)));
	isCheckmate = $derived(Boolean(this?.currentPosition?.san && sanMoveIsCheckmate(this?.currentPosition?.san)));
	isCapture = $derived(Boolean(this?.currentPosition?.san && sanMoveIsCapture(this?.currentPosition?.san)));
	hasIncrement = $derived<boolean>(this.gameTimeControl?.incrementMs !== 0);
	isWhiteTurn = $derived<boolean>(this.ply % 2 === 0);
	isBlackTurn = $derived<boolean>(!this.isWhiteTurn);
	currentTurn = $derived<Color>(this.isWhiteTurn ? Color.WHITE : Color.BLACK);
	isGameActive = $derived<boolean>(this.gameState === GameState.ACTIVE);
	isGameConcluded = $derived<boolean>(
		this.gameState === GameState.FINISHED || this.gameState === GameState.INTERRUPTED
	);
	gameConcludedText = $derived<string>(
		this.isGameConcluded ? formatGameConcludedMsg(this.gameResult, this.gameResultStatus) : ''
	);
	isMyTurnActive = $derived.by(() => {
		if (!this.isPlaying) {
			return false;
		}
		return (
			(this.gameState === GameState.ACTIVE && this.myColor === Color.WHITE && this.isWhiteTurn) ||
			(this.myColor === Color.BLACK && this.isBlackTurn)
		);
	});

	orientation = $state<Color>(Color.UNSPECIFIED);
	uiShowAbortButton = $derived<boolean>(
		this.gameState === GameState.ACTIVE &&
			((this.myColor === Color.WHITE && this.ply === 0) || (this.myColor === Color.BLACK && this.ply <= 1))
	);
	uiShowFlipBoardButton = $state<boolean>(true);
	uiShowResignButton = $derived<boolean>(this.gameState === GameState.ACTIVE && this.ply >= 2);
	uiShowOfferDrawButton = $derived<boolean>(this.gameState === GameState.ACTIVE && this.ply >= 2);
	uiShowDrawAcceptDeclineButtons = $derived.by(() => {
		if (this.gameState !== GameState.ACTIVE || !this.opponentPlayer?.userId) {
			return false;
		}
		const offer = this.pendingDrawOffers.get(this.opponentPlayer?.userId);
		return offer?.ply === this.ply;
	});
	uiShowChatButton = $derived<boolean>(false);

	moveDurationsMs = $derived.by(() => {
		if (!this.startTime || this.gameMoves.length < 1) {
			return [];
		}
		const startDate = timestampDate(this.startTime);
		return this.gameMoves.reduce<number[]>((acc, cur, i) => {
			if (cur.playedAt) {
				const prev = this.gameMoves[i - 1];
				const curDate = timestampDate(cur.playedAt);
				const prevDate = prev?.playedAt ? timestampDate(prev.playedAt) : startDate;
				acc.push(curDate.getTime() - prevDate.getTime());
			} else {
				acc.push(0);
			}
			return acc;
		}, []);
	});

	movesCanGoBack = $derived<boolean>(this.historyPointer > 0);
	movesCanGoForward = $derived<boolean>(this.historyPointer < this.gameMoves.length - 1);
	movesCanGotoStart = $derived<boolean>(this.historyPointer > 0);
	movesCanGotoEnd = $derived<boolean>(this.historyPointer < this.gameMoves.length - 1);

	movesGoBack() {
		if (this.movesCanGoBack) {
			this.historyPointer--;
		}
	}

	movesGoForward() {
		if (this.movesCanGoForward) {
			this.historyPointer++;
			soundManager.play(this.isCapture ? 'Capture' : 'Move');
		}
	}

	movesGotoStart() {
		if (this.movesCanGotoStart) {
			this.historyPointer = 0;
		}
	}

	movesGotoEnd() {
		if (this.movesCanGotoEnd) {
			this.historyPointer = this.gameMoves.length - 1;
			soundManager.play(this.isCapture ? 'Capture' : 'Move');
		}
	}

	movesGoto(move: number) {
		this.historyPointer = Math.min(Math.max(move, 0), this.gameMoves.length - 1);
	}

	incrementVersion(): void {
		this.version++;
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
		this.promotionPieceSymbol = this.myColor === Color.WHITE ? symbol.toUpperCase() : symbol;
		promotionPopoverElm.hidePopover();
		this.promotePiece(this.promotionPieceSymbol);
		this.playMoveUci(this.promotionSrcDest + this.promotionPieceSymbol.toLowerCase());
	}

	abortGame() {
		const abortGameMsg = create(MessageSchema, { event: { case: 'abortGame', value: { gameId: this.gameId } } });
		ws.send(abortGameMsg);
	}

	resignGame() {
		const resignGameMsg = create(MessageSchema, { event: { case: 'resignGame', value: { gameId: this.gameId } } });
		ws.send(resignGameMsg);
	}

	offerDraw() {
		const offerDrawMsg = create(MessageSchema, { event: { case: 'offerDraw', value: { gameId: this.gameId } } });
		ws.send(offerDrawMsg);
	}

	acceptDraw() {
		const acceptDrawMsg = create(MessageSchema, { event: { case: 'acceptDraw', value: { gameId: this.gameId } } });
		ws.send(acceptDrawMsg);
	}

	dedclineDraw() {
		const declineDrawMsg = create(MessageSchema, { event: { case: 'declineDraw', value: { gameId: this.gameId } } });
		ws.send(declineDrawMsg);
	}

	playMoveUci(uci: string) {
		this.ack++;
		const playMoveUciMsg = create(MessageSchema, {
			event: {
				case: 'playMoveUci',
				value: {
					ack: this.ack,
					gameId: this.gameId,
					uci
				}
			}
		});
		ws.send(playMoveUciMsg);
	}
}

function protoDurationToMs(dur: Duration): number {
	return Number(dur.seconds) * 1000 + dur.nanos / 1e6;
}

function formatGameConcludedMsg(gameResult: GameResult, gameResultStatus: GameResultStatus): string {
	if (gameResult === GameResult.UNSPECIFIED && gameResultStatus === GameResultStatus.UNSPECIFIED) {
		return '';
	}

	let msg = '';

	switch (gameResult) {
		case GameResult.DRAW:
			msg += 'Draw';
			break;
		case GameResult.WHITE_WON:
			msg += 'White won';
			break;
		case GameResult.BLACK_WON:
			msg += 'Black won';
			break;
		case GameResult.INTERRUPTED:
			msg += 'Interrupted';
			break;
		default:
			break;
	}

	switch (gameResultStatus) {
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
			msg += ' byresignation';
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
		case GameResultStatus.INTERRUPTED:
			msg += ', game was interrupted before it could finish';
			break;
		default:
			break;
	}

	return msg;
}

export function sanMoveIsCapture(san: string): boolean {
	return san.includes('x');
}

export function sanMoveIsCheckmate(san: string): boolean {
	return san.endsWith('#');
}

export function sanMoveIsCheck(san: string): boolean {
	return san.endsWith('+') || sanMoveIsCheckmate(san);
}

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

export function isPromotionUciMove(uci: string): boolean {
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

export function playedEnpassantMove(game: Game, move: string): Coord | undefined {
	if (move.length !== 4) {
		return;
	}
	const src = move.slice(0, 2) as Coord;
	const dest = move.slice(2, 4) as Coord;
	const pieceData = game.board.getPiece(src);
	if (!(pieceData?.piece === 'p' || pieceData?.piece === 'P')) {
		return;
	}
	const [srcRank, destRank] = [src[1], dest[1]];
	if (game.isWhiteTurn) {
		if (!(srcRank === '5' && destRank === '6')) {
			return;
		}
		if (game.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'p') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	} else {
		if (!(srcRank === '4' && destRank === '3')) {
			return;
		}
		if (game.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'P') {
			return `${dest[0]}${srcRank}` as Coord;
		}
	}
}

export function getPromotionLabelText(promotionSymbol: string): string {
	switch (promotionSymbol) {
		case 'q':
			return 'promote to queen';
		case 'r':
			return 'promote to rook';
		case 'n':
			return 'promote to knight';
		case 'b':
			return 'promote to bishop';
		default:
			return '';
	}
}
