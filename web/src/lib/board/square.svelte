<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->

<script lang="ts">
	import type { Square } from './square';

	export let square: Square;
	export let selected: boolean = false;
	export let bordered: boolean = false;
	export let highlighted: boolean = false;
</script>

<div
	class="square"
	style="--row:{square.row}; --col:{square.col};"
	data-sq={square.coord}
	data-selected={selected}
	data-bordered={bordered}
	data-highlighted={highlighted}
	data-color={square.color}
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
		background-color: rgba(165, 231, 59, 0.5);
	}

	.square[data-selected='true'] {
		background-color: rgba(116, 79, 202, 0.5);
	}

	.square[data-bordered='true'] {
		border: 3px solid rgb(201, 126, 77);
	}
</style>
