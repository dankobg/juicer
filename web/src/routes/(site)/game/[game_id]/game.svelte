<script lang="ts">
	import type { PageProps } from './$types';
	import { ws } from '$lib/ws/juicer-ws.svelte';
	import { onWsClose, onWsError, onWsMessage, onWsOpen } from '$lib/ws/ws-message-handler';
	import { Game, gameManager, getPromotionLabelText, PROMOS } from '$lib/gameplay/game-manager.svelte';
	import type { MoveCancelEvent, MoveFinishEvent, MoveStartEvent } from '@dankop/juicer-board';
	import { uiSettings } from '$lib/components/ui-settings/ui-settings-state.svelte';
	import { Color } from '$lib/gen/juicer_pb';
	import ChatBox, { type ChatMessage } from '$lib/components/chat-box/chat-box.svelte';
	import { presenceManager } from '$lib/gameplay/presence-manager.svelte';

	let { data, params }: PageProps = $props();

	let promotionPopoverElm!: HTMLDivElement;

	let game = $derived<Game | undefined>(gameManager.games?.get(Number(params.game_id)));

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

	let arr = $derived.by(() => {
		return game?.gameMoves.map(({ uci, playedAt }) => {
			return {
				m: uci || '-',
				ts: playedAt ? Number(playedAt.seconds || 0) : '-'
			};
		});
	});

	function onBoardMoveStart(event: MoveStartEvent) {
		console.log('onBoardMoveStart', event.data);
		// const isWhitePiece = event.data.pieceData.piece === event.data.pieceData.piece.toUpperCase();
		// if ((isWhitePiece && gameManager.color === Color.BLACK) || (!isWhitePiece && gameManager.color === Color.WHITE)) {
		// 	event.preventDefault();
		// 	console.debug('not your piece');
		// 	return;
		// }
		// if (!gameManager.hasActiveTurn) {
		// 	event.preventDefault();
		// 	console.debug('not your turn');
		// 	return;
		// }
	}

	function onBoardMoveCancel(event: MoveCancelEvent) {
		console.log('onBoardMoveCancel', event.data);
		// if (event.data.reason === 'out-of-bound') {
		// 	soundManager.play('OutOfBound');
		// }
	}

	function onBoardMoveFinish(event: MoveFinishEvent) {
		console.log('onBoardMoveFinish', event.data);
		// const isWhitePiece = event.data.pieceData.piece === event.data.pieceData.piece.toUpperCase();
		// if ((isWhitePiece && gameManager.color === Color.BLACK) || (!isWhitePiece && gameManager.color === Color.WHITE)) {
		// 	console.debug('not your piece');
		// 	event.preventDefault();
		// 	soundManager.play('Error');
		// 	return;
		// }
		// if (!gameManager.hasActiveTurn) {
		// 	console.debug('not your turn');
		// 	event.preventDefault();
		// 	soundManager.play('Error');
		// 	return;
		// }
		// const move = event.data.src + event.data.dest;
		// const isPromo = isPromotionMove(move, gameManager.legalMoves);
		// if (isPromo) {
		// 	gameManager.promotionSrcDest = move;
		// 	const dest = move.slice(2, 4);
		// 	const promoSquareElm = gameManager.board.shadowRoot?.querySelector(`juicer-square[coord='${dest}']`) ?? null;
		// 	if (!promoSquareElm) {
		// 		console.log('no promotion square element');
		// 		return;
		// 	}
		// 	const rect = promoSquareElm.getBoundingClientRect();
		// 	promotionPopoverElm.style.setProperty('--sq-size', `${rect.width}px`);
		// 	promotionPopoverElm.style.left = `${rect.left}px`;
		// 	promotionPopoverElm.style.top = `${rect.top}px`;
		// 	promotionPopoverElm.showPopover();
		// 	event.preventDefault();
		// 	return;
		// }
		// if (!gameManager.legalMoves.includes(move)) {
		// 	console.debug('invalid move attempt:', move, gameManager.legalMoves);
		// 	soundManager.play('Error');
		// 	event.preventDefault();
		// 	return;
		// }
		// const rookMove = getRookCastleMove(move);
		// if (rookMove) {
		// 	const rookSrc = rookMove.slice(0, 2) as Coord;
		// 	const rookDest = rookMove.slice(2, 4) as Coord;
		// 	gameManager.board.movePiece(rookSrc, rookDest);
		// }
		// const enpOppPieceCoordToDelete = playedEnpassantMove(move);
		// if (enpOppPieceCoordToDelete) {
		// 	gameManager.board.removePiece(enpOppPieceCoordToDelete);
		// }
		// const isCapture = Boolean(gameManager.board.getPiece(event.data.dest));
		// soundManager.play(isCapture ? 'Capture' : 'Move');
		// gameManager.playMoveUci(move);
	}

	function addBoardEventListeners() {
		if (game?.board) {
			game.board.addEventListener('movestart', onBoardMoveStart);
			game.board.addEventListener('movecancel', onBoardMoveCancel);
			game.board.addEventListener('movefinish', onBoardMoveFinish);
		}
	}

	function removeBoardEventListeners() {
		if (game?.board) {
			game.board.removeEventListener('movestart', onBoardMoveStart);
			game.board.removeEventListener('movecancel', onBoardMoveCancel);
			game.board.removeEventListener('movefinish', onBoardMoveFinish);
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
</script>

<pre class="l-20 t-20 fixed">{JSON.stringify(gameUserPresences?.values().toArray(), null, 2)}</pre>
<pre class="l-[1000px] t-20 fixed">{JSON.stringify(gameChatMessages, null, 2)}</pre>

{#if game?.gameId}
	<div class="game-layout">
		<div class="chat">
			<ChatBox
				title="Game chat"
				channel={`game.${game.gameId}`}
				chatUserId={game.currentPlayerId ?? ''}
				messages={[]}
				users={gameUserPresences}
				onMessage={msg => {
					gameChatMessages.push({ messageId: 1, message: msg, userId: game?.currentPlayerId ?? '', postedAt: '' });
					console.log('msg to send: ', msg);
				}}
			/>
		</div>
		<!-- <aside class="chat">chat</aside> -->
		<juicer-board
			board-theme={boardTheme}
			bind:this={game.board}
			class="order-1 flex-1 md:order-0"
			orientation={game.orientation}
			fen={game.gameMoves[0]?.fen || ''}
			coords-placement="inside"
			files-position="bottom"
			ranks-position="left"
			interactive
			show-ghost
			show-resizer={showResizer}
		></juicer-board>
		<aside class="controls">controls</aside>
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
				></button>
				<!-- onclick={() => game?.handlePromotionPiecePick(promotionPopoverElm, promo)} -->
			{/each}
		</div>
	</div>
{/if}

<style>
	.game-layout {
		position: relative;
		container-type: size;
		container-name: board-layout;
		width: 100dvw;
		height: calc(100dvh - var(--header-height) - 1px);
		display: grid;
		grid-template-columns: minmax(0, 300px) auto minmax(0, 300px);
		justify-content: center;
	}

	juicer-board {
		width: calc((min(100cqh, calc(100cqw - 300px - 300px))) * (var(--board-scale) / 100));
	}

	.chat,
	.controls {
		outline: 3px solid darkred;
	}

	/*##############################################################*/

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
