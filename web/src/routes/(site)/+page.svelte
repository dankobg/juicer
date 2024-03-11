<script lang="ts">
	import Board from '$lib/board/Board.svelte';
	import { FEN_START } from '$lib/board/model';

	let outerSize = '35rem';
</script>

<div style="width: {outerSize}; height: {outerSize};">
	<Board interactive fen={FEN_START} />
</div>
