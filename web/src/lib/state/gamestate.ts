import {
	GameTimeKind,
	GameVariant,
	GameState,
	GameTimeCategory,
	type GameTimeControl,
	GameResult,
	GameResultStatus,
	Color,
	type GameMove
} from '$lib/gen/juicer_pb';
import type { JuicerBoard, Coord, PieceFenSymbol } from '@dankop/juicer-board';

export type Player = {
	userId: string;
	username: string;
	guest: boolean;
	color: Color;
};

class GameManager {
	games = $state<Record<number, Game>>();
}

class Game {
	board!: JuicerBoard;
	gameId = $state<number | undefined>();
	userId = $state<string | undefined>();
	white = $state<Player | undefined>();
	black = $state<Player | undefined>();
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
	// pending draw offer
	// clocks and increment

	// ####################################

	color = $state<Color>(Color.UNSPECIFIED);
	orientation = $state<Color>(Color.UNSPECIFIED);
	fen = $state<string | undefined>();
	uci = $state<string | undefined>();
	san = $state<string | undefined>();
	lan = $state<string | undefined>();
	ply = $state<number>(0);

	gameHistoryIndex = $state<number>(0);

	// #######################################

	check: boolean = $derived<boolean>(Boolean(this?.san?.includes('+')));
	checkmate: boolean = $derived<boolean>(Boolean(this?.san?.includes('#')));
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
}

class LobbyManager {
	wsError = $state<string | undefined>();
	uiState = $state<'idle' | 'seeking' | 'playing'>('idle');
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

// const CASTLE_MOVES_PAIR = new Map<string, string>([
// 	['e1g1', 'h1f1'],
// 	['e1c1', 'a1d1'],
// 	['e8g8', 'h8f8'],
// 	['e8c8', 'a8d8']
// ]);

// export const PROMOS = ['q', 'r', 'n', 'b'];

// export function getRookCastleMove(uci: string): string | undefined {
// 	return CASTLE_MOVES_PAIR.get(uci);
// }

// export function isPromotionMove(move: string, legalMoves: string[]): boolean {
// 	if (move.length !== 4) {
// 		return false;
// 	}
// 	return legalMoves.some(m => m.length === 5 && m.startsWith(move) && PROMOS.includes(m[4]!));
// }

// function isPromotionUciMove(uci: string): boolean {
// 	return uci.length === 5 && PROMOS.includes(uci[4]!);
// }

// export function isEnpassantLanMove(lan: string, currentTurn: 'w' | 'b'): Coord | undefined {
// 	if (!lan.includes('x')) {
// 		return;
// 	}
// 	const [src, dest] = lan.split('x');
// 	if (src?.length !== 2 || dest?.length !== 2) {
// 		return;
// 	}
// 	const [srcRank, destRank] = [src[1], dest[1]];
// 	if (currentTurn === 'w') {
// 		if (srcRank === '5' && destRank === '6') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	} else {
// 		if (srcRank === '4' && destRank === '3') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	}
// }

// export function playedEnpassantMove(move: string): Coord | undefined {
// 	if (move.length !== 4) {
// 		return;
// 	}
// 	const src = move.slice(0, 2) as Coord;
// 	const dest = move.slice(2, 4) as Coord;
// 	const pieceData = gameManager.board.getPiece(src);
// 	if (!(pieceData?.piece === 'p' || pieceData?.piece === 'P')) {
// 		return;
// 	}
// 	const [srcRank, destRank] = [src[1], dest[1]];
// 	if (gameManager.isWhiteTurn) {
// 		if (!(srcRank === '5' && destRank === '6')) {
// 			return;
// 		}
// 		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'p') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	} else {
// 		if (!(srcRank === '4' && destRank === '3')) {
// 			return;
// 		}
// 		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'P') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	}
// }
