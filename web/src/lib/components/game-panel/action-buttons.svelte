<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import IconFlag from '@lucide/svelte/icons/flag';
	import IconHandshake from '@lucide/svelte/icons/handshake';
	import IconArrowUpDown from '@lucide/svelte/icons/arrow-up-down';
	import IconFootprints from '@lucide/svelte/icons/footprints';
	import IconMessageSquareText from '@lucide/svelte/icons/message-square-text';
	import { type Component } from 'svelte';
	import { goto } from '$app/navigation';
	import { GameState } from '$lib/gen/juicer_pb';
</script>

{#snippet gameActionButton(text: string, Icon: Component, onclick: () => void)}
	<Tooltip.Provider delayDuration={200}>
		<Tooltip.Root>
			<Tooltip.Trigger class={buttonVariants({ variant: 'outline', size: 'icon', class: 'rounded-full' })} {onclick}>
				<Icon />
			</Tooltip.Trigger>
			<Tooltip.Content>
				<span>{text}</span>
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/snippet}

<div class="flex flex-col gap-2 p-2">
	<div class="flex justify-center gap-4">
		{#if gameManager.showAbort}
			{@render gameActionButton('Abort', IconFootprints, () => gameManager.gameAbort())}
		{/if}
		{#if gameManager.showResign}
			{@render gameActionButton('Resign', IconFlag, () => gameManager.gameResign())}
		{/if}
		{#if gameManager.showOfferDraw}
			{@render gameActionButton('Offer draw', IconHandshake, () => gameManager.gameOfferDraw())}
		{/if}
		{@render gameActionButton('Flip board', IconArrowUpDown, () => gameManager.board.flip())}
		<div class="lg:hidden">
			{@render gameActionButton('Chat', IconMessageSquareText, () => gameManager.toggleChatDialog())}
		</div>
	</div>

	{#if gameManager.opponenOfferedDraw}
		<div class="flex justify-center gap-2">
			<p>Draw offered</p>
			<div class="flex flex-wrap items-center gap-2">
				<button
					onclick={() => gameManager.gameAcceptDraw()}
					class="rounded-full border border-green-800 bg-green-600/30 px-2 hover:bg-green-400/40 active:bg-green-400/60"
				>
					Accept
				</button>
				<button
					onclick={() => gameManager.gameDeclineDraw()}
					class="rounded-full border border-red-800 bg-red-600/30 px-2 hover:bg-red-400/40 active:bg-red-400/60"
				>
					Decline
				</button>
			</div>
		</div>
	{/if}

	{#if gameManager.gameState === GameState.FINISHED || (gameManager.gameState === GameState.INTERRUPTED && gameManager.poolLast.current.val)}
		<Button
			class="mt-4"
			onclick={() => {
				goto('/').finally(() => {
					const ctrl = gameManager.poolLast.current.val;
					if (!ctrl) {
						return;
					}
					const btn = document.querySelector<HTMLButtonElement>(
						`button[data-id='pool-${ctrl.clock}-${ctrl.increment}']`
					);
					btn?.click();
				});
			}}
		>
			New Game ({(gameManager.poolLast.current.val?.clock || 0) / 60}+{gameManager.poolLast.current.val?.increment})
		</Button>
	{/if}
</div>
