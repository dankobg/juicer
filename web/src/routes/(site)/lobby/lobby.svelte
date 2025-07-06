<script lang="ts">
	import type { PageProps } from '../$types';
	import * as Tabs from '$lib/components/ui/tabs';
	import { gameManager } from '$lib/state/game-manager.svelte';
	import QuickGame from './quick-game.svelte';
	import CustomGame from './custom-game.svelte';
	import PlayFriend from './play-friend.svelte';
	import { ws } from '$lib/state/ws-state.svelte';

	let { data }: PageProps = $props();

	$effect(() => {
		if (ws.readyState !== WebSocket.OPEN) {
			let params: URLSearchParams | undefined;
			if (gameManager.chatLastId.current) {
				params = new URLSearchParams();
				params.set('last_chat_id', gameManager.chatLastId.current);
			}
			ws.connect(params);
		}
		ws.onOpen = (event: Event) => {
			console.debug('ws open:', event);
		};
		ws.onClose = (event: CloseEvent) => {
			console.debug(`ws closed: code: ${event.code}, reason: ${event.reason}, wasClean: ${event.wasClean}`);
		};
		ws.onError = (event: Event) => {
			console.debug('ws error:', event);
		};
		ws.onMessage = (event: MessageEvent) => {
			gameManager.handleWebsocketMessage(event);
		};
		return () => {
			gameManager.cancelSeekGame();
		};
	});

	let tabsfeatureToggle = $state(false);
</script>

{#if tabsfeatureToggle}
	<div class="lobby mx-auto mt-8 max-w-screen-2xl">
		<Tabs.Root value="quick" class="mx-auto w-full max-w-xl">
			<Tabs.List class="grid w-full grid-cols-3">
				<Tabs.Trigger value="quick">Quick game</Tabs.Trigger>
				<Tabs.Trigger value="custom">Custom game</Tabs.Trigger>
				<Tabs.Trigger value="friend">Play a friend</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="quick" class="mt-4">
				<QuickGame gameTimeCategories={data?.gameTimeCategories ?? []} />
			</Tabs.Content>
			<Tabs.Content value="custom" class="mt-4">
				<CustomGame gameVariants={data?.gameVariants ?? []} gameTimeKinds={data?.gameTimeKinds ?? []} />
			</Tabs.Content>
			<Tabs.Content value="friend" class="mt-4">
				<PlayFriend />
			</Tabs.Content>
		</Tabs.Root>
	</div>
{:else}
	<div class="lobby mx-auto mt-8 max-w-screen-2xl">
		<div class="mx-auto w-full max-w-xl">
			<h1 class="mb-4 text-center text-2xl">Quick games</h1>
			<QuickGame gameTimeCategories={data?.gameTimeCategories ?? []} />
		</div>
	</div>
{/if}
