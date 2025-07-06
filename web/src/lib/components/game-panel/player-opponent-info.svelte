<script lang="ts">
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
	avatarAlt="opponent player avatar"
	online={false}
	hasActiveTurn={gameManager.opponentHasActiveTurn}
	showFirstMoveTimer={gameManager.showOpponentFirstMoveTimer}
	firstMoveTimerText="{gameManager.opponentFirstMoveTimer.time.seconds} seconds to play first move"
	showReconnectTimer={gameManager.showOpponentReconnectTimer}
	reconnectTimerText="{gameManager.opponentReconnectTimer.time.seconds} seconds to reconnect"
	clock={gameManager.opponentGameTimer?.time.formatted}
	progressMax={gameManager.gameTimeControl?.clock && Number(gameManager.gameTimeControl.clock.seconds) * 1000}
	progressValue={gameManager.opponentGameTimer && gameManager.opponentGameTimer.timeMs}
/>
