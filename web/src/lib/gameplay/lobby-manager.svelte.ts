import type { ChatMessage, ChatUser } from '$lib/components/chat-box/chat-box.svelte';
import { MessageSchema, type LobbyChat, type LobbyChatList } from '$lib/gen/juicer_pb';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { create } from '@bufbuild/protobuf';
import { presenceManager } from './presence-manager.svelte';

class LobbyManager {
	wsError = $state<string | undefined>();
	seekingQuickGame = $state<boolean>(false);

	lobbyChatMessages = $state<ChatMessage[]>([]);

	lobbyChatUsers = $derived.by(() => {
		const presences = presenceManager.lobbyChatPresence;
		const out = new Map<string, ChatUser>();
		for (const [k, v] of presences) {
			out.set(k, {
				userId: v.userId,
				username: v.username,
				guest: v.guest
			});
		}
		return out;
	});

	seekGame(clockMs: number, incrementMs: number): void {
		const seekGameMsg = create(MessageSchema, {
			event: {
				case: 'seekGame',
				value: {
					gameTimeControl: { clockMs, incrementMs }
				}
			}
		});
		ws.send(seekGameMsg);

		this.seekingQuickGame = true;
	}

	cancelSeekGame(): void {
		const cancelSeekGameMsg = create(MessageSchema, {
			event: { case: 'cancelSeekGame', value: {} }
		});
		ws.send(cancelSeekGameMsg);

		this.seekingQuickGame = false;
	}

	sendLobbyChat(message: string): void {
		const sendLobbyChatMsg = create(MessageSchema, {
			event: { case: 'sendLobbyChat', value: { message } }
		});
		ws.send(sendLobbyChatMsg);
	}

	onLobbyChat(lobbyChat: LobbyChat): void {
		this.lobbyChatMessages.push({
			userId: lobbyChat.userId,
			messageId: lobbyChat.messageId,
			message: lobbyChat.message,
			postedAt: lobbyChat.postedAt
		});
	}

	onLobbyChatList(lobbyChats: LobbyChatList): void {}
}

export const lobbyManager = new LobbyManager();
