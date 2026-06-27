<script lang="ts">
	import type { Game } from '$lib/gameplay/game.svelte';
	import { Color, GameTimeCategory, type GameTimeControl } from '$lib/gen/juicer_pb';
	import { gameTimeCategoryIcons } from '../quick-game-icons/quick-game-icons.svelte';

	let { game }: { game: Game } = $props();

	function formatGameTimeCategory(gtc: GameTimeCategory): string {
		switch (gtc) {
			case GameTimeCategory.HYPERBULLET:
				return 'Hyperbullet';
			case GameTimeCategory.BULLET:
				return 'Bullet';
			case GameTimeCategory.BLITZ:
				return 'Blitz';
			case GameTimeCategory.RAPID:
				return 'Rapid';
			case GameTimeCategory.CLASSICAL:
				return 'Classical';
			case GameTimeCategory.UNSPECIFIED:
				return 'unknown';
			default:
				return 'unknown';
		}
	}

	function formatGameTimeControl(gtc: GameTimeControl): string {
		return `${gtc.clockMs / 1000 / 60}+${gtc.incrementMs / 1000}`;
	}

	let activity = $derived.by(() => {
		if (game.isGameActive) {
			if (game.isSpectating) {
				if (game.whiteDisplayFirstMoveTimeMs !== undefined) {
					return 'Waiting for white first move';
				}
				if (game.blackDisplayFirstMoveTimeMs !== undefined) {
					return 'Waiting for black first move';
				}
				if (game.whiteDisconnectedAt) {
					return 'Waiting for white to reconnect';
				}
				if (game.blackDisconnectedAt) {
					return 'Waiting for black to reconnect';
				}
				return game.currentTurn === Color.WHITE ? 'White turn' : `Black turn`;
			}

			if (game.isPlaying) {
				if (
					(game.whiteDisplayFirstMoveTimeMs !== undefined && game.myColor === Color.WHITE) ||
					(game.blackDisplayFirstMoveTimeMs !== undefined && game.myColor === Color.BLACK)
				) {
					return 'Waiting for your first move';
				}
				if (
					(game.whiteDisplayFirstMoveTimeMs !== undefined && game.myColor === Color.BLACK) ||
					(game.blackDisplayFirstMoveTimeMs !== undefined && game.myColor === Color.WHITE)
				) {
					return `Waiting for opponent's first move`;
				}
				if (
					(game.whiteDisplayReconnectTimeMs !== undefined && game.myColor === Color.BLACK) ||
					(game.blackDisplayReconnectTimeMs !== undefined && game.myColor === Color.WHITE)
				) {
					return `Waiting for opponent to reconnect`;
				}
				return game.currentTurn === game.myColor ? 'Your turn' : `Opponent's turn`;
			}
		}

		return '';
	});
</script>

<div class="flex p-4 flex-col rounded-md border border-secondary gap-4">
	{#if game.gameTimeCategory}
		<div class="flex gap-1 items-center">
			{@render gameTimeCategoryIcons?.[game.gameTimeCategory]?.('h-8 w-8')}
			<span>{formatGameTimeCategory(game.gameTimeCategory)}</span>
			{#if game.gameTimeControl}
				<span>{formatGameTimeControl(game.gameTimeControl)}</span>
			{/if}
			<span>{game.rated ? 'Rated' : 'Unrated'}</span>
		</div>
	{/if}

	{#if game.isPlaying}
		<span>You are {game.myColor === Color.WHITE ? 'White' : 'Black'}</span>
	{/if}

	{#if game.isSpectating}
		<span>You are spectating</span>
	{/if}

	<div class="flex gap-1 items-center">
		<svg class="h-5 w-5 stroke-10 fill-black stroke-white" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
			<circle cx="50" cy="50" r="40" />
		</svg>
		<span>{game.black?.username}</span>
	</div>

	<div class="flex gap-1 items-center">
		<svg class="h-5 w-5 stroke-10 fill-white stroke-black" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
			<circle cx="50" cy="50" r="40" />
		</svg>
		<span>{game.white?.username}</span>
	</div>

	{#if activity}
		<span>{activity}</span>
	{/if}

	{#if game.gameConcludedText}
		<span>{game.gameConcludedText}</span>
	{/if}
</div>
