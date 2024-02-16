import { browser } from '$app/environment';

export type Row = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
export type Col = Row;
export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';
export type Coordinate = `${File}${Rank}`;

export type Color = 'w' | 'b';
export type PieceMovement = 'click' | 'drag' | 'either';

// export type BoardSetStyle = 'blue' | 'brown' | 'ic' | 'purple' | 'green';
export type BoardTheme = string;

export type PieceTheme =
	| 'alphacalifornia'
	| 'cburnett'
	| 'chess7'
	| 'companion'
	| 'dubrovny'
	| 'fresca'
	| 'governor'
	| 'icpieces'
	| 'leipzig'
	| 'libra'
	| 'merida'
	| 'pirouetti'
	| 'reillycraig'
	| 'shapes'
	| 'stauntyanarcandy'
	| 'cardinal'
	| 'celtic'
	| 'chessnut'
	| 'disguised'
	| 'fantasy'
	| 'gioco'
	| 'horsey'
	| 'kosal'
	| 'letter'
	| 'maestro'
	| 'mono'
	| 'pixel'
	| 'riohacha'
	| 'spatial'
	| 'tatiana';

export type BlackPieceFenSymbol = 'p' | 'n' | 'b' | 'r' | 'q' | 'k';
export type WhitePieceFenSymbol = 'P' | 'N' | 'B' | 'R' | 'Q' | 'K';
export type PieceFenSymbol = BlackPieceFenSymbol | WhitePieceFenSymbol;
export type PieceSymbol = Lowercase<PieceFenSymbol>;
export type PromotionPieceFenSymbol = Exclude<PieceFenSymbol, 'p' | 'k' | 'P' | 'K'>;
export type PromotionPieceSymbol = Lowercase<PromotionPieceFenSymbol>;

export type Position = {
	turn: Color;
	board: Board;
};

export type PieceCustomEvent = CustomEvent<{
	piece: Piece;
	squareIdx: number;
	elm: HTMLDivElement;
	clientX: number;
	clientY: number;
	offsetX: number;
	offsetY: number;
}>;

export type SquareCustomEvent = CustomEvent<{
	square: Square;
	elm: HTMLDivElement;
	clientX: number;
	clientY: number;
	offsetX: number;
	offsetY: number;
}>;

// ********************************************************************************
// ********************************************************************************
// ********************************************************************************

export const DEFAULT_APPEAR_SPEED = 200;
export const DEFAULT_MOVE_SPEED = 200;
export const DEFAULT_SNAPBACK_SPEED = 60;
export const DEFAULT_SNAP_SPEED = 30;
export const DEFAULT_TRASH_SPEED = 100;

export const DEFAULT_BOARD_THEME: BoardTheme = '/images/board/svg/brown.svg';
export const DEFAULT_PIECE_THEME: PieceTheme = 'cburnett';

export const WHITE: Color = 'w';
export const BLACK: Color = 'b';

export const EMPTY = '-';

export const PAWN: PieceSymbol = 'p';
export const KNIGHT: PieceSymbol = 'n';
export const BISHOP: PieceSymbol = 'b';
export const ROOK: PieceSymbol = 'r';
export const QUEEN: PieceSymbol = 'q';
export const KING: PieceSymbol = 'k';

export const WHITE_PAWN: WhitePieceFenSymbol = 'P';
export const WHITE_KNIGHT: WhitePieceFenSymbol = 'N';
export const WHITE_BISHOP: WhitePieceFenSymbol = 'B';
export const WHITE_ROOK: WhitePieceFenSymbol = 'R';
export const WHITE_QUEEN: WhitePieceFenSymbol = 'Q';
export const WHITE_KING: WhitePieceFenSymbol = 'K';

export const BLACK_PAWN: BlackPieceFenSymbol = 'p';
export const BLACK_KNIGHT: BlackPieceFenSymbol = 'n';
export const BLACK_BISHOP: BlackPieceFenSymbol = 'b';
export const BLACK_ROOK: BlackPieceFenSymbol = 'r';
export const BLACK_QUEEN: BlackPieceFenSymbol = 'q';
export const BLACK_KING: BlackPieceFenSymbol = 'k';

export const FEN_EMPTY_POSITION = '8/8/8/8/8/8/8/8';
export const FEN_STARTING_POSITION = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';

export const whitePiecesFen = ['R', 'N', 'B', 'Q', 'K', 'P'];
export const blackPiecesFen = ['r', 'n', 'b', 'q', 'k', 'p'];
export const allPiecesFen = whitePiecesFen.concat(blackPiecesFen);

const FILE_TO_COLUMN = new Map<File, Col>([
	['a', 0],
	['b', 1],
	['c', 2],
	['d', 3],
	['e', 4],
	['f', 5],
	['g', 6],
	['h', 7],
]);

