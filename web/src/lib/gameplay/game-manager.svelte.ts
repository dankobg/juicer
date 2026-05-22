import {
	GameTimeKind,
	GameVariant,
	GameState,
	GameTimeCategory,
	type GameTimeControl,
	GameResult,
	GameResultStatus,
	Color,
	type GameMove,
	MessageSchema,
	type GameFound,
	type GameInfo,
	type MoveSync,
	type MoveAck,
	type GameChat,
	type GameChatList,
	type AbortGame,
	type ResignGame,
	type OfferDraw,
	type AcceptDraw,
	type DeclineDraw,
	type GameFinished,
	type PlayerInfo
} from '$lib/gen/juicer_pb';
import { create } from '@bufbuild/protobuf';
import type { Coord, JuicerBoard } from '@dankop/juicer-board';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { SvelteMap } from 'svelte/reactivity';
import { goto } from '$app/navigation';

class GameManager {
	games = $state<SvelteMap<number, Game>>(new SvelteMap());

	sendGameChat(gameId: number, message: string): void {
		const sendGameChatMsg = create(MessageSchema, {
			event: { case: 'sendGameChat', value: { gameId, message } }
		});
		ws.send(sendGameChatMsg);
	}

	onGameFound(gameFound: GameFound): void {
		console.log('GameFound game_id: ', gameFound.gameId);
		const exists = this.games.has(gameFound.gameId);
		if (!exists) {
			this.games.set(gameFound.gameId, new Game(gameFound.gameId));
		}

		const game = this.games.get(gameFound.gameId);
		if (!game) {
			return;
		}

		goto(`/game/${gameFound.gameId}`);
	}

	onGameInfo(gameInfo: GameInfo): void {
		console.log('------------------------- GOT GAME INFO ------------------------', gameInfo);

		const exists = this.games.has(gameInfo.gameId);
		if (!exists) {
			this.games.set(gameInfo.gameId, new Game(gameInfo.gameId));
		}

		const game = this.games.get(gameInfo.gameId);
		if (!game) {
			return;
		}

		const gameOptions: GameOptions = {
			white: gameInfo.white,
			black: gameInfo.black,
			gameVariant: gameInfo.gameVariant,
			gameTimeKind: gameInfo.gameTimeKind,
			gameTimeCategory: gameInfo.gameTimeCategory,
			gameTimeControl: { $typeName: 'pb.GameTimeControl', clockMs: 423, incrementMs: 0 },
			gameState: gameInfo.gameState,
			gameResult: gameInfo.gameResult,
			gameResultStatus: gameInfo.gameResultStatus,
			reconnectTimeoutMs: gameInfo.reconnectTimeoutMs,
			firstMoveTimeoutMs: gameInfo.firstMoveTimeoutMs,
			// lastMove
			startTime: Number(gameInfo.startTime?.seconds), // @TODO: fix later
			endTime: Number(gameInfo.endTime?.seconds), // @TODO: fix later
			rated: gameInfo.rated,
			version: gameInfo.version,
			// repetitions: gameInfo.repetitions
			ply: gameInfo.ply,
			gameMoves: gameInfo.gameMoves,
			myColor: gameInfo.color,
			orientation: Color.WHITE, // @TODO: fix later
			gameHistoryIndex: 0 // @TODO: fix later
		};

		game.configure(gameOptions);
	}

	onAbortGame(abortGame: AbortGame): void {
		console.log('AbortGame: ', abortGame);
	}

	onResignGame(resignGame: ResignGame): void {
		console.log('ResignGame: ', resignGame);
	}

	onOfferDraw(offerDraw: OfferDraw): void {
		console.log('OfferDraw: ', offerDraw);
	}

	onAcceptDraw(acceptDraw: AcceptDraw): void {
		console.log('AcceptDraw: ', acceptDraw);
	}

	onDeclinedDraw(declineDraw: DeclineDraw): void {
		console.log('DeclineDraw: ', declineDraw);
	}

	onGameFinished(gameFinished: GameFinished): void {
		console.log('GameFinished: ', gameFinished);
	}

