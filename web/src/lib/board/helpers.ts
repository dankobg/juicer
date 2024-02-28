import type { Color, PieceSymbol } from './types';
import { BLACK, WHITE } from './common';
import { Piece } from './piece';
import { Square } from './square';
import { CastleRights } from './types';

export class CastleRightsHelper {
	constructor(public cr: CastleRights) {}

	toString(): string {
		let fenCastle = '';

		if (this.cr === CastleRights.None) {
			fenCastle += '-';
		}

		if ((this.cr & CastleRights.WhiteKingSide) !== 0) {
			fenCastle += 'K';
		}
		if ((this.cr & CastleRights.WhiteQueenSide) !== 0) {
			fenCastle += 'Q';
		}
		if ((this.cr & CastleRights.BlackKingSide) !== 0) {
			fenCastle += 'k';
		}
		if ((this.cr & CastleRights.BlackQueenSide) !== 0) {
			fenCastle += 'q';
		}

		return fenCastle;
	}
}

export function fenToBoard(fen: string): Square[] {
	const squares: Square[] = [];
	const [boardPosition] = fen.split(' ');

	let rowIndex = 0;
	let colIndex = 0;

	for (const char of boardPosition) {
		if (char === '/') {
			rowIndex++;
			colIndex = 0;
		} else if (Number.isNaN(Number(char))) {
			const squareIdx = rowIndex * 8 + colIndex;

			if (char === char.toUpperCase()) {
				squares[squareIdx] = new Square(squareIdx, new Piece(char.toUpperCase() as PieceSymbol, WHITE));
			} else {
				squares[squareIdx] = new Square(squareIdx, new Piece(char.toLowerCase() as PieceSymbol, BLACK));
			}

			colIndex++;
		} else {
			const numEmptySquares = Number.parseInt(char);

			for (let i = 0; i < numEmptySquares; i++) {
				const squareIdx = rowIndex * 8 + colIndex;
				squares[squareIdx] = new Square(squareIdx, null);
				colIndex++;
			}
		}
	}

	return squares;
}

export function boardToFen(board: Square[]): string {
	let fen = '';
	let emptySquareCount = 0;

	for (let i = 0; i < board.length; i++) {
		const sq = board[i];

		if (sq.isEmpty()) {
			emptySquareCount++;
		} else {
			if (emptySquareCount > 0) {
				fen += emptySquareCount.toString();
				emptySquareCount = 0;
			}

			fen += sq.piece?.toFenSymbol();
		}
		if ((i + 1) % 8 === 0 && i < 63) {
			if (emptySquareCount > 0) {
				fen += emptySquareCount.toString();
				emptySquareCount = 0;
			}

			fen += '/';
		}
	}

	return fen;
}

export function printBoard(squares: Square[], orientation: Color): string {
	let s = '   +------------------------+\n';

	for (let i = 0; i < 64; i++) {
		let sq: Square = squares[i];
		let rank: number = 8 - Math.floor(i / 8);

		if (orientation === BLACK) {
			const flippedBoard = squares.toReversed();
			sq = flippedBoard[i];
			rank = Math.floor(i / 8) + 1;
		}

		if (i % 8 === 0) {
			s += ` ${rank} |`;
		}

		s += ` ${sq.toString()} `;

		if (i % 8 === 7) {
			s += '| \n';
		}
	}

	s += '   +------------------------+\n';

	if (orientation === WHITE) {
		s += '     a  b  c  d  e  f  g  h';
	} else {
		s += '     h  g  f  e  d  c  b  a';
	}

	return s;
}

export function getPieceColorFromFenSymbol(pieceSymbol: PieceSymbol): Color {
	if (/^[prnbqk]$/.test(pieceSymbol)) {
		return 'b';
	}
	if (/^[PRNBQK]$/.test(pieceSymbol)) {
		return 'w';
	}
	throw new Error('invalid color');
}

export function generateRandomHexId(length = 32): string {
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
