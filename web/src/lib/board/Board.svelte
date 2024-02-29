<script lang="ts">
	import type { Coordinate, DragPosition, DropOffBoardAction } from './types';
	import { onMount, tick } from 'svelte';
	import JuicerSquare from './Square.svelte';
	import JuicerPiece from './Piece.svelte';
	import { browser } from '$app/environment';
	import { BLACK, BOARD_SIZE, FEN_STARTING_POSITION } from './common';
	import { Board } from './board';
	import { Square } from './square';
	import FileNotations from './FileNotations.svelte';
	import RankNotations from './RankNotations.svelte';

	export let fen: string = FEN_STARTING_POSITION;

	export let interactive: boolean = false;
	export let dropOffBoardAction: DropOffBoardAction = 'snapback';
	export let appearAnimationSpeedMs = 200;
	export let moveAnimationSpeedMs = 200;
	export let snapbackAnimationSpeedMs = 50;
	export let snapAnimationSpeedMs = 25;
	export let trashAnimationSpeedMs = 50;

	let board: Board | null = null;

	let boardElm: HTMLDivElement;
	let squareElements: NodeListOf<HTMLDivElement>;
	let pieceElements: NodeListOf<HTMLDivElement>;

	let boardWidth: number;
	let boardHeight: number;
	let boardSize: number = 480;

	let dragging: boolean = false;
	let draggedElm: HTMLDivElement | null = null;
	let dragOverSquare: number = -1;
	let activeSquare: number = -1;
	let dstSquare: number = -1;
	let srcSquare: number = -1;

	let offsetWidth: number = 0;
	let offsetHeight: number = 0;

	const dragPositionZero: DragPosition = { initialX: 0, initialY: 0, dx: 0, dy: 0 };
	let dragPos: Record<string, DragPosition> = {};
	let lastPos: DragPosition = dragPositionZero; // FIXME

	let offBoardBounds: boolean = false;

	function flipOrientation() {
		board?.flipOrientation();
		board = board;
	}

	function onPiecePointerDown(event: PointerEvent) {
		if (!interactive) {
			return;
		}

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

		const newPos = dragPos[pieceId];
		lastPos = { ...newPos }; // FIXME

		newPos.initialX = clientX - newPos.dx + (offsetWidth - offsetX);
		newPos.initialY = clientY - newPos.dy + (offsetHeight - offsetY);
		newPos.dx = clientX - newPos.initialX;
		newPos.dy = clientY - newPos.initialY;

		dragPos[pieceId] = { ...newPos };

		translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);
	}

	let rowDiff = 0,
		colDiff = 0,
		distDx = 0,
		distDy = 0,
		snapDx = 0,
		snapDy = 0;

	function onPiecePointerUp(event: PointerEvent) {
		if (!interactive) {
			return;
		}

		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;

		elm.style.userSelect = 'auto';

		const pieceId = elm.dataset.id ?? '';

		dragging = false;
		draggedElm = null;
		activeSquare = -1;

		const elements = document.elementsFromPoint(clientX, clientY);
		const dropSquare = elements.find(e => e.className.startsWith('square')) as HTMLDivElement | undefined;

		if (dropSquare) {
			// dropped inside board
			const sq = dropSquare.dataset.sq ?? '';
			const { squareIdx, row, col } = Square.fromCoord(sq as Coordinate);
			dstSquare = squareIdx;

			if (dstSquare === srcSquare) {
				// dropped to same starting square
				performSnapback(elm, pieceId);
			} else {
				// dropped to any other square except start
				const squareTopLeftX = col * (boardSize / BOARD_SIZE) + boardElm.offsetLeft;
				const squareTopLeftY = row * (boardSize / BOARD_SIZE) + boardElm.offsetTop;

				const snapDx = clientX - squareTopLeftX - 30;
				const snapDy = clientY - squareTopLeftY - 30;

				console.log({ row, col, squareTopLeftX, squareTopLeftY, snapDx, snapDy });

				dragPos[pieceId].dx -= snapDx;
				dragPos[pieceId].dy -= snapDy;
				translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);

				lastPos = { ...dragPos[pieceId] };

				// [rowDiff, colDiff] = calculateRowAndColDiffs(srcSquare, dstSquare);
				// distDx = colDiff * (boardSize / BOARD_SIZE);
				// distDy = rowDiff * (boardSize / BOARD_SIZE);
				// snapDx = distDx - dragPos[pieceId].dx;
				// snapDy = distDy - dragPos[pieceId].dy;
				// // Adjust signs based on direction
				// if (colDiff < 0) {
				// 	snapDx *= -1;
				// }
				// if (rowDiff < 0) {
				// 	snapDy *= -1;
				// }
				// dragPos[pieceId].dx += snapDx;
				// dragPos[pieceId].dy += snapDy;
				// translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);
			}
		} else {
			// dropped outside board
			performSnapback(elm, pieceId);
		}
	}

	function onPiecePointerCancel(event: PointerEvent) {
		if (!interactive) {
			return;
		}

		const { target } = event;
		const elm = target as HTMLDivElement;

		elm.style.userSelect = 'auto';

		dragging = false;
		draggedElm = null;
		activeSquare = -1;
	}

	function onPiecePointerMove(event: PointerEvent) {
		if (!interactive) {
			return;
		}

		if (!dragging) {
			return;
		}

		event.stopPropagation();

		const { target, clientX, clientY } = event;
		const elm = target as HTMLDivElement;

		const pieceId = elm.dataset.id ?? '';

		const sqIdx = getSquareIdxFromPointer(event);

		offBoardBounds =
			clientX < boardElm.offsetLeft ||
			clientX > boardElm.offsetLeft + boardElm.clientWidth ||
			clientY < boardElm.offsetTop ||
			clientY > boardElm.offsetTop + boardElm.clientHeight - 0;

		dragOverSquare = offBoardBounds ? -1 : sqIdx;

		const newPos = dragPos[pieceId];

		newPos.dx = clientX - newPos.initialX;
		newPos.dy = clientY - newPos.initialY;
		dragPos[pieceId] = { ...newPos };

		translateElm(elm, dragPos[pieceId].dx, dragPos[pieceId].dy);
	}

	function onPieceTouchStart(event: TouchEvent) {
		if (!interactive) {
			return;
		}

		event.preventDefault();
	}

	function onPieceDragStart(event: DragEvent) {
		if (!interactive) {
			return;
		}

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

	function translateElm(elm: HTMLDivElement, dx: number, dy: number): void {
		if (elm) {
			elm.style.translate = `${dx}px ${dy}px`;
		}
	}

	function performSnapback(elm: HTMLDivElement, pieceId: string, duration = 0): void {
		const keyFrames = { translate: '0' };
		const opts: KeyframeAnimationOptions = { easing: 'ease', duration };

		const anim = elm.animate(keyFrames, opts);
		anim.addEventListener('finish', () => {
			dragPos[pieceId] = lastPos;
			dstSquare = -1;
			translateElm(elm, lastPos.dx, lastPos.dy);
		});
	}

	function calculateRowAndColDiffs(srcIdx: number, destIdx: number) {
		const srcRow = Math.floor(srcIdx / BOARD_SIZE);
		const srcCol = srcIdx % BOARD_SIZE;

		const destRow = Math.floor(destIdx / BOARD_SIZE);
		const destCol = destIdx % BOARD_SIZE;

		const rowDiff = destRow - srcRow;
		const colDiff = destCol - srcCol;

		return [rowDiff, colDiff] as const;
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
		pieceElements = document.querySelectorAll('.piece');
		squareElements = document.querySelectorAll('.square');
	}

	function setInitialPieceAnimationState(squares: Square[]): void {
		for (const sq of squares) {
			if (sq.piece !== null) {
				dragPos[sq.piece.id] = dragPositionZero;
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
				const opts: KeyframeAnimationOptions = { fill: 'forwards', easing: 'ease-in', duration: 200 };

				keyframes.unshift(color === 'w' ? whiteKeyframesFrom : blackKeyframesFrom);

				const anim = p.animate(keyframes, opts);
				anim.addEventListener('finish', () => {
					dragPos[pieceId] = {
						...dragPositionZero,
						initialX: dx + boardElm.offsetLeft,
						initialY: dy + boardElm.offsetTop,
					};
				});
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
				setInitialPieceAnimationState(board.squares);
				animatePiecesMoves();
			}
		});
	});
</script>

<pre>{JSON.stringify(lastPos)}</pre>
<pre>{JSON.stringify(Object.values(dragPos).at(-1))}</pre>
<pre>{JSON.stringify({ rowDiff, colDiff, distDx, distDy, snapDx, snapDy })}</pre>

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
		<RankNotations orientation={board.orientation} />
		<FileNotations orientation={board.orientation} />

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

<button class="bg-blue-400 mt-5 p-2 rounded hover:bg-blue-500" on:click={flipOrientation}>FLIP ORIENTATION</button>

{#if board}
	<pre>{JSON.stringify({ dragging, offBoardBounds, dragOverSquare, activeSquare, dstSquare, srcSquare }, null, 2)}</pre>
	<pre style="margin-top:3rem;">{board.print()}</pre>
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
</style>
