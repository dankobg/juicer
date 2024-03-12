<script lang="ts">
	import { type Coordinate, Piece, type DragPosition, dragPositionZero } from './model';
	import { boardState } from './state';

	export let piece: Piece;
	export let square: Coordinate;
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	style="translate: {$boardState.dragPos[piece.id].dx}px {$boardState.dragPos[piece.id].dy}px;"
	class="piece"
	data-id={piece.id}
	data-sq={square}
	data-symbol={piece.toFenSymbol()}
	data-color={piece.color}
	data-dragging={$boardState.dragging && $boardState.draggedElm?.dataset.id === piece.id}
	on:pointerdown
	on:pointerup
	on:pointermove
	on:pointercancel
	on:touchstart
	on:dragstart
/>

<style>
	.piece {
		position: absolute;
		width: calc(var(--board-width) / var(--board-files));
		height: calc(var(--board-height) / var(--board-ranks));
		display: flex;
		justify-content: center;
		align-items: center;
		background-size: contain;
		background-position: center;
		background-repeat: no-repeat;
		cursor: pointer;
		z-index: 60;
		touch-action: none;
	}
	.piece[data-dragging='true'] {
		z-index: 65;
	}

	.piece[data-symbol='r'] {
		background-image: var(--br-theme);
	}
	.piece[data-symbol='b'] {
		background-image: var(--bb-theme);
	}
	.piece[data-symbol='n'] {
		background-image: var(--bn-theme);
	}
	.piece[data-symbol='q'] {
		background-image: var(--bq-theme);
	}
	.piece[data-symbol='k'] {
		background-image: var(--bk-theme);
	}
	.piece[data-symbol='p'] {
		background-image: var(--bp-theme);
	}
	.piece[data-symbol='R'] {
		background-image: var(--wr-theme);
	}
	.piece[data-symbol='B'] {
		background-image: var(--wb-theme);
	}
	.piece[data-symbol='N'] {
		background-image: var(--wn-theme);
	}
	.piece[data-symbol='Q'] {
		background-image: var(--wq-theme);
	}
	.piece[data-symbol='K'] {
		background-image: var(--wk-theme);
	}
	.piece[data-symbol='P'] {
		background-image: var(--wp-theme);
	}
</style>