const COLUMN_TO_FILE = new Map<Col, File>([
	[0, 'a'],
	[1, 'b'],
	[2, 'c'],
	[3, 'd'],
	[4, 'e'],
	[5, 'f'],
	[6, 'g'],
	[7, 'h'],
]);

const RANK_TO_ROW = new Map<Rank, Row>([
	[1, 7],
	[2, 6],
	[3, 5],
	[4, 4],
	[5, 3],
	[6, 2],
	[7, 1],
	[8, 0],
]);

const ROW_TO_RANK = new Map<Row, Rank>([
	[0, 8],
	[1, 7],
	[2, 6],
	[3, 5],
	[4, 4],
	[5, 3],
	[6, 2],
	[7, 1],
]);

export function swapColor(color: Color): Color {
	return color === WHITE ? BLACK : WHITE;
}

export function convertsquareIdxToCoordinate(squareIdx: number): Coordinate {
	const row = Math.floor(squareIdx / 8) as Row;
	const col = (squareIdx % 8) as Col;

	return convertRowAndColumnToCoordinate(row, col);
}

export function convertFileToColumn(file: File): Col {
	return FILE_TO_COLUMN.get(file) as Col;
}

export function convertColumnToFile(column: Col): File {
	return COLUMN_TO_FILE.get(column) as File;
}

export function convertRankToRow(rank: Rank): Row {
	return RANK_TO_ROW.get(rank) as Row;
}

export function convertRowToRank(row: Row): Rank {
	return ROW_TO_RANK.get(row) as Rank;
}

export function convertCoordinateToFileAndRank(coordinate: Coordinate): [File, Rank] {
	const chars = coordinate.split('');

	const file = chars[0] as File;
	const rank = Number.parseInt(chars[1]) as Rank;

	return [file, rank];
}

export function convertCoordinateToRowAndColumn(coordinate: Coordinate): [Row, Col] {
	const chars = coordinate.split('');

	const file = chars[0] as File;
	const rank = Number.parseInt(chars[1]) as Rank;

	const row = convertRankToRow(rank);
	const col = convertFileToColumn(file);

	return [row, col];
}

export function convertFileAndRankToCoordinate(file: File, rank: Rank): Coordinate {
	return `${file}${rank}`;
}

export function convertRowAndColumnToCoordinate(row: Row, col: Col): Coordinate {
	const file = convertColumnToFile(col);
	const rank = convertRowToRank(row);

	return convertFileAndRankToCoordinate(file, rank);
}

export function getRowAndCol(squareIdx: number): {
	row: Row;
	col: Col;
} {
	const row = Math.floor(squareIdx / 8) as Row;
	const col = (squareIdx % 8) as Col;

	return { row, col };
}

export function getSquareColor(squareIdx: number): Color {
	const { row, col } = getRowAndCol(squareIdx);
	return (row + col) % 2 === 0 ? WHITE : BLACK;
}

export function getPieceColorFromFenSymbol(pieceSymbol: PieceFenSymbol): Color {
	if (/^[prnbqk]$/.test(pieceSymbol)) {
		return 'b';
	}
	if (/^[PRNBQK]$/.test(pieceSymbol)) {
		return 'w';
	}
	throw new Error('invalid color');
}

function isDigit(c: string): boolean {
	return /^[0-9]$/.test(c);
}

