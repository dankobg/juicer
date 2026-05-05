import {
	Color,
	GameTimeControlSchema,
	MessageSchema,
	type Echo,
	type GameTimeControl,
	type GameMove,
	type Latency,
	type LobbyChat,
	type OpponentInfo,
	type Presence,
	type PresenceDiff,
	type PresenceState,
	type GameFound
} from '$lib/gen/juicer_pb';
import { create, fromJsonString } from '@bufbuild/protobuf';
import type { JuicerBoard, Coord, PieceFenSymbol } from '@dankop/juicer-board';
import { PersistedState } from 'runed';
import { ws } from '$lib/state/ws-state.svelte';
import { type Timestamp } from '@bufbuild/protobuf/wkt';
import { SvelteSet } from 'svelte/reactivity';
import { goto } from '$app/navigation';

export const RECONNECT_TIMEOUT_MS = 15_000;
export const FIRST_MOVE_TIMEOUT_MS = 10_000;

class GameManager {
	board!: JuicerBoard;
	wsError = $state<string | undefined>();
	uiState = $state<'idle' | 'seeking' | 'playing'>('idle');
	poolLast = new PersistedState<{ clockMs: number; incrementMs: number } | null>('juicer-pool-last', null);
	gameTimeControl = $state<GameTimeControl | null>(null);
	gameId = $state<string | undefined>();
	clientId = $state<string | undefined>();
	fen = $state<string | undefined>();
	color = $state<Color>(Color.UNSPECIFIED);
	orientation = $state<Color>(Color.UNSPECIFIED);
	uci = $state<string | undefined>();
	lan = $state<string | undefined>();
	san = $state<string | undefined>();
	ply = $state<number>(0);
	legalMoves = $state<string[]>([]);
	gameMoves = $state<GameMove[]>([]);
	historyIndex = $state<number>(0);
	opponentInfo = $state<OpponentInfo | null>(null);
	startTime = $state<Timestamp>();
	reconnectTimeoutMs = $state<number>(RECONNECT_TIMEOUT_MS);
	firstMoveTimeoutMs = $state<number>(FIRST_MOVE_TIMEOUT_MS);

	seekGame(clockMs: number, incrementMs: number): void {
		this.gameTimeControl = create(GameTimeControlSchema, { clockMs, incrementMs });

		const seekGameMsg = create(MessageSchema, {
			event: {
				case: 'seekGame',
				value: {
					timeControl: { clockMs, incrementMs }
				}
			}
		});
		ws.send(seekGameMsg);

		this.uiState = 'seeking';
		this.poolLast.current = { clockMs, incrementMs };
	}

	cancelSeekGame(): void {
		if (!this.gameTimeControl) {
			return;
		}

		const cancelSeekGameMsg = create(MessageSchema, {
			event: { case: 'cancelSeekGame', value: {} }
		});
		ws.send(cancelSeekGameMsg);

		this.uiState = 'idle';
		this.gameTimeControl = null;
		this.poolLast.current = null;
	}

	gameAbort() {
		const abortGameMsg = create(MessageSchema, { event: { case: 'abortGame', value: {} } });
		ws.send(abortGameMsg);
	}

	gameResign() {
		const resignGameMsg = create(MessageSchema, { event: { case: 'resignGame', value: {} } });
		ws.send(resignGameMsg);
	}

	gameOfferDraw() {
		const offerDrawMsg = create(MessageSchema, { event: { case: 'offerDraw', value: {} } });
		ws.send(offerDrawMsg);
	}

	gameDeclineDraw() {
		const declineDrawMsg = create(MessageSchema, { event: { case: 'declineDraw', value: {} } });
		ws.send(declineDrawMsg);
		// this.opponenOfferedDraw = false;
	}

	gameAcceptDraw() {
		const acceptDrawMsg = create(MessageSchema, { event: { case: 'acceptDraw', value: {} } });
		ws.send(acceptDrawMsg);
		// this.opponenOfferedDraw = false;
	}

	playMoveUci(uci: string) {
		const playMoveUciMsg = create(MessageSchema, { event: { case: 'playMoveUci', value: { move: uci } } });
		ws.send(playMoveUciMsg);
		// if (this.ply <= 1) {
		// 	this.clock?.setCurrentTurn(this.clock.currentTurn === 'w' ? 'b' : 'w');
		// }
		// this.ply++;
		// this.updateTimersState();
	}

	echo() {
		const echoMsg = create(MessageSchema, { event: { case: 'echo', value: { message: 'hello bozo' } } });
		ws.send(echoMsg);
	}

	sendLobbyChat() {
		const lobbyChatMsg = create(MessageSchema, {
			event: { case: 'sendLobbyChat', value: { message: `hello lobby ${Math.floor(Math.random() * 100) + 1}` } }
		});
		ws.send(lobbyChatMsg);
	}

	onEchoMsg(echoMsg: Echo) {
		console.log('got echo: ', echoMsg.message);
	}

	onLatencyMsg(latencyMsg: Latency) {
		console.log('got latency_ms: ', latencyMsg.latencyMs);
	}

	userPresences = $state<Record<string, Presence>>({});
	channelPresences = $state<Record<string, SvelteSet<string>>>({});

