<script lang="ts">
	import type { Coordinate, DragPosition } from './types';
	import { onMount, tick } from 'svelte';
	import JuicerSquare from './square.svelte';
	import JuicerPiece from './piece.svelte';
	import { browser } from '$app/environment';
	import { BLACK, FEN_STARTING_POSITION } from './common';
	import { Board } from './board';
	import { Square } from './square';

	export let fen: string = FEN_STARTING_POSITION;

	let board: Board | null = null;
	let boardElm: HTMLDivElement;
	let boardWidth: number;
	let boardHeight: number;
	let boardSize: number = 480;

	let transparentDragImage: HTMLImageElement | null = null;
	let squareElements: NodeListOf<HTMLDivElement>;
	let pieceElements: NodeListOf<HTMLDivElement>;

	let dragging: boolean = false;
	let draggedElm: HTMLDivElement | null = null;
	let dragOverSquare: number = -1;
	let activeSquare: number = -1;
	let dstSquare: number = -1;
	let srcSquare: number = -1;

	let offsetWidth: number = 0;
	let offsetHeight: number = 0;

	let dragPos: Record<string, DragPosition> = {};
	let initialDragPos: DragPosition = { initialX: 0, initialY: 0, dx: 0, dy: 0 };

	function onPiecePointerDown(event: PointerEvent) {
		if (event.button !== 0 || event.ctrlKey) {
			return;
		}

		const { target, clientX, clientY, offsetX, offsetY } = event;
		const elm = target as HTMLDivElement;
		const rect = elm.getBoundingClientRect();

		elm.setPointerCapture(event.pointerId);
		elm.style.userSelect = 'none';

		const pieceId = elm.dataset.id ?? '';

		const sqIdx = getSquareIdxFromPointer(event);
		srcSquare = sqIdx;
		activeSquare = sqIdx;

		dragging = true;
		draggedElm = elm;

		offsetWidth = rect.width / 2;
		offsetHeight = rect.height / 2;

		if (!dragPos[pieceId]) {
			dragPos[pieceId] = initialDragPos;
		}

		const newPos = dragPos[pieceId];

		newPos.initialX = clientX - newPos.dx + (offsetWidth - offsetX);
		newPos.initialY = clientY - newPos.dy + (offsetHeight - offsetY);
		newPos.dx = clientX - newPos.initialX;
		newPos.dy = clientY - newPos.initialY;

		dragPos[pieceId] = { ...newPos };

		translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);
	}

	function onPiecePointerUp(event: PointerEvent) {
		const { target } = event;
		const elm = target as HTMLDivElement;

		elm.style.userSelect = 'auto';

		dragging = false;
		draggedElm = null;
		activeSquare = -1;

		// TODO
		const elements = document.elementsFromPoint(event.clientX, event.clientY);
		const dropSquare = elements.find(e => e.className.startsWith('square')) as HTMLDivElement | undefined;
		if (dropSquare) {
			const sq = dropSquare.dataset.sq ?? '';
			const { squareIdx } = Square.fromCoord(sq as Coordinate);
			dstSquare = squareIdx;
		} else {
			dstSquare = -1;
		}
	}

	function onPiecePointerCancel(event: PointerEvent) {
		const { target } = event;
		const elm = target as HTMLDivElement;

		elm.style.userSelect = 'auto';

		dragging = false;
		draggedElm = null;
		activeSquare = -1;
	}

	function onPiecePointerMove(event: PointerEvent) {
		if (!dragging) {
			return;
		}

		event.stopPropagation();

		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;

		const pieceId = elm.dataset.id ?? '';

		const sqIdx = getSquareIdxFromPointer(event);
		dragOverSquare = sqIdx;

		const newPos = dragPos[pieceId];

		newPos.dx = clientX - newPos.initialX;
		newPos.dy = clientY - newPos.initialY;
		dragPos[pieceId] = { ...newPos };

		translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);
	}

	function onPieceTouchStart(event: TouchEvent) {
		event.preventDefault();
	}

	function onPieceDragStart(event: DragEvent) {
		event.preventDefault();
	}

	function initPieceTheme(pieceTheme: string | undefined) {
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

	function initBoardTheme(boardTheme: string | undefined) {
		if (boardTheme) {
			if (browser) {
				const rootElm = document.querySelector(':root') as HTMLElement;

				if (rootElm) {
					rootElm.style.setProperty('--board-theme', `url('${boardTheme}')`);
				}
			}
		}
	}

	function makeTransparentDragImage(): HTMLImageElement {
		const img = new Image();
		img.src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=';
		return img;
	}

	function translateElm(elm: HTMLDivElement, dx: number, dy: number): void {
		if (elm) {
			elm.style.translate = `${dx}px ${dy}px`;
		}
	}

	function getSquareIdxFromPointer(event: PointerEvent) {
		const squareSize = boardSize / 8;
		const dx = event.clientX - boardElm.offsetLeft;
		const dy = event.clientY - boardElm.offsetTop;

		const file = Math.max(0, Math.min(7, Math.floor(dx / squareSize)));
		const rank = Math.max(0, Math.min(7, Math.floor(dy / squareSize)));

		const sqIdx = rank * 8 + file;
		return sqIdx;
	}

	function setupDomElementsAndInitialVars(): void {
		transparentDragImage = makeTransparentDragImage();
		pieceElements = document.querySelectorAll('.piece');
		squareElements = document.querySelectorAll('.square');
	}

	function setInitialPieceAnimationState(squares: Square[]): void {
		for (const sq of squares) {
			if (sq.piece !== null) {
				dragPos[sq.piece.id] = { initialX: 0, initialY: 0, dx: 0, dy: 0 };
			}
		}
	}

	function animatePiecesMoves(): void {
		if (board) {
			for (let i = 0; i < pieceElements.length; i++) {
				const p = pieceElements[i];
				const psize = p.clientWidth;
				const sq = p.dataset.sq ?? '';
				const pieceId = p.dataset.id ?? '';
				const color = p.dataset.color ?? '';

				const { row, col } = Square.fromCoord(sq as Coordinate);

				let [dx, dy] = [col * psize, row * psize];
				if (board.orientation === BLACK) {
					dx = (7 - col) * (boardSize / 8);
					dy = (7 - row) * (boardSize / 8);
				}
				const blackKeyframesFrom = { opacity: 0, transform: `translate(${boardWidth / 2}px, ${-psize}px` };
				const whiteKeyframesFrom = { opacity: 0, transform: `translate(${boardWidth / 2}px, ${boardHeight}px)` };
				let keyframes: Keyframe[] = [{ opacity: 1, transform: `translate(${dx}px, ${dy}px)` }];

				keyframes.unshift(color === 'w' ? whiteKeyframesFrom : blackKeyframesFrom);

				const opts: KeyframeAnimationOptions = {
					fill: 'forwards',
					easing: 'ease-in',
					duration: 200,
				};

				const anim = p.animate(keyframes, opts);
				anim.onfinish = () => {
					dragPos[pieceId] = { ...initialDragPos, initialX: dx, initialY: dy };
				};
			}
		}
	}

	onMount(() => {
		initBoardTheme('/images/board/svg/brown.svg');
		initPieceTheme('cburnett');

		board = new Board(fen);

		tick().then(() => {
			if (board) {
				setupDomElementsAndInitialVars();
				// setInitialPieceAnimationState(board.squares);
				animatePiecesMoves();
			}
		});
	});
</script>

<pre>{JSON.stringify({ dragging, dragOverSquare, activeSquare, dstSquare, srcSquare }, null, 2)}</pre>

{#if board}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		bind:this={boardElm}
		bind:clientWidth={boardWidth}
		bind:clientHeight={boardHeight}
		class="board"
		style="--board-size: {boardSize}px;"
		data-orientation={board.orientation}
	>
		<div class="rank-notations" data-orientation={board.orientation}>
			<div>1</div>
			<div>2</div>
			<div>3</div>
			<div>4</div>
			<div>5</div>
			<div>6</div>
			<div>7</div>
			<div>8</div>
		</div>
		<div class="file-notations" data-orientation={board.orientation}>
			<div>a</div>
			<div>b</div>
			<div>c</div>
			<div>d</div>
			<div>e</div>
			<div>f</div>
			<div>g</div>
			<div>h</div>
		</div>

		{#each board.squares as sq}
			<JuicerSquare
				square={sq}
				selected={sq.squareIdx === activeSquare}
				highlighted={activeSquare === -1 &&
					srcSquare !== -1 &&
					dstSquare !== -1 &&
					(sq.squareIdx === srcSquare || sq.squareIdx === dstSquare)}
				bordered={activeSquare !== -1 && sq.squareIdx === dragOverSquare && sq.squareIdx !== srcSquare}
			/>

			{#if sq.piece !== null}
				<JuicerPiece
					square={sq.coord}
					piece={sq.piece}
					on:pointerdown={onPiecePointerDown}
					on:pointerup={onPiecePointerUp}
					on:pointermove={onPiecePointerMove}
					on:pointercancel={onPiecePointerCancel}
					on:touchstart{onPieceTouchStart}
					on:dragstart={onPieceDragStart}
				/>
			{/if}
		{/each}
	</div>
{/if}

<style>
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
