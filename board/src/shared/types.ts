export type ColorWhite = 'w';
export type ColorBlack = 'b';
export type Color = ColorWhite | ColorBlack;

export const WHITE: ColorWhite = 'w';
export const BLACK: ColorBlack = 'b';

export type RowCol = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
export type Rank = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8;
export type File = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';
export type Coordinate = `${File}${Rank}`;
export type Ranks = Rank[];
export type Files = File[];
export type BoardSize = 8;
export type PieceSetStyle =
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

export type PieceFENSymbol = 'p' | 'n' | 'b' | 'r' | 'q' | 'k' | 'P' | 'N' | 'B' | 'R' | 'Q' | 'K';
export type PieceSymbol = Lowercase<PieceFENSymbol>;
export type PromotionPiece = Exclude<PieceSymbol, 'p' | 'k'>;

type InfoFromFEN = {
  board: Square[][];
  activeColor: ColorBlack | ColorWhite;
  castlingRights: {
    whiteKingSide: boolean;
    whiteQueenSide: boolean;
    blackKingSide: boolean;
    blackQueenSide: boolean;
  };
  enPassantTargetSquare: Square | null;
  halfMoveClock: number;
  fullMoveClock: number;
};

export const BOARD_SIZE: BoardSize = 8;
export const FILES: Files = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
export const RANKS: Ranks = [1, 2, 3, 4, 5, 6, 7, 8];
export const STARTING_POSITION_FEN = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';
export const EMPTY_POSITION_FEN = '8/8/8/8/8/8/8/8';

export const PAWN: PieceSymbol = 'p';
export const KNIGHT: PieceSymbol = 'n';
export const BISHOP: PieceSymbol = 'b';
export const ROOK: PieceSymbol = 'r';
export const QUEEN: PieceSymbol = 'q';
export const KING: PieceSymbol = 'k';

const FILE_TO_COLUMN = new Map<File, RowCol>([
  ['a', 0],
  ['b', 1],
  ['c', 2],
  ['d', 3],
  ['e', 4],
  ['f', 5],
  ['g', 6],
  ['h', 7],
]);

const COLUMN_TO_FILE = new Map<RowCol, File>([
  [0, 'a'],
  [1, 'b'],
  [2, 'c'],
  [3, 'd'],
  [4, 'e'],
  [5, 'f'],
  [6, 'g'],
  [7, 'h'],
]);

const RANK_TO_ROW = new Map<Rank, RowCol>([
  [1, 7],
  [2, 6],
  [3, 5],
  [4, 4],
  [5, 3],
  [6, 2],
  [7, 1],
  [8, 0],
]);