function validateFEN(fen: string): { ok: boolean; err: string } {
	// 1st criterion: 6 space-seperated fields?
	const tokens = fen.split(' ');
	if (tokens.length !== 6) {
		return {
			ok: false,
			err: 'Invalid FEN: must contain six space-delimited fields',
		};
	}

	// 2nd criterion: full move clock number is an integer value >= 1?
	const fullMoveClock = Number.parseInt(tokens[5], 10);
	if (Number.isNaN(fullMoveClock) || fullMoveClock <= 0) {
		return {
			ok: false,
			err: 'Invalid FEN: full move clock must be a positive integer',
		};
	}

	// 3rd criterion: half move clock is an integer >= 0?
	const halfMoveClock = Number.parseInt(tokens[4], 10);
	if (Number.isNaN(halfMoveClock) || halfMoveClock < 0) {
		return {
			ok: false,
			err: 'Invalid FEN: half move counter number must be a non-negative integer',
		};
	}

	// 4th criterion: 4th field is a valid en-passant square target or `-` if empty?
	if (!/^(-|[abcdefgh][36])$/.test(tokens[3])) {
		return { ok: false, err: 'Invalid FEN: en-passant square is invalid' };
	}

	// 5th criterion: 3th field is a valid castle-string?
	if (/[^kKqQ-]/.test(tokens[2])) {
		return { ok: false, err: 'Invalid FEN: castling rights are invalid' };
	}

	// 6th criterion: 2nd field is "w" (white) or "b" (black)?
	if (!/^(w|b)$/.test(tokens[1])) {
		return { ok: false, err: 'Invalid FEN: active color is invalid' };
	}

	// 7th criterion: 1st field contains 8 rows?
	const rows = tokens[0].split('/');
	if (rows.length !== 8) {
		return {
			ok: false,
			err: "Invalid FEN: piece data does not contain 8 '/'-delimited rows",
		};
	}

	// 8th criterion: every row is valid?
	for (let i = 0; i < rows.length; i++) {
		// check for right sum of fields AND not two numbers in succession
		let sumFields = 0;
		let previousWasNumber = false;

		for (let k = 0; k < rows[i].length; k++) {
			if (isDigit(rows[i][k])) {
				if (previousWasNumber) {
					return {
						ok: false,
						err: 'Invalid FEN: piece data is invalid (consecutive number)',
					};
				}
				sumFields += parseInt(rows[i][k], 10);
				previousWasNumber = true;
			} else {
				if (!/^[prnbqkPRNBQK]$/.test(rows[i][k])) {
					return {
						ok: false,
						err: 'Invalid FEN: piece data is invalid (invalid piece)',
					};
				}
				sumFields += 1;
				previousWasNumber = false;
			}
		}
		if (sumFields !== 8) {
			return {
				ok: false,
				err: 'Invalid FEN: piece data is invalid (too many squares in rank)',
			};
		}
	}

	if ((tokens[3][1] === '3' && tokens[1] === 'w') || (tokens[3][1] === '6' && tokens[1] === 'b')) {
		return { ok: false, err: 'Invalid FEN: illegal en-passant target square' };
	}

	const kings = [
		{ color: 'white', regex: /K/g },
		{ color: 'black', regex: /k/g },
	];

	for (const { color, regex } of kings) {
		if (!regex.test(tokens[0])) {
			return { ok: false, err: `Invalid FEN: missing ${color} king` };
		}

		if ((tokens[0].match(regex) || []).length > 1) {
			return { ok: false, err: `Invalid FEN: too many ${color} kings` };
		}
	}

	return { ok: true, err: '' };
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

			fen += sq.piece?.toPieceFenSymbol();
		}
		if ((i + 1) % 8 === 0 && i < 63) {
			if (emptySquareCount > 0) {
				fen += emptySquareCount.toString();
				emptySquareCount = 0;
			}

			fen += '/';
		}
	}

	fen += ' w KQkq - 0 1';

	return fen;
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

