<script lang="ts">
	import IconX from '@lucide/svelte/icons/x';
	import Button from '#lib/components/ui/button/button.svelte';
	import DataTableViewOptions from '#lib/components/data-grid-shared/data-table-view-options.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { useTableContext } from '#lib/components/data-grid-shared/app-table.ts';
	import type { components } from '#lib/gen/juicer_openapi.ts';
	import { capitalize } from '#lib/utils.ts';
	import DataTableFacetedFilter from '#lib/components/data-grid-shared/data-table-faceted-filter.svelte';

	const {
		gameVariants,
		gameTimeKinds,
		gameTimeCategories,
		gameResults,
		gameResultStatuses,
		gameStates
	}: {
		gameVariants: components['schemas']['GameVariant'][];
		gameTimeKinds: components['schemas']['GameTimeKind'][];
		gameTimeCategories: components['schemas']['GameTimeCategory'][];
		gameResults: components['schemas']['GameResult'][];
		gameResultStatuses: components['schemas']['GameResultStatus'][];
		gameStates: components['schemas']['GameState'][];
	} = $props();

	const table = useTableContext();

	const isFiltered = $derived(table.atoms.columnFilters.get().length > 0);
	let gameVariantIdCol = $derived(table.getColumn('game_variant_id'));
	let gameTimeKindIdCol = $derived(table.getColumn('game_time_kind_id'));
	let gameTimeCategoryIdCol = $derived(table.getColumn('game_time_category_id'));
	// let gameResultIdCol = $derived(table.getColumn('game_result_id'));
	let userGameResultCol = $derived(table.getColumn('user_game_result'));
	let gameResultStatusIdCol = $derived(table.getColumn('game_result_status_id'));
	let gameStateIdCol = $derived(table.getColumn('game_state_id'));
	let ratedCol = $derived(table.getColumn('rated'));

	let gameVariantIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('game_variant_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	let gameTimeKindIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('game_time_kind_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	let gameTimeCategoryIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('game_time_category_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	// let gameResultIds = $derived.by(() => {
	// 	const params = page.url.searchParams.getAll('game_result_id');
	// 	if (params.length > 0) {
	// 		const ids = params.map(Number).filter(Number.isFinite);
	// 		if (ids.length > 0) {
	// 			return ids;
	// 		}
	// 	}
	// });

	let userGameResults = $derived.by(() => {
		const params = page.url.searchParams.getAll('user_game_result');
		if (params.length > 0) {
			return params;
		}
	});

	let gameResultStatusIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('game_result_status_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	let gameStateIds = $derived.by(() => {
		const params = page.url.searchParams.getAll('game_state_id');
		if (params.length > 0) {
			const ids = params.map(Number).filter(Number.isFinite);
			if (ids.length > 0) {
				return ids;
			}
		}
	});

	let rated = $derived.by(() => {
		const param = page.url.searchParams.get('rated');
		if (param) {
			return param === 'true';
		}
	});

	let gameVariantOptions = gameVariants.map(x => ({ label: capitalize(x.name), value: String(x.id) }));
	let gameTimeKindOptions = gameTimeKinds.map(x => ({ label: capitalize(x.name), value: String(x.id) }));
	let gameTimeCategoryOptions = gameTimeCategories.map(x => ({ label: capitalize(x.name), value: String(x.id) }));
	// let gameResultOptions = gameResults.map(x => ({ label: capitalize(x.name), value: String(x.id) }));

	const userGameResultOptions = [
		{ label: 'Win', value: 'win' },
		{ label: 'Loss', value: 'loss' },
		{ label: 'Draw', value: 'draw' },
		{ label: 'Interrupted', value: 'interrupted' }
	];

	let gameResultStatusOptions = gameResultStatuses.map(x => ({ label: capitalize(x.name), value: String(x.id) }));
	let gameStateOptions = gameStates.map(x => ({ label: capitalize(x.name), value: String(x.id) }));

	function gotoWithFilters(params: URLSearchParams) {
		goto(page.url.pathname + params.size ? `?${params}` : '', { reset: false });
	}

	onMount(() => {
		if (gameVariantIdCol && gameVariantIds) {
			gameVariantIdCol?.setFilterValue(gameVariantIds);
		}
		if (gameTimeKindIdCol && gameTimeKindIds) {
			gameTimeKindIdCol.setFilterValue(gameTimeKindIds);
		}
		if (gameTimeCategoryIdCol && gameTimeCategoryIds) {
			gameTimeCategoryIdCol.setFilterValue(gameTimeCategoryIds);
		}
		// if (gameResultIdCol && gameResultIds) {
		// 	gameResultIdCol.setFilterValue(gameResultIds);
		// }
		if (userGameResultCol && userGameResults) {
			userGameResultCol.setFilterValue(userGameResults);
		}
		if (gameResultStatusIdCol && gameResultStatusIds) {
			gameResultStatusIdCol.setFilterValue(gameResultStatusIds);
		}
		if (gameStateIdCol && gameStateIds) {
			gameStateIdCol.setFilterValue(gameStateIds);
		}
		if (ratedCol && rated) {
			ratedCol.setFilterValue(rated);
		}
	});
</script>

<div class="flex items-center justify-between">
	<div class="flex flex-1 flex-wrap items-center space-x-2 gap-y-2">
		{#if gameVariantIdCol}
			<DataTableFacetedFilter
				column={gameVariantIdCol}
				title="Variant"
				options={gameVariantOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_variant_id', opt.value);
					} else {
						sp.delete('game_variant_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if gameTimeKindIdCol}
			<DataTableFacetedFilter
				column={gameTimeKindIdCol}
				title="Time kind"
				options={gameTimeKindOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_time_kind_id', opt.value);
					} else {
						sp.delete('game_time_kind_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if gameTimeCategoryIdCol}
			<DataTableFacetedFilter
				column={gameTimeCategoryIdCol}
				title="Time category"
				options={gameTimeCategoryOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_time_category_id', opt.value);
					} else {
						sp.delete('game_time_category_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		<!-- {#if gameResultIdCol}
			<DataTableFacetedFilter
				column={gameResultIdCol}
				title="Result"
				options={gameResultOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_result_id', opt.value);
					} else {
						sp.delete('game_result_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if} -->

		{#if userGameResultCol}
			<DataTableFacetedFilter
				column={userGameResultCol}
				title="Result"
				options={userGameResultOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('user_game_result', opt.value);
					} else {
						sp.delete('user_game_result', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if gameResultStatusIdCol}
			<DataTableFacetedFilter
				column={gameResultStatusIdCol}
				title="Result status"
				options={gameResultStatusOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_result_status_id', opt.value);
					} else {
						sp.delete('game_result_status_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if gameStateIdCol}
			<DataTableFacetedFilter
				column={gameStateIdCol}
				title="State"
				options={gameStateOptions}
				onFacetChange={(selected, opt) => {
					const sp = new URLSearchParams(new URL(page.url.href).searchParams);
					if (selected) {
						sp.append('game_state_id', opt.value);
					} else {
						sp.delete('game_state_id', opt.value);
					}
					gotoWithFilters(sp);
				}}
			/>
		{/if}

		{#if isFiltered}
			<Button
				variant="ghost"
				onclick={() => {
					table.resetColumnFilters();
					gotoWithFilters(new URLSearchParams());
				}}
				class="h-8 px-2 lg:px-3"
			>
				Reset
				<IconX />
			</Button>
		{/if}
	</div>

	<DataTableViewOptions />
</div>
