<script lang="ts">
	import type { PageProps } from './$types';
	import { ws } from '$lib/ws/juicer-ws.svelte';
	import { onWsClose, onWsError, onWsMessage, onWsOpen } from '$lib/ws/ws-message-handler';
	import { gameManager } from '$lib/gameplay/game-manager.svelte';
	import {
		Game,
		getPromotionLabelText,
		getRookCastleMove,
		isPromotionMove,
		playedEnpassantMove,
		PROMOS
	} from '$lib/gameplay/game.svelte';
	import type {
		Coord,
		MoveCancelEvent,
		MoveFinishEvent,
		MoveStartEvent,
		ResizerScaleChangeEvent
	} from '@dankop/juicer-board';
	import { uiSettings } from '$lib/components/ui-settings/ui-settings-state.svelte';
	import { Color, GameState } from '$lib/gen/juicer_pb';
	import ChatBox, { type ChatMessage } from '$lib/components/chat-box/chat-box.svelte';
	import { presenceManager } from '$lib/gameplay/presence-manager.svelte';
	import { soundManager } from '$lib/sound/sound-manager.svelte';
	import GameControls from '$lib/components/game/game-controls.svelte';
	import PlayerInfo from '$lib/components/game/player-info.svelte';
	import MovesList from '$lib/components/game/moves-list.svelte';

	let { data, params }: PageProps = $props();

	let game = $derived<Game | undefined>(gameManager.games?.get(Number(params.game_id)));

	let promotionPopoverElm!: HTMLDivElement;

	let boardTheme = $derived.by(() => {
		if (!uiSettings.boardActiveTheme.current) {
			return;
		}
		return `/images/board/${uiSettings.boardActiveTheme.current.src}`;
	});

	let showResizer = $derived.by(() => {
		switch (uiSettings.resizer.current) {
			case 'disabled':
				return false;
			case 'always':
				return true;
			case 'first-move':
				return (game?.myColor === Color.WHITE && game.ply === 0) || (game?.myColor === Color.BLACK && game.ply <= 1);
			default:
				return false;
		}
	});

	function onBoardMoveStart(event: MoveStartEvent) {
		if (!game?.isViewingLatestPosition) {
			event.preventDefault();
			return;
		}

		const isWhitePiece = event.data.pieceData.piece === event.data.pieceData.piece.toUpperCase();
		if ((isWhitePiece && game?.myColor === Color.BLACK) || (!isWhitePiece && game?.myColor === Color.WHITE)) {
			event.preventDefault();
			console.debug('not your piece');
			return;
		}
		if (!game?.isMyTurnActive) {
			event.preventDefault();
			console.debug('not your turn');
			return;
		}
	}

	function onBoardMoveCancel(event: MoveCancelEvent) {
		if (event.data.reason === 'out-of-bound') {
			soundManager.play('OutOfBound');
		}
	}

	function onBoardMoveFinish(event: MoveFinishEvent) {
		const isWhitePiece = event.data.pieceData.piece === event.data.pieceData.piece.toUpperCase();
		if ((isWhitePiece && game?.myColor === Color.BLACK) || (!isWhitePiece && game?.myColor === Color.WHITE)) {
			console.debug('not your piece');
			event.preventDefault();
			soundManager.play('Error');
			return;
		}
		if (!game?.isMyTurnActive) {
			console.debug('not your turn');
			event.preventDefault();
			soundManager.play('Error');
			return;
		}
		const move = event.data.src + event.data.dest;
		const isPromo = isPromotionMove(move, game.legalMoves);
		if (isPromo) {
			game.promotionSrcDest = move;
			const dest = move.slice(2, 4);
			const promoSquareElm = game.board.shadowRoot?.querySelector(`juicer-square[coord='${dest}']`) ?? null;
			if (!promoSquareElm) {
				console.log('no promotion square element');
				return;
			}
			const rect = promoSquareElm.getBoundingClientRect();
			promotionPopoverElm.style.setProperty('--sq-size', `${rect.width}px`);
			promotionPopoverElm.style.left = `${rect.left}px`;
			promotionPopoverElm.style.top = `${rect.top}px`;
			promotionPopoverElm.showPopover();
			event.preventDefault();
			return;
		}
		if (!game.legalMoves.includes(move)) {
			console.debug('invalid move attempt:', move, game.legalMoves);
			soundManager.play('Error');
			event.preventDefault();
			return;
		}
		const rookMove = getRookCastleMove(move);
		if (rookMove) {
			const rookSrc = rookMove.slice(0, 2) as Coord;
			const rookDest = rookMove.slice(2, 4) as Coord;
			game.board.movePiece(rookSrc, rookDest);
		}
		const enpOppPieceCoordToDelete = playedEnpassantMove(game, move);
		if (enpOppPieceCoordToDelete) {
			game.board.removePiece(enpOppPieceCoordToDelete);
		}
		const isCapture = Boolean(game.board.getPiece(event.data.dest));
		soundManager.play(isCapture ? 'Capture' : 'Move');
		game.playMoveUci(move);
	}

	function onResizerScaleChanged(event: ResizerScaleChangeEvent) {
		event.stopPropagation();
		const clamped = Math.max(10, Math.min(event.data.scale, 100));
		document.documentElement.style.setProperty('--board-scale', `${clamped}`);
	}

	function addBoardEventListeners() {
		if (game?.board) {
			game.board.addEventListener('movestart', onBoardMoveStart);
			game.board.addEventListener('movecancel', onBoardMoveCancel);
			game.board.addEventListener('movefinish', onBoardMoveFinish);
			game.board.addEventListener('resizer:scale-changed', onResizerScaleChanged);
		}
	}

	function removeBoardEventListeners() {
		if (game?.board) {
			game.board.removeEventListener('movestart', onBoardMoveStart);
			game.board.removeEventListener('movecancel', onBoardMoveCancel);
			game.board.removeEventListener('movefinish', onBoardMoveFinish);
			game.board.removeEventListener('resizer:scale-changed', onResizerScaleChanged);
		}
	}

	$effect(() => {
		const params = new URLSearchParams();
		params.set('path', window.location.pathname);
		ws.connect(params);

		ws.onOpen = onWsOpen;
		ws.onError = onWsError;
		ws.onClose = onWsClose;
		ws.onMessage = onWsMessage;

		return () => {
			ws.close();
		};
	});

	$effect(() => {
		if (game?.board) {
			addBoardEventListeners();
		}
		return () => {
			removeBoardEventListeners();
		};
	});

	$effect(() => {
		if (uiSettings.pieceActiveTheme.current && game?.board) {
			game.board.pieceTheme = (piece: string) => {
				const clr = piece === piece.toUpperCase() ? 'w' : 'b';
				return `/images/piece/${uiSettings.pieceActiveTheme.current}/${clr}${piece.toUpperCase()}.svg`;
			};
		}
	});

	let gameUserPresences = $derived(presenceManager.getPresenceInChannel(`game.${game?.gameId}`));

	let gameChatMessages = $state<ChatMessage[]>([]);

	function clockPrecision(ms: number): 'deciseconds' | 'centiseconds' | null {
		if (ms <= 10_000) return 'centiseconds';
		if (ms <= 60_000) return 'deciseconds';
		return null;
	}
