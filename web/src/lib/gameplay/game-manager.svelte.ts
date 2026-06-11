import {
	Color,
	MessageSchema,
	type GameFound,
	type GameInfo,
	type MoveSync,
	type MoveAck,
	type GameChat,
	type GameChatList,
	type AbortGame,
	type ResignGame,
	type AcceptDraw,
	type DeclineDraw,
	type GameFinished,
	type DrawOffer,
	type DrawDeclined,
	GameResult,
	type PlayerRejoined,
	type PlayerLeft
} from '$lib/gen/juicer_pb';
import { create } from '@bufbuild/protobuf';
import type { Coord } from '@dankop/juicer-board';
import { ws } from '$lib/ws/juicer-ws.svelte';
import { SvelteMap } from 'svelte/reactivity';
import { goto } from '$app/navigation';
import type { ChatMessage } from '$lib/components/chat-box/chat-box.svelte';
import { soundManager } from '$lib/sound/sound-manager.svelte';
import { Game, getRookCastleMove, isEnpassantLanMove, isPromotionUciMove, type GameOptions } from './game.svelte';

export class GameManager {
	games = $state<SvelteMap<number, Game>>(new SvelteMap());
	gameChatMessages = $state<SvelteMap<number, ChatMessage[]>>(new SvelteMap());

	sendGameChat(gameId: number, message: string): void {
		const sendGameChatMsg = create(MessageSchema, {
			event: { case: 'sendGameChat', value: { gameId, message } }
		});
		ws.send(sendGameChatMsg);
	}

	onGameFound(gameFound: GameFound): void {
		const exists = this.games.has(gameFound.gameId);
		if (!exists) {
			this.games.set(gameFound.gameId, new Game({ gameId: gameFound.gameId }));
		}
		goto(`/game/${gameFound.gameId}`);
		soundManager.play('NewChallenge');
	}

	onGameInfo(gameInfo: GameInfo): void {
		const gameOpts: GameOptions = {
			gameId: gameInfo.gameId,
			white: gameInfo.white,
			black: gameInfo.black,
			gameVariant: gameInfo.gameVariant,
			gameTimeKind: gameInfo.gameTimeKind,
			gameTimeCategory: gameInfo.gameTimeCategory,
			gameTimeControl: gameInfo.gameTimeControl,
			gameState: gameInfo.gameState,
			gameResult: gameInfo.gameResult,
			gameResultStatus: gameInfo.gameResultStatus,
			reconnectTimeoutMs: gameInfo.reconnectTimeoutMs,
			firstMoveTimeoutMs: gameInfo.firstMoveTimeoutMs,
			lastMove: gameInfo.lastMove,
			startTime: gameInfo.startTime,
			endTime: gameInfo.endTime,
			rated: gameInfo.rated,
			version: gameInfo.version,
			ply: gameInfo.ply,
			gameMoves: gameInfo.gameMoves,
			myColor: gameInfo.color,
			orientation: gameInfo.color,
			legalMoves: gameInfo.legalMoves,
			whiteRemainingGameTime: gameInfo.clocks?.white,
			blackRemainingGameTime: gameInfo.clocks?.black,
			ack: gameInfo.version,
			pendingDrawOffers: gameInfo.pendingDrawOffers,
			whiteDisconnectedAt: gameInfo.whiteDisconnectedAt,
			blackDisconnectedAt: gameInfo.blackDisconnectedAt
		};

		if (gameInfo.gameMoves.length > 0) {
			gameOpts.historyPointer = gameInfo.gameMoves.length - 1;
		}

		if (this.games.has(gameInfo.gameId)) {
			const game = this.games.get(gameInfo.gameId);
			game?.configure(gameOpts);
			game?.startLoop();
		} else {
			const game = new Game(gameOpts);
			game.startLoop();
			this.games.set(gameInfo.gameId, game);
		}
	}

	onAbortGame(abortGame: AbortGame): void {}

	onResignGame(resignGame: ResignGame): void {}

	onDrawOffer(drawOffer: DrawOffer): void {
		const game = this.games.get(drawOffer.gameId);
		if (game) {
			game.pendingDrawOffers.set(drawOffer.offeredBy, drawOffer);
		}
	}

	onDrawDeclined(drawDeclined: DrawDeclined): void {
		const game = this.games.get(drawDeclined.gameId);
		if (game) {
			game.pendingDrawOffers.clear();
		}
	}

	onAcceptDraw(acceptDraw: AcceptDraw): void {}

	onDeclinedDraw(declineDraw: DeclineDraw): void {}

