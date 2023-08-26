<script lang="ts">
  import JuicerPiece from './piece.svelte';
  import JuicerSquare from './square.svelte';
  import { type Color, type PieceMovement, Board, FEN_STARTING_POSITION, WHITE } from '$lib/board/board';
  import { onMount } from 'svelte';
  import { boardStore } from './store';

  export let orientation: Color = WHITE;
  export let pieceMovement: PieceMovement = 'click';
  export let fen: string = FEN_STARTING_POSITION;
  export let boardInitialSize = '35rem';

  let draggedElm: HTMLDivElement | null = null;
  let startX: number = 0;
  let startY: number = 0;
  let dx: number = 0;
  let dy: number = 0;

  function onPieceClick(event: MouseEvent) {}

  function onPieceMouseDown(event: MouseEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;
    const parent = elm.parentElement as HTMLDivElement;

    startX = elm.clientWidth / 2;
    startY = elm.clientHeight / 2;

    dx = clientX - parent.offsetLeft - elm.offsetLeft - startX;
    dy = clientY - parent.offsetTop - elm.offsetTop - startY;

    elm.style.translate = `${dx}px ${dy}px`;
  }

  function onPieceDragStart(event: DragEvent) {
    const { target } = event;
    const elm = target as HTMLDivElement;

    const dragImage = new Image();
    dragImage.src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=';
    event.dataTransfer!.setDragImage(dragImage, 0, 0);

    draggedElm = elm;
  }

  function onPieceDragEnd(event: DragEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;
    const parent = elm.parentElement as HTMLDivElement;

    dx = clientX - parent.offsetLeft - elm.offsetLeft - startX;
    dy = clientY - parent.offsetTop - elm.offsetTop - startY;

    elm.style.translate = `${dx}px ${dy}px`;

    draggedElm = null;
  }

  function onPieceDrag(event: DragEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;
    const parent = elm.parentElement as HTMLDivElement;

    dx = clientX - parent.offsetLeft - elm.offsetLeft - startX;
    dy = clientY - parent.offsetTop - elm.offsetTop - startY;

    elm.style.translate = `${dx}px ${dy}px`;
  }

  function onPieceDragOver(event: DragEvent) {
    event.preventDefault();
  }

  function onPieceDragEnter(event: DragEvent) {
    event.preventDefault();
  }

  function onPieceDragLeave(event: DragEvent) {
    event.preventDefault();
  }

  function onPieceDrop(event: DragEvent) {
    event.preventDefault();
    boardStore.selectSquare(-1);
  }

  function onSquareClick(event: MouseEvent) {
    // const { square } = event.detail;
    // boardStore.move(62, 45);
    // if (pieceMovement !== 'click') {
    //   return;
    // }
    // if (selectedSquares.length === 1) {
    //   const from = selectedSquares[0];
    //   const to = event.detail.square;
    //   move(from, to);
    //   selectedSquares = [];
    // }
  }

  function onSquareDragEnter(event: DragEvent) {
    event.preventDefault();
  }

  function onSquareDragLeave(event: DragEvent) {
    event.preventDefault();
  }

  function onSquareDragOver(event: DragEvent) {
    event.preventDefault();
  }

  function onSquareDrop(event: DragEvent) {
    event.preventDefault();
    boardStore.selectSquare(-1);
  }

  onMount(() => {
    const board = new Board(fen);
    boardStore.init(board);
  });
</script>

<div class="board" style="--board-size: {boardInitialSize};" data-orientation={orientation}>
  <div class="rank-notations" data-orientation={orientation}>
    <div>1</div>
    <div>2</div>
    <div>3</div>
    <div>4</div>
    <div>5</div>
    <div>6</div>
    <div>7</div>
    <div>8</div>
  </div>
  <div class="file-notations" data-orientation={orientation}>
    <div>a</div>
    <div>b</div>
    <div>c</div>
    <div>d</div>
    <div>e</div>
    <div>f</div>
    <div>g</div>
    <div>h</div>
  </div>

  {#if $boardStore.board}
    {#each $boardStore.board.squares as sq}
      {#if sq.piece !== null}
        <JuicerSquare
          square={sq}
          selected={$boardStore.selectedSquare === sq.squareIdx}
          highlighted={$boardStore.highlightedSquares.includes(sq.squareIdx)}
        />

        <JuicerPiece
          piece={sq.piece}
          squareIdx={sq.squareIdx}
          on:click={onPieceClick}
          on:mousedown={onPieceMouseDown}
          on:drag={onPieceDrag}
          on:dragstart={onPieceDragStart}
          on:dragend={onPieceDragEnd}
          on:dragenter={onPieceDragEnter}
          on:dragleave={onPieceDragLeave}
          on:dragover={onPieceDragOver}
          on:drop={onPieceDrop}
        />
      {:else}
        <JuicerSquare
          square={sq}
          selected={$boardStore.selectedSquare === sq.squareIdx}
          highlighted={$boardStore.highlightedSquares.includes(sq.squareIdx)}
          on:click={onSquareClick}
          on:dragenter={onSquareDragEnter}
          on:dragleave={onSquareDragLeave}
          on:dragover={onSquareDragOver}
          on:drop={onSquareDrop}
        />
      {/if}
    {/each}
  {/if}
</div>

<style>
  .board {
    width: var(--board-size);
    height: var(--board-size);
    background-image: var(--board-theme);
    background-size: cover;
    position: relative;
    user-select: none;
    z-index: 50;
  }

  .file-notations {
    font-size: 1em;
    font-weight: 500;
    display: flex;
    flex-direction: row;
    position: absolute;
    width: 100%;
    bottom: 0;
    left: 0;
    margin-inline-start: -2px;
    margin-block-end: 2px;
    z-index: 55;
  }
  .file-notations[data-orientation='b'] {
    flex-direction: row-reverse;
  }
  .file-notations div {
    flex: 1;
    display: flex;
    justify-content: end;
  }

  .rank-notations {
    font-size: 1em;
    font-weight: 500;
    display: flex;
    flex-direction: column-reverse;
    position: absolute;
    height: 100%;
    bottom: 0;
    left: 0;
    z-index: 55;
  }
  .rank-notations[data-orientation='b'] {
    flex-direction: column;
  }
  .rank-notations div {
    flex: 1;
    margin-inline-start: 2px;
    margin-block-start: 2px;
  }

  .file-notations[data-orientation='w'] div:nth-child(odd),
  .rank-notations[data-orientation='w'] div:nth-child(odd) {
    color: #f6f6f6;
  }
  .file-notations[data-orientation='w'] div:nth-child(even),
  .rank-notations[data-orientation='w'] div:nth-child(even) {
    color: #333;
  }
  .file-notations[data-orientation='b'] div:nth-child(odd),
  .rank-notations[data-orientation='b'] div:nth-child(odd) {
    color: #333;
  }
  .file-notations[data-orientation='b'] div:nth-child(even),
  .rank-notations[data-orientation='b'] div:nth-child(even) {
    color: #f6f6f6;
  }
</style>
