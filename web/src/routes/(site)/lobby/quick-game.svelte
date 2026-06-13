<script lang="ts">
	import { page } from '$app/state';
	import { quickGameIcons } from '$lib/components/quick-game-icons/quick-game-icons.svelte';
	import { lobbyManager } from '$lib/gameplay/lobby-manager.svelte';
	import type { components } from '$lib/gen/juicer_openapi';
	import { ws } from '$lib/ws/juicer-ws.svelte';

	let {
		gameTimeCategories,
		quickGames
	}: {
		gameTimeCategories: components['schemas']['GameTimeCategory'][];
		quickGames: components['schemas']['QuickGame'][];
	} = $props();

	function determineGameTimeCategoryFromTimeControl(
		clockSecs: number,
		incrementSecs: number,
		categories: components['schemas']['GameTimeCategory'][]
	): components['schemas']['GameTimeCategory'] | undefined {
		if (clockSecs <= 0) {
			throw new Error('clock must be > 0');
		}

		if (incrementSecs < 0) {
			throw new Error(`increment must be >= 0`);
		}

		const avgMovesEstimate = 40;
		const totalTime = clockSecs + incrementSecs * avgMovesEstimate;

		const category = categories
			.toSorted((x, y) => {
				const a = x.upper_time_limit_secs === undefined ? Infinity : x.upper_time_limit_secs;
				const b = y.upper_time_limit_secs === undefined ? Infinity : y.upper_time_limit_secs;
				return a - b;
			})
			.find(gtc => {
				if (gtc.upper_time_limit_secs !== undefined) {
					return totalTime < gtc.upper_time_limit_secs;
				}

				return true;
			});

		return category;
	}

	function formatPresetTimeControl(clockSecs: number, incrementSecs: number): string {
		return `${clockSecs / 60}+${incrementSecs}`;
	}

	function formatCategoryNameFromControl(clockSecs: number, incrementSecs: number): string {
		const category = determineGameTimeCategoryFromTimeControl(clockSecs, incrementSecs, gameTimeCategories);
		if (!category) {
			return '';
		}
		return category.name[0]?.toUpperCase() + category.name.slice(1);
	}

	function onSeekGame(quickGame: components['schemas']['QuickGame']) {
		const updater = () => {
			lobbyManager.seekGame(quickGame.clock_secs * 1000, quickGame.increment_secs * 1000);
		};

		if (!document.startViewTransition) {
			updater();
			return;
		}

		document.startViewTransition(updater);
	}

	function onCancelSeekGame() {
		const updater = () => {
			lobbyManager.cancelSeekGame();
		};

		if (!document.startViewTransition) {
			updater();
			return;
		}

		document.startViewTransition(updater);
	}

	function onReseekGame(clockMs: number, incrementMs: number) {
		const updater = () => {
			lobbyManager.seekGame(clockMs, incrementMs);
		};

		if (!document.startViewTransition) {
			updater();
			return;
		}

		document.startViewTransition(updater);
	}

	let reseekDone = $state(false);

	$effect(() => {
		if (reseekDone || lobbyManager.seekingQuickGame || ws.readyState !== WebSocket.OPEN) {
			return;
		}
		const reseek = page.state.reseek;
		if (!reseek) {
			return;
		}
		onReseekGame(reseek.clockMs, reseek.incrementMs);
		reseekDone = true;
	});
</script>

<div class="rounded-xl border border-yellow-400/20">
	<div class={['quick-game', { seeking: lobbyManager.seekingQuickGame }]}>
		<div class="[view-transition-name:quick-game-content]">
			{#if lobbyManager.seekingQuickGame}
				<div class="grid gap-4 p-4">
					<p>Searching for game...</p>
					<p>Blitz</p>
					<p>Active players: 69420</p>
					{@render quickGameIcons?.['Blitz']?.('h-8 w-8')}
					<button class="rounded-md bg-yellow-600 px-3 py-1 text-black" onclick={onCancelSeekGame}>cancel seek</button>
				</div>
			{:else}
				<div class="grid w-full grid-cols-[repeat(auto-fill,minmax(min(10rem,100%),1fr))] gap-4">
					{#each quickGames as quickGame (`${quickGame.name}-${quickGame.clock_secs}-${quickGame.increment_secs}`)}
						<button
							class="flex h-full items-center justify-center gap-1 rounded-lg border-2 border-yellow-400 bg-yellow-600/10 p-2 text-xl tracking-wide text-primary transition duration-100 ease-in-out hover:bg-yellow-600/20 hover:text-yellow-400 active:bg-yellow-600/30"
							onclick={() => onSeekGame(quickGame)}
						>
							<div class="flex flex-1 flex-col gap-2 text-start">
								<p>{formatPresetTimeControl(quickGame.clock_secs, quickGame.increment_secs)}</p>
								<p>
									{#if gameTimeCategories.length > 0}
										{formatCategoryNameFromControl(quickGame.clock_secs, quickGame.increment_secs)}
									{:else}
										{quickGame.name}
									{/if}
								</p>
							</div>
							{@render quickGameIcons?.[quickGame.name as keyof typeof quickGameIcons]?.('h-8 w-8')}
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	::view-transition-group(quick-game-content) {
		animation-duration: 0.35s;
	}
	::view-transition-old(quick-game-content) {
		animation-name: scale-out;
		animation-timing-function: ease-in-out;
	}
	::view-transition-new(quick-game-content) {
		animation-name: scale-in;
		animation-timing-function: ease-in-out;
	}

	.quick-game {
		--c1: oklch(0.85 0.1973 86.82);
		--c2: oklch(0.7118 0.2097 36);
		--c3: oklch(0.74 0.1973 51.88);
		padding: 1rem;
		border-radius: 1rem;
		position: relative;
		background-color: var(--card);

		&.seeking {
			border: 0.5rem solid transparent;
			background:
				linear-gradient(var(--card), var(--card)) padding-box,
				conic-gradient(from var(--gradient-angle), transparent var(--gradient-percent), var(--c1), var(--c2), var(--c3))
					border-box;
			animation:
				spin-once 2.3s linear 1,
				spin 2.3s linear infinite;
		}
	}

	@keyframes scale-in {
		from {
			scale: 0;
			opacity: 0;
		}
		to {
			scale: 1;
			opacity: 1;
		}
	}

	@keyframes scale-out {
		from {
			scale: 1;
			opacity: 1;
		}
		to {
			scale: 0;
			opacity: 0;
		}
	}

	@property --gradient-angle {
		syntax: '<angle>';
		initial-value: 0deg;
		inherits: false;
	}

	@property --gradient-percent {
		syntax: '<percentage>';
		initial-value: 50%;
		inherits: false;
	}

	@keyframes spin-once {
		from {
			--gradient-angle: 0deg;
			--gradient-percent: 95%;
		}
		to {
			--gradient-angle: 360deg;
			--gradient-percent: 50%;
		}
	}

	@keyframes spin {
		from {
			--gradient-angle: 0deg;
		}
		to {
			--gradient-angle: 360deg;
		}
	}
</style>
