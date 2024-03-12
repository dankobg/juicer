<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { browser } from '$app/environment';
	import { validateFen } from 'chess.js';
	import {
		FEN_EMPTY,
		BOARD_RANKS,
		BOARD_FILES,
		type EventData,
		Square,
		type Row,
		type Col,
		type Coordinate,
		NO_SQUARE,
	} from './model';
	import { boardState } from './state';
	import RankNotations from './RankNotations.svelte';
	import FileNotations from './FileNotations.svelte';
	import JuicerSquare from './Square.svelte';
	import JuicerPiece from './Piece.svelte';

	export let interactive = false;
	export let fen = FEN_EMPTY;
	export let boardTheme = 'brown';
	export let pieceTheme = 'cburnett';

	let pieceElements: NodeListOf<HTMLDivElement>;
	let boardElement: HTMLDivElement;
	let clientWidth: number;
	let clientHeight: number;

	$: squareWidth = clientWidth / BOARD_FILES;
	$: squareHeight = clientHeight / BOARD_RANKS;
	$: appearAnimationFinished = $boardState.piecesAnimated === $boardState.piecesCount;

	function translateElm(elm: HTMLDivElement, dx: number, dy: number): void {
		elm.style.translate = `${dx}px ${dy}px`;
	}

	function getEventData(event: PointerEvent): EventData {
		const { target, clientX, clientY, offsetX, offsetY } = event;
		const elm = target as HTMLDivElement;
		return { elm, clientX, clientY, offsetX, offsetY };
	}

	function getSquareIndexFromPointer(event: PointerEvent): number {
		const dx = event.clientX - boardElement.offsetLeft;
		const dy = event.clientY - boardElement.offsetTop;

		const col = Math.floor(dx / squareWidth);
		const row = Math.floor(dy / squareHeight);

		if (row < 0 || row > 7 || col < 0 || col > 7) {
			return -1;
		}

		return Square.calcIndex(row as Row, col as Col);
	}

	// --------------------------------------------------------------------------------------------------------------------

	function onPiecePointerDown(event: PointerEvent): void {
		if (event.button !== 0 || event.ctrlKey) {
			return;
		}

		const eventData = getEventData(event);
		const { elm } = eventData;
		const rect = elm.getBoundingClientRect();
		const pieceId = elm.dataset.id!;
		const squareIndex = getSquareIndexFromPointer(event);

		elm.setPointerCapture(event.pointerId);
		elm.style.userSelect = 'none';
		boardState.setDragStart(pieceId, elm, rect, squareIndex, eventData);
		translateElm(elm, $boardState.dragPos[pieceId].dx, $boardState.dragPos[pieceId].dy);
	}

	function onPiecePointerUp(event: PointerEvent): void {
		if (event.button !== 0 || event.ctrlKey) {
			return;
		}

		const eventData = getEventData(event);
		const { elm, clientX, clientY } = eventData;
		const pieceId = elm.dataset.id!;

		elm.style.userSelect = 'auto';
		boardState.setDragEnd();

		const elements = document.elementsFromPoint(clientX, clientY);
		const destSquareElement = elements.find(e => e.className.startsWith('square')) as HTMLDivElement | undefined;

		if (destSquareElement) {
			const destSquareCoord = destSquareElement.dataset.sq!;
			const { index, row, col } = Square.getDataFromCoord(destSquareCoord as Coordinate);
			boardState.setDestSquareIndex(index);

			if ($boardState.srcSquareIndex === $boardState.destSquareIndex) {
				snapBack(elm, pieceId);
			} else {
				snapToSquare(elm, pieceId, row, col, clientX, clientY);
			}
		} else {
			snapBack(elm, pieceId);
		}
	}

	function onPiecePointerMove(event: PointerEvent): void {
		event.stopPropagation();

		if (!$boardState.dragging) {
			return;
		}

		const eventData = getEventData(event);
		const { elm } = eventData;
		const pieceId = elm.dataset.id!;
		const squareIndex = getSquareIndexFromPointer(event);

		boardState.setDragMove(pieceId, squareIndex, eventData);
		translateElm(elm, $boardState.dragPos[pieceId].dx, $boardState.dragPos[pieceId].dy);
	}

	function onPiecePointerCancel(event: PointerEvent): void {
		const { elm } = getEventData(event);
		elm.style.userSelect = 'auto';
		boardState.setDragEnd();
	}

	function onPieceTouchStart(event: TouchEvent): void {
		event.preventDefault();
	}

	function onPieceDragStart(event: DragEvent): void {
		event.preventDefault();
	}

	function onBoardContextMenu(event: Event): void {
		event.preventDefault();
	}

	// --------------------------------------------------------------------------------------------------------------------

	function snapToSquare(
		elm: HTMLDivElement,
		pieceId: string,
		row: Row,
		col: Col,
		clientX: number,
		clientY: number
	): void {
		boardState.setSnapToSquareState(pieceId, boardElement, squareWidth, squareHeight, row, col, clientX, clientY);
		translateElm(elm, $boardState.dragPos[pieceId].dx, $boardState.dragPos[pieceId].dy);
	}

	function snapBack(elm: HTMLDivElement, pieceId: string, duration = 0): void {
		const keyFrames = { translate: '0' };
		const opts: KeyframeAnimationOptions = { easing: 'ease', duration };

		const anim = elm.animate(keyFrames, opts);
		anim.addEventListener('finish', () => {
			boardState.setSnapbackState(pieceId);
			translateElm(elm, $boardState.lastPos[pieceId].dx, $boardState.lastPos[pieceId].dy);
		});
	}

	function applyBoardTheme() {
		if (browser) {
			const rootElm = document.querySelector(':root') as HTMLElement;

			if (rootElm && $boardState.boardTheme) {
				rootElm.style.setProperty('--board-theme', `url('${$boardState.boardTheme}')`);
			}
		}
	}

	function applyPieceTheme() {
		if (browser) {
			const rootElm = document.querySelector(':root') as HTMLElement;

			if (rootElm && $boardState.pieceTheme) {
				for (const [k, v] of $boardState.pieceTheme.entries()) {
					const prefix = k === k.toUpperCase() ? 'w' : 'b';
					const themeVar = `--${prefix}${k.toLowerCase()}-theme`;
					const value = `url('${v}')`;
					rootElm.style.setProperty(themeVar, value);
				}
			}
		}
	}

	function animateElm(elm: HTMLDivElement, dx: number, dy: number, cb?: () => void) {
		const pieceId = elm.dataset.id ?? '';
		const opts = { easing: 'ease', duration: 500 };

		const keyframes = {
			translate: `${$boardState.lastPos[pieceId].dx + dx}px ${$boardState.lastPos[pieceId].dy + dy}px`,
		};

		const anim = elm.animate(keyframes, opts);
		anim.addEventListener('finish', () => {
			boardState.update(s => {
				const updatedDx = s.lastPos[pieceId].dx + dx;
				const updatedDy = s.lastPos[pieceId].dy + dy;

				s.dragPos[pieceId] = { initialX: 0, initialY: 0, dx: updatedDx, dy: updatedDy };
				s.lastPos[pieceId] = { ...s.dragPos[pieceId] };
				return s;
			});

			translateElm(elm, $boardState.dragPos[pieceId].dx, $boardState.dragPos[pieceId].dy);
			cb?.();
		});
	}

	function appearPiecesAnimate() {
		for (const elm of pieceElements) {
			const sq = elm.dataset.sq ?? '';
			const { row, col } = Square.getDataFromCoord(sq as Coordinate);
			const [dx, dy] = [col * squareWidth, row * squareHeight];

			if (!appearAnimationFinished) {
				animateElm(elm, dx, dy, boardState.incrementPiecesAnimated);
			}
		}
	}

	onMount(() => {
		boardState.setBoardTheme(`/images/board/svg/${boardTheme}.svg`);
		boardState.setPieceTheme(symbol => {
			const upper = symbol.toUpperCase();
			const name = symbol === upper ? `w${upper}` : `b${upper}`;
			return `/images/piece/${pieceTheme}/${name}.svg`;
		});

		applyBoardTheme();
		applyPieceTheme();

		const { ok, error } = validateFen(fen);
		boardState.setupFenState(error ?? null);

		if (!ok) {
			console.warn(error);
			return;
		}

		boardState.init(fen);

		tick().then(() => {
			pieceElements = document.querySelectorAll('.piece');
			appearPiecesAnimate();
		});
	});
