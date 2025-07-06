<script lang="ts">
	import Progress from '$lib/components/ui/progress/progress.svelte';

	type Props = {
		username: string;
		avatar?: string;
		avatarAlt?: string;
		hasActiveTurn: boolean;
		showFirstMoveTimer: boolean;
		showReconnectTimer?: boolean;
		firstMoveTimerText: string;
		reconnectTimerText?: string;
		clock?: string;
		progressMax?: number;
		progressValue?: number;
		online: boolean;
	};

	let {
		username,
		avatar,
		avatarAlt,
		hasActiveTurn = false,
		showFirstMoveTimer = false,
		showReconnectTimer = false,
		firstMoveTimerText,
		reconnectTimerText,
		clock,
		progressMax = 100,
		progressValue = 0,
		online
	}: Props = $props();
</script>

<div class="flex flex-col gap-2 p-2">
	{#if showFirstMoveTimer}
		<p class="bg-sky-700 px-2 text-center font-medium">
			{firstMoveTimerText}
		</p>
	{/if}
	{#if showReconnectTimer}
		<p class="bg-sky-700 px-2 text-center font-medium">
			{reconnectTimerText}
		</p>
	{/if}
	<div class="flex items-center justify-between">
		<span
			class={[
				'rounded-sm bg-green-700 p-1 text-5xl',
				hasActiveTurn ? 'bg-green-700' : 'bg-neutral-200 dark:bg-neutral-700'
			]}
		>
			{clock || '00:00'}
		</span>
		<div class="grid justify-items-end">
			<div class="grid [grid-template-areas:'img-stack']">
				<img
					src={avatar || '/images/empty-avatar.svg'}
					alt={avatarAlt}
					class="m-0 h-12 w-12 max-w-full object-contain p-0 [grid-area:img-stack]"
				/>
				<svg
					class={['h-3 w-3 [grid-area:img-stack]', online ? 'fill-green-400' : 'fill-gray-400']}
					viewBox="0 0 100 100"
					xmlns="http://www.w3.org/2000/svg"
				>
					<circle cx="50" cy="50" r="50" />
				</svg>
			</div>
			<div>{username}</div>
		</div>
	</div>
	<Progress
		max={progressMax}
		value={progressValue}
		class="h-2 rounded-full [&>[data-slot='progress-indicator']]:bg-green-400"
	/>
</div>
