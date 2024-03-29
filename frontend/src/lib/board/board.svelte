<script lang="ts">
  import JuicerPiece from './piece.svelte';
  import JuicerSquare from './square.svelte';
  import {
    type Color,
    Board,
    FEN_STARTING_POSITION,
    WHITE,
    makeTransparentDragImage,
    whiteSparePieces,
    blackSparePieces,
    translateElm,
    getRowAndCol,
    getSquareIdxFromDragPos,
    BLACK,
  } from '$lib/board/board';
  import { onMount } from 'svelte';
  import { boardStore } from './store';

  export let fen: string = FEN_STARTING_POSITION;
  export let boardInitialSize = '30rem';
  export let orientation: Color = WHITE;
  // export let pieceMovement: PieceMovement = 'drag';
  export let showSparePieces: boolean = false;

  // let board = new Board(fen);
  let board = new Board('r7/2k5/6b1/8/2P5/8/5P2/3Q2N1 w - - 0 1');
  let transparentDragImage: HTMLImageElement | null = null;
  let squareElements: NodeListOf<HTMLDivElement>;
  let pieceElements: NodeListOf<HTMLDivElement>;
  let boardElm: HTMLDivElement;
  let boardWidth: number;
  let boardHeight: number;

  let draggedElm: HTMLDivElement | null = null;
  let draggedElmId: string = '';
  let dragging: boolean = false;
  let dragOverSquare: number = -1;
  let activeSquare: number = -1;
  let dstSquare: number = -1;
  let srcSquare: number = -1;

  let xOffset: number = 0;
  let yOffset: number = 0;
  let pos: Record<string, { initialX: number; initialY: number; dx: number; dy: number }> = {};

  let firstAnimationPlayed: boolean = false;

  function onBoardDragEnter(event: DragEvent) {
    event.preventDefault();
  }

  function onBoardDragLeave(event: DragEvent) {}

  function onBoardDragOver(event: DragEvent) {
    event.preventDefault();
    dragOverSquare = getSquareIdxFromDragPos(boardElm, event);
  }

  function onBoardDrop(event: DragEvent) {
    event.preventDefault();

    const { target, clientX, clientY } = event;
    const pieceElm = target as HTMLDivElement;

    const sqIdx = Number.parseInt(pieceElm.dataset.square ?? '-1');
    const pieceId = pieceElm.dataset.id ?? '';
    const symbol = pieceElm.dataset.symbol ?? '';

    pos[pieceId].initialX = clientX - xOffset;
    pos[pieceId].initialY = clientY - yOffset;
    draggedElm = null;
    draggedElmId = '';
    dragging = false;
    dragOverSquare = -1;
    activeSquare = -1;
    dstSquare = getSquareIdxFromDragPos(boardElm, event);
    pieceElm.dataset.square = dstSquare.toString();
  }

  function onPieceClick(event: MouseEvent) {}

  function onPiecePointerDown(event: MouseEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;
    const rect = elm.getBoundingClientRect();

    const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
    const pieceId = elm.dataset.id ?? '';
    const symbol = elm.dataset.symbol ?? '';

    xOffset = elm.clientWidth / 2;
    yOffset = elm.clientHeight / 2;

    pos[pieceId].dx = clientX + (pos?.[pieceId]?.dx ?? 0) - rect.left - xOffset;
    pos[pieceId].dy = clientY + (pos?.[pieceId]?.dy ?? 0) - rect.top - yOffset;

    translateElm(elm, pos[pieceId].dx, pos[pieceId].dy);
  }

  function onPieceDragStart(event: DragEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;

    const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
    const pieceId = elm.dataset.id ?? '';
    const symbol = elm.dataset.symbol ?? '';

    pos[pieceId].initialX = clientX - (pos?.[pieceId]?.dx ?? 0);
    pos[pieceId].initialY = clientY - (pos?.[pieceId]?.dy ?? 0);

    dragging = true;
    draggedElm = elm;
    draggedElmId = pieceId;

    if (transparentDragImage) {
      event.dataTransfer!.setDragImage(transparentDragImage, 0, 0);
    }

    if (elm.dataset.square) {
      srcSquare = Number.parseInt(elm.dataset.square);
    }
  }

  function onPieceDragEnd(event: DragEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;

    const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
    const pieceId = elm.dataset.id ?? '';
    const symbol = elm.dataset.symbol ?? '';

    pos[pieceId].initialX = clientX;
    pos[pieceId].initialY = clientY;
    dragging = false;
    draggedElm = null;
    draggedElmId = '';
  }

  function onPieceDrag(event: DragEvent) {
    const { target, clientX, clientY } = event;
    const elm = target as HTMLDivElement;

    const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
    const pieceId = elm.dataset.id ?? '';
    const symbol = elm.dataset.symbol ?? '';

    pos[pieceId].dx = clientX - (pos?.[pieceId]?.initialX ?? 0);
    pos[pieceId].dy = clientY - (pos?.[pieceId]?.initialY ?? 0);

    translateElm(elm, pos[pieceId].dx, pos[pieceId].dy);
  }

  function setupDomElementsAndInitialVars(): void {
    transparentDragImage = makeTransparentDragImage();
    pieceElements = document.querySelectorAll('.piece');
    squareElements = document.querySelectorAll('.square');
  }

  function setInitialPiecesPositions(): void {
    if (board) {
      for (const sq of board.squares) {
        if (sq.piece !== null) {
          pos[sq.piece.id] = { initialX: 0, initialY: 0, dx: 0, dy: 0 };
        }
      }
    }
  }

  function animatePiecesMoves(): void {
    if (board) {
      for (let i = 0; i < pieceElements.length; i++) {
        const p = pieceElements[i];
        const psize = p.clientWidth;
        const sqIdx = Number.parseInt(p.dataset.square ?? '-1');
        const pieceId = p.dataset.id ?? '';
        const symbol = p.dataset.symbol ?? '';
        const color = p.dataset.color ?? '';
        const { row, col } = getRowAndCol(sqIdx);
        let [dx, dy] = [col * psize, row * psize];
        if (orientation === BLACK) {
          dx = (7 - col) * 60;
          dy = (7 - row) * 60;
        }
        const blackKeyframesFrom = { opacity: 0, transform: `translate(${boardWidth / 2}px, ${-psize}px` };
        const whiteKeyframesFrom = { opacity: 0, transform: `translate(${boardWidth / 2}px, ${boardHeight}px)` };
        let keyframes: Keyframe[] = [{ opacity: 1, transform: `translate(${dx}px, ${dy}px)` }];

        keyframes.unshift(color === 'w' ? whiteKeyframesFrom : blackKeyframesFrom);

        const opts: KeyframeAnimationOptions = {
          fill: 'forwards',
          easing: 'ease-in',
          duration: !firstAnimationPlayed ? 200 : 0,
        };

        const anim = p.animate(keyframes, opts);
        anim.onfinish = () => {
          pos[pieceId].initialX = dx;
          pos[pieceId].initialY = dy;
        };
      }
    }

    firstAnimationPlayed = true;
  }

  onMount(() => {
    setupDomElementsAndInitialVars();
    setInitialPiecesPositions();
    // animatePiecesMoves();
  });

  $: pieceElements && orientation && animatePiecesMoves();
