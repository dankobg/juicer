<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->

<script lang="ts">
	import type { Piece } from './piece';
	import type { Coordinate } from './types';

	export let piece: Piece;
	export let square: Coordinate;
</script>

<div
	class="piece"
	data-id={piece.id}
	data-sq={square}
	data-symbol={piece.toFenSymbol()}
	data-color={piece.color}
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
		width: calc(var(--board-size) / 8);
		height: calc(var(--board-size) / 8);
		display: flex;
		justify-content: center;
		align-items: center;
		background-size: contain;
		background-position: center;
		background-repeat: no-repeat;
		cursor: pointer;
		z-index: 60;
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
