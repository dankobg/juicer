<script lang="ts">
	import type { Game } from '$lib/gameplay/game.svelte';
	import IconFlag from '@lucide/svelte/icons/flag';
	import IconHandshake from '@lucide/svelte/icons/handshake';
	import IconArrowUpDown from '@lucide/svelte/icons/arrow-up-down';
	import IconFootprints from '@lucide/svelte/icons/footprints';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import IconMessageSquareText from '@lucide/svelte/icons/message-square-text';
	import IconCirclePlay from '@lucide/svelte/icons/circle-play';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { buttonVariants } from '../ui/button';
	import type { Component } from 'svelte';
	import type { GameTimeControl } from '$lib/gen/juicer_pb';
	import Button from '../ui/button/button.svelte';

	let { game }: { game: Game } = $props();

	function formatTimeControl(gameTimeControl: GameTimeControl): string {
		return `${Math.floor(gameTimeControl.clockMs / 1000) / 60}+${Math.floor(gameTimeControl.incrementMs / 1000)}`;
	}
</script>

<div class="flex gap-2 items-center justify-center">
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

		{#if game.uiShowDrawAcceptDeclineButtons}
			{@render btn('Accept draw', IconCheck, () => game.acceptDraw())}
			{@render btn('Decline draw', IconX, () => game.dedclineDraw())}
		{/if}

		{#if game.uiShowRequeueQuickGame && game?.gameTimeControl}
			<Button class="font-medium" onclick={() => game?.requeueQuickGame()}>
				<IconCirclePlay />
				<span>Play ({formatTimeControl(game.gameTimeControl)}) again</span>
			</Button>
		{/if}
	{/if}
</div>

{#snippet btn(text: string, Icon: Component, onclick?: VoidFunction)}
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
