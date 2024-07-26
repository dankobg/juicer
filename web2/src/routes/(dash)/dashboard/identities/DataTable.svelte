<script lang="ts">
	import { Render, Subscribe, createRender, createTable } from 'svelte-headless-table';
	import {
		addHiddenColumns,
		addPagination,
		addSelectedRows,
		addSortBy,
		addTableFilter
	} from 'svelte-headless-table/plugins';
	import { readable } from 'svelte/store';
	import CaretSort from 'svelte-radix/CaretSort.svelte';
	import ChevronDown from 'svelte-radix/ChevronDown.svelte';
	import DataTableActions from './DataTableActions.svelte';
	import DataTableCheckbox from './DataTableCheckbox.svelte';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { cn } from '$lib/utils.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Pagination from '$lib/components/ui/pagination';
	import * as Select from '$lib/components/ui/select';

	const data = Array.from({ length: 500 })
		.fill(0)
		.map((_, i) => ({
			id: crypto.randomUUID(),
			schema_id: Math.random() < 0.5 ? 'customer' : 'employee',
			email: 'test@test.com',
			first_name: `First ${i + 1}`,
			last_name: `Last ${i + 1}`,
			avatar_url: 'https://lh3.googleusercontent.com/a/ACg8ocLRu3aq_ujewp1GIpwyc6L-8C4jPw6nUCqbP70lHVUV1TcZYko=s96-c',
			created_at: new Date(),
			updated_at: new Date(),
			nid: crypto.randomUUID(),
			state: Math.random() < 0.3 ? 'inactive' : 'active',
			state_changed_at: new Date(),
			metadata_public: null,
			metadata_admin: null,
			available_aal: Math.random() < 0.1 ? 'aal2' : 'aal1',
			organization_id: null
		}));

	let perPageOptions = [20, 40, 50, 60, 75, 100];
	let perPageSelected = { value: perPageOptions[2], label: perPageOptions[2].toString() };

	const table = createTable(readable(data), {
		sort: addSortBy({ disableMultiSort: true }),
		page: addPagination({ initialPageSize: perPageSelected.value }),
		filter: addTableFilter({
			fn: ({ filterValue, value }) => value.toLowerCase().includes(filterValue.toLowerCase())
		}),
		select: addSelectedRows(),
		hide: addHiddenColumns()
	});

	const columns = table.createColumns([
		table.column({
			header: (_, { pluginStates }) => {
				const { allPageRowsSelected } = pluginStates.select;
				return createRender(DataTableCheckbox, {
					checked: allPageRowsSelected
				});
			},
			accessor: 'id',
			cell: ({ row }, { pluginStates }) => {
				const { getRowState } = pluginStates.select;
				const { isSelected } = getRowState(row);

				return createRender(DataTableCheckbox, {
					checked: isSelected
				});
			},
			plugins: { sort: { disable: true }, filter: { exclude: true } }
		}),
		table.column({
			header: 'Schema ID',
			accessor: 'schema_id',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'First name',
			accessor: 'first_name',
			plugins: { sort: { disable: false }, filter: { exclude: false } }
		}),
		table.column({
			header: 'Last name',
			accessor: 'last_name',
			plugins: { sort: { disable: false }, filter: { exclude: false } }
		}),
		table.column({
			header: 'E-Mail',
			accessor: 'email',
			plugins: { sort: { disable: false }, filter: { exclude: false } }
		}),
		table.column({
			header: 'NID',
			accessor: 'nid',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'State',
			accessor: 'state'
		}),
		table.column({
			header: 'State changed at',
			accessor: 'state_changed_at',
			cell: ({ value }) => value?.toLocaleDateString(undefined, { hourCycle: 'h23' }) ?? '',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'Available aal',
			accessor: 'available_aal',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'Created at',
			accessor: 'created_at',
			cell: ({ value }) => value?.toLocaleDateString(undefined, { hourCycle: 'h23' }) ?? '',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'Updated at',
			accessor: 'updated_at',
			cell: ({ value }) => value?.toLocaleDateString(undefined, { hourCycle: 'h23' }) ?? '',
			plugins: { sort: { disable: false }, filter: { exclude: true } }
		}),
		table.column({
			header: 'Actions',
			accessor: ({ id }) => id,
			cell: item => {
				return createRender(DataTableActions, { id: item.value });
			},
			plugins: { sort: { disable: true }, filter: { exclude: true } }
		})
	]);

	const { headerRows, pageRows, tableAttrs, tableBodyAttrs, flatColumns, pluginStates, rows } =
		table.createViewModel(columns);

	$: if (perPageSelected) {
		$pageSize = perPageSelected.value;
	}

	const { sortKeys } = pluginStates.sort;

	const { hiddenColumnIds } = pluginStates.hide;
	const ids = flatColumns.map(c => c.id);
	let hideForId = Object.fromEntries(ids.map(id => [id, true]));

	$: $hiddenColumnIds = Object.entries(hideForId)
		.filter(([, hide]) => !hide)
		.map(([id]) => id);

	const { hasNextPage, hasPreviousPage, pageIndex, pageCount, pageSize } = pluginStates.page;
	const { filterValue } = pluginStates.filter;
	const { selectedDataIds } = pluginStates.select;

	const hideableCols = [
		'schema_id',
		'email',
		'first_name',
		'last_name',
		'avatar_url',
		'nid',
		'state',
		'state_changed_at',
		'available_aal',
		'created_at',
		'updated_at'
	];
