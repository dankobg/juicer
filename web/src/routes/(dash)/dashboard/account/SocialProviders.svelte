<script lang="ts">
	import type { UiNode } from '@ory/client';
	import type { PageData } from './$types';
	import { Button, Card } from 'flowbite-svelte';

	export let data: PageData;

	const filterBy = (n: UiNode, action: 'link' | 'unlink') =>
		n.group === 'oidc' && n.type === 'input' && n.attributes.node_type === 'input' && n.attributes.name === action;

	let providersToLink = data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'link')) ?? [];
	let providersToUnlink = data?.flow?.ui?.nodes?.filter(n => filterBy(n, 'unlink')) ?? [];
</script>

{#if providersToLink.length > 0}
	<Card>
		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Link social auth providers</h3>

		{#each providersToLink as provider}
			{#if provider.attributes.node_type === 'input'}
				<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
					<input type="hidden" name="link" value={provider.attributes.value} readonly required />
					<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

					<Button color="alternative" class="w-full font-semibold">
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
		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Unlink social auth providers</h3>

		{#each providersToUnlink as provider}
			{#if provider.attributes.node_type === 'input'}
				<form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded" class="space-y-6">
					<input type="hidden" name="unlink" value={provider.attributes.value} readonly required />
					<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />

					<Button color="alternative" class="w-full font-semibold">
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
{/if}
