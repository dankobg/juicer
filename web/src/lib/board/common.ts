import type { BlackPieceSymbol, Color, PieceSymbol, WhitePieceSymbol } from './types';

export const WHITE: Color = 'w';
export const BLACK: Color = 'b';

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

export const FEN_EMPTY_POSITION = '8/8/8/8/8/8/8/8';
export const FEN_STARTING_POSITION = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';

export const whitePiecesFen = ['R', 'N', 'B', 'Q', 'K', 'P'];
export const blackPiecesFen = ['r', 'n', 'b', 'q', 'k', 'p'];
export const allPiecesFen = whitePiecesFen.concat(blackPiecesFen);

export const BOARD_SIZE = 8;
export const BOARD_TOTAL_SQUARES = 64;

export const FILE_CHARS = 'abcdefgh';
export const RANK_CHARS = '12345678';
