<script lang="ts">
	import type { Color } from './types';

	export let orientation: Color = 'w';
</script>

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

<style>
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

	.file-notations[data-orientation='w'] div:nth-child(odd) {
		color: #f6f6f6;
	}
	.file-notations[data-orientation='w'] div:nth-child(even) {
		color: #333;
	}
	.file-notations[data-orientation='b'] div:nth-child(odd) {
		color: #333;
	}
	.file-notations[data-orientation='b'] div:nth-child(even) {
		color: #f6f6f6;
	}
</style>
