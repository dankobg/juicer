export type Color = 'w' | 'b';

export type Row = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
export type Col = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;

export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';

export type Coordinate = `${File}${Rank}`;

export type PieceSymbol = 'p' | 'n' | 'b' | 'r' | 'q' | 'k';
export type WhitePieceSymbol = Uppercase<PieceSymbol>;
export type BlackPieceSymbol = Lowercase<PieceSymbol>;
export type PieceFenSymbol = WhitePieceSymbol | BlackPieceSymbol;

export type PromotionPieceSymbol = Exclude<PieceSymbol, 'p' | 'k'>;
export type PromotionWhitePieceSymbol = Uppercase<PromotionPieceSymbol>;
export type PromotionBlackPieceSymbol = Lowercase<PromotionPieceSymbol>;
export type PromotionPieceFenSymbol = PromotionWhitePieceSymbol | PromotionBlackPieceSymbol;

export type DropOffBoardAction = 'trash' | 'snapback';
export type DragPosition = { initialX: number; initialY: number; dx: number; dy: number };

export type Vec = [number, number];
export type SquareData = { index: number; coord: Coordinate; file: File; rank: Rank; row: Row; col: Col };

export type RowCol = { row: Row; col: Col };
export type RankFile = { rank: Rank; file: File };

export type EventData = { clientX: number; clientY: number; offsetX: number; offsetY: number; elm: HTMLDivElement };

// ----------------------------------------------------------------------------------------------------------------

export const WHITE: Color = 'w';
export const BLACK: Color = 'b';
export const COLORS: [Color, Color] = [WHITE, BLACK];

export const PAWN: PieceSymbol = 'p';
export const KNIGHT: PieceSymbol = 'n';
export const BISHOP: PieceSymbol = 'b';
export const ROOK: PieceSymbol = 'r';
export const QUEEN: PieceSymbol = 'q';
export const KING: PieceSymbol = 'k';

export const WHITE_PAWN: WhitePieceSymbol = 'P';
export const WHITE_KNIGHT: WhitePieceSymbol = 'N';
export const WHITE_BISHOP: WhitePieceSymbol = 'B';
export const WHITE_ROOK: WhitePieceSymbol = 'R';
export const WHITE_QUEEN: WhitePieceSymbol = 'Q';
export const WHITE_KING: WhitePieceSymbol = 'K';

export const BLACK_PAWN: BlackPieceSymbol = 'p';
export const BLACK_KNIGHT: BlackPieceSymbol = 'n';
export const BLACK_BISHOP: BlackPieceSymbol = 'b';
export const BLACK_ROOK: BlackPieceSymbol = 'r';
export const BLACK_QUEEN: BlackPieceSymbol = 'q';
export const BLACK_KING: BlackPieceSymbol = 'k';

export const WHITE_PIECES_FEN: WhitePieceSymbol[] = ['R', 'N', 'B', 'Q', 'K', 'P'];
export const BLACK_PIECES_FEN: BlackPieceSymbol[] = ['r', 'n', 'b', 'q', 'k', 'p'];
export const ALL_PIECES_FEN: PieceFenSymbol[] = ['R', 'N', 'B', 'Q', 'K', 'P', 'r', 'n', 'b', 'q', 'k', 'p'];

