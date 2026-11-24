<script lang="ts">
	import { Color } from '#lib/gen/juicer_pb.ts';

	let { username, color }: { username: string; color: Color } = $props();
</script>

<div class="flex items-center justify-center gap-2">
	{username}
	<svg
		class={['h-5 w-5 stroke-10', color === Color.WHITE ? 'fill-white stroke-black' : 'fill-black stroke-white']}
		viewBox="0 0 100 100"
		xmlns="http://www.w3.org/2000/svg"
	>
		<circle cx="50" cy="50" r="40" />
	</svg>
</div>
