<script lang="ts" module>
	type TData = unknown;
</script>

<script lang="ts" generics="TData">
	import IconChevronRight from '@lucide/svelte/icons/chevron-right';
	import IconChevronLeft from '@lucide/svelte/icons/chevron-left';
	import IconChevronsRight from '@lucide/svelte/icons/chevrons-right';
	import IconChevronsLeft from '@lucide/svelte/icons/chevrons-left';
	import type { Table } from '@tanstack/table-core';
	import * as Select from '$lib/components/ui/select/index';
	import { Button } from '$lib/components/ui/button/index';

	let { table }: { table: Table<TData> } = $props();
</script>

<div class="flex items-center justify-between px-2">
	<div class="text-muted-foreground flex-1 text-sm">
		{table.getFilteredSelectedRowModel().rows.length} of
		{table.getFilteredRowModel().rows.length} row(s) selected.
	</div>
	<div class="flex items-center space-x-6 lg:space-x-8">
		<div class="flex items-center space-x-2">
			<p class="text-sm font-medium">Rows per page</p>
			<Select.Root
				allowDeselect={false}
				type="single"
				value={`${table.getState().pagination.pageSize}`}
				onValueChange={value => {
					table.setPageSize(Number(value));
				}}
			>
				<Select.Trigger class="h-8 w-[70px]">
					{String(table.getState().pagination.pageSize)}
				</Select.Trigger>
				<Select.Content side="top">
					{#each [10, 20, 30, 50, 100] as pageSize (pageSize)}
						<Select.Item value={`${pageSize}`}>
							{pageSize}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>
		<div class="flex w-[100px] items-center justify-center text-sm font-medium">
			Page {table.getState().pagination.pageIndex + 1} of
			{table.getPageCount()}
		</div>
		<div class="flex items-center space-x-2">
			<Button
				variant="outline"
				class="hidden size-8 p-0 lg:flex"
				onclick={() => table.setPageIndex(0)}
				disabled={!table.getCanPreviousPage()}
			>
				<span class="sr-only">Go to first page</span>
				<IconChevronsLeft />
			</Button>
			<Button
				variant="outline"
				class="size-8 p-0"
				onclick={() => table.previousPage()}
				disabled={!table.getCanPreviousPage()}
			>
				<span class="sr-only">Go to previous page</span>
				<IconChevronLeft />
			</Button>
			<Button variant="outline" class="size-8 p-0" onclick={() => table.nextPage()} disabled={!table.getCanNextPage()}>
				<span class="sr-only">Go to next page</span>
				<IconChevronRight />
			</Button>
			<Button
				variant="outline"
				class="hidden size-8 p-0 lg:flex"
				onclick={() => table.setPageIndex(table.getPageCount() - 1)}
				disabled={!table.getCanNextPage()}
			>
				<span class="sr-only">Go to last page</span>
				<IconChevronsRight />
			</Button>
		</div>
	</div>
</div>