export function printBoardAscii(board: Square[], orientation: Color): string {
	let s = '   +------------------------+\n';

	for (let i = 0; i < 64; i++) {
		let sq: Square = board[i];
		let rank: number = 8 - Math.floor(i / 8);

		if (orientation === BLACK) {
			const flippedBoard = board.toReversed();
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

export class Board {
	squares!: Square[];
	orientation: Color = WHITE;

	constructor(
		public fen: string = FEN_STARTING_POSITION,
		opts?: {
			orientation: Color;
			boardTheme: BoardTheme;
			pieceTheme: PieceTheme;
		}
	) {
		if (opts?.orientation === BLACK) {
			this.orientation = BLACK;
		}

		this.initBoardTheme(opts?.boardTheme || DEFAULT_BOARD_THEME);
		this.initPieceTheme(opts?.pieceTheme || DEFAULT_PIECE_THEME);
		this.loadFromFen(fen);
	}

	initBoardTheme(boardTheme: BoardTheme): void {
		if (browser) {
			const rootElm = document.querySelector(':root') as HTMLElement;

			if (rootElm) {
				rootElm.style.setProperty('--board-theme', `url('${boardTheme}')`);
			}
		}
	}

	initPieceTheme(pieceTheme: PieceTheme): void {
		if (browser) {
			const rootElm = document.querySelector(':root') as HTMLElement;

			const sets = new Map<string, string>([
				['--br-theme', `url('/images/piece/${pieceTheme}/bR.svg')`],
				['--bb-theme', `url('/images/piece/${pieceTheme}/bB.svg')`],
				['--bn-theme', `url('/images/piece/${pieceTheme}/bN.svg')`],
				['--bq-theme', `url('/images/piece/${pieceTheme}/bQ.svg')`],
				['--bk-theme', `url('/images/piece/${pieceTheme}/bK.svg')`],
				['--bp-theme', `url('/images/piece/${pieceTheme}/bP.svg')`],
				['--wr-theme', `url('/images/piece/${pieceTheme}/wR.svg')`],
				['--wb-theme', `url('/images/piece/${pieceTheme}/wB.svg')`],
				['--wn-theme', `url('/images/piece/${pieceTheme}/wN.svg')`],
				['--wq-theme', `url('/images/piece/${pieceTheme}/wQ.svg')`],
				['--wk-theme', `url('/images/piece/${pieceTheme}/wK.svg')`],
				['--wp-theme', `url('/images/piece/${pieceTheme}/wP.svg')`],
			]);

			if (rootElm) {
				for (const [k, v] of sets.entries()) {
					rootElm.style.setProperty(k, v);
				}
			}
		}
	}

	loadFromFen(fen: string): void {
		this.squares = fenToBoard(fen);
	}

	clear() {
		for (let i = 0; i < this.squares.length; i++) {
			this.squares[i].piece = null;
		}

		console.log(this.squares.map(sq => sq.toString()));
	}

	print(): string {
		return printBoardAscii(this.squares, this.orientation);
	}

	setOrientation(orientation: Color): void {
		this.orientation = orientation;
	}

	setOrientationWhite(): void {
		this.setOrientation(WHITE);
	}

	setOrientationBlack(): void {
		this.setOrientation(BLACK);
	}

	flipOrientation(): void {
		if (this.orientation === WHITE) {
			this.setOrientationBlack();
		} else {
			this.setOrientationWhite();
		}
	}

	copy() {
		return new Board(this.fen);
	}
}

export class Square {
	color: Color;

	constructor(
		public squareIdx: number,
		public piece: Piece | null
	) {
		this.color = getSquareColor(squareIdx);
	}

	private get rowCol(): { row: Row; col: Col } {
		return getRowAndCol(this.squareIdx);
	}

	get row(): Row {
		return this.rowCol.row;
	}

	get col(): Col {
		return this.rowCol.col;
	}

	get file(): File {
		return convertColumnToFile(this.col);
	}

	get rank(): Rank {
		return convertRowToRank(this.row);
	}

	get coordinate(): Coordinate {
		return convertRowAndColumnToCoordinate(this.row, this.col);
	}

	isWhite(): boolean {
		return this.color === WHITE;
	}

	isBlack(): boolean {
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

	colorEquals(square: Square): boolean {
		return this.color === square.color;
	}

	toString(): string {
		if (this.hasPiece()) {
			return this.piece?.toPieceFenSymbol() as string;
		}

		return EMPTY;
	}

	copy(): Square {
		return new Square(this.squareIdx, this.piece);
	}
}

export class Piece {
	alive: boolean = true;
	id: string = '';

	constructor(
		public symbol: PieceSymbol,
		public color: Color
	) {
		this.id = generateRandomHexId();
	}

	static fromPieceFenSymbol(symbol: PieceFenSymbol): Piece | null {
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

	toPieceFenSymbol(): string {
		if (this.color === WHITE) {
			return this.symbol.toUpperCase();
		}

		return this.symbol;
	}

	copy(): Piece {
		return new Piece(this.symbol, this.color);
	}
}

export function makeTransparentDragImage(): HTMLImageElement {
	const img = new Image();
	img.src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=';
	return img;
}

export function calculateDistanceBetweenSquares(
	srcIdx: number,
	dstIdx: number,
	boardSize: number
): { dx: number; dy: number } {
	const { row: srcRow, col: srcCol } = getRowAndCol(srcIdx);
	const { row: dstRow, col: dstCol } = getRowAndCol(dstIdx);

	const c = dstCol - srcCol;
	const r = dstRow - srcRow;

	const size = boardSize / 8;

	const dx = c * size;
	const dy = r * size;

	return { dx, dy };
}

export const blackSparePieces: Piece[] = [
	new Piece('r', 'b'),
	new Piece('n', 'b'),
	new Piece('b', 'b'),
	new Piece('q', 'b'),
	new Piece('k', 'b'),
	new Piece('p', 'b'),
];

export const whiteSparePieces: Piece[] = [
	new Piece('r', 'w'),
	new Piece('n', 'w'),
	new Piece('b', 'w'),
	new Piece('q', 'w'),
	new Piece('k', 'w'),
	new Piece('p', 'w'),
];

export const allSparePieces = blackSparePieces.concat(whiteSparePieces);

export function translateElm(elm: HTMLDivElement, dx: number, dy: number): void {
	if (elm) {
		elm.style.translate = `${dx}px ${dy}px`;
	}
}

export function getSquareIdxFromDragPos(boardElm: HTMLDivElement, e: MouseEvent | DragEvent) {
	const boardSize = boardElm.clientWidth;
	const squareSize = boardSize / 8;
	const dx = e.clientX - boardElm.offsetLeft;
	const dy = e.clientY - boardElm.offsetTop;

	const file = Math.max(0, Math.min(7, Math.floor(dx / squareSize)));
	const rank = Math.max(0, Math.min(7, Math.floor(dy / squareSize)));

	const dstSquare = rank * 8 + file;
	return dstSquare;
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
