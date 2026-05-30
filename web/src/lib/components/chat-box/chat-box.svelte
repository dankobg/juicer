<script lang="ts">
	import { tick } from 'svelte';
	import { Button } from '$lib/components/ui/button/index';
	import { Input } from '$lib/components/ui/input';
	import IconArrowDown from '@lucide/svelte/icons/arrow-down';

	export type ChatMessage = {
		messageId: number;
		userId: string;
		message: string;
		postedAt: string;
	};

	export type ChatUser = {
		userId: string;
		username: string;
		avatarUrl?: string;
		guest?: boolean;
	};

	type Props = {
		title: string;
		channel: string;
		chatUserId: string;
		messages: ChatMessage[];
		users: Map<string, ChatUser>;
		onSend?: (text: string) => void;
	};

	let { title, channel, chatUserId, messages, users, onSend }: Props = $props();

	let messagesContainer: HTMLDivElement;
	let scrollPointElm: HTMLDivElement;
	let allowedToScrollToLatest = $state<boolean>(true);
	let text = $state<string>('');

	export function sendMessage() {
		if (!text) {
			return;
		}
		onSend?.(text);
		text = '';
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
		if (allowedToScrollToLatest && messages.length > 0 && messagesContainer) {
			scrollToLatestMessage();
		}
	});
</script>

<div class="relative flex h-full w-full flex-1 flex-col rounded-md border border-secondary">
	<div class="rounded-t-md bg-secondary p-1">
		<p class="text-center">{title}</p>
	</div>
	<div class="flex flex-1 flex-col overflow-hidden p-4">
		<div class="grid flex-1 content-start gap-1 overflow-y-auto" bind:this={messagesContainer} onscroll={onScroll}>
			{#each messages as msg (msg.messageId)}
				{@const chatUser = users.get(msg.userId)}

				<div class="flex flex-wrap items-center gap-2">
					<div class="flex items-center justify-center gap-1">
						<img
							class="aspect-square h-[22px] w-[22px] max-w-full object-cover"
							src={chatUser?.avatarUrl || '/images/empty-avatar.svg'}
							alt={chatUser?.username + ' avatar'}
						/>
						<span
							class={['rounded px-2 text-sm text-black', chatUserId === chatUser?.userId ? 'bg-primary' : 'bg-sky-700']}
						>
							{chatUser?.username}
						</span>
					</div>
					<p>{msg.message}</p>
				</div>
			{/each}

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
