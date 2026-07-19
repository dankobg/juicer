import { MessageSchema, type LobbyChat, type LobbyChatList } from '$lib/gen/juicer_pb';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { create } from '@bufbuild/protobuf';
import { chatManager, LOBBY_CHAT_CHANNEL, type ChatMessage } from './chat-manager.svelte';

class LobbyManager {
	wsError = $state<string | undefined>();
	seekingQuickGame = $state<boolean>(false);
	seekingGameTimeControl = $state<{ clockMs: number; incrementMs: number } | null>(null);

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
		this.seekingGameTimeControl = { clockMs, incrementMs };
	}

	cancelSeekGame(): void {
		const cancelSeekGameMsg = create(MessageSchema, {
			event: { case: 'cancelSeekGame', value: {} }
		});
		ws.send(cancelSeekGameMsg);

		this.seekingQuickGame = false;
		this.seekingGameTimeControl = null;
	}

	sendLobbyChat(message: string): void {
		const sendLobbyChatMsg = create(MessageSchema, {
			event: { case: 'sendLobbyChat', value: { message } }
		});
		ws.send(sendLobbyChatMsg);
	}

	fetchOlderLobbyChatMessages(): void {
		const cursor = chatManager.getChatCursor(LOBBY_CHAT_CHANNEL);
		const listLobbyChatsMsg = create(MessageSchema, {
			event: { case: 'listLobbyChats', value: { cursor } }
		});
		ws.send(listLobbyChatsMsg);
	}

	onLobbyChat(lobbyChat: LobbyChat): void {
		chatManager.addChatMessage('lobby.chat', {
			messageId: lobbyChat.messageId,
			message: lobbyChat.message,
			user: {
				id: lobbyChat.user?.id ?? '',
				username: lobbyChat.user?.username ?? ''
			},
			postedAt: lobbyChat.postedAt
		});
	}

	onLobbyChatList(lobbyChats: LobbyChatList): void {
		const messages: ChatMessage[] = [];
		for (const msg of lobbyChats.lobbyChats) {
			messages.push({
				messageId: msg.messageId,
				message: msg.message,
				user: {
					id: msg.user?.id ?? '',
					username: msg.user?.username ?? ''
				},
				postedAt: msg.postedAt
			});
		}
		chatManager.prependChatMessages(LOBBY_CHAT_CHANNEL, messages, lobbyChats.hasMore);
	}
}

export const lobbyManager = new LobbyManager();