	onMoveSync(moveSync: MoveSync): void {
		console.log('MoveSync: ', moveSync);
	}

	onMoveAck(moveAck: MoveAck): void {
		console.log('MoveAck: ', moveAck);
	}

	onGameChat(gameChat: GameChat): void {
		console.log('GameChat: ', gameChat);
	}

	onGameChatList(gameChats: GameChatList): void {
		console.log('GameChatList: ', gameChats);
	}
}

type GameOptions = {
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
	lastMove?: number;
	startTime?: number;
	endTime?: number;
	rated?: boolean;
	version?: number;
	repetitions?: number;
	ply?: number;
	gameMoves?: GameMove[];
	myColor?: Color;
	orientation?: Color;
	gameHistoryIndex?: number;
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
	lastMove = $state<number | undefined>();
	startTime = $state<number | undefined>();
	endTime = $state<number | undefined>();
	rated = $state<boolean | undefined>();
	version = $state<number>(0);
	repetitions = $state<number>(0);
	gameMoves = $state<GameMove[]>([]);
	legalMoves = $state<string[]>([]);
	ply = $state<number>(0);
	// pending draw offer
	// clocks and increment

	// ####################################

	myColor = $state<Color>(Color.UNSPECIFIED);
	orientation = $state<Color>(Color.UNSPECIFIED);
	gameHistoryIndex = $state<number>(0);

	constructor(gameId: number, options?: GameOptions) {
		this.gameId = gameId;

		if (options) {
			this.configure(options);
		}
	}

	configure(options: GameOptions): void {
		const defaultOptions: GameOptions = {
			gameVariant: GameVariant.UNSPECIFIED,
			gameTimeKind: GameTimeKind.UNSPECIFIED,
			gameTimeCategory: GameTimeCategory.UNSPECIFIED,
			gameState: GameState.UNSPECIFIED,
			gameResult: GameResult.UNSPECIFIED,
			gameResultStatus: GameResultStatus.UNSPECIFIED,
			version: 0,
			repetitions: 0,
			ply: 0,
			gameMoves: [],
			myColor: Color.UNSPECIFIED,
			orientation: Color.UNSPECIFIED,
			gameHistoryIndex: 0
		};

		const opts = { ...defaultOptions, ...options };

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
		this.repetitions = opts?.repetitions ?? 0;
		this.ply = opts?.ply ?? 0;
		this.gameMoves = opts?.gameMoves ?? [];
		this.myColor = opts?.myColor ?? Color.UNSPECIFIED;
		this.orientation = opts?.orientation ?? Color.UNSPECIFIED;
		this.gameHistoryIndex = opts?.gameHistoryIndex ?? 0;
	}

	// #######################################

	view = $derived<'spectating' | 'playing'>(this.myColor === Color.UNSPECIFIED ? 'spectating' : 'playing');
	currentPlayerId = $derived<string | undefined>(
		this.myColor === Color.WHITE ? this.white?.userId : this.black?.userId
	);

	isCheck = $state<boolean>(Boolean(this.gameMoves.at(-1)?.check));
	isCheckmate = $state<boolean>(Boolean(this.gameMoves.at(-1)?.san?.includes('#')));
	hasIncrement: boolean = $derived<boolean>(this.gameTimeControl?.incrementMs !== 0);
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

	incrementVersion(): void {
		this.version++;
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
		const playMoveUciMsg = create(MessageSchema, {
			event: {
				case: 'playMoveUci',
				value: {
					gameId: this.gameId,
					uci
				}
			}
		});
		ws.send(playMoveUciMsg);
		// if (this.ply <= 1) {
		// 	this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
		// }
		// this.ply++;
		// this.updateTimersState();
	}

	handlePromotionPiecePick() {}
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

export const gameManager = new GameManager();

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
	if (promotionSymbol === 'q') {
		return 'promote to queen';
	}
	if (promotionSymbol === 'r') {
		return 'promote to rook';
	}
	if (promotionSymbol === 'n') {
		return 'promote to knight';
	}
	if (promotionSymbol === 'b') {
		return 'promote to bishop';
	}
	return '';
}
