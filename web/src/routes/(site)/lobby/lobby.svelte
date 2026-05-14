<script lang="ts">
	import { ws } from '$lib/ws/juicer-ws.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import type { PageProps } from '../$types';
	import QuickGame from './quick-game.svelte';
	import CustomGame from './custom-game.svelte';
	import PlayFriend from './play-friend.svelte';
	import { onWsClose, onWsError, onWsMessage, onWsOpen } from '$lib/ws/ws-message-handler';

	let { data }: PageProps = $props();

	$effect(() => {
		const params = new URLSearchParams();
		params.set('path', window.location.pathname);
		ws.connect(params);

		ws.onOpen = onWsOpen;
		ws.onError = onWsError;
		ws.onClose = onWsClose;
		ws.onMessage = onWsMessage;

		return () => {
			ws.close();
		};
	});
</script>

<div class="lobby mx-auto mt-8 w-full max-w-screen-2xl">
	<Tabs.Root value="quick-game" class="mx-auto w-full max-w-xl">
		<Tabs.List class="grid w-full grid-cols-3">
			<Tabs.Trigger value="quick-game">Quick game</Tabs.Trigger>
			<Tabs.Trigger value="custom-game">Custom game</Tabs.Trigger>
			<Tabs.Trigger value="play-friend">Play a friend</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="quick-game" class="mt-4">
			{#if data?.gameTimeCategoriesResult?.data?.data && data?.quickGamesResult?.data}
				<QuickGame
					gameTimeCategories={data.gameTimeCategoriesResult.data.data}
					quickGames={data.quickGamesResult.data}
				/>
			{/if}
		</Tabs.Content>
		<Tabs.Content value="custom-game" class="mt-4">
			{#if data?.gameVariantsResult?.data?.data && data?.gameTimeKindsResult?.data?.data}
				<CustomGame
					gameVariants={data.gameVariantsResult.data.data}
					gameTimeKinds={data.gameTimeKindsResult.data.data}
				/>
			{/if}
		</Tabs.Content>
		<Tabs.Content value="play-friend" class="mt-4">
			<PlayFriend />
		</Tabs.Content>
	</Tabs.Root>
</div>