const ROW_TO_RANK = new Map<RowCol, Rank>([
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

export function convertFileToColumn(file: File): RowCol {
  return FILE_TO_COLUMN.get(file) as RowCol;
}

export function convertColumnToFile(column: RowCol): File {
  return COLUMN_TO_FILE.get(column) as File;
}

export function convertRankToRow(rank: Rank): RowCol {
  return RANK_TO_ROW.get(rank) as RowCol;
}

export function convertRowToRank(row: RowCol): Rank {
  return ROW_TO_RANK.get(row) as Rank;
}

export function convertCoordinateToFileAndRank(coordinate: Coordinate): [File, Rank] {
  const chars = coordinate.split('');

  const file = chars[0] as File;
  const rank = Number.parseInt(chars[1]) as Rank;

  return [file, rank];
}

export function convertCoordinateToRowAndColumn(coordinate: Coordinate): [RowCol, RowCol] {
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

export function convertRowAndColumnToCoordinate(row: RowCol, col: RowCol): Coordinate {
  const file = convertColumnToFile(col);
  const rank = convertRowToRank(row);

  return convertFileAndRankToCoordinate(file, rank);
}

function isDigit(c: string): boolean {
  return /^[0-9]$/.test(c);
}

export function validateFEN(fen: string): { ok: boolean; err: string } {
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

export function convertBoardToFen(board: Square[][]): string {
  let fen = '';

  for (let i = 0; i < board.length; i++) {
    let empty = 0;

    for (let j = 0; j < board[i].length; j++) {
      if (board[i][j].isEmpty()) {
        empty += 1;
      }

      if (board[i][j].hasPiece()) {
        if (empty > 0) {
          fen += empty.toString();
          empty = 0;
        }

        fen += board[i][j].piece?.toFENPieceSymbol();
      }
    }

    if (empty > 0) {
      fen += empty.toString();
    }

    if (i < board.length - 1) {
      fen += '/';
    }
  }

  fen += ' w KQkq - 0 1';

  return fen;
}

function convertFenPiecePlacementTokenToBoard(piecePlacmentToken: string): Square[][] {
  const board: Square[][] = [];

  piecePlacmentToken.split('/').forEach((row, i) => {
    const squares: Square[] = [];

    row.split('').forEach((char, j) => {
      if (isDigit(char)) {
        for (let x = 0; x < Number.parseInt(char); x++) {
          const square = new Square(i as RowCol, (j + x) as RowCol, null);
          squares.push(square);
          board[i] = squares;
        }
      } else {
        const piece = Piece.fromFENSymbol(char as PieceFENSymbol);
        const square = new Square(i as RowCol, j as RowCol, piece);
        squares.push(square);
        board[i] = squares;
      }
    });
  });

  return board;
}

export function parseInfoFromFEN(fen: string): InfoFromFEN | never {
  const res = validateFEN(fen);
  if (!res.ok) {
    throw new Error(res.err);
  }

  const tokens = fen.split(' ');
  const [piecePlacement, activeColor, castlingRights, enPassantTargetSquare, halfMoveClock, fullMoveClock] = tokens;

  const calcEnPassantTargetSquare = (): Square | null => {
    if (enPassantTargetSquare === '-') {
      return null;
    }

    const [row, col] = convertCoordinateToRowAndColumn(enPassantTargetSquare as Coordinate);
    return new Square(row, col, null);
  };

  const info: InfoFromFEN = {
    board: convertFenPiecePlacementTokenToBoard(piecePlacement),
    activeColor: activeColor as Color,
    castlingRights: {
      whiteKingSide: castlingRights.includes('K'),
      whiteQueenSide: castlingRights.includes('Q'),
      blackKingSide: castlingRights.includes('k'),
      blackQueenSide: castlingRights.includes('q'),
    },
    enPassantTargetSquare: calcEnPassantTargetSquare(),
    halfMoveClock: Number.parseInt(halfMoveClock) ?? 0,
    fullMoveClock: Number.parseInt(fullMoveClock) ?? 1,
  };

  return info;
}

export class Piece {
  alive: boolean = true;

  constructor(public symbol: PieceSymbol, public color: Color) {}

  static fromFENSymbol(symbol: PieceFENSymbol): Piece | null {
    const color = () => {
      if (/^[prnbqk]$/.test(symbol)) {
        return 'b';
      }
      if (/^[PRNBQK]$/.test(symbol)) {
        return 'w';
      }
      throw new Error('invalid color');
    };

    return new Piece(symbol.toLowerCase() as PieceSymbol, color());
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

  showLegalMoves(): Square[] | null {
    return null;
  }

  symbolEquals(piece: Piece): boolean {
    return this.symbol === piece.symbol;
  }

  toFENPieceSymbol(): string {
    if (this.color === WHITE) {
      return this.symbol.toUpperCase();
    }

    return this.symbol;
  }
}

export class Square {
  color: Color;

  constructor(public row: RowCol, public column: RowCol, public piece: Piece | null) {
    this.color = this.calculateSquareColor(row, column);
  }

  static indexInRange(idx: number): boolean {
    return idx >= 0 && idx <= 7;
  }

  static indicesInRange(...indices: number[]): boolean {
    return indices.every(idx => this.indexInRange(idx));
  }

  get file(): File {
    return convertColumnToFile(this.column);
  }

  get rank(): Rank {
    return convertRowToRank(this.row);
  }

  get coordinate(): Coordinate {
    return convertRowAndColumnToCoordinate(this.row, this.column);
  }

  hasPiece(): boolean {
    return this.piece !== null;
  }

  isWhite(): boolean {
    return this.color === WHITE;
  }

  isBlack(): boolean {
    return this.color === BLACK;
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

  toFENSymbol(): string | undefined {
    return this.piece?.toFENPieceSymbol();
  }

  private calculateSquareColor(x: number, y: number): Color {
    return (x + y) % 2 === 0 ? WHITE : BLACK;
  }

  equals(square: Square): boolean {
    return this.coordinate === square.coordinate;
  }
}

export class Board {
  private boardSize: BoardSize = BOARD_SIZE;
  squares: Square[][] = [];
  pieceSet: PieceSetStyle = 'cburnett';

  constructor(fen = STARTING_POSITION_FEN) {
    this.initializeBoard();
    this.initializePieceSet();
    this.loadFromFEN(fen);
  }

  fen(): string {
    return convertBoardToFen(this.squares);
  }

  ascii(): string {
    let s = '   +------------------------+\n';

    for (let i = 0; i < this.squares.length; i++) {
      for (let j = 0; j < this.squares[i].length; j++) {
        if (j % 8 === 0) {
          s += ' ' + convertRowToRank(i as RowCol) + ' |';
        }

        if (this.squares[i][j].piece) {
          s += ' ' + this.squares[i][j].piece?.toFENPieceSymbol() + ' ';
        } else {
          s += ' - ';
        }

        if ((j + 1) % 8 === 0) {
          s += '| \n';
        }
      }
    }

    s += '   +------------------------+\n';
    s += '     a  b  c  d  e  f  g  h';

    return s;
  }

  loadFromFEN(fen = STARTING_POSITION_FEN): void {
    const info = parseInfoFromFEN(fen);

    for (let i = 0; i < info.board.length; i++) {
      for (let j = 0; j < info.board[i].length; j++) {
        this.squares[i][j].piece = info.board[i][j].piece;
      }
    }
  }

  initializeBoard(): void {
    for (let i = 0; i < this.boardSize; i++) {
      const squares: Square[] = [];

      for (let j = 0; j < this.boardSize; j++) {
        const square = new Square(i as RowCol, j as RowCol, null);
        squares.push(square);
        this.squares[i] = squares;
      }
    }
  }

  initializePieceSet(pieceSetName?: string): void {
    const name = pieceSetName ?? this.pieceSet;

    const rootElm = document.querySelector(':root') as HTMLElement;

    const sets = new Map<string, string>([
      ['--br-piece-set', `url('/piece/${name}/bR.svg')`],
      ['--bb-piece-set', `url('/piece/${name}/bB.svg')`],
      ['--bn-piece-set', `url('/piece/${name}/bN.svg')`],
      ['--bq-piece-set', `url('/piece/${name}/bQ.svg')`],
      ['--bk-piece-set', `url('/piece/${name}/bK.svg')`],
      ['--bp-piece-set', `url('/piece/${name}/bP.svg')`],
      ['--wr-piece-set', `url('/piece/${name}/wR.svg')`],
      ['--wb-piece-set', `url('/piece/${name}/wB.svg')`],
      ['--wn-piece-set', `url('/piece/${name}/wN.svg')`],
      ['--wq-piece-set', `url('/piece/${name}/wQ.svg')`],
      ['--wk-piece-set', `url('/piece/${name}/wK.svg')`],
      ['--wp-piece-set', `url('/piece/${name}/wP.svg')`],
    ]);

    if (rootElm) {
      for (const [k, v] of sets.entries()) {
        rootElm.style.setProperty(k, v);
      }
    }
  }

  clone() {
    return new Board(this.fen());
  }
}
