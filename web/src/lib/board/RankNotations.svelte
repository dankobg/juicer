<script lang="ts">
	import { WHITE, type Color } from './model';

	export let orientation: Color = WHITE;
</script>

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

<style>
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

	.rank-notations[data-orientation='w'] div:nth-child(odd) {
		color: #f6f6f6;
	}
	.rank-notations[data-orientation='w'] div:nth-child(even) {
		color: #333;
	}
	.rank-notations[data-orientation='b'] div:nth-child(odd) {
		color: #333;
	}
	.rank-notations[data-orientation='b'] div:nth-child(even) {
		color: #f6f6f6;
	}
</style>
