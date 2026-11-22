<script lang="ts">
	import { Badge } from '#lib/components/ui/badge/index.ts';

	let { value }: { value?: 'win' | 'loss' | 'draw' | 'interrupted' } = $props();

	let color = $derived.by(() => {
		switch (value) {
			case 'win':
				return 'text-green-400';
			case 'loss':
				return 'text-red-400';
			case 'draw':
				return 'text-yellow-400';
			case 'interrupted':
				return 'text-orange-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<span>{value ?? ''}</span>
</Badge>
