<script lang="ts">
	import type { PageProps } from './$types';
	import * as Card from '$lib/components/ui/card/index';
	import { PieChart, Text } from 'layerchart';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import IconUsers from '@lucide/svelte/icons/users';
	import IconCat from '@lucide/svelte/icons/cat';
	import IconBuilding2 from '@lucide/svelte/icons/building-2';
	import IconHandHeart from '@lucide/svelte/icons/hand-heart';

	let { data }: PageProps = $props();

	const chartConfigGameStats = {
		draw: { label: 'Draw', color: 'var(--color-yellow-400)' },
		interrupted: { label: 'Interrupted', color: 'var(--color-purple-600)' },
		loss: { label: 'Loss', color: 'var(--color-red-600)' },
		win: { label: 'Win', color: 'var(--color-green-600)' }
	} satisfies Chart.ChartConfig;

	let chartDataAll = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.all.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.all.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.all.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.all.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});

	let chartDataHyperbullet = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.hyperbullet.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.hyperbullet.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.hyperbullet.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.hyperbullet.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});

	let chartDataBullet = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.bullet.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.bullet.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.bullet.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.bullet.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});

	let chartDataBlitz = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.blitz.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.blitz.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.blitz.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.blitz.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});

	let chartDataRapid = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.rapid.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.rapid.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.rapid.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.rapid.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});

	let chartDataClassical = $derived.by(() => {
		if (!data?.gameStatsResult?.data) {
			return [];
		}
		return [
			{
				label: chartConfigGameStats.draw.label,
				value: data.gameStatsResult.data.classical.draw,
				color: chartConfigGameStats.draw.color
			},
			{
				label: chartConfigGameStats.interrupted.label,
				value: data.gameStatsResult.data.classical.interrupted,
				color: chartConfigGameStats.interrupted.color
			},
			{
				label: chartConfigGameStats.win.label,
				value: data.gameStatsResult.data.classical.win,
				color: chartConfigGameStats.win.color
			},
			{
				label: chartConfigGameStats.loss.label,
				value: data.gameStatsResult.data.classical.loss,
				color: chartConfigGameStats.loss.color
			}
		];
	});
</script>

{#if data?.gameStatsResult?.data}
	<!-- <div class="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Adoptions</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.gameStatsResult?.data?.stats?.adoptions.total}
				</Card.Title>
				<Card.Action>
					<IconHandHeart />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.gameStatsResult?.data?.stats?.adoptions?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.gameStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending adoptions'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Favorites</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.gameStatsResult?.data?.stats?.animals.favorites}
				</Card.Title>
				<Card.Action>
					<IconUsers />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">{favoriteAnimalsDesc}</div>
				<div class="text-muted-foreground">
					{favoriteAnimalsText}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Animals</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.gameStatsResult?.data?.stats?.animals.posted.total}
				</Card.Title>
				<Card.Action>
					<IconCat />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.gameStatsResult?.data?.stats?.animals?.posted.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.gameStatsResult?.data?.stats?.animals?.posted.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending animals'}
				</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="@container/card">
			<Card.Header>
				<Card.Description>Organizations</Card.Description>
				<Card.Title class="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
					Total {data?.gameStatsResult?.data?.stats?.organizations?.total}
				</Card.Title>
				<Card.Action>
					<IconBuilding2 />
				</Card.Action>
			</Card.Header>
			<Card.Footer class="flex-col items-start gap-1.5 text-sm">
				<div class="line-clamp-1 flex gap-2 font-medium">
					Pending approval {data?.gameStatsResult?.data?.stats?.organizations?.by_status.pending}
				</div>
				<div class="text-muted-foreground">
					{data?.gameStatsResult?.data?.stats?.organizations?.by_status.pending > 0
						? 'Waiting for admin approval'
						: 'No pending organizations'}
				</div>
			</Card.Footer>
		</Card.Root>
	</div> -->

	<div class="grid grid-cols-[repeat(auto-fill,minmax(min(28rem,100%),1fr))] justify-center gap-4">
		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>All</Card.Title>
				<Card.Description>All categories stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataAll}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.all?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">All categories stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Hyperbullet</Card.Title>
				<Card.Description>Hyperbullet stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataHyperbullet}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.hyperbullet?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Hyperbullet stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Bullet</Card.Title>
				<Card.Description>Bullet stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataBullet}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.bullet?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Bullet stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Blitz</Card.Title>
				<Card.Description>Blitz stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataBlitz}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.blitz?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Blitz stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Rapid</Card.Title>
				<Card.Description>Rapid stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataRapid}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.rapid?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Rapid stats</div>
			</Card.Footer>
		</Card.Root>

		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Classical</Card.Title>
				<Card.Description>Classical stats</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={chartConfigGameStats} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={chartDataClassical}
						key="label"
						value="value"
						c="color"
						innerRadius={60}
						padding={28}
						props={{
							pie: {
								motion: 'tween'
							}
						}}
						legend
					>
						{#snippet aboveMarks()}
							<Text
								value={data.gameStatsResult.data?.classical?.total}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
			<Card.Footer class="flex-col gap-2 text-sm">
				<div class="leading-none text-muted-foreground">Classical stats</div>
			</Card.Footer>
		</Card.Root>
	</div>
{/if}
