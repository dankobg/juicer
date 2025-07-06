<script lang="ts">
	import { online } from 'svelte/reactivity/window';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import PlayerInfo from './player-info.svelte';

	type Props = {
		username: string;
		avatar?: string;
	};

	let { username, avatar }: Props = $props();
</script>

<PlayerInfo
	{username}
	{avatar}
	avatarAlt="my player logo"
	online={online.current ?? false}
	hasActiveTurn={gameManager.hasActiveTurn}
	showFirstMoveTimer={gameManager.showOwnFirstMoveTimer}
	firstMoveTimerText="{gameManager.ownFirstMoveTimer.time.seconds} seconds to play first move"
	clock={gameManager.ownGameTimer?.time.formatted}
	progressMax={gameManager.gameTimeControl?.clock && Number(gameManager.gameTimeControl.clock.seconds) * 1000}
	progressValue={gameManager.ownGameTimer && gameManager.ownGameTimer.timeMs}
/>