	lobbyUserPresence = $derived.by(() => {
		if (this.channelPresences['lobby']?.size === 0) {
			return [];
		}

		return [...(this.channelPresences['lobby']?.values() ?? [])].reduce((acc, userId) => {
			const presence = this.userPresences[userId];
			if (presence) {
				acc.push(presence);
			}
			return acc;
		}, [] as Presence[]);
	});

	onPresenceState(presenceState: PresenceState) {
		console.log('snapshot', presenceState.presences);

		for (const presence of presenceState.presences) {
			this.userPresences[presence.userId] = presence;
			this.channelPresences[presence.channel] ||= new SvelteSet();
			this.channelPresences[presence.channel]?.add(presence.userId);
		}
	}

	onPresenceDiff(presenceDiff: PresenceDiff) {
		if (presenceDiff.joined.length > 0) {
			console.log('joined', presenceDiff.joined);
		} else {
			console.log('left', presenceDiff.left);
		}

		for (const presence of presenceDiff.joined) {
			this.channelPresences[presence.channel]?.add(presence.userId);
			this.userPresences[presence.userId] = presence;
		}

		for (const presence of presenceDiff.left) {
			this.channelPresences[presence.channel]?.delete(presence.userId);

			if (this.channelPresences[presence.channel]?.size === 0) {
				delete this.channelPresences[presence.channel];
			}

			let stillInAnyChannel = false;

			for (const ch of Object.keys(this.channelPresences)) {
				if (this.channelPresences[ch]?.has(presence.userId)) {
					stillInAnyChannel = true;
					break;
				}
			}

			if (!stillInAnyChannel) {
				delete this.userPresences[presence.userId];
			}
		}
	}

	onLobbyChat(lobbyChat: LobbyChat) {}

	onGameFound(gameFound: GameFound) {
		console.log('GameFound game_id: ', gameFound.gameId);
		goto(`/game/${gameFound.gameId}`);
	}

	handleWebsocketMessage(event: MessageEvent) {
		try {
			const msg = fromJsonString(MessageSchema, event.data);

			switch (msg.event.case) {
				case 'echo':
					this.onEchoMsg(msg.event.value);
					break;
				case 'latency':
					this.onLatencyMsg(msg.event.value);
					break;
				case 'presenceState':
					this.onPresenceState(msg.event.value);
					break;
				case 'presenceDiff':
					this.onPresenceDiff(msg.event.value);
					break;
				case 'lobbyChat':
					this.onLobbyChat(msg.event.value);
					break;
				case 'gameFound':
					this.onGameFound(msg.event.value);
					break;
				default:
					console.error('unknown message', msg.event.case, msg.event.value);
					break;
			}
		} catch (error) {
			console.error('json parse msg data', error);
		}
	}
}

export const gameManager = new GameManager();

// const CASTLE_MOVES_PAIR = new Map<string, string>([
// 	['e1g1', 'h1f1'],
// 	['e1c1', 'a1d1'],
// 	['e8g8', 'h8f8'],
// 	['e8c8', 'a8d8']
// ]);

// export const PROMOS = ['q', 'r', 'n', 'b'];

// export function getRookCastleMove(uci: string): string | undefined {
// 	return CASTLE_MOVES_PAIR.get(uci);
// }

// export function isPromotionMove(move: string, legalMoves: string[]): boolean {
// 	if (move.length !== 4) {
// 		return false;
// 	}
// 	return legalMoves.some(m => m.length === 5 && m.startsWith(move) && PROMOS.includes(m[4]!));
// }

// function isPromotionUciMove(uci: string): boolean {
// 	return uci.length === 5 && PROMOS.includes(uci[4]!);
// }

// export function isEnpassantLanMove(lan: string, currentTurn: 'w' | 'b'): Coord | undefined {
// 	if (!lan.includes('x')) {
// 		return;
// 	}
// 	const [src, dest] = lan.split('x');
// 	if (src?.length !== 2 || dest?.length !== 2) {
// 		return;
// 	}
// 	const [srcRank, destRank] = [src[1], dest[1]];
// 	if (currentTurn === 'w') {
// 		if (srcRank === '5' && destRank === '6') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	} else {
// 		if (srcRank === '4' && destRank === '3') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	}
// }

// export function playedEnpassantMove(move: string): Coord | undefined {
// 	if (move.length !== 4) {
// 		return;
// 	}
// 	const src = move.slice(0, 2) as Coord;
// 	const dest = move.slice(2, 4) as Coord;
// 	const pieceData = gameManager.board.getPiece(src);
// 	if (!(pieceData?.piece === 'p' || pieceData?.piece === 'P')) {
// 		return;
// 	}
// 	const [srcRank, destRank] = [src[1], dest[1]];
// 	if (gameManager.isWhiteTurn) {
// 		if (!(srcRank === '5' && destRank === '6')) {
// 			return;
// 		}
// 		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'p') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	} else {
// 		if (!(srcRank === '4' && destRank === '3')) {
// 			return;
// 		}
// 		if (gameManager.board.getPiece(`${dest[0]}${srcRank}` as Coord)?.piece === 'P') {
// 			return `${dest[0]}${srcRank}` as Coord;
// 		}
// 	}
// }
