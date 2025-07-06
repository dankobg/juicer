<script module>
	export { blitzIcon, bulletIcon, rapidIcon, classicalIcon, hyperBulletIcon };
</script>

<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { gameManager, type ClockControl } from '$lib/state/game-manager.svelte';
	import type { GameTimeCategory } from '$lib/gen/juicer_openapi';

	// svelte-ignore non_reactive_update
	let cardRootElm: HTMLElement | null = null;
	let { gameTimeCategories }: { gameTimeCategories: GameTimeCategory[] } = $props();

	let searchingGameLabel = $state('');

	const quickGames = [
		{ defaultLabel: 'Hyperbullet', control: { clock: 30, increment: 0 } },
		{ defaultLabel: 'Bullet', control: { clock: 60, increment: 0 } },
		{ defaultLabel: 'Blitz', control: { clock: 180, increment: 0 } },
		{ defaultLabel: 'Blitz', control: { clock: 180, increment: 1 } },
		{ defaultLabel: 'Blitz', control: { clock: 300, increment: 0 } },
		{ defaultLabel: 'Blitz', control: { clock: 300, increment: 2 } },
		{ defaultLabel: 'Rapid', control: { clock: 600, increment: 0 } },
		{ defaultLabel: 'Rapid', control: { clock: 600, increment: 5 } },
		{ defaultLabel: 'Rapid', control: { clock: 900, increment: 0 } },
		{ defaultLabel: 'Rapid', control: { clock: 900, increment: 5 } },
		{ defaultLabel: 'Classical', control: { clock: 1800, increment: 0 } },
		{ defaultLabel: 'Classical', control: { clock: 2700, increment: 10 } }
	];

	const quickGameIcons = {
		Hyperbullet: hyperBulletIcon,
		Bullet: bulletIcon,
		Blitz: blitzIcon,
		Rapid: rapidIcon,
		Classical: classicalIcon
	};

	function determineGameTimeCategoryFromTimeControl(control: ClockControl): GameTimeCategory | undefined {
		if (control.clock <= 0) {
			throw new Error('clock must be >= 0');
		}
		if (control.increment < 0) {
			throw new Error('increment must be > 0');
		}
		const avgMovesEstimate = 40;
		const totalTime = control.clock + control.increment * avgMovesEstimate;
		const category = gameTimeCategories.find(c => {
			if (c.upperTimeLimitSecs) {
				return totalTime < c.upperTimeLimitSecs;
			}
			return true;
		});
		return category;
	}

	function formatPresetTimeControl(control: ClockControl): string {
		return `${control.clock / 60}+${control.increment}`;
	}

	function formatCategoryNameFromControl(control: ClockControl): string {
		const cat = determineGameTimeCategoryFromTimeControl(control);
		if (!cat) {
			return '';
		}
		return cat.name[0]?.toUpperCase() + cat.name.slice(1);
	}

	function onSeekPress(quickGame: (typeof quickGames)[number]) {
		const supportsVt = document.startViewTransition;

		const updater = () => {
			searchingGameLabel = quickGame.defaultLabel;
			gameManager.seekGame(quickGame.control);
			if (!supportsVt) {
				cardRootElm?.classList.toggle('seeking');
			}
		};
		if (!supportsVt) {
			updater();
		} else {
			const vt = document.startViewTransition(updater);
			vt.finished.then(() => {
				cardRootElm?.classList.toggle('seeking');
			});
		}
	}

	function onCancelSeekPress() {
		const updater = () => {
			searchingGameLabel = '';
			gameManager.cancelSeekGame();
			cardRootElm?.classList.toggle('seeking');
		};
		if (!document.startViewTransition) {
			updater();
		} else {
			document.startViewTransition(updater);
		}
	}

	let infoFeatureToggle = $state(false);
</script>

