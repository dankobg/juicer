<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconX from '@lucide/svelte/icons/x';
	import type { Table } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import DataTableFacetedFilter from '$lib/components/data-grid-shared/data-table-faceted-filter.svelte';
	import DataTableViewOptions from '$lib/components/data-grid-shared/data-table-view-options.svelte';
	import { aals } from './data';
	import { confirmation } from '$lib/components/confirmation-dialog/confirmation-dialog-state.svelte';
	import type { Identity } from '$lib/gen/juicer_openapi';
	import { juicer } from '$lib/juicer/client';
	import { toast } from 'svelte-sonner';
	import { invalidate } from '$app/navigation';

	let { table, identity }: { table: Table<TData>; identity: Identity | undefined } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	const aalCol = $derived(table.getColumn('authenticatorAssuranceLevel'));

	async function onConfirmDeleteIdentitySessions() {
		if (!identity) {
			return;
		}
		try {
			await juicer.deleteIdentitySessions({ id: identity.id });
			toast.success('sessions deleted');
			invalidate(`data:identity-sessions-${identity.id}`);
		} catch (error) {
			console.log('err', error);
			toast.error('sessions delete failed');
		} finally {
			confirmation.closeDialog();
		}
	}

	function onDeleteIdentitySessionsClick() {
		confirmation.openDialog({
			description: deleteIdentitySessionsDescriptionSnippet,
			onConfirm: onConfirmDeleteIdentitySessions
		});
	}
</script>

{#snippet deleteIdentitySessionsDescriptionSnippet()}
	This action cannot be undone. This will delete all sessions for user: <strong>{identity?.traits?.['email']}</strong>
{/snippet}

<div class="flex gap-4">
	<Button onclick={onDeleteIdentitySessionsClick} disabled={table.getRowModel().rows.length === 0}>
		Delete all sessions
	</Button>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 items-center space-x-2">
		{#if aalCol}
			<DataTableFacetedFilter column={aalCol} title="Aal" options={aals} />
		{/if}

		{#if isFiltered}
			<Button variant="ghost" onclick={() => table.resetColumnFilters()} class="h-8 px-2 lg:px-3">
				Reset
				<IconX />
			</Button>
		{/if}
	</div>
	<DataTableViewOptions {table} />
</div>
