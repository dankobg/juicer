import { MessageSchema } from '$lib/gen/juicer_pb';
import { fromJsonString } from '@bufbuild/protobuf';
import { latencyStats } from './latency.svelte';
import { presenceManager } from '$lib/gameplay/presence-manager.svelte';
import { lobbyManager } from '$lib/gameplay/lobby-manager.svelte';
import { gameManager } from '$lib/gameplay/game-manager.svelte';

export function onWsOpen(event: Event): void {
	console.debug('ws open:', event);
}

export function onWsClose(event: CloseEvent): void {
	console.debug(`ws closed: code: ${event.code}, reason: ${event.reason}, wasClean: ${event.wasClean}`);
}

export function onWsError(event: Event): void {
	console.debug('ws error:', event);
}

export function onWsMessage(event: MessageEvent): void {
	try {
		const msg = fromJsonString(MessageSchema, event.data);

		switch (msg.event.case) {
			case 'problem':
				// handle problem err msgs
				break;
			case 'echo':
				presenceManager.onEcho(msg.event.value);
				break;
			case 'latency':
				latencyStats.onLatency(msg.event.value);
				break;
			case 'presenceState':
				presenceManager.onPresenceState(msg.event.value);
				break;
			case 'presenceDiff':
				presenceManager.onPresenceDiff(msg.event.value);
				break;
			case 'lobbyChat':
				lobbyManager.onLobbyChat(msg.event.value);
				break;
			case 'lobbyChats':
				lobbyManager.onLobbyChatList(msg.event.value);
				break;
			case 'gameChat':
				gameManager.onGameChat(msg.event.value);
				break;
			case 'gameChats':
				gameManager.onGameChatList(msg.event.value);
				break;
			case 'gameFound':
				gameManager.onGameFound(msg.event.value);
				break;
			case 'gameInfo':
				gameManager.onGameInfo(msg.event.value);
				break;
			case 'moveAck':
				gameManager.onMoveAck(msg.event.value);
				break;
			case 'moveSync':
				gameManager.onMoveSync(msg.event.value);
				break;
			case 'abortGame':
				gameManager.onAbortGame(msg.event.value);
				break;
			case 'resignGame':
				gameManager.onResignGame(msg.event.value);
				break;
			case 'offerDraw':
				gameManager.onOfferDraw(msg.event.value);
				break;
			case 'acceptDraw':
				gameManager.onAcceptDraw(msg.event.value);
				break;
			case 'declineDraw':
				gameManager.onDeclinedDraw(msg.event.value);
				break;
			case 'gameFinished':
				gameManager.onGameFinished(msg.event.value);
				break;
			default:
				console.error('unknown message', msg.event.case, msg.event.value);
				break;
		}
	} catch (error: unknown) {
		console.error('json parse msg data', error);
	}
}
