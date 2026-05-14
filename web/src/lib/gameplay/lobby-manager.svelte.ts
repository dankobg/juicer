import { MessageSchema, type LobbyChat, type LobbyChatList } from '$lib/gen/juicer_pb';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { create } from '@bufbuild/protobuf';

class LobbyManager {
	wsError = $state<string | undefined>();
	seekingQuickGame = $state<boolean>(false);

	seekGame(clockMs: number, incrementMs: number): void {
		const seekGameMsg = create(MessageSchema, {
			event: {
				case: 'seekGame',
				value: {
					timeControl: { clockMs, incrementMs }
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

	onLobbyChat(lobbyChat: LobbyChat): void {}

	onLobbyChatList(lobbyChats: LobbyChatList): void {}
}

export const lobbyManager = new LobbyManager();
