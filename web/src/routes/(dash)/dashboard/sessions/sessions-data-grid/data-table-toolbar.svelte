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
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import {
		type ListSessionsExpandEnum,
		ListSessionsExpandEnum as ListSessionsExpandOptions
	} from '$lib/gen/juicer_openapi';
	import Label from '$lib/components/ui/label/label.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import { expandSessions } from '../expandSessions.svelte';
	import Input from '$lib/components/ui/input/input.svelte';

	let { table }: { table: Table<TData> } = $props();

	const isFiltered = $derived(table.getState().columnFilters.length > 0);
	const aalCol = $derived(table.getColumn('authenticatorAssuranceLevel'));
	const expandOptions = Object.values(ListSessionsExpandOptions);

	const searchParams = new SvelteURLSearchParams(page.url.searchParams);

	function onExpandCheckedChange(val: boolean, opt: ListSessionsExpandEnum) {
		const all = searchParams.getAll('expand') as ListSessionsExpandEnum[];
		let newUrl = page.url.pathname;
		if (val) {
			if (!all.includes(opt)) {
				searchParams.append('expand', opt);
				expandSessions.expanded = [...all, opt];
			}
		} else {
			searchParams.delete('expand');
			for (const p of all.filter(x => x !== opt)) {
				searchParams.append('expand', p);
			}
			expandSessions.expanded = expandSessions.expanded.filter(x => x !== opt);
		}
		if (searchParams.size > 0) {
			newUrl += `?${searchParams}`;
		}
		goto(newUrl);
	}

	$effect(() => {
		expandSessions.expanded = searchParams.getAll('expand') as ListSessionsExpandEnum[];
	});
</script>

<div class="flex gap-4">
	<p>Expand</p>
	<div class="flex gap-2">
		{#each expandOptions as opt, i}
			{@const checked = expandSessions.expanded.includes(opt)}
			<div>
				<Checkbox
					id="expand-{opt}"
					aria-labelledby="expand-{opt}-label"
					{checked}
					value={opt}
					onCheckedChange={val => onExpandCheckedChange(val, opt)}
				/>
				<Label
					id="expand-{opt}-label"
					for="expand-{opt}"
					class="text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
				>
					{opt}
				</Label>
			</div>
		{/each}
	</div>
</div>

<div class="flex items-center justify-between">
	<div class="flex flex-1 items-center space-x-2">
		{#if expandSessions.expanded.includes('identity')}
			<Input
				placeholder="Filter by email..."
				value={(table.getColumn('email')?.getFilterValue() as string) ?? ''}
				oninput={e => {
					table.getColumn('email')?.setFilterValue(e.currentTarget.value);
				}}
				onchange={e => {
					table.getColumn('email')?.setFilterValue(e.currentTarget.value);
				}}
				class="h-8 w-[150px] lg:w-[250px]"
			/>
			<Input
				placeholder="Filter by username..."
				value={(table.getColumn('username')?.getFilterValue() as string) ?? ''}
				oninput={e => {
					table.getColumn('username')?.setFilterValue(e.currentTarget.value);
				}}
				onchange={e => {
					table.getColumn('username')?.setFilterValue(e.currentTarget.value);
				}}
				class="h-8 w-[150px] lg:w-[250px]"
			/>
			<Input
				placeholder="Filter by full name..."
				value={(table.getColumn('fullName')?.getFilterValue() as string) ?? ''}
				oninput={e => {
					table.getColumn('fullName')?.setFilterValue(e.currentTarget.value);
				}}
				onchange={e => {
					table.getColumn('fullName')?.setFilterValue(e.currentTarget.value);
				}}
				class="h-8 w-[150px] lg:w-[250px]"
			/>
		{/if}

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
