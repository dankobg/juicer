<script lang="ts">
	import { gameManager } from '$lib/state/game-manager.svelte';
</script>

{#if gameManager.showGameResultMessage}
	<p class="bg-rose-400 px-4 py-2 text-center font-medium text-black">{gameManager.gameResultMessage}</p>
{/if}
