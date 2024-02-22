export type Color = 'w' | 'b';

export type Row = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
export type Col = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;

export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';

export type Coordinate = `${File}${Rank}`;

export type PieceSymbol = 'p' | 'n' | 'b' | 'r' | 'q' | 'k';
export type WhitePieceSymbol = Uppercase<PieceSymbol>;
export type BlackPieceSymbol = Lowercase<PieceSymbol>;

export type PromotionPieceSymbol = Exclude<PieceSymbol, 'p' | 'k'>;
export type PromotionWhitePieceSymbol = Uppercase<PromotionPieceSymbol>;
export type PromotionBlackPieceSymbol = Lowercase<PromotionPieceSymbol>;

export enum CastleRights {
	None = 0,
	WhiteKingSide = 1,
	WhiteQueenSide = 2,
	BlackKingSide = 4,
	BlackQueenSide = 8,
}

export type DragPosition = {
	initialX: number;
	initialY: number;
	dx: number;
	dy: number;
};

// export class CastleRightsHelper {
// 	constructor(public cr: CastleRights) {}

// 	toString(): string {
// 		let s = '';

// 		if (this.cr === CastleRights.None) {
// 			s += '-';
// 		}

// 		if ((this.cr & CastleRights.WhiteKingSide) !== 0) {
// 			s += 'K';
// 		}
// 		if ((this.cr & CastleRights.WhiteQueenSide) !== 0) {
// 			s += 'Q';
// 		}
// 		if ((this.cr & CastleRights.BlackKingSide) !== 0) {
// 			s += 'k';
// 		}
// 		if ((this.cr & CastleRights.BlackQueenSide) !== 0) {
// 			s += 'q';
// 		}

// 		return s;
// 	}
// }

// export function getRowAndCol(squareIdx: number): { row: Row; col: Col } {
// 	const row = Math.floor(squareIdx / 8) as Row;
// 	const col = (squareIdx % 8) as Col;
// 	return { row, col };
// }

// export function printBoard(squares: Square[], orientation: Color): string {
// 	let s = '   +------------------------+\n';

// 	for (let i = 0; i < 64; i++) {
// 		let sq: Square = squares[i];
// 		let rank: number = 8 - Math.floor(i / 8);

// 		if (orientation === BLACK) {
// 			const flippedBoard = squares.toReversed();
// 			sq = flippedBoard[i];
// 			rank = Math.floor(i / 8) + 1;
// 		}

// 		if (i % 8 === 0) {
// 			s += ` ${rank} |`;
// 		}

// 		s += ` ${sq.toString()} `;

// 		if (i % 8 === 7) {
// 			s += '| \n';
// 		}
// 	}

// 	s += '   +------------------------+\n';

// 	if (orientation === WHITE) {
// 		s += '     a  b  c  d  e  f  g  h';
// 	} else {
// 		s += '     h  g  f  e  d  c  b  a';
// 	}

// 	return s;
// }

// export function fenToBoard(fen: string): Square[] {
// 	const squares: Square[] = [];
// 	const [boardPosition] = fen.split(' ');

// 	let rowIndex = 0;
// 	let colIndex = 0;

// 	for (const char of boardPosition) {
// 		if (char === '/') {
// 			rowIndex++;
// 			colIndex = 0;
// 		} else if (Number.isNaN(Number(char))) {
// 			const squareIdx = rowIndex * 8 + colIndex;

// 			if (char === char.toUpperCase()) {
// 				squares[squareIdx] = new Square(squareIdx, new Piece(char.toUpperCase() as PieceSymbol, WHITE));
// 			} else {
// 				squares[squareIdx] = new Square(squareIdx, new Piece(char.toLowerCase() as PieceSymbol, BLACK));
// 			}

// 			colIndex++;
// 		} else {
// 			const numEmptySquares = Number.parseInt(char);

// 			for (let i = 0; i < numEmptySquares; i++) {
// 				const squareIdx = rowIndex * 8 + colIndex;
// 				squares[squareIdx] = new Square(squareIdx, null);
// 				colIndex++;
// 			}
// 		}
// 	}

// 	return squares;
// }

// export function boardToFen(board: Square[]): string {
// 	let fen = '';
// 	let emptySquareCount = 0;

// 	for (let i = 0; i < board.length; i++) {
// 		const sq = board[i];

// 		if (sq.isEmpty()) {
// 			emptySquareCount++;
// 		} else {
// 			if (emptySquareCount > 0) {
// 				fen += emptySquareCount.toString();
// 				emptySquareCount = 0;
// 			}

// 			fen += sq.piece?.toFenSymbol();
// 		}
// 		if ((i + 1) % 8 === 0 && i < 63) {
// 			if (emptySquareCount > 0) {
// 				fen += emptySquareCount.toString();
// 				emptySquareCount = 0;
// 			}

