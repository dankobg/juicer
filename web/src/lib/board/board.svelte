<script lang="ts">
	import { onMount, tick } from 'svelte';
	import JuicerSquare from './square.svelte';
	import JuicerPiece from './piece.svelte';
	import { browser } from '$app/environment';
	import { BLACK, FEN_STARTING_POSITION } from './common';
	import { Board } from './board';
	import { Square } from './square';

	export let fen: string = FEN_STARTING_POSITION;

	let board: Board | null = null;
	let boardSize = 480;

	let boardElm: HTMLDivElement;
	let boardWidth: number;
	let boardHeight: number;

	let transparentDragImage: HTMLImageElement | null = null;
	let squareElements: NodeListOf<HTMLDivElement>;
	let pieceElements: NodeListOf<HTMLDivElement>;

	let draggedElm: HTMLDivElement | null = null;
	let draggedElmId: string = '';
	let pieceState: 'dragging' | 'idle' | 'waitingSecond' = 'idle';
	let dragOverSquare: number = -1;
	let activeSquare: number = -1;
	let dstSquare: number = -1;
	let srcSquare: number = -1;

	let xOffset: number = 0;
	let yOffset: number = 0;
	let piecesAnimationState: Record<string, { initialX: number; initialY: number; dx: number; dy: number }> = {};

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

		piecesAnimationState[pieceId].initialX = clientX - xOffset;
		piecesAnimationState[pieceId].initialY = clientY - yOffset;
		draggedElm = null;
		draggedElmId = '';
		pieceState = 'idle';
		dragOverSquare = -1;
		activeSquare = -1;
		dstSquare = getSquareIdxFromDragPos(boardElm, event);
		pieceElm.dataset.square = dstSquare.toString();
	}

	function onPiecePointerUp(event: PointerEvent) {
		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;
		const rect = elm.getBoundingClientRect();

		const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
		const pieceId = elm.dataset.id ?? '';
		const symbol = elm.dataset.symbol ?? '';
	}

	function onPiecePointerDown(event: PointerEvent) {
		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;
		const rect = elm.getBoundingClientRect();

		const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
		const pieceId = elm.dataset.id ?? '';
		const symbol = elm.dataset.symbol ?? '';

		xOffset = elm.clientWidth / 2;
		yOffset = elm.clientHeight / 2;

		piecesAnimationState[pieceId].dx = clientX + (piecesAnimationState?.[pieceId]?.dx ?? 0) - rect.left - xOffset;
		piecesAnimationState[pieceId].dy = clientY + (piecesAnimationState?.[pieceId]?.dy ?? 0) - rect.top - yOffset;

		translateElm(elm, piecesAnimationState[pieceId].dx, piecesAnimationState[pieceId].dy);
		pieceState = 'waitingSecond';
	}

	function onPieceDrag(event: DragEvent) {
		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;

		const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
		const pieceId = elm.dataset.id ?? '';
		const symbol = elm.dataset.symbol ?? '';

		piecesAnimationState[pieceId].dx = clientX - (piecesAnimationState?.[pieceId]?.initialX ?? 0);
		piecesAnimationState[pieceId].dy = clientY - (piecesAnimationState?.[pieceId]?.initialY ?? 0);

		translateElm(elm, piecesAnimationState[pieceId].dx, piecesAnimationState[pieceId].dy);
	}

	function onPieceDragStart(event: DragEvent) {
		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;

		const sqIdx = Number.parseInt(elm.dataset.square ?? '-1');
		const pieceId = elm.dataset.id ?? '';
		const symbol = elm.dataset.symbol ?? '';

		piecesAnimationState[pieceId].initialX = clientX - (piecesAnimationState?.[pieceId]?.dx ?? 0);
		piecesAnimationState[pieceId].initialY = clientY - (piecesAnimationState?.[pieceId]?.dy ?? 0);

		pieceState = 'dragging';
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

		piecesAnimationState[pieceId].initialX = clientX;
		piecesAnimationState[pieceId].initialY = clientY;
		pieceState = 'idle';
		draggedElm = null;
		draggedElmId = '';
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

	function getSquareIdxFromDragPos(boardElm: HTMLDivElement, e: MouseEvent | DragEvent) {
		const boardSize = boardElm.clientWidth;
		const squareSize = boardSize / 8;
		const dx = e.clientX - boardElm.offsetLeft;
		const dy = e.clientY - boardElm.offsetTop;

		const file = Math.max(0, Math.min(7, Math.floor(dx / squareSize)));
		const rank = Math.max(0, Math.min(7, Math.floor(dy / squareSize)));

		const dstSquare = rank * 8 + file;
		return dstSquare;
	}

	function setupDomElementsAndInitialVars(): void {
		transparentDragImage = makeTransparentDragImage();
		pieceElements = document.querySelectorAll('.piece');
		squareElements = document.querySelectorAll('.square');
	}

	function setInitialPieceAnimationState(squares: Square[]): void {
		for (const sq of squares) {
			if (sq.piece !== null) {
				piecesAnimationState[sq.piece.id] = { initialX: 0, initialY: 0, dx: 0, dy: 0 };
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

				const { row, col } = new Square(sqIdx, null);

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
					duration: !firstAnimationPlayed ? 200 : 0,
				};

				const anim = p.animate(keyframes, opts);
				anim.onfinish = () => {
					piecesAnimationState[pieceId].initialX = dx;
					piecesAnimationState[pieceId].initialY = dy;
				};
			}
		}

		firstAnimationPlayed = true;
	}

	onMount(() => {
		initBoardTheme('/images/board/svg/brown.svg');
		initPieceTheme('cburnett');

		board = new Board(fen);

		tick().then(() => {
			if (board) {
				setupDomElementsAndInitialVars();
				setInitialPieceAnimationState(board.squares);
				animatePiecesMoves();
			}
		});
	});
</script>

{#if board}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		bind:this={boardElm}
		bind:clientWidth={boardWidth}
		bind:clientHeight={boardHeight}
		class="board"
		style="--board-size: {boardSize}px;"
		data-orientation={board.orientation}
		on:dragenter={onBoardDragEnter}
		on:dragleave={onBoardDragLeave}
		on:dragover={onBoardDragOver}
		on:drop={onBoardDrop}
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
			<JuicerSquare square={sq} selected={false} highlighted={false} />

			{#if sq.piece !== null}
				<JuicerPiece
					square={sq.squareIdx}
					piece={sq.piece}
					on:pointerdown={onPiecePointerDown}
					on:pointerup={onPiecePointerUp}
					on:drag={onPieceDrag}
					on:dragstart={onPieceDragStart}
					on:dragend={onPieceDragEnd}
				/>
			{/if}
		{/each}
	</div>
{/if}

<pre>{JSON.stringify(
		{
			draggedElmId,
			boardState: pieceState,
			dragOverSquare,
			activeSquare,
			dstSquare,
			srcSquare,
			xOffset,
			yOffset,
		},
		null,
		2
	)}</pre>

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
