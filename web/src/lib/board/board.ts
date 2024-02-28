import { BLACK, FEN_EMPTY_POSITION, WHITE } from './common';
import { validateFen } from './fen';
import { fenToBoard, printBoard } from './helpers';
import { Square } from './square';
import type { Color } from './types';

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

		this.squares.reverse();
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
