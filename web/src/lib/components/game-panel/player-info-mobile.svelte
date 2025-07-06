<script lang="ts">
	import Progress from '$lib/components/ui/progress/progress.svelte';

	type Props = {
		username: string;
		avatar?: string;
		avatarAlt?: string;
		hasActiveTurn: boolean;
		showFirstMoveTimer: boolean;
		showReconnectTimer: boolean;
		firstMoveTimerText: string;
		reconnectTimerText: string;
		clock?: string;
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
		online
	}: Props = $props();
</script>

<div class="grid w-full gap-2 p-2 sm:hidden">
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
	<div class="flex flex-wrap items-center justify-between gap-2">
		<div class="flex items-center gap-2">
			<svg
				class={['h-3 w-3', online ? 'fill-green-400' : 'fill-gray-400']}
				viewBox="0 0 100 100"
				xmlns="http://www.w3.org/2000/svg"
			>
				<circle cx="50" cy="50" r="50" />
			</svg>
			<span>
				{username}
			</span>
		</div>
		<div>
			<span
				class={[
					'rounded-sm bg-green-700 p-1 text-2xl',
					hasActiveTurn ? 'bg-green-700' : 'bg-neutral-200 dark:bg-neutral-700'
				]}
			>
				{clock || '00:00'}
			</span>
		</div>
	</div>
</div>
