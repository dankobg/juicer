<script lang="ts">
	import { online } from 'svelte/reactivity/window';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import PlayerInfoMobile from './player-info-mobile.svelte';

	type Props = {
		username: string;
		avatar?: string;
	};

	let { username, avatar }: Props = $props();
</script>

<PlayerInfoMobile
	{username}
	{avatar}
	avatarAlt="my player avatar"
	online={online.current ?? false}
	hasActiveTurn={gameManager.hasActiveTurn}
	showFirstMoveTimer={gameManager.showOwnFirstMoveTimer}
	firstMoveTimerText="{gameManager.ownFirstMoveTimer.time.seconds} seconds to play first move"
	showReconnectTimer={gameManager.showOpponentReconnectTimer}
	reconnectTimerText="{gameManager.opponentReconnectTimer.time.seconds} seconds to reconnect"
	clock={gameManager.ownGameTimer?.time.formatted}
/>
