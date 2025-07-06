<script lang="ts">
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
	avatarAlt="opponent player avatar"
	online={false}
	hasActiveTurn={gameManager.opponentHasActiveTurn}
	showFirstMoveTimer={gameManager.showOpponentFirstMoveTimer}
	firstMoveTimerText="{gameManager.opponentFirstMoveTimer.time.seconds} seconds to play first move"
	showReconnectTimer={gameManager.showOpponentReconnectTimer}
	reconnectTimerText="{gameManager.opponentReconnectTimer.time.seconds} seconds to reconnect"
	clock={gameManager.opponentGameTimer?.time.formatted}
/>
