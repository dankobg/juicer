<script lang="ts">
	import { Color, type PlayerInfo } from '$lib/gen/juicer_pb';

	type Props = {
		player: PlayerInfo;
		color: Color;
		online?: boolean;
		clock?: string;
		active?: boolean;
	};

	let { player, color, online, clock, active }: Props = $props();
</script>

<div class="flex items-center justify-between">
	<span
		class={['rounded-sm bg-green-700 p-1 text-5xl', active ? 'bg-green-700' : 'bg-neutral-200 dark:bg-neutral-700']}
	>
		{clock || '00:00'}
	</span>
	<div class="flex items-start justify-center gap-1">
		<div class="flex flex-wrap items-center gap-1">
			<svg
				class={['h-5 w-5 stroke-10', color === Color.WHITE ? 'fill-white stroke-black' : 'fill-black stroke-white']}
				viewBox="0 0 100 100"
				xmlns="http://www.w3.org/2000/svg"
			>
				<circle cx="50" cy="50" r="40" />
			</svg>
			<span>{player?.username}</span>
		</div>
		<div class="grid [grid-template-areas:'img-stack']">
			<img
				src={player?.avatarUrl || '/images/empty-avatar.svg'}
				alt={`${player?.username} avatar`}
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
	</div>
</div>