</script>

<div style="--board-size: {boardInitialSize};">
  {#if showSparePieces}
    <div class="spare-pieces black">
      {#each blackSparePieces as bSpare}
        <JuicerPiece
          square={-1}
          piece={bSpare}
          spare
          on:click={onPieceClick}
          on:pointerdown={onPiecePointerDown}
          on:drag={onPieceDrag}
          on:dragstart={onPieceDragStart}
          on:dragend={onPieceDragEnd}
        />
      {/each}
    </div>
  {/if}

  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    bind:this={boardElm}
    bind:clientWidth={boardWidth}
    bind:clientHeight={boardHeight}
    class="board"
    style="--board-size: {boardInitialSize};"
    data-orientation={orientation}
    on:dragenter={onBoardDragEnter}
    on:dragleave={onBoardDragLeave}
    on:dragover={onBoardDragOver}
    on:drop={onBoardDrop}
  >
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

    <!-- {#if $boardStore.board}
      {#each $boardStore.board.squares as sq} -->
    {#if board}
      {#each board.squares as sq}
        <JuicerSquare
          square={sq}
          selected={$boardStore.activeSquare === sq.squareIdx}
          highlighted={$boardStore.highlightedSquares.includes(sq.squareIdx)}
        />

        {#if sq.piece !== null}
          <JuicerPiece
            square={sq.squareIdx}
            piece={sq.piece}
            on:click={onPieceClick}
            on:pointerdown={onPiecePointerDown}
            on:drag={onPieceDrag}
            on:dragstart={onPieceDragStart}
            on:dragend={onPieceDragEnd}
          />
        {/if}
      {/each}
    {/if}
  </div>

  {#if showSparePieces}
    <div class="spare-pieces white">
      {#each whiteSparePieces as wSpare}
        <JuicerPiece
          square={-1}
          piece={wSpare}
          spare
          on:click={onPieceClick}
          on:pointerdown={onPiecePointerDown}
          on:drag={onPieceDrag}
          on:dragstart={onPieceDragStart}
          on:dragend={onPieceDragEnd}
        />
      {/each}
    </div>
  {/if}
</div>

<button
  style="margin-block: 3rem;"
  on:click={() => {
    board.orientation = board.orientation === 'w' ? 'b' : 'w';
    board = board;
    orientation = orientation === 'w' ? 'b' : 'w';
    // board.squares = [...board.squares].reverse();
  }}>rofl</button
>

<pre>{JSON.stringify(board.orientation, null, 2)}</pre>
<pre style="font-family: Courier;">{board.print()}</pre>

<style>
  .spare-pieces {
    width: var(--board-size);
    position: relative;
    display: flex;
    justify-content: space-evenly;
    gap: 1rem;
  }

  .board {
    width: var(--board-size);
    height: var(--board-size);
    background-image: var(--board-theme);
    background-size: cover;
    position: relative;
    user-select: none;
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
    z-index: 51;
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
    z-index: 51;
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
