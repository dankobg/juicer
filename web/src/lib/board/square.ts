import { BLACK, FILES, RANKS, WHITE } from './common';
import type { Piece } from './piece';
import type { Col, Color, Coordinate, Rank, File, Row } from './types';

export class Square {
	color: Color;
	row: Row;
	col: Col;

	constructor(
		public squareIdx: number,
		public piece: Piece | null
	) {
		this.color = getSquareColor(squareIdx);
		const { row, col } = getRowAndCol(squareIdx);
		this.row = row;
		this.col = col;
	}

	static fromCoord(coord: Coordinate): Square {
		const f = coord[0] as File;
		const r = Number.parseInt(coord[1]) as Rank;

		const rowIndex = 7 - RANKS.indexOf(r);
		const colIndex = FILES.indexOf(f);
		const squareIdx = rowIndex * 8 + colIndex;
		return new Square(squareIdx, null);
	}

	get file(): File {
		return FILES[this.col];
	}

	get rank(): Rank {
		return RANKS[7 - this.row];
	}

	get coord(): Coordinate {
		return `${this.file}${this.rank}`;
	}

	isLight(): boolean {
		return this.color === WHITE;
	}

	isDark(): boolean {
		return this.color === BLACK;
	}

	hasPiece(): boolean {
		return this.piece !== null;
	}

	isEmpty(): boolean {
		return this.piece === null;
	}

	hasFriendlyPiece(currentTurn: Color): boolean {
		return this.hasPiece() && this.piece?.color === currentTurn;
	}

	hasEnemyPiece(currentTurn: Color): boolean {
		return this.hasPiece() && this.piece?.color !== currentTurn;
	}

	isEmptyOrHasEnemyPiece(currentTurn: Color): boolean {
		return this.isEmpty() || this.hasEnemyPiece(currentTurn);
	}

	equals(square: Square): boolean {
		return this.squareIdx === square.squareIdx;
	}

	toString(): string {
		if (this.piece !== null) {
			return this.piece.toFenSymbol();
		}

		return '-';
	}

	copy(): Square {
		return new Square(this.squareIdx, this.piece);
	}
}

function getSquareColor(squareIdx: number): Color {
	if ((squareIdx / 8) % 2 === squareIdx % 2) {
		return BLACK;
	}
	return WHITE;
}

function getRowAndCol(squareIdx: number): { row: Row; col: Col } {
	const row = Math.floor(squareIdx / 8) as Row;
	const col = (squareIdx % 8) as Col;
	return { row, col };
}
