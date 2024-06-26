<script lang="ts">
	import { createEventDispatcher, tick } from 'svelte';
	import { Input, Button, GradientButton } from 'flowbite-svelte';
	import { chatMessages } from './messages';

	const dispatch = createEventDispatcher();

	let messagesContainer: HTMLDivElement;
	let scrollPointElm: HTMLDivElement;
	let allowedToScrollToLatest = true;

	let newMsg: string = '';

	export function sendMessage() {
		if (!newMsg) {
			return;
		}
		dispatch('message', { text: newMsg });
		chatMessages.update(msgs => [...msgs, { own: true, text: newMsg }]);
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

	$: if (allowedToScrollToLatest && $chatMessages && messagesContainer) {
		scrollToLatestMessage();
	}
</script>

<div class="chat">
	<div class="header">
		<p>Game chat</p>
	</div>
	<div class="box">
		<div class="messages-container" bind:this={messagesContainer} on:scroll={onScroll}>
			{#each $chatMessages as m}
				<div class="message" class:own={m.own} class:opp={!m.own}>
					<div class="author">
						<img class="avatar" src="/images/logo.svg" alt="avatar" />
						<span class="name">bozo69</span>
					</div>
					<div class="content">
						<p class="text">{m.text}</p>
					</div>
				</div>
			{/each}

			{#if !allowedToScrollToLatest}
				<div class="scroll-btn">
					<GradientButton size="xs" outline pill color="pinkToOrange" on:click={scrollToLatestMessage}>
						Scroll down
					</GradientButton>
				</div>
			{/if}

			<div bind:this={scrollPointElm}></div>
		</div>
		<div class="controls">
			<Input
				type="text"
				bind:value={newMsg}
				on:keydown={e => e.key === 'Enter' && sendMessage()}
				placeholder="Send a message"
			/>
			<Button on:click={sendMessage}>Send</Button>
		</div>
	</div>
</div>

<style>
	.chat {
		width: 30rem;
		height: 30rem;
		display: flex;
		flex-direction: column;
		border: 1px solid #eb4f27;
		border-top: none;
		border-radius: 0.5rem;
		position: relative;
	}

	.header {
		background-color: #f3937a;
		text-align: center;
		padding: 0.2rem;
		border-top-left-radius: 0.5rem;
		border-top-right-radius: 0.5rem;
	}

	.box {
		flex: 1;
		padding: 1rem;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.messages-container {
		flex: 1;
		display: grid;
		align-content: start;
		gap: 4px;
		overflow-y: scroll;
	}

	.message {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	.message.own .author .name {
		background-color: rgb(196, 232, 255);
	}
	.message.opp .author .name {
		background-color: rgb(255, 196, 199);
	}

	.author {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 4px;
	}

	.avatar {
		aspect-ratio: 1;
		width: 22px;
		height: 22px;
		max-width: 100%;
		object-fit: cover;
	}

	.name {
		font-size: 13px;
		border-radius: 0.5rem;
		padding-inline: 0.5rem;
	}

	.content {
		border-radius: 0.5rem;
	}

	.text {
		font-weight: 500;
	}

	.controls {
		flex: 0;
		display: flex;
		gap: 0.3rem;
		margin-top: auto;
		margin-top: 1rem;
	}

	.scroll-btn {
		position: absolute;
		bottom: 4rem;
		left: 50%;
		transform: translateX(-50%);
	}
</style>
