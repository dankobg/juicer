import { KING, QUEEN, KNIGHT, BISHOP, ROOK, PAWN, WHITE, BLACK } from './common';
import type { Color, PieceSymbol } from './types';

export class Piece {
	id!: string;

	constructor(
		public symbol: PieceSymbol,
		public color: Color
	) {
		this.id = generateRandomHexId();
	}

	static fromPieceFenSymbol(symbol: PieceSymbol): Piece {
		return new Piece(symbol.toLowerCase() as PieceSymbol, getPieceColorFromFenSymbol(symbol));
	}

	isKing(): boolean {
		return this.symbol === KING;
	}

	isQueen(): boolean {
		return this.symbol === QUEEN;
	}

	isKnight(): boolean {
		return this.symbol === KNIGHT;
	}

	isBishop(): boolean {
		return this.symbol === BISHOP;
	}

	isRook(): boolean {
		return this.symbol === ROOK;
	}

	isPawn(): boolean {
		return this.symbol === PAWN;
	}

	isWhite(): boolean {
		return this.color === WHITE;
	}

	isBlack(): boolean {
		return this.color === BLACK;
	}

	symbolEquals(piece: Piece): boolean {
		return this.symbol === piece.symbol;
	}

	colorEquals(piece: Piece): boolean {
		return this.color === piece.color;
	}

	toFenSymbol(): string {
		if (this.color === WHITE) {
			return this.symbol.toUpperCase();
		}

		return this.symbol;
	}

	copy(): Piece {
		return new Piece(this.symbol, this.color);
	}
}

function generateRandomHexId(length = 32): string {
	if (length % 2 !== 0) {
		throw new Error('Hex ID length must be even.');
	}

	const characters = '0123456789ABCDEF';
	let result = '';

	for (let i = 0; i < length; i++) {
		const randomIndex = Math.floor(Math.random() * characters.length);
		result += characters.charAt(randomIndex);
	}

	return result;
}

function getPieceColorFromFenSymbol(pieceSymbol: PieceSymbol): Color {
	if (/^[prnbqk]$/.test(pieceSymbol)) {
		return 'b';
	}
	if (/^[PRNBQK]$/.test(pieceSymbol)) {
		return 'w';
	}
	throw new Error('invalid color');
}