</script>

{#if game}
	<div class="game-layout">
		<div class="chat">
			<ChatBox
				title="Game chat"
				channel={`game.${game.gameId}.chat`}
				chatUserId={data?.auth?.user?.id ?? ''}
				messages={[]}
				users={new Map()}
				onSend={msg => {
					gameManager.sendGameChat(game.gameId!, msg);
				}}
			/>
		</div>

		<div class="moves">
			<MovesList {game} />
		</div>

		<div class="player opp">
			{#if game.opponentPlayer}
				<PlayerInfo
					player={game.opponentPlayer}
					color={game.opponentColor}
					active={game.gameState === GameState.ACTIVE && game.currentTurn === game.opponentColor}
					online={gameUserPresences.has(game.opponentPlayer.userId)}
					clockMs={game.opponentColor === Color.WHITE ? game.whiteDisplayTimeMs : game.blackDisplayTimeMs}
					{clockPrecision}
				/>
			{/if}
		</div>

		<div class="board-wrapper">
			<juicer-board
				bind:this={() => game.board, v => (game.board = v)}
				board-theme={boardTheme}
				fen={game?.currentPosition?.fen || 'start'}
				orientation={game.orientation === Color.WHITE ? 'w' : 'b'}
				coords-placement={uiSettings.boardCoordinates.current.placement}
				ranks-position={uiSettings.boardCoordinates.current.ranksPosition}
				files-position={uiSettings.boardCoordinates.current.filesPosition}
				interactive
				show-ghost={uiSettings.showGhost.current}
				show-resizer={showResizer}
				check-square={game?.checkSquare}
			></juicer-board>
		</div>

		<div class="controls">
			<GameControls {game} />
		</div>

		<div class="player me">
			{#if game.mePlayer}
				<PlayerInfo
					player={game.mePlayer}
					color={game.myColor}
					active={game.gameState === GameState.ACTIVE && game.currentTurn === game.myColor}
					online={gameUserPresences.has(game.mePlayer.userId)}
					clockMs={game.myColor === Color.WHITE ? game.whiteDisplayTimeMs : game.blackDisplayTimeMs}
					{clockPrecision}
				/>
			{/if}
		</div>
	</div>

	<div id="promotion-popover" popover="auto" bind:this={promotionPopoverElm}>
		<div>
			{#each PROMOS as promo (promo)}
				<button
					class="promotion-btn"
					aria-label={getPromotionLabelText(promo)}
					popovertarget="promotion-popover"
					popovertargetaction="hide"
					data-promo={promo}
					onclick={() => game?.handlePromotionPiecePick(promotionPopoverElm, promo)}
				></button>
			{/each}
		</div>
	</div>
{/if}

<style>
	.game-layout {
		container-name: game-layout;
		container-type: size;
		min-height: 0;
		min-width: 0;
		width: 100dvw;
		height: calc(100dvh - var(--header-height) - 1px);
		display: grid;
		gap: var(--game-layout-gap);
		justify-content: center;
		overflow: hidden;
		grid-template-rows: auto auto minmax(0, 1fr) auto auto;
		grid-template-areas:
			'moves'
			'player-opp'
			'board'
			'controls'
			'player-me';

		@media screen and (width > 60rem) {
			--game-layout-gap: 1rem;
			grid-template-rows: unset;
			grid-template-columns:
				minmax(var(--chat-min-width), var(--chat-max-width))
				auto
				minmax(var(--panel-min-width), var(--panel-max-width));
			grid-template-areas:
				'chat board player-opp'
				'chat board moves'
				'chat board controls'
				'chat board player-me';
		}
	}

	.chat {
		grid-area: chat;
		min-height: 20rem;
		display: none;

		@media screen and (width > 60rem) {
			display: block;
		}
	}

	.moves {
		grid-area: moves;
	}

	.controls {
		grid-area: controls;
	}

	.board-wrapper {
		grid-area: board;
		container-name: board-wrapper;
		container-type: size;
		display: grid;
		place-content: center;
		width: calc(100cqmin * (var(--board-scale) / 100));

		@media screen and (width > 60rem) {
			place-content: unset;
			width: calc(
				(min(100cqh, calc(100cqw - var(--chat-min-width) - var(--panel-min-width) - 2 * var(--game-layout-gap)))) *
					(var(--board-scale) / 100)
			);
		}
	}

	juicer-board {
		min-width: 0;
		min-height: 0;
		width: 100cqmin;
		height: 100cqmin;
	}

	.player.opp {
		grid-area: player-opp;
	}
	.player.me {
		grid-area: player-me;
	}

	#promotion-popover {
		margin: 0;
		padding: 0;
		border: none;
		z-index: 70;
		background-color: transparent;
	}
	#promotion-popover::backdrop {
		background: rgba(0, 0, 0, 0.5);
		position: fixed;
		inset: 0;
	}
	#promotion-popover > div {
		display: grid;
	}
	.promotion-btn {
		width: var(--sq-size);
		height: var(--sq-size);
		background-color: whitesmoke;
		box-shadow: none;
		border: 0;
		border-radius: 1000vmax;
		background-position: center;
		background-repeat: no-repeat;
		background-size: 75%;
		cursor: pointer;
	}
	.promotion-btn[data-promo='Q'] {
		background-image: var(--wq-theme);
	}
	.promotion-btn[data-promo='R'] {
		background-image: var(--wr-theme);
	}
	.promotion-btn[data-promo='N'] {
		background-image: var(--wn-theme);
	}
	.promotion-btn[data-promo='B'] {
		background-image: var(--wb-theme);
	}
	.promotion-btn[data-promo='q'] {
		background-image: var(--bq-theme);
	}
	.promotion-btn[data-promo='r'] {
		background-image: var(--br-theme);
	}
	.promotion-btn[data-promo='n'] {
		background-image: var(--bn-theme);
	}
	.promotion-btn[data-promo='b'] {
		background-image: var(--bb-theme);
	}
	.promotion-btn:hover {
		background-color: rgb(230, 230, 230);
	}
</style>
