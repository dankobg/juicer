import type { Timestamp } from '@bufbuild/protobuf/wkt';

export type ChatMessage = {
	messageId: string;
	user: ChatUserSnapshot;
	message: string;
	postedAt?: Timestamp;
};

export type ChatUserSnapshot = {
	id: string;
	username: string;
	avatarUrl?: string;
};

export type ChatState = {
	messages: ChatMessage[];
	hasMore?: boolean;
	cursor?: string;
};

export const LOBBY_CHAT_CHANNEL = 'lobby.chat';

export function gameChatChannel(gameId: number): string {
	return `game.${gameId}.chat`;
}

export function gameTvChatChannel(gameId: number): string {
	return `gametv.${gameId}.chat`;
}

export class ChatManager {
	channelChats = $state<Record<string, ChatState>>({});
	lobbyChats = $derived(this.getChatsForChannel('lobby.chat'));

	getChatsForChannel(channel: string): ChatState {
		return this.channelChats[channel] ?? { messages: [] };
	}

	getChatCursor(channel: string): string | undefined {
		const chat = this.getChatsForChannel(channel);
		return chat?.cursor ?? chat.messages?.at(0)?.messageId;
	}

	addChatMessage(channel: string, msg: ChatMessage): void {
		this.channelChats[channel] ||= { messages: [] };
		this.channelChats[channel].messages.push(msg);
	}

	loadChatMessages(channel: string, messages: ChatMessage[], hasMore?: boolean): void {
		this.channelChats[channel] = { messages, hasMore, cursor: messages?.at(0)?.messageId };
	}

	prependChatMessages(channel: string, messages: ChatMessage[], hasMore?: boolean): void {
		const prevMsgs = this.channelChats[channel]?.messages ?? [];
		this.channelChats[channel] = {
			messages: [...messages, ...prevMsgs],
			hasMore,
			cursor: messages.at(0)?.messageId
		};
	}
}

export const chatManager = new ChatManager();
