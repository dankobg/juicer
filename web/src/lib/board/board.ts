import { BLACK, FEN_EMPTY_POSITION, WHITE } from './common';
import { validateFen } from './fen';
import { Square } from './square';
import { Piece } from './piece';
import type { Color, PieceSymbol } from './types';

export class Board {
	squares!: Square[];

	constructor(
		public fen: string = FEN_EMPTY_POSITION,
		public orientation: Color = WHITE
	) {
		this.loadFromFen(fen);
	}

	loadFromFen(fen = FEN_EMPTY_POSITION): void {
		const meta = validateFen(fen);
		this.squares = meta.squares;
	}

	print(): string {
		return printBoard(this.squares, this.orientation);
	}

	setOrientation(orientation: Color): void {
		this.orientation = orientation;
	}

	flipOrientation(): void {
		if (this.orientation === WHITE) {
			this.setOrientation(BLACK);
		} else {
			this.setOrientation(WHITE);
		}
	}

	clear(): void {
		this.fen = FEN_EMPTY_POSITION;
		this.orientation = WHITE;
		this.squares = fenToBoard(FEN_EMPTY_POSITION);
	}

	copy(): Board {
		return new Board(this.fen);
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

	// fen += ' w KQkq - 0 1';

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
