<script lang="ts">
	import type { PageProps } from './$types';
	import * as Table from '$lib/components/ui/table/index';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';
	import IconEye from '@lucide/svelte/icons/eye';
	import { IdentityStateEnum, type Session } from '$lib/gen/juicer_openapi';
	import { stateIcons } from '../../identities/identities-data-grid/data';
	import Button from '$lib/components/ui/button/button.svelte';
	import { juicer } from '$lib/juicer/client';
	import { toast } from 'svelte-sonner';
	import { invalidate } from '$app/navigation';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let ActiveIcon = $derived(data.session?.active ? IconCheck : IconX);
	let activeIconClasses = $derived(data.session?.active ? 'text-green-400' : 'text-yellow-400');

	let IdentityStateIcon = $derived(data.session?.identity?.state && stateIcons.get(data.session?.identity.state));
	let identityStateIconClasses = $derived.by(() => {
		switch (data.session?.identity?.state as IdentityStateEnum) {
			case IdentityStateEnum.Active:
				return 'text-green-400';
			case IdentityStateEnum.Inactive:
				return 'text-red-400';
			default:
				return '';
		}
	});

	async function onConfirmDeactivateSession() {
		try {
			await juicer.disableSession({ id: data.session?.id ?? '' });
			toast.success('session deactivated');
			invalidate(`data:identity-sessions-${data.session?.identity?.id ?? ''}`);
		} catch (error) {
			console.log('err', error);
			toast.error('session deactivation failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	async function onConfirmExtendSession() {
		try {
			await juicer.extendSession({ id: data.session?.id ?? '' });
			toast.success('session extended');
			invalidate(`data:identity-sessions-${data.session?.identity?.id ?? ''}`);
		} catch (error) {
			console.log('err', error);
			toast.error('session extend failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeactivateSessionClick(sess: Session) {
		confirmation.openDialog({
			description: deactivateSessionDescriptionSnippet,
			onConfirm: onConfirmDeactivateSession
		});
	}

	function onExtendSessionClick(sess: Session) {
		confirmation.openDialog({
			description: extendSessionDescriptionSnippet,
			onConfirm: onConfirmExtendSession
		});
	}
</script>

{#snippet deactivateSessionDescriptionSnippet()}
	{@const email = data.session?.identity?.traits?.['email'] ?? ''}
	This action cannot be undone. This will deactive (invalidate) the session
	{#if email}
		for user: <strong>{email}</strong>
	{:else}
		<strong>{data?.session?.id ?? ''}</strong>
	{/if}
	so they will have to login again.
{/snippet}

{#snippet extendSessionDescriptionSnippet()}
	{@const email = data.session?.identity?.traits?.['email'] ?? ''}
	This will extend the session
	{#if email}
		for user: <strong>{email}</strong>
	{:else}
		<strong>{data?.session?.id ?? ''}</strong>
	{/if}
	so they will have to login again.
{/snippet}

{#if data.session}
	<h1 class="mb-6 text-2xl font-bold">Session</h1>

	<section class="mb-6 gap-4">
		<p class="mb-6 text-lg">Actions</p>
		<div class="flex gap-4">
			<Button
				variant="destructive"
				disabled={!data.session.active}
				onclick={() => onDeactivateSessionClick(data.session)}
			>
				Deactivate
			</Button>
			<Button disabled={!data.session.active} onclick={() => onExtendSessionClick(data.session)}>Extend</Button>
		</div>
	</section>

	<p class="mb-6 text-lg">Details</p>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.session.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Active</span>
			<span class="flex gap-2 font-medium">
				{data.session.active}
				<ActiveIcon class={activeIconClasses} />
			</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Authenticated time</span>
			<time class="font-medium">{fmt.format(data.session.authenticatedAt)}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Authenticator assurance level</span>
			<time class="font-medium">{data.session.authenticatorAssuranceLevel}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Expires time</span>
			<time class="font-medium">{fmt.format(data.session.expiresAt)}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Issued time</span>
			<time class="font-medium">{fmt.format(data.session.issuedAt)}</time>
		</div>
	</div>

	{#if data.session.authenticationMethods && data.session.authenticationMethods.length > 0}
		<p class="mt-8 text-lg">Authentication methods</p>
		<Table.Root>
			<Table.Caption>A list of authentication methods</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>Method</Table.Head>
					<Table.Head>Aal</Table.Head>
					<Table.Head>Completed time</Table.Head>
					<Table.Head>Organization</Table.Head>
					<Table.Head>Provider</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.session.authenticationMethods as method (method)}
					<Table.Row>
						<Table.Cell class="font-medium">{method.method}</Table.Cell>
						<Table.Cell>{method.aal}</Table.Cell>
						<Table.Cell>{fmt.format(method.completedAt)}</Table.Cell>
						<Table.Cell>{method.organization}</Table.Cell>
						<Table.Cell>{method.provider}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}

	{#if data.session.devices && data.session.devices.length > 0}
		<p class="mt-8 text-lg">Devices</p>
		<Table.Root>
			<Table.Caption>A list of devices</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Ip address</Table.Head>
					<Table.Head>Location</Table.Head>
					<Table.Head>User agent</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.session.devices as device (device)}
					<Table.Row>
						<Table.Cell class="font-medium">{device.id}</Table.Cell>
						<Table.Cell>{device.ipAddress}</Table.Cell>
						<Table.Cell>{device.location}</Table.Cell>
						<Table.Cell>{device.userAgent}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}

	{#if data.session.identity}
		<h2 class="mb-6 text-xl font-bold">Session identity</h2>
		<p class="mb-6 text-lg">Details</p>
		<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">ID</span>
				<div class="flex items-center gap-2 font-medium">
					<span>{data.session.identity.id}</span>
					<Button variant="ghost" size="icon" href="/dashboard/identities/{data.session.identity.id}">
						<IconEye class="text-sky-400" />
					</Button>
				</div>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">E-Mail</span>
				<span class="font-medium">{data.session.identity.traits['email'] ?? ''}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">First name</span>
				<span class="font-medium">{data.session.identity.traits['first_name'] ?? ''}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Last name</span>
				<span class="font-medium">{data.session.identity.traits['last_name'] ?? ''}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Avatar URL</span>
				<span class="font-medium">{data.session.identity.traits['avatar_url'] ?? ''}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Schema ID</span>
				<span class="font-medium">{data.session.identity.schemaId}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Schema URL</span>
				<span class="font-medium">{data.session.identity.schemaUrl}</span>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">State</span>
				<span class="flex gap-2 font-medium"
					>{data.session.identity.state} <IdentityStateIcon class={identityStateIconClasses} /></span
				>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">State changed time</span>
				<time class="font-medium">{fmt.format(data.session.identity.stateChangedAt ?? undefined)}</time>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Created time</span>
				<time class="font-medium">{fmt.format(data.session.identity.createdAt)}</time>
			</div>
			<div class="flex flex-col justify-center">
				<span class="text-muted-foreground">Updated time</span>
				<time class="font-medium">{fmt.format(data.session.identity.updatedAt)}</time>
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
				{#each Object.values(data.session.identity.credentials ?? {}) as credential}
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

		{#if data.session.identity.recoveryAddresses && data.session.identity.recoveryAddresses.length > 0}
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
					{#each data.session.identity.recoveryAddresses as recAddr (recAddr)}
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

		{#if data.session.identity.verifiableAddresses && data.session.identity.verifiableAddresses.length > 0}
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
					{#each data.session.identity.verifiableAddresses as verAddr (verAddr)}
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
{/if}
