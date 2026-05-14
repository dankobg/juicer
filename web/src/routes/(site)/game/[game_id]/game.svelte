<script lang="ts">
	import type { PageProps } from './$types';
	import { ws } from '$lib/ws/juicer-ws.svelte';
	import { onWsClose, onWsError, onWsMessage, onWsOpen } from '$lib/ws/ws-message-handler';
	import { Game, gameManager, PROMOS } from '$lib/gameplay/game-manager.svelte';
	import type { MoveCancelEvent, MoveFinishEvent, MoveStartEvent } from '@dankop/juicer-board';

	let { data, params }: PageProps = $props();

	let game = $derived<Game | undefined>(gameManager.games?.get(Number(params.game_id)));

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
</script>

picetina