// 			fen += '/';
// 		}
// 	}

// 	// fen += ' w KQkq - 0 1';

// 	return fen;
// }

// function getSquareColor(squareIdx: number): Color {
// 	if ((squareIdx / 8) % 2 === squareIdx % 2) {
// 		return BLACK;
// 	}
// 	return WHITE;
// }

// function generateRandomHexId(length = 32): string {
// 	if (length % 2 !== 0) {
// 		throw new Error('Hex ID length must be even.');
// 	}

// 	const characters = '0123456789ABCDEF';
// 	let result = '';

// 	for (let i = 0; i < length; i++) {
// 		const randomIndex = Math.floor(Math.random() * characters.length);
// 		result += characters.charAt(randomIndex);
// 	}

// 	return result;
// }

// function getPieceColorFromFenSymbol(pieceSymbol: PieceFenSymbol): Color {
// 	if (/^[prnbqk]$/.test(pieceSymbol)) {
// 		return 'b';
// 	}
// 	if (/^[PRNBQK]$/.test(pieceSymbol)) {
// 		return 'w';
// 	}
// 	throw new Error('invalid color');
// }

// export class Board {
// 	squares!: Square[];

// 	constructor(
// 		public fen: string = FEN_EMPTY_POSITION,
// 		public orientation: Color = WHITE
// 	) {
// 		this.loadFromFen(fen);
// 	}

// 	loadFromFen(fen = FEN_EMPTY_POSITION): void {
// 		const meta = validateFen(fen);
// 		this.squares = meta.squares;
// 	}

// 	print(): string {
// 		return printBoard(this.squares, this.orientation);
// 	}

// 	setOrientation(orientation: Color): void {
// 		this.orientation = orientation;
// 	}

// 	flipOrientation(): void {
// 		if (this.orientation === WHITE) {
// 			this.setOrientation(BLACK);
// 		} else {
// 			this.setOrientation(WHITE);
// 		}
// 	}

// 	clear(): void {
// 		this.fen = FEN_EMPTY_POSITION;
// 		this.orientation = WHITE;
// 		this.squares = fenToBoard(FEN_EMPTY_POSITION);
// 	}

// 	copy(): Board {
// 		return new Board(this.fen);
// 	}
// }

// export class Square {
// 	color: Color;
// 	row: Row;
// 	col: Col;

// 	constructor(
// 		public squareIdx: number,
// 		public piece: Piece | null
// 	) {
// 		this.color = getSquareColor(squareIdx);
// 		const { row, col } = getRowAndCol(squareIdx);
// 		this.row = row;
// 		this.col = col;
// 	}

// 	static fromCoord(coord: Coordinate): Square {
// 		const rowIndex = RANK_CHARS.indexOf(coord[0]);
// 		const colIndex = FILE_CHARS.indexOf(coord[1]);
// 		const squareIdx = rowIndex * 8 + colIndex;
// 		return new Square(squareIdx, null);
// 	}

// 	get file(): File {
// 		const n = this.squareIdx % BOARD_SIZE;
// 		return FILE_CHARS.slice(n, n + 1) as File;
// 	}

// 	get rank(): Rank {
// 		const n = this.squareIdx / BOARD_SIZE;
// 		return Number.parseInt(RANK_CHARS.slice(n, n + 1)) as Rank;
// 	}

// 	get coord(): Coordinate {
// 		return `${this.file}${this.rank}`;
// 	}

// 	isLight(): boolean {
// 		return this.color === WHITE;
// 	}

// 	isDark(): boolean {
// 		return this.color === BLACK;
// 	}

// 	hasPiece(): boolean {
// 		return this.piece !== null;
// 	}

// 	isEmpty(): boolean {
// 		return this.piece === null;
// 	}

// 	hasFriendlyPiece(currentTurn: Color): boolean {
// 		return this.hasPiece() && this.piece?.color === currentTurn;
// 	}

// 	hasEnemyPiece(currentTurn: Color): boolean {
// 		return this.hasPiece() && this.piece?.color !== currentTurn;
// 	}

// 	isEmptyOrHasEnemyPiece(currentTurn: Color): boolean {
// 		return this.isEmpty() || this.hasEnemyPiece(currentTurn);
// 	}

// 	equals(square: Square): boolean {
// 		return this.squareIdx === square.squareIdx;
// 	}

// 	toString(): string {
// 		if (this.piece !== null) {
// 			return this.piece.toFenSymbol();
// 		}

// 		return '-';
// 	}

// 	copy(): Square {
// 		return new Square(this.squareIdx, this.piece);
// 	}
// }
