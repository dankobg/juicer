import {
	BLACK,
	BLACK_KING,
	BLACK_PAWN,
	BLACK_QUEEN,
	BOARD_SIZE,
	FEN_EMPTY_POSITION,
	WHITE,
	WHITE_KING,
	WHITE_PAWN,
	WHITE_QUEEN,
} from './common';
import { fenToBoard } from './helpers';
import { Piece } from './piece';
import { Square } from './square';
import { CastleRights, type Color, type Coordinate, type PieceSymbol } from './types';

type FenToken = {
	halfMoveClock: number;
	fullMoveClock: number;
	enpSquare: Square | null;
	castleRights: CastleRights;
	turnColor: Color;
	position: string;
};

const FEN_PARTS_LENGTH = 6;
const FEN_NONE_SYMBOL = '-';
const FEN_SEPARATOR = ' ';
const FEN_POSITION_SEPARATOR = '/';

const RE_IS_DIGIT = /^[0-9]$/;
const RE_ENP_SQUARE = /^(-|[abcdefgh][36])$/;
const RE_CASTLE_RIGHTS = /^(-|\bK?Q?k?q?)$/m;
const RE_TURN_COLOR = /^(w|b)$/;
const RE_FEN_PIECE_SYMBOL = /^[prnbqkPRNBQK]$/;

export function validateFenMetadataParts(fen: string): FenToken {
	// tokens length must be 6 after splitting the fen by a single space delimiter
	const tokens = fen.split(FEN_SEPARATOR);
	if (tokens.length !== FEN_PARTS_LENGTH) {
		throw new Error('invalid FEN: length must be exactly 6 after splitting by a single space delimiter');
	}

	const [positionToken, turnColorToken, castleRightsToken, enpSquareToken, halfMoveClockToken, fullMoveClockToken] =
		tokens;

	// turn color must be either `w` | `b`
	if (!RE_TURN_COLOR.test(turnColorToken)) {
		throw new Error('invalid FEN: invalid active turn color');
	}

	const turn = turnColorToken === 'b' ? BLACK : WHITE;

	// full move clock must be a number >= 1
	const fullMoveClock = Number.parseInt(fullMoveClockToken);
	if (Number.isNaN(fullMoveClock) || fullMoveClock < 1) {
		throw new Error('invalid FEN: full move clock must be a number >= 1');
	}

	// half move clock must be a number >= 0
	const halfMoveClock = Number.parseInt(halfMoveClockToken);
	if (Number.isNaN(halfMoveClock) || halfMoveClock < 0) {
		throw new Error('invalid FEN: half move clock must be a number >= 0');
	}

	const n = turn === BLACK ? 1 : 0;

	// half move clock must be within the limit
	if (!(halfMoveClock <= (fullMoveClock - 1) * 2 + n)) {
		throw new Error('invalid FEN: half move clock must be whithin the valid limit');
	}

	// in case of an en-passant square, the half move clock must be equal to 0
	if (enpSquareToken !== FEN_NONE_SYMBOL && halfMoveClock !== 0) {
		throw new Error('invalid FEN: half move clock must be 0 if en-passant square exists');
	}

	// en-passant square must be a valid square or `-` if empty
	if (!RE_ENP_SQUARE.test(enpSquareToken)) {
		throw new Error('invalid FEN: en-passant target square is invalid');
	}

	let enpSquare: Square | null = null;
	if (enpSquareToken !== FEN_NONE_SYMBOL) {
		enpSquare = Square.fromCoord(enpSquareToken as Coordinate);
	}

	if (enpSquare !== null) {
		if ((turn === WHITE && enpSquare.rank === 3) || (turn === BLACK && enpSquare.rank === 6)) {
			throw new Error('invalid FEN: en-passant target square coordinate is invalid');
		}
	}

	// castle rights string must be of valid fen castle string format
	if (!RE_CASTLE_RIGHTS.test(castleRightsToken)) {
		throw new Error('invalid FEN: invalid castling rights string');
	}

	let cr: CastleRights = CastleRights.None;

	if (castleRightsToken.includes(WHITE_KING)) {
		cr |= CastleRights.WhiteKingSide;
	}
	if (castleRightsToken.includes(WHITE_QUEEN)) {
		cr |= CastleRights.WhiteQueenSide;
	}
	if (castleRightsToken.includes(BLACK_KING)) {
		cr |= CastleRights.BlackKingSide;
	}
	if (castleRightsToken.includes(BLACK_QUEEN)) {
		cr |= CastleRights.BlackQueenSide;
	}

	return {
		halfMoveClock,
		fullMoveClock,
		enpSquare,
		castleRights: cr,
		turnColor: turn,
		position: positionToken,
	};
}

