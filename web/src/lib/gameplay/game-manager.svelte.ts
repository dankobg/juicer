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
import type { Coord, JuicerBoard, PieceFenSymbol } from '@dankop/juicer-board';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { SvelteMap } from 'svelte/reactivity';
import { goto } from '$app/navigation';
import type { ChatMessage } from '$lib/components/chat-box/chat-box.svelte';
import { soundManager } from '$lib/sound/sound-manager.svelte';

export class GameManager {
	games = $state<SvelteMap<number, Game>>(new SvelteMap());
	gameChatMessages = $state<SvelteMap<number, ChatMessage[]>>(new SvelteMap());

	sendGameChat(gameId: number, message: string): void {
		const sendGameChatMsg = create(MessageSchema, {
			event: { case: 'sendGameChat', value: { gameId, message } }
		});
		ws.send(sendGameChatMsg);
	}

	onGameFound(gameFound: GameFound): void {
		const exists = this.games.has(gameFound.gameId);
		if (!exists) {
			this.games.set(gameFound.gameId, new Game({ gameId: gameFound.gameId }));
		}
		goto(`/game/${gameFound.gameId}`);
		soundManager.play('NewChallenge');
	}

	onGameInfo(gameInfo: GameInfo): void {
		const gameOpts: GameOptions = {
			gameId: gameInfo.gameId,
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
			orientation: gameInfo.color, // @TODO: fix later
			gameHistoryIndex: 0, // @TODO: fix later
			legalMoves: gameInfo.legalMoves
		};

		if (this.games.has(gameInfo.gameId)) {
			const game = this.games.get(gameInfo.gameId);
			game?.configure(gameOpts);
		} else {
			this.games.set(gameInfo.gameId, new Game(gameOpts));
		}
	}

	onAbortGame(abortGame: AbortGame): void {}

	onResignGame(resignGame: ResignGame): void {}

	onOfferDraw(offerDraw: OfferDraw): void {}

	onAcceptDraw(acceptDraw: AcceptDraw): void {}

	onDeclinedDraw(declineDraw: DeclineDraw): void {}

	onGameFinished(gameFinished: GameFinished): void {
		const game = this.games.get(gameFinished.gameId);
		if (game) {
			game.gameResult = gameFinished.gameResult;
			game.gameResultStatus = gameFinished.gameResultStatus;
			game.gameState = gameFinished.gameState;
		}
	}

	onMoveSync(moveSync: MoveSync): void {
		const game = this.games.get(moveSync.gameId);
		if (!game) {
			throw new Error('movesync: no game found');
		}

		const myMove = moveSync.ply === game.ply;
		game.version = moveSync.version;
		game.legalMoves = moveSync.legalMoves;
		game.ply = moveSync.ply;
		game.gameMoves.push({
			$typeName: 'pb.GameMove',
			uci: moveSync.uci,
			san: moveSync.san,
			lan: moveSync.lan,
			check: moveSync.san.includes('+'),
			fen: moveSync.fen
		});

		// handleeeeeeeeeeeeeee moveSync.clocks

		if (!myMove) {
			// if (game.ply <= 2) {
			//   this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
			// }
			// 		this.updateTimersState();
			const src = moveSync.uci.slice(0, 2) as Coord;
			const dest = moveSync.uci.slice(2, 4) as Coord;
			const isPromo = isPromotionUciMove(moveSync.uci);
			if (isPromo) {
				game.promotionSrcDest = moveSync.uci.slice(0, 4);
				game.promotionPieceSymbol = game.myColor === Color.WHITE ? moveSync.uci[4]! : moveSync.uci[4]!.toUpperCase();
				game.promotePiece(game.promotionPieceSymbol);
				return;
			}
			const rookMove = getRookCastleMove(moveSync.uci);
			if (rookMove) {
				const rookSrc = rookMove.slice(0, 2) as Coord;
				const rookDest = rookMove.slice(2, 4) as Coord;
				const pos = new Map(game.board.position);
				const pieceData = game.board.getPiece(src)!;
				const rookPieceData = game.board.getPiece(rookSrc)!;
				pos.delete(src);
				pos.set(dest, pieceData);
				pos.delete(rookSrc);
				pos.set(rookDest, rookPieceData);
				game.board.setPosition(pos);
				return;
			}
			const enpOppPieceCoordToDelete = isEnpassantLanMove(moveSync.lan, game.opponentColor === Color.WHITE ? 'w' : 'b');
			if (enpOppPieceCoordToDelete) {
				const pos = new Map(game.board.position);
				const pieceData = game.board.getPiece(src)!;
				pos.delete(enpOppPieceCoordToDelete);
				pos.delete(src);
				pos.set(dest, pieceData);
				game.board.setPosition(pos);
				return;
			}
			game.board.movePiece(src, dest);
		}
	}

	onMoveAck(moveAck: MoveAck): void {}

	onGameChat(gameChat: GameChat): void {
		const gameChats = this.gameChatMessages.get(gameChat.gameId);
		if (!gameChats) {
			this.gameChatMessages.set(gameChat.gameId, []);
		}

		gameChats!.push({
			userId: gameChat.userId,
			messageId: gameChat.messageId,
			message: gameChat.message,
			postedAt: gameChat.postedAt
		});
	}

	onGameChatList(gameChats: GameChatList): void {}
}

type GameOptions = {
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
	lastMove?: number;
	startTime?: number;
	endTime?: number;
	rated?: boolean;
	version?: number;
	ack?: number;
	repetitions?: number;
	ply?: number;
	gameMoves?: GameMove[];
	myColor?: Color;
	orientation?: Color;
	gameHistoryIndex?: number;
	legalMoves?: string[];
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
	ack = $state<number>(0);
	repetitions = $state<number>(0);
	gameMoves = $state<GameMove[]>([]);
	legalMoves = $state<string[]>([]);
	ply = $state<number>(0);
	// pending draw offer
	// clocks and increment

	// ####################################
	myColor = $state<Color>(Color.UNSPECIFIED);
	opponentColor = $derived(
		this.myColor === Color.UNSPECIFIED ? Color.UNSPECIFIED : this.myColor === Color.WHITE ? Color.BLACK : Color.WHITE
	);
	orientation = $state<Color>(Color.UNSPECIFIED);
	gameHistoryIndex = $state<number>(0);

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
			gameHistoryIndex: 0,
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
		this.gameHistoryIndex = opts?.gameHistoryIndex ?? 0;
		this.legalMoves = opts?.legalMoves ?? [];
	}

	// #######################################
	promotionSrcDest = $state<string>('');
	promotionPieceSymbol = $state<string>('');

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
	isMyTurnActive = $derived(
		(this.gameState === GameState.ACTIVE && this.myColor === Color.WHITE && this.isWhiteTurn) ||
			(this.myColor === Color.BLACK && this.isBlackTurn)
	);

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
		// if (this.ply <= 1) {
		// 	this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
		// }
		// this.ply++;
		// this.updateTimersState();
	}
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
