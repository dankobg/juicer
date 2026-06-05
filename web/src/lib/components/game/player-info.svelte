<script lang="ts">
	import { Color, type PlayerInfo } from '$lib/gen/juicer_pb';

	type Props = {
		player: PlayerInfo;
		color: Color;
		online?: boolean;
		active?: boolean;
		clockMs?: number;
		clockPrecision?: (ms: number) => 'deciseconds' | 'centiseconds' | null;
	};

	let { player, color, online, active, clockMs, clockPrecision }: Props = $props();

	function formatChessTime(totalMs: number, precisionFn?: Props['clockPrecision']): string {
		const precision = precisionFn?.(totalMs) ?? null;

		if (!Number.isFinite(totalMs) || totalMs <= 0) {
			if (precision === 'deciseconds') return '00:00.0';
			if (precision === 'centiseconds') return '00:00.00';
			return '00:00';
		}

		const totalSeconds = Math.floor(totalMs / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = Math.floor(totalSeconds % 60);

		if (precision !== null) {
			const baseTime = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;

			if (precision === 'deciseconds') {
				const tenths = Math.floor((totalMs % 1000) / 100);
				return `${baseTime}.${tenths}`;
			}

			if (precision === 'centiseconds') {
				const hundredths = Math.floor((totalMs % 1000) / 10);
				return `${baseTime}.${hundredths.toString().padStart(2, '0')}`;
			}
		}

		if (hours > 0) {
			return `${hours}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
		}

		return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
	}
</script>

<div class="flex items-center justify-between">
	<span
		class={['rounded-sm bg-green-700 p-1 text-5xl', active ? 'bg-green-700' : 'bg-neutral-200 dark:bg-neutral-700']}
	>
		{formatChessTime(clockMs ?? 0, clockPrecision)}
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