</script>

{#if $boardState.fenInfo.state === 'error'}
	<span class="text-red-600">{$boardState.fenInfo.error}</span>
{/if}

{#if $boardState.fenInfo.state === 'success' && $boardState.board}
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		bind:this={boardElement}
		bind:clientWidth
		bind:clientHeight
		class="board"
		style="--board-width: {clientWidth}px; --board-height: {clientHeight}px; --board-ranks: {BOARD_RANKS}; --board-files: {BOARD_FILES};"
		data-orientation={$boardState.orientation}
		on:contextmenu={onBoardContextMenu}
	>
		<RankNotations orientation={$boardState.orientation} />
		<FileNotations orientation={$boardState.orientation} />

		{#each $boardState.board.squares ?? [] as sq}
			<JuicerSquare
				square={sq}
				selected={sq.index === $boardState.activeSquareIndex}
				highlighted={$boardState.activeSquareIndex === NO_SQUARE &&
					$boardState.srcSquareIndex !== NO_SQUARE &&
					$boardState.destSquareIndex !== NO_SQUARE &&
					(sq.index === $boardState.srcSquareIndex || sq.index === $boardState.destSquareIndex)}
				bordered={$boardState.activeSquareIndex !== NO_SQUARE &&
					sq.index === $boardState.dragOverSquareIndex &&
					sq.index !== $boardState.srcSquareIndex}
			/>

			{#if sq.piece !== null}
				{#if interactive}
					<JuicerPiece
						square={sq.coordinate}
						piece={sq.piece}
						on:pointerdown={onPiecePointerDown}
						on:pointerup={onPiecePointerUp}
						on:pointermove={onPiecePointerMove}
						on:pointercancel={onPiecePointerCancel}
						on:touchstart{onPieceTouchStart}
						on:dragstart={onPieceDragStart}
					/>
				{:else}
					<JuicerPiece square={sq.coordinate} piece={sq.piece} />
				{/if}
			{/if}
		{/each}
	</div>
{/if}

<pre>
	{$boardState.board?.print()}
</pre>

<style>
	.board {
		width: 100%;
		height: 100%;
		background-image: var(--board-theme);
		background-size: cover;
		position: relative;
		user-select: none;
	}
</style>