<Card.Root class="max-auto min-[37.5rem]:m-0 m-4" bind:ref={cardRootElm}>
	<Card.Content class="space-y-2 [view-transition-name:quick-game-content]">
		{#if gameManager.uiState === 'seeking'}
			<div class="flex flex-col items-center justify-center gap-4">
				<p>Searching for a game...</p>

				<div class="flex flex-wrap items-center justify-center gap-2">
					<p>{searchingGameLabel}</p>
					<p>
						{(Number(gameManager.gameTimeControl?.clock?.seconds) || 0) / 60}+{Number(
							gameManager.gameTimeControl?.increment?.seconds || 0
						)}
					</p>
					{@render quickGameIcons[searchingGameLabel as keyof typeof quickGameIcons]('h-8 w-8')}
				</div>

				{#if infoFeatureToggle}
					<p>in lobby: {gameManager.lobbyCount}</p>
					<p>playing: {gameManager.playingCount}</p>
				{/if}

				<Button
					class="bg-orange-400 hover:bg-orange-500/80 active:bg-orange-700/70"
					size="lg"
					onclick={onCancelSeekPress}>Cancel Seek</Button
				>
			</div>
		{:else}
			<div class="grid grid-cols-[repeat(auto-fit,minmax(10rem,1fr))] gap-4">
				{#each quickGames as quickGame}
					<button
						data-id="pool-{quickGame.control.clock}-{quickGame.control.increment}"
						class="text-primary flex h-full items-center gap-1 rounded-md border border-2 border-orange-400 bg-orange-600/10 p-2 text-xl tracking-wide transition duration-100 ease-in-out hover:bg-orange-600/20 hover:text-orange-400 active:bg-orange-600/30"
						onclick={() => onSeekPress(quickGame)}
					>
						<div class="flex flex-1 flex-col gap-2 text-start">
							<p>
								{formatPresetTimeControl(quickGame.control)}
							</p>
							<p>
								{#if gameTimeCategories.length > 0}
									{formatCategoryNameFromControl(quickGame.control)}
								{:else}
									{quickGame.defaultLabel}
								{/if}
							</p>
						</div>
						{@render quickGameIcons[quickGame.defaultLabel as keyof typeof quickGameIcons]('h-8 w-8')}
					</button>
				{/each}
			</div>
		{/if}
	</Card.Content>
</Card.Root>

{#snippet blitzIcon(className: string)}
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" class={className}>
		<path
			fill="currentColor"
			d="M5.52.359A.5.5 0 0 1 6 0h4a.5.5 0 0 1 .474.658L8.694 6H12.5a.5.5 0 0 1 .395.807l-7 9a.5.5 0 0 1-.873-.454L6.823 9.5H3.5a.5.5 0 0 1-.48-.641zM6.374 1L4.168 8.5H7.5a.5.5 0 0 1 .478.647L6.78 13.04L11.478 7H8a.5.5 0 0 1-.474-.658L9.306 1z"
		/>
	</svg>
{/snippet}

{#snippet bulletIcon(className: string)}
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class={className}>
		<path
			fill="currentColor"
			d="M265.8 18.05c-4.7 38.56-4.7 38.56-38.4 57.92c38.6 4.73 38.6 4.73 58 38.43c4.7-38.58 4.7-38.58 38.4-57.95c-38.6-4.73-38.6-4.73-58-38.4m206.3 20.59c-3.8 1.14-9 3.12-15.2 6.04c-14.1 6.57-32.6 17.05-51.9 29c-38.5 23.86-80.5 54.32-96.1 70.42l-.8.8l-42 24.4c3.6 2.2 7 4.6 10.5 7.3c12.8 9.9 25.3 22.6 32 28.9l-12.2 13.2c-7.5-7-19.4-19.1-30.8-27.9c-5.6-4.3-11.2-7.8-15-9.3c-2.2-.8-3.3-1-3.8-1l-.8.5L60.57 366.2c3.35.5 6.73 1.4 10.09 2.5c14.85 4.9 30.54 14.9 44.84 29.2c14.2 14.2 24.2 29.9 29.2 44.7c.6 1.9 1.2 3.8 1.6 5.8l183.3-183.3l36-58.6l.7-.8c17.8-17.7 48.1-60.4 71.6-99.3c11.8-19.41 22-38.06 28.3-52.18c2.9-6.4 4.8-11.71 5.9-15.58M438 153.2c4.1 31.3 4.1 31.3-18.4 53.5c31.4-4.2 31.4-4.2 53.5 18.2c-4.2-31.2-4.2-31.2 18.1-53.4c-31.1 4.1-31.1 4.1-53.2-18.3M85.47 185.4c-16.43 30.2-16.43 30.2-50.41 35.3c30.18 16.5 30.18 16.5 35.3 50.4C86.79 241 86.79 241 120.7 235.8c-30.14-16.5-30.14-16.5-35.23-50.4m333.03 55.2c-25.1 52-25.1 52-81.9 63.1c52.1 25.1 52.1 25.1 63.2 81.9c25.1-52.1 25.1-52.1 81.8-63.1c-52-25.2-52-25.2-63.1-81.9M52.38 383.5c-4.41 0-7.54 1.2-9.37 3c-3.25 3.3-4.52 10.6-.78 22c3.82 11.3 12.45 25.2 24.89 37.7c12.45 12.4 26.31 21 37.68 24.8c11.4 3.8 18.7 2.5 22-.7c3.2-3.3 4.5-10.6.8-22c-3.9-11.3-12.5-25.2-25-37.6c-12.42-12.5-26.28-21.1-37.6-24.9c-4.98-1.7-9.19-2.3-12.62-2.3m166.12 28.4c3 25.2 3 25.2-15.4 42.9c25.3-3.1 25.3-3.1 43 15.3c-3-25.3-3-25.3 15.2-42.9c-25.2 3-25.2 3-42.8-15.3M69.32 421a20.66 7.804 45 0 1 16.83 10.1a20.66 7.804 45 0 1 9.09 20.1a20.66 7.804 45 0 1-20.13-9.1a20.66 7.804 45 0 1-9.09-20.1a20.66 7.804 45 0 1 3.3-1"
		/>
	</svg>
{/snippet}

{#snippet rapidIcon(className: string)}
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={className}>
		<path
			fill="currentColor"
			d="M7.75 19a2.75 2.75 0 0 1-2.745-2.582L5 16.25V15a3 3 0 1 1 2.562-4.56a4.5 4.5 0 0 1 1.68-.432L9.5 10h4l.248.007l.245.02l.127.017l.11-.167l.131-.178l-1.777-1.777c-.7-.7-.773-1.788-.22-2.57l.103-.134l.117-.128c.74-.74 1.918-.78 2.705-.117l.127.117l5.412 5.413a4 4 0 0 1-2.642 6.824l-.225.004l-.182-.007l-.026.064a2.75 2.75 0 0 1-2.138 1.588l-.192.019l-.174.005zm6.605-12.849a.502.502 0 0 0-.768.641l.058.07l2.808 2.808l-.638.642q-.316.317-.523.696l-.097.192l-.266.584l-.618-.173a3 3 0 0 0-.604-.104l-.208-.007H9.5c-.7 0-1.343.24-1.853.64l-.143.12l-.165.16l-.093.1l-.111.134l-.121.167l-.042.064q-.062.095-.115.196l-.058.112l-.062.137l-.034.084l-.051.142l-.055.184a3 3 0 0 0-.062.297l-.011.083l-.018.189l-.006.191v1.75c0 .648.492 1.18 1.122 1.244l.128.006h4.251v-.246a1.25 1.25 0 0 0-1.121-1.244l-.128-.006h-1a.75.75 0 0 1-.102-1.494l.102-.006h1a2.75 2.75 0 0 1 2.745 2.582l.005.168l-.001.246h1.749c.591 0 1.094-.414 1.218-.975l.02-.122l.1-.84l.822.198a2.5 2.5 0 0 0 2.48-4.067l-.122-.13zM5 10.501a1.5 1.5 0 0 0-.145 2.992l.145.006l.114-.005l.005-.025q.053-.224.126-.438l.097-.255l.106-.236l.057-.113l.13-.234l.088-.14l.145-.21l.074-.098l.108-.134l.129-.147l.15-.157c-.21-.4-.59-.69-1.037-.778l-.151-.022z"
		/>
	</svg>
{/snippet}

{#snippet classicalIcon(className: string)}
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={className}>
		<path
			fill="currentColor"
			d="M10.997 5.998a6.14 6.14 0 0 1 5.8 4.126l.075.233l.044.144h2.33a2.75 2.75 0 0 1 2.744 2.582l.005.167v1a1.75 1.75 0 0 1-1.606 1.743l-.143.005H18.62l.241.584a1.75 1.75 0 0 1-.813 2.22l-.137.064a1.8 1.8 0 0 1-.496.124l-.171.008h-1.787a1.75 1.75 0 0 1-1.51-.867l-.072-.137l-.539-1.143l.054-.007c-1.4.186-2.817.208-4.221.066l-.497-.057l-.535 1.136a1.75 1.75 0 0 1-1.583 1.005H4.75a1.75 1.75 0 0 1-1.618-2.415l.433-1.05a3.24 3.24 0 0 1-1.57-2.78a.75.75 0 0 1 .648-.742L2.745 12h1.88l.497-1.643a6.14 6.14 0 0 1 5.875-4.359m6.777 9.693q-1.158.465-2.356.765l-.549.129l.362.77a.25.25 0 0 0 .117.119l.053.018l.056.007h1.787a.25.25 0 0 0 .248-.28l-.017-.065l-.478-1.156h-.043l.411-.148zm-13.552 0l.39.152l.388.141l-.482 1.166a.25.25 0 0 0 .232.345h1.804l.057-.007a.25.25 0 0 0 .17-.137l.359-.763l.044.01a18 18 0 0 1-2.962-.906m6.775-8.194a4.64 4.64 0 0 0-4.371 3.087l-.068.207l-1.136 3.75l.163.059a16.67 16.67 0 0 0 10.42.133l.406-.133l.162-.059l-1.136-3.75a4.64 4.64 0 0 0-4.006-3.273l-.216-.016zm-6.977 6.5l.151-.5l-.507.002l.025.052q.13.25.33.446M17.37 12l.756 2.498l2.12.001a.25.25 0 0 0 .243-.192l.007-.058v-.999a1.25 1.25 0 0 0-1.122-1.243l-.128-.006z"
		/>
	</svg>
{/snippet}

{#snippet hyperBulletIcon(className: string)}
	<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class={className}>
		<path
			fill="currentColor"
			d="M13.03 14.156V60.47l119.782 144a916 916 0 0 1 34.282-36.22a948 948 0 0 1 33.53-31.938L54.065 14.156zm432.533 5.97c-2.307.043-4.7.183-7.188.405c-19.907 1.777-44.893 9.52-72.656 22.782c-45.372 21.676-98.133 57.952-150.564 105.126l-.03-.032c-.96.864-1.918 1.754-2.876 2.625a901 901 0 0 0-5.78 5.282q-.049.047-.095.094a925 925 0 0 0-8.375 7.813c-.107.1-.205.21-.313.31c-2.9 2.75-5.796 5.562-8.687 8.376a931 931 0 0 0-8.688 8.563c-.078.077-.17.14-.25.218l-.812.812C116.164 245.746 68.015 312.14 41.5 367.53c-13.316 27.82-21.125 52.866-22.938 72.814s2.15 34.025 10.97 42.844c8.818 8.818 22.895 12.78 42.843 10.968s44.995-9.59 72.813-22.906c36.475-17.46 77.708-44.312 119.687-78.625l-13-15.625c-76.125 63.634-142.623 97.127-161.97 77.78c-21.25-21.25 21.226-99.45 97.407-184.75l.344.408c12.673-14.077 26.176-28.306 40.438-42.563a1075 1075 0 0 1 38.47-36.594l-.408-.343c86.176-77.464 165.56-120.875 187-99.437c19.556 19.554-14.89 87.342-79.875 164.5l15.658 13.03c35.244-42.798 62.73-84.904 80.468-122.03c13.264-27.763 21.037-52.75 22.813-72.656s-2.235-33.953-11.064-42.78c-7.725-7.726-19.446-11.746-35.594-11.44zM281.03 203.343a1058 1058 0 0 0-39.75 37.75c-14.714 14.71-28.594 29.393-41.56 43.875l66.436 79.874c-.017.014-.045.016-.062.03l13.125 15.75l.03-.03l46.03 55.344c-25.77 6.714-52.722 5.31-77.03-7.657c4.94 6.544 9.707 13.083 15.72 19.095c58.928 58.93 146.78 66.75 196.092 17.438c49.314-49.314 41.523-137.165-17.406-196.094c-3.683-3.685-6.796-7.407-10.687-10.69c-2.463-2.075-5.342-3.71-7.876-5.624c14.742 25.24 16.597 52.502 9.625 78.22z"
		/>
	</svg>
{/snippet}

<style>
	::view-transition-old(quick-game-content) {
		animation: scale-out 0.35s ease-in-out;
	}
	::view-transition-new(quick-game-content) {
		animation: scale-in 0.35s ease-in-out;
	}

	:global([data-slot='card']) {
		position: relative;
	}

	:global([data-slot='card'].seeking)::after {
		content: '';
		position: absolute;
		display: block;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		width: calc(100% + 1rem);
		height: calc(100% + 1rem);
		z-index: -1;
		border-radius: 14px;
		background: conic-gradient(from var(--angle), transparent var(--gradient-percent), orange, orangered, yellow);
		animation:
			spin-once 2.5s linear 1,
			spin 2.5s linear infinite;
	}

	@property --angle {
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
			--angle: 0deg;
			--gradient-percent: 95%;
		}
		to {
			--angle: 360deg;
			--gradient-percent: 50%;
		}
	}

	@keyframes spin {
		from {
			--angle: 0deg;
		}
		to {
			--angle: 360deg;
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
</style>
