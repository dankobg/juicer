<script lang="ts">
	import type { Game } from '$lib/gameplay/game.svelte';

	let { game }: { game: Game } = $props();
</script>

{#if game.isGameConcluded}
	<p class="bg-rose-400 px-4 py-2 text-center font-medium text-black">{game.gameConcludedText}</p>
{/if}