</script>

<div class="w-full">
	<div class="mb-4 flex items-center gap-4">
		<Input class="max-w-sm" placeholder="Filter emails..." type="text" bind:value={$filterValue} />
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild let:builder>
				<Button variant="outline" class="ml-auto" builders={[builder]}>
					Columns <ChevronDown class="ml-2 h-4 w-4" />
				</Button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content>
				{#each flatColumns as col}
					{#if hideableCols.includes(col.id)}
						<DropdownMenu.CheckboxItem bind:checked={hideForId[col.id]}>
							{col.header}
						</DropdownMenu.CheckboxItem>
					{/if}
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
	<div class="rounded-md border">
		<Table.Root {...$tableAttrs}>
			<Table.Header>
				{#each $headerRows as headerRow}
					<Subscribe rowAttrs={headerRow.attrs()}>
						<Table.Row>
							{#each headerRow.cells as cell (cell.id)}
								<Subscribe attrs={cell.attrs()} let:attrs props={cell.props()} let:props>
									<Table.Head {...attrs} class={cn('[&:has([role=checkbox])]:pl-3')}>
										{#if props.sort.disabled}
											<Render of={cell.render()} />
										{:else}
											<Button variant="ghost" on:click={props.sort.toggle}>
												<Render of={cell.render()} />
												<CaretSort class={cn($sortKeys[0]?.id === cell.id && 'text-foreground', 'ml-2 h-4 w-4')} />
											</Button>
										{/if}
									</Table.Head>
								</Subscribe>
							{/each}
						</Table.Row>
					</Subscribe>
				{/each}
			</Table.Header>
			<Table.Body {...$tableBodyAttrs}>
				{#each $pageRows as row (row.id)}
					<Subscribe rowAttrs={row.attrs()} let:rowAttrs>
						<Table.Row {...rowAttrs} data-state={$selectedDataIds[row.id] && 'selected'}>
							{#each row.cells as cell (cell.id)}
								<Subscribe attrs={cell.attrs()} let:attrs>
									<Table.Cell class="[&:has([role=checkbox])]:pl-3" {...attrs}>
										<Render of={cell.render()} />
									</Table.Cell>
								</Subscribe>
							{/each}
						</Table.Row>
					</Subscribe>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>

	<div class="flex w-full flex-wrap items-center py-4">
		<div class="flex flex-1 items-center">
			<span class="text-muted-foreground mr-4 text-sm">
				selected {Object.keys($selectedDataIds).length} of {$rows.length} row(s)
			</span>

			<Select.Root bind:selected={perPageSelected}>
				<Select.Label>Per page</Select.Label>
				<Select.Trigger class="w-fit">
					<Select.Value placeholder="Per page" />
				</Select.Trigger>
				<Select.Content>
					{#each perPageOptions as opt}
						<Select.Item value={opt}>{opt}</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</div>

		<Pagination.Root count={$rows.length} perPage={$pageSize} let:pages let:currentPage class="w-auto items-end">
			<Pagination.Content>
				<Pagination.Item>
					<Pagination.PrevButton on:click={() => ($pageIndex = $pageIndex - 1)} disabled={!$hasPreviousPage} />
				</Pagination.Item>
				{#each pages as page (page.key)}
					{#if page.type === 'ellipsis'}
						<Pagination.Item>
							<Pagination.Ellipsis />
						</Pagination.Item>
					{:else}
						<Pagination.Item>
							<Pagination.Link
								{page}
								isActive={currentPage == page.value}
								on:click={() => ($pageIndex = page.value - 1)}
							>
								{page.value}
							</Pagination.Link>
						</Pagination.Item>
					{/if}
				{/each}
				<Pagination.Item>
					<Pagination.NextButton on:click={() => ($pageIndex = $pageIndex + 1)} disabled={!$hasNextPage} />
				</Pagination.Item>
			</Pagination.Content>
		</Pagination.Root>
	</div>
</div>
