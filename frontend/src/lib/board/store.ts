import { writable } from 'svelte/store';
import { WHITE, type Board, type Position, findRowAndCol } from '$lib/board/board';

interface BoardStore {
  board: Board | null;
  positions: Position[];
  highlightedSquares: number[];
  selectedSquare: number;
}

function createBoardStore() {
  const { subscribe, set, update } = writable<BoardStore>({
    board: null,
    positions: [],
    highlightedSquares: [],
    selectedSquare: -1,
  });

  return {
    subscribe,
    init: (board: Board) =>
      update(s => {
        s.board = board;
        s.positions = [{ board: board, turn: WHITE }];

        return s;
      }),
    highlightSquare: (squareIdx: number) =>
      update(s => {
        s.highlightedSquares = [...s.highlightedSquares, squareIdx];

        return s;
      }),
    selectSquare: (squareIdx: number) =>
      update(s => {
        s.selectedSquare = squareIdx;

        return s;
      }),
    move: (fromSquareIdx: number, toSquareIdx: number) =>
      update(s => {
        const p = document.querySelector(`.piece[data-square='${fromSquareIdx}']`) as HTMLElement | null;

        console.log(p);

        if (p) {
          const { row, col } = findRowAndCol(toSquareIdx);
          const moveAnimation = [
            { translate: `calc(${col} * (var(--board-size) / 8)) calc(${row} * (var(--board-size) / 8))` },
          ];

          const anim = p.animate(moveAnimation, {
            duration: 80,
            fill: 'forwards',
          });

          anim.onfinish = () => {
            if (s.board) {
              s.board.squares[toSquareIdx] = s.board.squares[fromSquareIdx];
              s.board.squares[fromSquareIdx].piece = null;
              s.positions = [...s.positions, { board: s.board, turn: 'w' }];
            }

            return s;
          };
        }
        return s;
      }),
  };
}

export const boardStore = createBoardStore();