export const FILES: File[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
export const RANKS: Rank[] = [1, 2, 3, 4, 5, 6, 7, 8];

export const FEN_EMPTY = '8/8/8/8/8/8/8/8';
export const FEN_START = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';

export const BOARD_RANKS: number = 8;
export const BOARD_FILES: number = 8;
export const BOARD_TOTAL_SQUARES: number = BOARD_RANKS * BOARD_FILES;

export const dragPositionZero: DragPosition = { initialX: 0, initialY: 0, dx: 0, dy: 0 };

// prettier-ignore
export const COORDINATES: Coordinate[] = [
	'a8','b8','c8','d8','e8','f8','g8','h8',
	'a7','b7','c7','d7','e7','f7','g7','h7',
	'a6','b6','c6','d6','e6','f6','g6','h6',
	'a5','b5','c5','d5','e5','f5','g5','h5',
	'a4','b4','c4','d4','e4','f4','g4','h4',
	'a3','b3','c3','d3','e3','f3','g3','h3',
	'a2','b2','c2','d2','e2','f2','g2','h2',
	'a1','b1','c1','d1','e1','f1','g1','h1',
];

export const NO_SQUARE = -1;

export class Board {
	squares!: Square[];

	constructor(
		public fen: string = FEN_EMPTY,
		public orientation: Color = WHITE
	) {
		this.loadFromFen(fen);
	}

	static indexInBounds(squareIndex: number): boolean {
		return squareIndex >= 0 && squareIndex <= 63;
	}

	static fenToSquares(fen: string): Square[] {
		const squares: Square[] = [];
		const [boardPosition] = fen.split(' ');

		let row = 0;
		let col = 0;

		for (const char of boardPosition) {
			if (char === '/') {
				row++;
				col = 0;
			} else if (Number.isNaN(Number.parseInt(char))) {
				const squareIndex = Square.calcIndex(row as Row, col as Col);
				squares[squareIndex] = Square.newFromIndex(squareIndex, new Piece(char as PieceFenSymbol));
				col++;
			} else {
				const numEmptySquares = Number.parseInt(char);

				for (let i = 0; i < numEmptySquares; i++) {
					const squareIndex = Square.calcIndex(row as Row, col as Col);
					squares[squareIndex] = Square.newFromIndex(squareIndex, null);
					col++;
				}
			}
		}

		return squares;
	}

	static squaresToFen(board: Square[]): string {
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

	loadFromFen(fen = FEN_EMPTY): void {
		this.squares = Board.fenToSquares(fen);
	}

	getFen(): string {
		return Board.squaresToFen(this.squares);
	}

	setOrientation(orientation: Color): void {
		this.orientation = orientation;
	}

	flipOrientation(): void {
		this.orientation = this.orientation === WHITE ? BLACK : WHITE;
		this.squares = this.squares.toReversed();
	}

	print(): string {
		let s = '   +------------------------+\n';

		for (let i = 0; i < this.squares.length; i++) {
			let rank: number = 8 - Math.floor(i / 8);

			if (this.orientation === BLACK) {
				rank = Math.floor(i / BOARD_RANKS) + 1;
			}

			if (i % BOARD_RANKS === 0) {
				s += ` ${rank} |`;
			}

			s += ` ${this.squares[i].toFenSymbol()} `;

			if (i % BOARD_FILES === 7) {
				s += '| \n';
			}
		}

		s += '   +------------------------+\n';

		const letters = this.orientation === WHITE ? FILES : FILES.toReversed();
		s += '     ' + letters.join('  ');

		return s;
	}

	setPiece(squareIndex: number, piece: Piece): void {
		if (this.squares[squareIndex].isEmpty()) {
			this.squares[squareIndex].piece = piece;
		}
	}

	setOrReplacePiece(squareIndex: number, piece: Piece): void {
		this.squares[squareIndex].piece = piece;
	}

	getPiece(squareIndex: number): Piece | null {
		return this.squares[squareIndex]?.piece ?? null;
	}

	removePiece(squareIndex: number): void {
		if (this.getPiece(squareIndex)) {
			this.squares[squareIndex].piece = null;
		}
	}

	movePiece(srcSquareIndex: number, destSquareIndex: number): void {
		if (!this.getPiece(srcSquareIndex)) {
			console.warn(`no piece on src square: ${srcSquareIndex}`);
			return;
		}

		this.squares[destSquareIndex].piece = this.squares[srcSquareIndex].piece;
		this.squares[srcSquareIndex].piece = null;
	}

	getPiecesCount(): number {
		return this.squares.filter(sq => sq.hasPiece()).length;
	}

	clear(): void {
		this.fen = FEN_EMPTY;
		this.orientation = WHITE;
		this.squares = Board.fenToSquares(FEN_EMPTY);
	}

	clone(): Board {
		return new Board(this.fen);
	}
}

export class Square {
	index: number;
	row: Row;
	col: Col;
	file: File;
	rank: Rank;
	color: Color;

	constructor(
		public coordinate: Coordinate,
		public piece: Piece | null
	) {
		const { index, row, col, file, rank } = Square.getDataFromCoord(coordinate);
		this.index = index;
		this.row = row;
		this.col = col;
		this.file = file;
		this.rank = rank;
		this.color = Square.getSquareColor(index);
	}

	static coordToIndex(coord: Coordinate): number {
		return COORDINATES.indexOf(coord);
	}

	static indexToCoord(index: number): Coordinate | null {
		return COORDINATES[index] ?? null;
	}

	static calcIndex(row: Row, col: Col): number {
		return row * BOARD_RANKS + col;
	}

	static rankFileToRowCol(rank: Rank, file: File): RowCol {
		const row = (7 - RANKS.indexOf(rank)) as Row;
		const col = FILES.indexOf(file) as Col;
		return { row, col };
	}

	static rowColToRankFile(row: Row, col: Col): RankFile {
		const file = FILES[col];
		const rank = RANKS[7 - row];
		return { file, rank };
	}

	static getDataFromCoord(coord: Coordinate): SquareData {
		const file = coord[0] as File;
		const rank = Number.parseInt(coord[1]) as Rank;
		const { row, col } = Square.rankFileToRowCol(rank, file);
		const index = Square.calcIndex(row, col);
		return { index, coord, file, rank, row, col };
	}

	static getDataFromIndex(index: number): SquareData {
		const row = Math.floor(index / 8) as Row;
		const col = (index % 8) as Col;
		const { rank, file } = Square.rowColToRankFile(row, col);
		const coord = `${file}${rank}` as const;
		return { index, coord, file, rank, row, col };
	}

	static getDistanceRowCol(srcSquareIndex: number, destSquareIndex: number): Vec {
		const src = Square.getDataFromIndex(srcSquareIndex);
		const dest = Square.getDataFromIndex(destSquareIndex);
		return [src.row - dest.row, src.col - dest.col];
	}

	static getSquareColor(index: number): Color {
		if ((index / 8) % 2 === index % 2) {
			return BLACK;
		}
		return WHITE;
	}

	static newFromIndex(squareIndex: number, piece: Piece | null): Square {
		const coord = Square.indexToCoord(squareIndex);
		if (!coord) {
			throw new Error('invalid square index');
		}
		return new Square(coord, piece);
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

	equals(square: Square): boolean {
		return this.coordinate === square.coordinate;
	}

	toFenSymbol(): string {
		return this.piece?.toFenSymbol() ?? '-';
	}

	toString(): string {
		return this.toFenSymbol();
	}

	clone(): Square {
		return new Square(this.coordinate, this.piece);
	}
}

export class Piece {
	id: string;
	symbol: PieceSymbol;
	color: Color;

	constructor(public fenSymbol: PieceFenSymbol) {
		this.id = this.generateId();
		this.symbol = this.fenSymbol.toLowerCase() as PieceSymbol;
		this.color = this.fenSymbol === this.fenSymbol.toUpperCase() ? WHITE : BLACK;
	}

	private generateId(length = 32): string {
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

	static pairToFenSymbol(symbol: PieceSymbol, color: Color): PieceFenSymbol {
		return (color === WHITE ? symbol.toUpperCase() : symbol) as PieceFenSymbol;
	}

	static fenSymbolToPair(fenSymbol: PieceFenSymbol): [PieceSymbol, Color] {
		return [fenSymbol.toLowerCase() as PieceSymbol, fenSymbol === fenSymbol.toUpperCase() ? WHITE : BLACK];
	}

	static newFromPair(symbol: PieceSymbol, color: Color): Piece {
		return new Piece(Piece.pairToFenSymbol(symbol, color));
	}

	equals(piece: Piece): boolean {
		return this.id === piece.id;
	}

	toFenSymbol(): string {
		return this.color === WHITE ? this.symbol.toUpperCase() : this.symbol;
	}

	toString(): string {
		return this.toFenSymbol();
	}

	clone(): Piece {
		const clone = new Piece(this.fenSymbol);
		clone.id = this.id;
		return clone;
	}
}
