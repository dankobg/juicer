<script lang="ts">
	import type { Game } from '$lib/gameplay/game-manager.svelte';
	import IconFlag from '@lucide/svelte/icons/flag';
	import IconHandshake from '@lucide/svelte/icons/handshake';
	import IconArrowUpDown from '@lucide/svelte/icons/arrow-up-down';
	import IconFootprints from '@lucide/svelte/icons/footprints';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import IconMessageSquareText from '@lucide/svelte/icons/message-square-text';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { buttonVariants } from '../ui/button';
	import type { Component } from 'svelte';

	let { game }: { game: Game } = $props();
</script>

<div>
	{#if game}
		{#if game.uiShowAbortButton}
			{@render btn('Abort', IconFootprints, () => game.abortGame())}
		{/if}

		{#if game.uiShowResignButton}
			{@render btn('Resign', IconFlag, () => game.resignGame())}
		{/if}

		{#if game.uiShowOfferDrawButton}
			{@render btn('Offer draw', IconHandshake, () => game.offerDraw())}
		{/if}

		{#if game.uiShowFlipBoardButton}
			{@render btn('Flip board', IconArrowUpDown, () => game?.board?.flip())}
		{/if}

		{#if game.uiShowChatButton}
			{@render btn('Toggle chat', IconMessageSquareText, () => console.log('chat'))}
		{/if}

		{#if game.uiShowDrawOfferResponseButtons}
			{@render btn('Accept draw', IconCheck, () => game.acceptDraw())}
			{@render btn('Decline draw', IconX, () => game.dedclineDraw())}
		{/if}

		<!-- @TODO: new game button -->
	{/if}
</div>

{#snippet btn(text: string, Icon: Component, onclick?: () => void)}
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
