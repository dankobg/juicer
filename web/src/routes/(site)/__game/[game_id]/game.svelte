<script lang="ts">
	import { innerHeight } from 'svelte/reactivity/window';
	import { Color, MessageSchema } from '$lib/gen/juicer_pb';
	import {
		gameManager,
		getRookCastleMove,
		isPromotionMove,
		playedEnpassantMove,
		PROMOS
	} from '$lib/state/game-manager.svelte';
	import ChatBox from '$lib/components/chat-box/chat-box.svelte';
	import { create } from '@bufbuild/protobuf';
	import type { Coord, MoveCancelEvent, MoveFinishEvent, MoveStartEvent } from '@dankop/juicer-board';
	import type { User } from '$lib/kratos/service';
	import GamePanel from '$lib/components/game-panel/game-panel.svelte';
	import { ws } from '$lib/state/ws-state.svelte';
	import ChatBoxDialog from '$lib/components/chat-box/chat-box-dialog.svelte';
	import PlayerOpponentMobile from '$lib/components/game-panel/player-opponent-mobile.svelte';
	import PlayerMeMobile from '$lib/components/game-panel/player-me-mobile.svelte';
	import ActionButtons from '$lib/components/game-panel/action-buttons.svelte';
	import MovesListMobile from '$lib/components/game-panel/moves-list-mobile.svelte';
	import { settings } from '$lib/components/settings/settings-state.svelte';
	import { soundManager } from '$lib/state/sound-manager.svelte';

	let { authenticated, user }: { authenticated: boolean; user: User | null } = $props();

	let promotionPopoverElm!: HTMLDivElement;
	let resizeTargetElm!: HTMLDivElement;

	function onChatMessage(msg: string) {
		if (!msg) {
			return;
		}
		const chatMsg = create(MessageSchema, { event: { case: 'gameChat', value: { message: msg } } });
		ws.send(chatMsg);
	}

	function onBoardMoveStart(event: MoveStartEvent) {
		const isWhitePiece = event.data.pieceData.piece === event.data.pieceData.piece.toUpperCase();
		if ((isWhitePiece && gameManager.color === Color.BLACK) || (!isWhitePiece && gameManager.color === Color.WHITE)) {
			event.preventDefault();
			console.debug('not your piece');
			return;
		}
		if (!gameManager.hasActiveTurn) {
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
		if ((isWhitePiece && gameManager.color === Color.BLACK) || (!isWhitePiece && gameManager.color === Color.WHITE)) {
			console.debug('not your piece');
			event.preventDefault();
			soundManager.play('Error');
			return;
		}
		if (!gameManager.hasActiveTurn) {
			console.debug('not your turn');
			event.preventDefault();
			soundManager.play('Error');
			return;
		}

		const move = event.data.src + event.data.dest;

		const isPromo = isPromotionMove(move, gameManager.legalMoves);
		if (isPromo) {
			gameManager.promotionSrcDest = move;
			const dest = move.slice(2, 4);
			const promoSquareElm = gameManager.board.shadowRoot?.querySelector(`juicer-square[coord='${dest}']`) ?? null;
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
		if (!gameManager.legalMoves.includes(move)) {
			console.debug('invalid move attempt:', move, gameManager.legalMoves);
			soundManager.play('Error');
			event.preventDefault();
			return;
		}
		const rookMove = getRookCastleMove(move);
		if (rookMove) {
			const rookSrc = rookMove.slice(0, 2) as Coord;
			const rookDest = rookMove.slice(2, 4) as Coord;
			gameManager.board.movePiece(rookSrc, rookDest);
		}
		const enpOppPieceCoordToDelete = playedEnpassantMove(move);
		if (enpOppPieceCoordToDelete) {
			gameManager.board.removePiece(enpOppPieceCoordToDelete);
		}
		const isCapture = Boolean(gameManager.board.getPiece(event.data.dest));
		soundManager.play(isCapture ? 'Capture' : 'Move');
		gameManager.playMoveUci(move);
	}

	function addBoardEventListeners() {
		if (gameManager.board) {
			gameManager.board.addEventListener('movestart', onBoardMoveStart);
			gameManager.board.addEventListener('movecancel', onBoardMoveCancel);
			gameManager.board.addEventListener('movefinish', onBoardMoveFinish);
		}
	}

	function removeBoardEventListeners() {
		if (gameManager.board) {
			gameManager.board.removeEventListener('movestart', onBoardMoveStart);
			gameManager.board.removeEventListener('movecancel', onBoardMoveCancel);
			gameManager.board.removeEventListener('movefinish', onBoardMoveFinish);
		}
	}

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
	});

	$effect(() => {
		if (gameManager.board) {
			addBoardEventListeners();
		}
		return () => {
			removeBoardEventListeners();
		};
	});

	$effect(() => {
		if (gameManager.board && resizeTargetElm) {
			gameManager.board.resizeTarget = resizeTargetElm;
		}
	});

	function getPromotionLabelText(promotionSymbol: string): string {
		if (promotionSymbol === 'q') {
			return 'promote to queen';
		}
		if (promotionSymbol === 'r') {
			return 'promote to rook';
		}
		if (promotionSymbol === 'n') {
			return 'promote to knight';
		}
		if (promotionSymbol === 'b') {
			return 'promote to bishop';
		}
		return '';
	}

	let showResizer = $derived.by(() => {
		switch (settings.resizer.current) {
			case 'disabled':
				return false;
			case 'always':
				return true;
			case 'first-move':
				return (
					(gameManager.color === Color.WHITE && gameManager.ply === 0) ||
					(gameManager.color === Color.BLACK && gameManager.ply <= 1)
				);
			default:
				return false;
		}
	});

	$effect(() => {
		if (settings.pieceActiveTheme.current && gameManager.board) {
			gameManager.board.pieceTheme = (piece: string) => {
				const clr = piece === piece.toUpperCase() ? 'w' : 'b';
				return `/images/piece/${settings.pieceActiveTheme.current}/${clr}${piece.toUpperCase()}.svg`;
			};
		}
	});
	let boardTheme = $derived.by(() => {
		if (!settings.boardActiveTheme.current) {
			return;
		}
		return `/images/board/${settings.boardActiveTheme.current.src}`;
	});

	let arr = $derived.by(() => {
		return gameManager.historyMovesInfo.map(({ move, playedAt }) => {
			return {
				m: move?.uci || '-',
				ts: playedAt ? Number(playedAt.seconds || 0) : '-'
			};
		});
	});
