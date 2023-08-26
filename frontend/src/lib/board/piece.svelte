<script lang="ts">
  import { findRowAndCol, type Piece } from './board';

  export let piece: Piece;
  export let squareIdx: number;

  let { row, col } = findRowAndCol(squareIdx);
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  class="piece"
  style="--row:{row}; --col:{col};"
  data-piece-fen={piece.toPieceFenSymbol()}
  data-square={squareIdx}
  data-row={row}
  data-col={col}
  on:click
  on:mousedown
  on:drag
  on:dragstart
  on:dragend
  on:dragover
  on:dragenter
  on:dragleave
  on:drop
  draggable="true"
/>

<style>
  .piece {
    position: absolute;
    width: calc(var(--board-size) / 8);
    height: calc(var(--board-size) / 8);
    display: flex;
    justify-content: center;
    align-items: center;
    background-size: contain;
    background-position: center;
    background-repeat: no-repeat;
    overflow: hidden;
    z-index: 60;
    top: calc(var(--row) * (var(--board-size) / 8));
    left: calc(var(--col) * (var(--board-size) / 8));
  }

  .piece[data-piece-fen='r'] {
    background-image: var(--br-theme);
  }
  .piece[data-piece-fen='b'] {
    background-image: var(--bb-theme);
  }
  .piece[data-piece-fen='n'] {
    background-image: var(--bn-theme);
  }
  .piece[data-piece-fen='q'] {
    background-image: var(--bq-theme);
  }
  .piece[data-piece-fen='k'] {
    background-image: var(--bk-theme);
  }
  .piece[data-piece-fen='p'] {
    background-image: var(--bp-theme);
  }
  .piece[data-piece-fen='R'] {
    background-image: var(--wr-theme);
  }
  .piece[data-piece-fen='B'] {
    background-image: var(--wb-theme);
  }
  .piece[data-piece-fen='N'] {
    background-image: var(--wn-theme);
  }
  .piece[data-piece-fen='Q'] {
    background-image: var(--wq-theme);
  }
  .piece[data-piece-fen='K'] {
    background-image: var(--wk-theme);
  }
  .piece[data-piece-fen='P'] {
    background-image: var(--wp-theme);
  }
</style>
