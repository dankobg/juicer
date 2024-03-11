import { writable } from 'svelte/store';
import {
	WHITE,
	Board,
	type Color,
	type DragPosition,
	type PieceFenSymbol,
	ALL_PIECES_FEN,
	WHITE_ROOK,
	WHITE_BISHOP,
	WHITE_KNIGHT,
	WHITE_QUEEN,
	WHITE_KING,
	WHITE_PAWN,
	BLACK_ROOK,
	BLACK_BISHOP,
	BLACK_KNIGHT,
	BLACK_QUEEN,
	BLACK_KING,
	BLACK_PAWN,
	NO_SQUARE,
	dragPositionZero,
	type EventData,
	type Row,
	type Col,
} from './model';

type FenState = { state: 'success' } | { state: 'error'; error: string };

type BoardState = {
	fenInfo: FenState;
	board: Board | null;
	orientation: Color;
	boardTheme: string;
	pieceTheme: Map<PieceFenSymbol, string>;
	dragging: boolean;
	draggedElm: HTMLDivElement | null;
	dragOverSquareIndex: number;
	activeSquareIndex: number;
	srcSquareIndex: number;
	destSquareIndex: number;
	pieceOffsetWidth: number;
	pieceOffsetHeight: number;
	dragPos: Record<string, DragPosition>;
	lastPos: Record<string, DragPosition>;
	offBoardBounds: boolean;
	piecesCount: number;
	piecesAnimated: number;
};

const defaultBoardState: BoardState = {
	fenInfo: { state: 'success' },
	board: null,
	orientation: WHITE,
	boardTheme: '/images/board/svg/brown.svg',
	pieceTheme: new Map([
		[WHITE_ROOK, '/images/piece/cburnett/wR.svg'],
		[WHITE_BISHOP, '/images/piece/cburnett/wB.svg'],
		[WHITE_KNIGHT, '/images/piece/cburnett/wN.svg'],
		[WHITE_QUEEN, '/images/piece/cburnett/wQ.svg'],
		[WHITE_KING, '/images/piece/cburnett/wK.svg'],
		[WHITE_PAWN, '/images/piece/cburnett/wP.svg'],
		[BLACK_ROOK, '/images/piece/cburnett/bR.svg'],
		[BLACK_BISHOP, '/images/piece/cburnett/bB.svg'],
		[BLACK_KNIGHT, '/images/piece/cburnett/bN.svg'],
		[BLACK_QUEEN, '/images/piece/cburnett/bQ.svg'],
		[BLACK_KING, '/images/piece/cburnett/bK.svg'],
		[BLACK_PAWN, '/images/piece/cburnett/bP.svg'],
	]),
	dragging: false,
	draggedElm: null,
	dragOverSquareIndex: NO_SQUARE,
	activeSquareIndex: NO_SQUARE,
	srcSquareIndex: NO_SQUARE,
	destSquareIndex: NO_SQUARE,
	pieceOffsetWidth: 0,
	pieceOffsetHeight: 0,
	dragPos: {},
	lastPos: {},
	offBoardBounds: false,
	piecesCount: 0,
	piecesAnimated: 0,
};

function createBoardState() {
	const { subscribe, update } = writable<BoardState>(defaultBoardState);

	return {
		subscribe,
		setupFenState: (error: string | null) =>
			update(s => {
				s.fenInfo = error ? { state: 'error', error: error } : { state: 'success' };
				return s;
			}),
		init: (fen: string) =>
			update(s => {
				s.board = new Board(fen);
				for (const sq of s.board?.squares ?? []) {
					if (sq.piece !== null) {
						s.dragPos[sq.piece.id] = dragPositionZero;
						s.piecesCount++;
					}
				}
				s.lastPos = { ...s.dragPos };
				return s;
			}),
		clear: () =>
			update(s => {
				s.fenInfo = { state: 'success' };
				s.board = null;
				return s;
			}),
		setBoardTheme: (boardTheme: string) =>
			update(s => {
				s.boardTheme = boardTheme;
				return s;
			}),
		setPieceTheme: (resolver: (pieceFenSymbol: PieceFenSymbol) => string) =>
			update(s => {
				const pieceTheme = new Map<PieceFenSymbol, string>();
				for (const pfs of ALL_PIECES_FEN) {
					pieceTheme.set(pfs, resolver(pfs));
				}
				s.pieceTheme = pieceTheme;
				return s;
			}),
		setDragStart: (
			pieceId: string,
			draggedElm: HTMLDivElement,
			draggedElmRect: DOMRect,
			srcSquareIndex: number,
			edata: EventData
		) =>
			update(s => {
				s.dragging = true;
				s.draggedElm = draggedElm;
				s.srcSquareIndex = srcSquareIndex;
				s.activeSquareIndex = srcSquareIndex;
				s.pieceOffsetWidth = draggedElmRect.width / 2;
				s.pieceOffsetHeight = draggedElmRect.height / 2;
				s.lastPos[pieceId] = { ...s.dragPos[pieceId] };
				s.dragPos[pieceId].initialX = edata.clientX - s.dragPos[pieceId].dx + (s.pieceOffsetWidth - edata.offsetX);
				s.dragPos[pieceId].initialY = edata.clientY - s.dragPos[pieceId].dy + (s.pieceOffsetHeight - edata.offsetY);
				s.dragPos[pieceId].dx = edata.clientX - s.dragPos[pieceId].initialX;
				s.dragPos[pieceId].dy = edata.clientY - s.dragPos[pieceId].initialY;
				return s;
			}),
		setDragEnd: () =>
			update(s => {
				s.dragging = false;
				s.draggedElm = null;
				s.activeSquareIndex = NO_SQUARE;
				return s;
			}),
		setDragMove: (pieceId: string, squareIndex: number, edata: EventData) =>
			update(s => {
				s.dragOverSquareIndex = squareIndex;
				s.dragPos[pieceId].dx = edata.clientX - s.dragPos[pieceId].initialX;
				s.dragPos[pieceId].dy = edata.clientY - s.dragPos[pieceId].initialY;
				return s;
			}),
		setDestSquareIndex: (squareIndex: number) =>
			update(s => {
				s.destSquareIndex = squareIndex;
				return s;
			}),
		setSnapbackState: (pieceId: string) =>
			update(s => {
				s.dragPos[pieceId] = { ...s.lastPos[pieceId] };
				s.destSquareIndex = NO_SQUARE;
				return s;
			}),
		setSnapToSquareState: (
			pieceId: string,
			boardElement: HTMLDivElement,
			squareWidth: number,
			squareHeight: number,
			row: Row,
			col: Col,
			clientX: number,
			clientY: number
		) =>
			update(s => {
				const squareTopLeftX = col * squareWidth + boardElement.offsetLeft;
				const squareTopLeftY = row * squareHeight + boardElement.offsetTop;
				const snapDx = clientX - squareTopLeftX - s.pieceOffsetWidth;
				const snapDy = clientY - squareTopLeftY - s.pieceOffsetHeight;

				s.dragPos[pieceId].dx -= snapDx;
				s.dragPos[pieceId].dy -= snapDy;
				s.lastPos[pieceId] = { ...s.dragPos[pieceId] };
				s.board!.squares[s.destSquareIndex].piece = s.board!.squares[s.srcSquareIndex].piece;
				s.board!.squares[s.srcSquareIndex].piece = null;
				return s;
			}),
		incrementPiecesAnimated: () =>
			update(s => {
				s.piecesAnimated++;
				return s;
			}),
		update,
	};
}

export const boardState = createBoardState();
