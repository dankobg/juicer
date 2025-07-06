<script lang="ts">
	import { tick } from 'svelte';
	import { Button } from '$lib/components/ui/button/index';
	import { Input } from '$lib/components/ui/input';
	import IconArrowDown from '@lucide/svelte/icons/arrow-down';
	import IconX from '@lucide/svelte/icons/x';
	import { gameManager } from '$lib/state/game-manager.svelte';

	type Props = {
		gameId: string;
		username: string;
		avatarUrl?: string;
		opponentUsername: string;
		opponentAvatarUrl?: string;
		onMessage?: (msg: string) => void;
	};
	let props: Props = $props();

	let messagesContainer: HTMLDivElement;
	let scrollPointElm: HTMLDivElement;
	let allowedToScrollToLatest: boolean = $state(true);
	let newMsg: string = $state('');

	export function sendMessage() {
		if (!newMsg) {
			return;
		}
		props.onMessage?.(newMsg);
		if (gameManager.chatGameId.current !== gameManager.gameId) {
			gameManager.chatMessages.current = [];
			gameManager.chatGameId.current = props.gameId;
		}
		gameManager.chatMessages.current.push({ text: newMsg });
		newMsg = '';
	}

	function onScroll(event: Event) {
		const elm = event.target as HTMLDivElement;
		const threshold = 50;
		if (elm.scrollTop + elm.clientHeight >= elm.scrollHeight - threshold) {
			allowedToScrollToLatest = true;
		} else {
			allowedToScrollToLatest = false;
		}
	}

	async function scrollToLatestMessage() {
		await tick();
		scrollPointElm.scrollIntoView({ behavior: 'smooth', block: 'end' });
	}

	$effect(() => {
		if (
			allowedToScrollToLatest &&
			gameManager.gameId === gameManager.chatGameId.current &&
			gameManager.chatMessages.current.length > 0 &&
			messagesContainer
		) {
			scrollToLatestMessage();
		}
	});
</script>

<dialog
	bind:this={gameManager.chatDialogElm}
	onclose={() => gameManager.closeChatDialog()}
	class="top-1/2 left-1/2 w-full translate-[-50%] rounded-sm border border-3 border-orange-700 sm:max-w-sm"
>
	<div class="h-[30rem]">
		<div class="border-secondary relative order-3 flex h-full w-full flex-1 flex-col rounded border">
			<div class="bg-secondary grid grid-cols-3 items-center justify-center rounded-t-md p-1">
				<p class="col-2 text-center">Game chat</p>
				<form method="dialog" class="col-3 flex">
					<Button variant="ghost" size="icon" type="submit" class="ml-auto">
						<IconX />
					</Button>
				</form>
			</div>
			<div class="flex flex-1 flex-col overflow-hidden p-4">
				<div
					class="grid flex-1 content-start gap-1 overflow-y-scroll"
					bind:this={messagesContainer}
					onscroll={onScroll}
				>
					{#if gameManager.chatGameId.current === gameManager.gameId}
						{#each gameManager.chatMessages.current as m}
							{@const username = m.received ? props.opponentUsername : props.username}
							{@const avatar = m.received ? props.opponentAvatarUrl : props.avatarUrl}
							<div class="flex flex-wrap items-center gap-2">
								<div class="flex items-center justify-center gap-1">
									<img
										class="aspect-square h-[22px] w-[22px] max-w-full object-cover"
										src={avatar || '/images/empty-avatar.svg'}
										alt={username + ' avatar'}
									/>
									<span class={['rounded px-2 text-sm', m.received ? 'bg-primary' : 'bg-sky-700']}>{username}</span>
								</div>
								<p>{m.text}</p>
							</div>
						{/each}
					{/if}

					{#if !allowedToScrollToLatest}
						<div class="absolute bottom-16 left-1/2 -translate-x-1/2 transform">
							<Button variant="default" onclick={scrollToLatestMessage} class="bg-primary">
								Jump to latest
								<IconArrowDown />
							</Button>
						</div>
					{/if}
					<div bind:this={scrollPointElm}></div>
				</div>
				<div class="mt-2 flex flex-shrink-0 gap-[0.3rem]">
					<Input
						type="text"
						bind:value={newMsg}
						onkeydown={e => e.key === 'Enter' && sendMessage()}
						placeholder="Send a message"
					/>
					<Button onclick={sendMessage}>Send</Button>
				</div>
			</div>
		</div>
	</div>
</dialog>
