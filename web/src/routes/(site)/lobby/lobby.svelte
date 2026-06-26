<script lang="ts">
	import { ws } from '$lib/ws/juicer-ws.svelte';
	import * as Tabs from '$lib/components/ui/tabs';
	import type { PageProps } from '../$types';
	import QuickGame from './quick-game.svelte';
	import CustomGame from './custom-game.svelte';
	import PlayFriend from './play-friend.svelte';
	import { onWsClose, onWsError, onWsMessage, onWsOpen } from '$lib/ws/ws-message-handler';
	import ChatBox from '$lib/components/chat-box/chat-box.svelte';
	import { lobbyManager } from '$lib/gameplay/lobby-manager.svelte';
	import { presenceManager } from '$lib/gameplay/presence-manager.svelte';
	import { chatManager } from '$lib/gameplay/chat-manager.svelte';

	let { data }: PageProps = $props();

	$effect(() => {
		ws.onOpen = onWsOpen;
		ws.onError = onWsError;
		ws.onClose = onWsClose;
		ws.onMessage = onWsMessage;

		const params = new URLSearchParams();
		params.set('path', window.location.pathname);
		ws.connect(params);

		return () => {
			ws.close();
			lobbyManager.seekingQuickGame = false;
		};
	});
</script>

<div
	class="mx-auto mt-8 grid w-full max-w-screen-2xl grid-cols-1 justify-center gap-8 p-4 px-4 lg:grid-cols-[minmax(20rem,30rem)_minmax(20rem,42rem)]"
>
	<div class="max-h-[45rem] min-h-[30rem] w-full max-w-[30rem] justify-self-center">
		<ChatBox
			title="Lobby chat"
			channel="lobby.chat"
			chatUserId={data?.auth?.user?.id ?? ''}
			messages={chatManager.lobbyChats.messages}
			hasMore={chatManager.lobbyChats.hasMore}
			presences={presenceManager.lobbyChatPresence}
			onSend={msg => {
				lobbyManager.sendLobbyChat(msg);
			}}
			onLoadMore={() => lobbyManager.fetchOlderLobbyChatMessages()}
		/>
	</div>

	<Tabs.Root value="quick-game" class="w-full max-w-[42rem] justify-self-center">
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
