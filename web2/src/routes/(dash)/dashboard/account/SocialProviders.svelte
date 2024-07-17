<script lang="ts">
	import type { UiNode } from '@ory/client-fetch';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import BellRing from 'lucide-svelte/icons/bell-ring';
	import Check from 'lucide-svelte/icons/check';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';

	export let data: PageData;
	let socialsAction: 'link' | 'unlink' | undefined;

	const filterBy = (n: UiNode, action: 'link' | 'unlink') =>
		n.group === 'oidc' && n.type === 'input' && n.attributes.node_type === 'input' && n.attributes.name === action;

	let providersToLink = data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'link')) ?? [];
	let providersToUnlink = data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'unlink')) ?? [];
	// let providers = data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'link') || filterBy(n, 'unlink')) ?? [];

	let providers = ['a', 'b', 'c', 'd', 'e', 'f'];

	onMount(() => {
		const val = window.sessionStorage.getItem('socialsAction') as 'link' | 'unlink' | undefined;

		if (val) {
			socialsAction = val;
			toast.success(`Your account has been ${val}ed`);
		}

		return () => {
			socialsAction = undefined;
			sessionStorage.removeItem('socialsAction');
		};
	});
</script>

<svelte:window
	on:beforeunload={() => {
		if (socialsAction) {
			sessionStorage.removeItem('socialsAction');
		}
	}}
/>

<Card.Root class="w-[380px]">
	<Card.Header>
		<Card.Title>Social providers</Card.Title>
		<Card.Description>Link/Unlink auth social providers</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-4">
		{#each providers as provider}
		{#if provider.attributes.node_type === 'input'}
		<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
			<input type="hidden" name="link" value="kurac" readonly required />
			<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
			
			<div class="items-center flex gap-4">
				<img class="ml-4 inline-flex h-6 w-6 object-cover" src="/images/providers/{alt={provider.attributes.value}}.svg" alt={provider.attributes.value} />
					Link {provider.attributes.value} account
					
					<Switch
						on:click={() => {
							window.sessionStorage.setItem('socialsAction', 'link');
						}}
					/>
				</div>
			</form>
			{/if}
		{/each}

		<!-- {#each providers as provider}
			{#if provider.attributes.node_type === 'input'}
				<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
					<input type="hidden" name="link" value={provider.attributes.value} readonly required />
					<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

					<Button
						type="submit"
						color="alternative"
						class="w-full font-semibold"
						on:click={() => {
							window.sessionStorage.setItem('socialsAction', 'link');
						}}
					>
						Link {provider.attributes.value} account
						<img
							class="ml-4 inline-flex h-6 w-6 object-cover"
							src="/images/providers/{provider.attributes.value}.svg"
							alt={provider.attributes.value}
						/>
					</Button>
				</form>
			{/if}
		{/each} -->
		<!-- <div class=" flex items-center space-x-4 rounded-md border p-4">
			<BellRing />
			<div class="flex-1 space-y-1">
				<p class="text-sm font-medium leading-none">Social providers</p>
				<p class="text-muted-foreground text-sm">Link/Unlink auth social providers</p>
			</div>
			<Switch />
		</div>
		<div>
			{#each providers as provider, idx (idx)}
				<div class="mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0">
					<span class="flex h-2 w-2 translate-y-1 rounded-full bg-sky-500" />
					<div class="space-y-1">
						{#if provider.attributes.node_type === 'input'}
							<p class="text-sm font-medium leading-none">
								{provider.attributes.value}
							</p>
							<p class="text-muted-foreground text-sm">
								{provider.meta.label}
							</p>
						{/if}
					</div>
				</div>
			{/each}
		</div> -->
	</Card.Content>
</Card.Root>

<!-- {#if providersToLink.length > 0}
	<Card>
		{#if socialsAction === 'link'}
			{#each data?.flow?.ui?.messages ?? [] as msg}
				{@const err = msg.type === 'error'}
				<SimpleAlert kind={msg.type} title={err ? 'Unable to link account' : ''} text={msg.text} />
			{/each}
		{/if}

		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Link social auth providers</h3>

		{#each providersToLink as provider}
			{#if provider.attributes.node_type === 'input'}
				<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
					<input type="hidden" name="link" value={provider.attributes.value} readonly required />
					<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

					<Button
						type="submit"
						color="alternative"
						class="w-full font-semibold"
						on:click={() => {
							window.sessionStorage.setItem('socialsAction', 'link');
						}}
					>
						Link {provider.attributes.value} account
						<img
							class="w-6 h-6 object-cover inline-flex ml-4"
							src="/images/providers/{provider.attributes.value}.svg"
							alt={provider.attributes.value}
						/>
					</Button>
				</form>
			{/if}
		{/each}
	</Card>
{/if}

{#if providersToUnlink.length > 0}
	<Card>
		{#if socialsAction === 'unlink'}
			{#each data?.flow?.ui?.messages ?? [] as msg}
				{@const err = msg.type === 'error'}
				<SimpleAlert kind={msg.type} title={err ? 'Unable to unlink account' : ''} text={msg.text} />
			{/each}
		{/if}

		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Unlink social auth providers</h3>

		{#each providersToUnlink as provider}
			{#if provider.attributes.node_type === 'input'}
				<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
					<input type="hidden" name="unlink" value={provider.attributes.value} readonly required />
					<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

					<Button
						type="submit"
						color="alternative"
						class="w-full font-semibold"
						on:click={() => {
							window.sessionStorage.setItem('socialsAction', 'unlink');
						}}
					>
						Unlink {provider.attributes.value} account
						<img
							class="w-6 h-6 object-cover inline-flex ml-4"
							src="/images/providers/{provider.attributes.value}.svg"
							alt={provider.attributes.value}
						/>
					</Button>
				</form>
			{/if}
		{/each}
	</Card>
{/if} -->
