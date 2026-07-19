<script lang="ts">
	import ChatBox from './chat-box.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import type { Game } from '$lib/gameplay/game.svelte';
	import { presenceManager } from '$lib/gameplay/presence-manager.svelte';
	import { chatManager } from '$lib/gameplay/chat-manager.svelte';
	import { gameManager } from '$lib/gameplay/game-manager.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import IconX from '@lucide/svelte/icons/x';
	import { gameChatDialog, toggleChatDialog } from './game-chat-dialog-state.svelte';

	let { chatUserId, game }: { chatUserId: string; game: Game } = $props();

	let gameUserPresences = $derived(presenceManager.getPresenceInChannel(`game.${game?.gameId}`));
	let gameChatMessages = $derived(chatManager.channelChats[`game.${game?.gameId}.chat`]?.messages ?? []);
</script>

<Dialog.Root
	open={gameChatDialog.open}
	onOpenChange={() => {
		toggleChatDialog();
	}}
>
	<Dialog.Content class="z-999 min-h-[70vh] p-0" showCloseButton={false}>
		<div class="relative">
			<Button size="icon" variant="ghost" class="absolute top-0 right-3 z-999" onclick={() => toggleChatDialog()}>
				<IconX />
			</Button>
			<ChatBox
				title="Game chat"
				channel={`game.${game.gameId}.chat`}
				{chatUserId}
				messages={gameChatMessages}
				presences={gameUserPresences}
				onSend={msg => {
					gameManager.sendGameChat(game.gameId!, msg);
				}}
				onLoadMore={() => gameManager.fetchOlderGameChatMessages(game.gameId!)}
			/>
		</div>
	</Dialog.Content>
</Dialog.Root>