function validatePositionPart(ft: FenToken): void {
	const ranks = ft.position.split(FEN_POSITION_SEPARATOR);
	if (ranks.length !== BOARD_SIZE) {
		throw new Error(`invalid FEN: it does not contain 8 ranks delimited by ${FEN_POSITION_SEPARATOR} character`);
	}

	const piecesCount: Record<string, number> = {};

	for (let r = 0; r < ranks.length; r++) {
		let sumSquaresInRank = 0;
		let previousWasNumber = false;

		for (let f = 0; f < ranks[r].length; f++) {
			if (RE_IS_DIGIT.test(ranks[r][f])) {
				if (previousWasNumber) {
					throw new Error('invalid FEN: position string is invalid, it has consecutive numbers');
				}

				const n = Number.parseInt(ranks[r][f]);
				if (Number.isNaN(n) || n < 0) {
					throw new Error('invalid FEN: failed to parse row number');
				}

				sumSquaresInRank += n;
				previousWasNumber = true;
			} else {
				if (!RE_FEN_PIECE_SYMBOL.test(ranks[r][f])) {
					throw new Error('invalid FEN: invalid piece symbol');
				}
				const piece = Piece.fromPieceFenSymbol(ranks[r][f] as PieceSymbol);
				piecesCount[piece.toFenSymbol()]++;
				sumSquaresInRank++;
				previousWasNumber = false;
			}
		}

		if (sumSquaresInRank != BOARD_SIZE) {
			throw new Error('invalid FEN: position string is invalid, too many squares in rank');
		}
	}

	if (piecesCount[WHITE_KING] === 0) {
		throw new Error('invalid FEN: position is missing white king');
	}
	if (piecesCount[BLACK_KING] === 0) {
		throw new Error('invalid FEN: position is missing black king');
	}

	const wkc = piecesCount[WHITE_KING];
	if (wkc > 1) {
		throw new Error(`invalid FEN: position is having too many white kings ${wkc}`);
	}

	const bkc = piecesCount[BLACK_KING];
	if (bkc > 1) {
		throw new Error(`invalid FEN: position is having too many black kings ${bkc}`);
	}

	for (const char of ranks[0]) {
		if (char === WHITE_PAWN) {
			throw new Error('invalid FEN: white pawn is on 8th rank');
		}
	}

	for (const char of ranks[7]) {
		if (char === BLACK_PAWN) {
			throw new Error('invalid FEN: black pawn is on 1st rank');
		}
	}
}

export type PositionMeta = {
	fenToken: FenToken;
	squares: Square[];
};

export function validateFen(fen: string): PositionMeta {
	if (fen === FEN_EMPTY_POSITION) {
		const squares = fenToBoard(fen);

		return {
			fenToken: {
				castleRights: CastleRights.None,
				enpSquare: null,
				fullMoveClock: 1,
				halfMoveClock: 0,
				position: fen.split(' ')[0],
				turnColor: WHITE,
			},
			squares,
		};
	}

	const fenToken = validateFenMetadataParts(fen);
	validatePositionPart(fenToken);

	const squares = fenToBoard(fenToken.position);

	return {
		fenToken,
		squares,
	};
}
