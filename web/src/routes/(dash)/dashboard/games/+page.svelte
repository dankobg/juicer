<script lang="ts">
	import DataTable from './games-data-grid/data-table.svelte';
	import { makeColumns } from './games-data-grid/columns';
	import type { PageProps } from './$types';
	import { gameResults, gameResultStatuses, gameStates } from './games-data-grid/data.ts';

	let { data }: PageProps & {} = $props();

	const gameVariants = $derived(data.gameVariantsResult?.data?.data ?? []);
	const gameTimeKinds = $derived(data.gameTimeKindsResult?.data?.data ?? []);
	const gameTimeCategories = $derived(data.gameTimeCategoriesResult?.data?.data ?? []);
	const columns = $derived(makeColumns(data?.auth?.user?.id, gameVariants, gameTimeKinds, gameTimeCategories));
</script>

{#if data.gamesResult?.data}
	<div class="flex h-full flex-1 flex-col space-y-8 p-8">
		<div class="flex items-center justify-between space-y-2">
			<div>
				<h2 class="text-2xl font-bold tracking-tight">Games</h2>
				<p class="text-muted-foreground">List of Games</p>
			</div>
		</div>
		<DataTable
			data={data.gamesResult?.data.data ?? []}
			meta={data.gamesResult?.data.meta}
			{columns}
			{gameVariants}
			{gameTimeKinds}
			{gameTimeCategories}
			{gameResults}
			{gameResultStatuses}
			{gameStates}
		/>
	</div>
{/if}
