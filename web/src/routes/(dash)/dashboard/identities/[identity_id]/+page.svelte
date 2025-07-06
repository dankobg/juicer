<script lang="ts">
	import type { PageProps } from './$types';
	import * as Table from '$lib/components/ui/table/index';
	import { stateIcons } from '../identities-data-grid/data';
	import { IdentityStateEnum } from '$lib/gen/juicer_openapi';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import Button from '$lib/components/ui/button/button.svelte';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import { juicer } from '$lib/juicer/client';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let StateIcon = $derived(data.identity?.state && stateIcons.get(data.identity.state));
	let stateIconClasses = $derived.by(() => {
		switch (data.identity?.state as IdentityStateEnum) {
			case IdentityStateEnum.Active:
				return 'text-green-400';
			case IdentityStateEnum.Inactive:
				return 'text-red-400';
			default:
				return '';
		}
	});

	async function onConfirmDeleteIdentity() {
		if (!data.identity) {
			return;
		}
		try {
			await juicer.deleteIdentity({ id: data.identity.id });
			goto('/dashboard/identities');
		} catch (error) {
			console.log('err', error);
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeleteIdentityClick() {
		confirmation.openDialog({
			description: deleteIdentityDescriptionSnippet,
			onConfirm: onConfirmDeleteIdentity,
			destructive: true
		});
	}
</script>

{#snippet deleteIdentityDescriptionSnippet()}
	{@const email = data?.identity?.traits?.['email']}
	This action cannot be undone. This will delete the identity <strong>{email}</strong> completely.
{/snippet}

{#if data.identity}
	<section class="mb-6 gap-4">
		<p class="mb-6 text-lg">Actions</p>
		<div class="flex gap-4">
			<Button href="/dashboard/identities/{data.identity.id}/sessions">View sessions</Button>
			<Button href="/dashboard/identities/{data.identity.id}/edit">Edit identity</Button>
			<Button variant="destructive" onclick={onDeleteIdentityClick}>Delete identity</Button>
		</div>
	</section>

	<h1 class="mb-6 text-2xl font-bold">Identity</h1>
	<p class="mb-6 text-lg">Details</p>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.identity.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">E-Mail</span>
			<span class="font-medium">{data.identity.traits['email'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">First name</span>
			<span class="font-medium">{data.identity.traits['first_name'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Last name</span>
			<span class="font-medium">{data.identity.traits['last_name'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Avatar URL</span>
			<span class="font-medium">{data.identity.traits['avatar_url'] ?? ''}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Schema ID</span>
			<span class="font-medium">{data.identity.schemaId}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Schema URL</span>
			<span class="font-medium">{data.identity.schemaUrl}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">State</span>
			<span class="flex gap-2 font-medium">{data.identity.state} <StateIcon class={stateIconClasses} /></span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">State changed time</span>
			<time class="font-medium">{fmt.format(data.identity.stateChangedAt ?? undefined)}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(data.identity.createdAt)}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(data.identity.updatedAt)}</time>
		</div>
	</div>

	<p class="mt-8 text-lg">Credentials</p>
	<Table.Root>
		<Table.Caption>A list of credentials</Table.Caption>
		<Table.Header>
			<Table.Row>
				<Table.Head>Type</Table.Head>
				<Table.Head>Version</Table.Head>
				<Table.Head>Config</Table.Head>
				<Table.Head>Identifiers</Table.Head>
				<Table.Head>Created time</Table.Head>
				<Table.Head>Update time</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each Object.values(data.identity.credentials ?? {}) as credential}
				<Table.Row>
					<Table.Cell class="font-medium">{credential.type}</Table.Cell>
					<Table.Cell class="font-medium">{credential.version}</Table.Cell>
					<Table.Cell class="font-medium"><pre>{JSON.stringify(credential.config, null, 2)}</pre></Table.Cell>
					<Table.Cell class="font-medium">{credential.identifiers?.join(', ')}</Table.Cell>
					<Table.Cell>{fmt.format(credential.createdAt)}</Table.Cell>
					<Table.Cell>{fmt.format(credential.updatedAt)}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>

	{#if data.identity.recoveryAddresses && data.identity.recoveryAddresses.length > 0}
		<p class="mt-8 text-lg">Recovery addresses</p>
		<Table.Root>
			<Table.Caption>A list of recovery addresses</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Value</Table.Head>
					<Table.Head>Via</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.identity.recoveryAddresses as recAddr (recAddr)}
					<Table.Row>
						<Table.Cell class="font-medium">{recAddr.id}</Table.Cell>
						<Table.Cell>{recAddr.value}</Table.Cell>
						<Table.Cell>{recAddr.via}</Table.Cell>
						<Table.Cell>{fmt.format(recAddr.createdAt)}</Table.Cell>
						<Table.Cell>{fmt.format(recAddr.updatedAt)}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}

	{#if data.identity.verifiableAddresses && data.identity.verifiableAddresses.length > 0}
		<p class="mt-8 text-lg">Verifiable addresses:</p>
		<Table.Root>
			<Table.Caption>A list of verifiable addresses</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Value</Table.Head>
					<Table.Head>Via</Table.Head>
					<Table.Head>Status</Table.Head>
					<Table.Head>Verfiied</Table.Head>
					<Table.Head>Verified time</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.identity.verifiableAddresses as verAddr (verAddr)}
					<Table.Row>
						<Table.Cell class="font-medium">{verAddr.id}</Table.Cell>
						<Table.Cell>{verAddr.value}</Table.Cell>
						<Table.Cell>{verAddr.via}</Table.Cell>
						<Table.Cell>{verAddr.status}</Table.Cell>
						<Table.Cell>
							<div class="flex gap-2">
								{verAddr.verified}
								{#if verAddr.verified}
									<IconCheck class="text-green-400" />
								{:else}
									<IconX class="text-red-400" />
								{/if}
							</div>
						</Table.Cell>
						<Table.Cell>{fmt.format(verAddr.verifiedAt ?? undefined)}</Table.Cell>
						<Table.Cell>{fmt.format(verAddr.createdAt)}</Table.Cell>
						<Table.Cell>{fmt.format(verAddr.updatedAt)}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
{/if}
