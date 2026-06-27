<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { Button } from '$lib/components/ui/button/index';
	import { Input } from '$lib/components/ui/input';
	import IconArrowDown from '@lucide/svelte/icons/arrow-down';
	import type { Presence } from '$lib/gen/juicer_pb';
	import { colorFromUserId } from '$lib/utils';
	import { type ChatMessage } from '$lib/gameplay/chat-manager.svelte';

	type Props = {
		title: string;
		channel: string;
		chatUserId: string;
		messages: ChatMessage[];
		hasMore?: boolean;
		presences: Record<string, Presence>;
		onSend?: (text: string) => void;
		onLoadMore?: VoidFunction;
	};

	let { title, channel, chatUserId, messages, hasMore, presences, onSend, onLoadMore }: Props = $props();

	let text = $state<string>('');
	let messagesContainer: HTMLDivElement;
	let scrollTopSentinelElm: HTMLDivElement;
	let scrollBottomSentinelElm: HTMLDivElement;
	let canScrollToLatest = $state<boolean>(true);

	let loadingMore = false;
	let restoreScroll = false;
	let oldHeight = 0;

	export function sendMessage() {
		if (!text) {
			return;
		}
		onSend?.(text);
		text = '';
	}

	async function scrollToLatestMessage() {
		await tick();
		scrollBottomSentinelElm.scrollIntoView({ behavior: 'smooth', block: 'end' });
	}

	function onScroll(event: Event) {
		const elm = event.target as HTMLDivElement;
		canScrollToLatest = elm.scrollTop + elm.clientHeight >= elm.scrollHeight - 50;
	}

	async function loadMore() {
		if (loadingMore) {
			return;
		}
		loadingMore = true;
		oldHeight = messagesContainer.scrollHeight;
		restoreScroll = true;
		if (hasMore && onLoadMore) {
			onLoadMore();
		}
	}

	onMount(() => {
		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry?.isIntersecting) {
					loadMore();
				}
			},
			{
				root: messagesContainer,
				threshold: 0
			}
		);
		observer.observe(scrollTopSentinelElm);

		return () => observer.disconnect();
	});

	$effect(() => {
		if (canScrollToLatest && messages.length > 0 && messagesContainer) {
			scrollToLatestMessage();
		}
	});

	$effect(() => {
		if (!messages.length || !restoreScroll) {
			return;
		}
		tick().then(() => {
			const newHeight = messagesContainer.scrollHeight;
			messagesContainer.scrollTop += newHeight - oldHeight;
			restoreScroll = false;
			loadingMore = false;
		});
	});
</script>

<div class="relative flex h-full w-full flex-1 flex-col rounded-md border border-secondary">
	<div class="rounded-t-md bg-secondary p-1">
		<p class="text-center">{title}</p>
	</div>
	<div class="flex flex-1 flex-col overflow-hidden p-4">
		<div class="grid flex-1 content-start gap-1 overflow-y-auto" bind:this={messagesContainer} onscroll={onScroll}>
			<div bind:this={scrollTopSentinelElm}></div>
			{#each messages as msg (msg.messageId)}
				<div class="flex flex-wrap items-center gap-2">
					<div class="flex items-center justify-center gap-1">
						<img
							class="aspect-square h-[22px] w-[22px] max-w-full object-cover"
							src={msg?.user?.avatarUrl || '/images/empty-avatar.svg'}
							alt={msg?.user?.username + ' avatar'}
						/>
						<span
							style:--chat-bg={chatUserId === msg?.user?.id ? 'var(--primary)' : colorFromUserId(msg?.user?.id ?? '')}
							class="rounded px-2 text-sm bg-(--chat-bg) text-[contrast-color(var(--chat-bg))]"
						>
							{msg?.user?.username}
						</span>
					</div>
					<p>{msg.message}</p>
				</div>
			{/each}

			{#if !canScrollToLatest}
				<div class="absolute bottom-16 left-1/2 -translate-x-1/2 transform">
					<Button variant="default" onclick={scrollToLatestMessage} class="bg-primary">
						Jump to latest
						<IconArrowDown />
					</Button>
				</div>
			{/if}
			<div bind:this={scrollBottomSentinelElm}></div>
		</div>
		<div class="mt-2 flex shrink-0 gap-[0.3rem]">
			<Input
				type="text"
				bind:value={text}
				onkeydown={e => e.key === 'Enter' && sendMessage()}
				placeholder="Send a message"
			/>
			<Button onclick={sendMessage}>Send</Button>
		</div>
	</div>
</div>