	onGamePlayerLeft(playerLeft: PlayerLeft): void {
		const game = this.games.get(playerLeft.gameId);
		if (game) {
			const color = game?.getPlayerColor(playerLeft.userId);
			if (color === Color.WHITE) {
				game.whiteDisconnectedAt = playerLeft.leftAt;
			} else if (color === Color.BLACK) {
				game.blackDisconnectedAt = playerLeft.leftAt;
			}
		}
	}

	onGamePlayerRejoined(playerRejoined: PlayerRejoined): void {
		const game = this.games.get(playerRejoined.gameId);
		if (game) {
			const color = game?.getPlayerColor(playerRejoined.userId);
			if (color === Color.WHITE) {
				game.whiteDisconnectedAt = undefined;
			} else if (color === Color.BLACK) {
				game.blackDisconnectedAt = undefined;
			}
		}
	}

	onGameFinished(gameFinished: GameFinished): void {
		const game = this.games.get(gameFinished.gameId);
		if (game) {
			game.gameResult = gameFinished.gameResult;
			game.gameResultStatus = gameFinished.gameResultStatus;
			game.gameState = gameFinished.gameState;
			game.stopLoop();

			switch (game?.gameResult) {
				case GameResult.DRAW:
					soundManager.play('Draw');
					break;
				case GameResult.WHITE_WON:
					soundManager.play(game?.myColor === Color.WHITE ? 'Victory' : 'Defeat');
					break;
				case GameResult.BLACK_WON:
					soundManager.play(game?.myColor === Color.BLACK ? 'Victory' : 'Defeat');
					break;
			}
		}
	}

	onMoveSync(moveSync: MoveSync): void {
		const game = this.games.get(moveSync.gameId);
		if (!game) {
			throw new Error('movesync: no game found');
		}

		const wasViewingLatestPosition = game.historyPointer === game.ply;

		const myMove = moveSync.ply === game.ply;
		game.version = moveSync.version;
		game.legalMoves = moveSync.legalMoves;
		game.ply = moveSync.ply;
		game.whiteRemainingGameTime = moveSync.clocks?.white;
		game.blackRemainingGameTime = moveSync.clocks?.black;
		game.lastMove = moveSync.playedAt;
		game.gameMoves.push({
			$typeName: 'pb.GameMove',
			uci: moveSync.uci,
			san: moveSync.san,
			lan: moveSync.lan,
			fen: moveSync.fen,
			playedAt: moveSync.playedAt
		});

		if (wasViewingLatestPosition) {
			game.historyPointer = game.gameMoves.length - 1;
		}

		game.startLoop();

		if (!myMove && wasViewingLatestPosition) {
			const src = moveSync.uci.slice(0, 2) as Coord;
			const dest = moveSync.uci.slice(2, 4) as Coord;
			const isPromo = isPromotionUciMove(moveSync.uci);
			if (isPromo) {
				game.promotionSrcDest = moveSync.uci.slice(0, 4);
				game.promotionPieceSymbol = game.myColor === Color.WHITE ? moveSync.uci[4]! : moveSync.uci[4]!.toUpperCase();
				game.promotePiece(game.promotionPieceSymbol);
				return;
			}
			const rookMove = getRookCastleMove(moveSync.uci);
			if (rookMove) {
				const rookSrc = rookMove.slice(0, 2) as Coord;
				const rookDest = rookMove.slice(2, 4) as Coord;
				const pos = new Map(game.board.position);
				const pieceData = game.board.getPiece(src)!;
				const rookPieceData = game.board.getPiece(rookSrc)!;
				pos.delete(src);
				pos.set(dest, pieceData);
				pos.delete(rookSrc);
				pos.set(rookDest, rookPieceData);
				game.board.setPosition(pos);
				return;
			}
			const enpOppPieceCoordToDelete = isEnpassantLanMove(moveSync.lan, game.opponentColor === Color.WHITE ? 'w' : 'b');
			if (enpOppPieceCoordToDelete) {
				const pos = new Map(game.board.position);
				const pieceData = game.board.getPiece(src)!;
				pos.delete(enpOppPieceCoordToDelete);
				pos.delete(src);
				pos.set(dest, pieceData);
				game.board.setPosition(pos);
				return;
			}
			game.board.movePiece(src, dest);
		}
	}

	onMoveAck(moveAck: MoveAck): void {}

	onGameChat(gameChat: GameChat): void {
		const gameChats = this.gameChatMessages.get(gameChat.gameId);
		if (!gameChats) {
			this.gameChatMessages.set(gameChat.gameId, []);
		}

		gameChats!.push({
			userId: gameChat.userId,
			messageId: gameChat.messageId,
			message: gameChat.message,
			postedAt: gameChat.postedAt
		});
	}

	onGameChatList(gameChats: GameChatList): void {}
}

export const gameManager = new GameManager();
