<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->

<script lang="ts">
  import { getRowAndCol, getSquareColor, type Square } from '$lib/board/board';

  export let square: Square;
  export let bordered: boolean = false;
  export let selected: boolean = false;
  export let highlighted: boolean = false;

  let { row, col } = getRowAndCol(square.squareIdx);

  let squareColor = getSquareColor(square.squareIdx);
</script>

<div
  class="square"
  style="--row:{row}; --col:{col};"
  data-square={square.squareIdx}
  data-highlighted={highlighted}
  data-selected={selected}
  data-bordered={bordered}
  data-color={squareColor}
  on:click
  on:dragenter
  on:dragleave
  on:dragover
  on:drop
/>

<style>
  .square {
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
    top: calc(var(--row) * (var(--board-size) / 8));
    left: calc(var(--col) * (var(--board-size) / 8));
    z-index: 55;
  }

  .square[data-highlighted='true'] {
    background-color: rgba(164, 206, 74, 0.5);
  }

  .square[data-selected='true'] {
    background-color: rgba(173, 90, 194, 0.5);
  }

  .square[data-bordered='true'] {
    border: 3px solid goldenrod;
  }
</style>
