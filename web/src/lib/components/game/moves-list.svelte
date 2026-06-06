<script lang="ts">
	import IconStepBack from '@lucide/svelte/icons/step-back';
	import IconStepForward from '@lucide/svelte/icons/step-forward';
	import IconSkipBack from '@lucide/svelte/icons/skip-back';
	import IconSkipForward from '@lucide/svelte/icons/skip-forward';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { buttonVariants } from '$lib/components/ui/button/index';
	import IconArrowDown from '@lucide/svelte/icons/arrow-down-to-line';
	import { tick, type Component } from 'svelte';
	import type { Game } from '$lib/gameplay/game-manager.svelte';

	let { game }: { game: Game } = $props();

	let scrollPointElm: HTMLDivElement;
	let allowedToScrollToLatest = $state<boolean>(true);

	function onScroll(event: Event) {
		const elm = event.target as HTMLDivElement;
		const threshold = 50;
		if (elm.scrollTop + elm.clientHeight >= elm.scrollHeight - threshold) {
			allowedToScrollToLatest = true;
		} else {
			allowedToScrollToLatest = false;
		}
	}

	async function scrollToLatestMessage() {
		await tick();
		scrollPointElm.scrollIntoView({ behavior: 'smooth', block: 'end' });
	}

	$effect(() => {
		if (allowedToScrollToLatest && game?.gameMoves?.length > 0) {
			scrollToLatestMessage();
		}
	});

	function onKeydown(e: KeyboardEvent) {
		if (e.key === 'ArrowLeft') {
			e.preventDefault();
			game.movesStepBack();
		} else if (e.key === 'ArrowRight') {
			e.preventDefault();
			game.movesStepForward();
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			game.movesSkipToStart();
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			game.movesSkipToEnd();
		}
	}
</script>

{#snippet btn(text: string, Icon: Component, onclick?: () => void)}
	<Tooltip.Provider delayDuration={200}>
		<Tooltip.Root>
			<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon', class: 'rounded-full' })} {onclick}>
				<Icon />
			</Tooltip.Trigger>
			<Tooltip.Content>
				<span>{text}</span>
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/snippet}

<svelte:window onkeydown={onKeydown} />

<div class="flex flex-col">
	<div class="">
		<div class="mt-auto flex justify-center gap-4 p-2 pb-0">
			{@render btn('Skip to start', IconSkipBack, () => game.movesSkipToStart())}
			{@render btn('Step back', IconStepBack, () => game.movesStepBack())}
			{@render btn('Step forward', IconStepForward, () => game.movesStepForward())}
			{@render btn('Skip to end', IconSkipForward, () => game.movesSkipToEnd())}
			{#if !allowedToScrollToLatest}
				{@render btn('Scroll to last', IconArrowDown, () => scrollToLatestMessage())}
			{/if}
		</div>

		<div class="flex max-h-30 min-h-40 flex-col gap-1 overflow-y-auto p-2" onscroll={onScroll}>
			{#each Array(Math.ceil((game.gameMoves.length - 1) / 2)) as _, i}
				{@const w = i * 2 + 1}
				{@const b = i * 2 + 2}
				{@const h1 = game.gameMoves[w]}
				{@const h2 = game.gameMoves[b]}
				{@const d1 = game.moveDurationsMs[w]}
				{@const d2 = game.moveDurationsMs[b]}
				{@const num = i + 1}

				<div class="grid grid-cols-[10%_1fr_1fr_15%] items-center justify-between gap-4 border-b border-stone-300/10">
					<div class="text-start">{num}.</div>
					<button
						class={['hover:bg-secondary rounded-sm text-center', game.historyIndex === w && 'bg-secondary']}
						onclick={() => game.movesJumpTo(w)}
					>
						{h1?.san}
					</button>
					<button
						class={['hover:bg-secondary rounded-sm text-center', game.historyIndex === b && 'bg-secondary']}
						onclick={() => game.movesJumpTo(b)}
					>
						{h2?.san}
					</button>
					{#if game.gameMoves.length > 0}
						<div>
							{#if d1}
								<div class="flex h-full items-center justify-between gap-1 text-xs">
									<svg viewBox="0 0 5 10" xmlns="http://www.w3.org/2000/svg" class="h-full w-[6px]">
										<rect width="100%" height="100%" fill="#fff" />
									</svg>
									{(d1 / 1000).toFixed(1)}s
								</div>
							{/if}
							{#if d2}
								<div class="flex h-full items-center justify-between gap-1 text-xs">
									<svg viewBox="0 0 5 10" xmlns="http://www.w3.org/2000/svg" class="h-full w-[6px]">
										<rect width="100%" height="100%" fill="#000" stroke="#fff" />
									</svg>
									{(d2 / 1000).toFixed(1)}s
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}

			<div bind:this={scrollPointElm}></div>
		</div>
	</div>
</div>
