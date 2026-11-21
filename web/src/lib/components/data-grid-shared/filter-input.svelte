<script lang="ts" generics="TData extends RowData">
	import Input from '#lib/components/ui/input/input.svelte';
	import type { RowData, Column } from '@tanstack/svelte-table';
	import type { AppFeatures } from './app-table';

	let {
		title,
		column,
		value,
		onFilterChange
	}: {
		title?: string;
		column: Column<AppFeatures, TData>;
		value?: string;
		onFilterChange?: (value: string, column: Column<AppFeatures, TData>) => void;
	} = $props();

	const placeholder = $derived(column.columnDef.meta?.label ?? title ?? column.id);
</script>

<Input
	placeholder="filter by {placeholder}..."
	value={value ?? (column?.getFilterValue() as string) ?? ''}
	oninput={e => {
		column?.setFilterValue(e.currentTarget.value);
		onFilterChange?.(e.currentTarget.value, column);
	}}
	onchange={e => {
		column?.setFilterValue(e.currentTarget.value);
		onFilterChange?.(e.currentTarget.value, column);
	}}
	class="h-8 w-[150px] lg:w-[250px]"
/>
