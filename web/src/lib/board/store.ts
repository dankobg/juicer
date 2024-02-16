import { writable } from 'svelte/store';
import { type Board } from '$lib/board/board';

interface BoardStore {
	board: Board | null;
	// positions: Position[];
	highlightedSquares: number[];
	activeSquare: number;
}

function createBoardStore() {
	const { subscribe, update } = writable<BoardStore>({
		board: null,
		// positions: [],
		highlightedSquares: [],
		activeSquare: -1,
	});

	return {
		subscribe,
		init: (board: Board) =>
			update(s => {
				s.board = board;
				// s.positions = [{ board: board, turn: WHITE }];

				return s;
			}),
		highlightSquare: (squareIndices: number[]) =>
			update(s => {
				s.highlightedSquares = [...s.highlightedSquares, ...squareIndices];

				return s;
			}),
		setActiveSquare: (squareIdx: number) =>
			update(s => {
				s.activeSquare = squareIdx;

				return s;
			}),
		move: (fromSquareIdx: number, toSquareIdx: number) =>
			update(s => {
				if (s.board) {
					s.board.squares[toSquareIdx] = s.board.squares[fromSquareIdx];
					s.board.squares[fromSquareIdx].piece = null;
					// s.positions = [...s.positions, { board: s.board, turn: 'w' }];
				}

				return s;
			}),
	};
}

export const boardStore = createBoardStore();

/**
 * init(board): Board
 * move([from, to, cfg]): void
 * selectSquare(e2): void
 * highlightSquares([a1, c4, h8]): void
 * setSquarePiece([e5, 'Q']): void - or maybe addPiece and removePiece
 * toggleOrientation(): void
 * setOrientation(color: Color): void
 * getFen(): void
 * setFen(fen: string): void
 * clearBoard(): void
 */
