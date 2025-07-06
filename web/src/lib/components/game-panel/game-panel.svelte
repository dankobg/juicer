<script lang="ts">
	import type { User } from '$lib/kratos/service';
	import MovesList from './moves-list.svelte';
	import ActionButtons from './action-buttons.svelte';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import GameResultMessage from './game-result-message.svelte';
	import PlayerOpponentInfo from './player-opponent-info.svelte';
	import PlayerMeInfo from './player-me-info.svelte';

	let { authenticated, user }: { authenticated: boolean; user: User | null } = $props();
</script>

<div
	class="border-secondary order-2 hidden h-full max-w-[22rem] min-w-[18rem] flex-1 flex-col justify-center rounded border sm:flex md:order-none"
>
	<GameResultMessage />
	<PlayerOpponentInfo username={authenticated ? gameManager.opponentInfo?.username || '' : 'Opp (Guest)'} />
	<MovesList />
	<ActionButtons />
	<PlayerMeInfo username={authenticated ? (user?.username ?? '') : 'Me (Guest)'} avatar={user?.avatarUrl} />
</div>