</script>

<main class="max-w-[1920px]">
	<section class="justify-betwee flex h-[calc(100dvh-4rem)] flex-col items-center">
		<PlayerOpponentMobile username={authenticated ? gameManager.opponentInfo?.username || '' : 'Opp (Guest)'} />
		<div
			class="[container-type:size] flex h-full w-full flex-1 items-center justify-center gap-2 [container-name:wrapper-container]"
		>
			{@render gameChat({
				gameId: gameManager.gameId,
				username: authenticated ? (user?.username ?? '') : 'Me (Guest)',
				avatarUrl: user?.avatarUrl,
				opponentUsername: authenticated ? (gameManager.opponentInfo?.username ?? '') : 'Opp (Guest)',
				opponentAvatarUrl: gameManager.opponentInfo?.avatarUrl,
				onMessage: onChatMessage
			})}
			<div class="board-inner aspect-square w-full place-content-center overflow-hidden" bind:this={resizeTargetElm}>
				<juicer-board
					board-theme={boardTheme}
					bind:this={gameManager.board}
					class="order-1 flex-1 md:order-none"
					orientation={gameManager.orientation}
					fen={gameManager.currentPreview?.fen || undefined}
					coords-placement="inside"
					files-position="bottom"
					ranks-position="left"
					interactive={true}
					show-ghost={true}
					show-resizer={showResizer || undefined}
					min-size="100"
					max-size={innerHeight?.current ? innerHeight.current - 80 + 16 : undefined}
					check-square={gameManager.checkSquare}
				></juicer-board>
			</div>
			<GamePanel {authenticated} {user} />
		</div>
		<div class="sm:hidden">
			<ActionButtons />
			<MovesListMobile />
		</div>
		<PlayerMeMobile username={authenticated ? (user?.username ?? '') : 'Me (Guest)'} avatar={user?.avatarUrl} />
	</section>
</main>

<div id="promotion-popover" popover="auto" bind:this={promotionPopoverElm}>
	<div>
		{#each PROMOS as promo}
			<button
				class="promotion-btn"
				aria-label={getPromotionLabelText(promo)}
				popovertarget="promotion-popover"
				popovertargetaction="hide"
				data-promo={promo}
				onclick={() => gameManager.handlePromotionPiecePick(promotionPopoverElm, promo)}
			></button>
		{/each}
	</div>
</div>

{#snippet gameChat(props: {
	gameId: string;
	username: string;
	avatarUrl?: string;
	opponentUsername: string;
	opponentAvatarUrl?: string;
	onMessage?: (msg: string) => void;
})}
	<ChatBox {...props} />
	<ChatBoxDialog {...props} />
{/snippet}

<style>
	@container wrapper-container (aspect-ratio > 1) {
		.board-inner {
			width: auto;
			height: 100%;
		}
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
