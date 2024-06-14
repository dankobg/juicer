<script lang="ts">
	import { Input, Button } from 'flowbite-svelte';

	export let messages: string[] = [];
	export let msg: string = '';

	export function sendMessage() {
		if (!msg) {
			return;
		}
		messages = [...messages, msg];
		msg = '';
	}
</script>

<div class="chat">
	<div class="box">
		<div class="header">
			<p class="title">Game room chat</p>
		</div>
		<div class="messages">
			{#each messages as m}
				<div class="message">
					<div>
						<small>by @user69</small>
						<p class="text">{m}</p>
					</div>
				</div>
			{/each}
		</div>

		<div class="controls">
			<Input
				type="text"
				bind:value={msg}
				on:keydown={e => e.key === 'Enter' && sendMessage()}
				placeholder="Send a message"
			/>
			<Button on:click={sendMessage} disabled={!msg}>Send</Button>
		</div>
	</div>
</div>

<style>
	.chat {
		border: 1px solid purple;
		padding: 0.625rem;
		width: 30rem;
		height: 30rem;
		background-color: #fff;
		border-radius: 0.375rem;
	}

	.box {
		display: flex;
		flex-direction: column;
		height: 100%;
	}

	.header {
		margin-bottom: 1rem;
		background-color: rgb(248, 205, 180);
		text-align: center;
		border-radius: 0.5rem;
	}

	.messages {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: start;
		gap: 0.125rem;
		overflow-y: scroll;
	}

	.message {
		border-radius: 0.625rem;
		padding-inline: 0.5rem;
	}

	.message:nth-child(even) {
		background-color: #fff;
	}
	.message:nth-child(odd) {
		background-color: #eee;
	}

	.text {
		margin: 0;
	}

	.controls {
		display: flex;
		gap: 0.5rem;
		margin-top: 1rem;
	}
</style>
