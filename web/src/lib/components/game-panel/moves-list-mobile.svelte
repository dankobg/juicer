<script lang="ts">
	import { gameManager } from '$lib/state/game-manager.svelte';
	import IconSkipBack from '@lucide/svelte/icons/skip-back';
	import IconSkipForward from '@lucide/svelte/icons/skip-forward';
	import Button from '$lib/components/ui/button/button.svelte';
</script>

<div class="flex w-screen items-center justify-between gap-x-2 p-2">
	<Button variant="outline" size="icon" onclick={() => gameManager.movesSkipToStart()}>
		<IconSkipBack />
	</Button>
	<div class="flex flex-1 items-start gap-x-2 overflow-x-auto">
		{#each gameManager.historyMovesInfo.slice(1) as hist, idx}
			{@const num = idx + 1}
			<Button
				size="sm"
				variant={num === gameManager.historyIndex ? 'secondary' : 'outline'}
				onclick={() => gameManager.movesJumpTo(num)}
			>
				{num}. {hist.move?.san}
			</Button>
		{/each}
	</div>
	<Button variant="outline" size="icon" onclick={() => gameManager.movesSkipToEnd()}>
		<IconSkipForward />
	</Button>
</div>
