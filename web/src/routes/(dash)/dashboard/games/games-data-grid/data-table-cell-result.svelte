<script lang="ts">
	import { Badge } from '#lib/components/ui/badge/index.ts';

	let { value }: { value?: string } = $props();

	let color = $derived.by(() => {
		switch (value) {
			case 'win':
			case 'won':
				return 'text-green-400';
			case 'loss':
			case 'lost':
				return 'text-red-400';
			default:
				return 'text-orange-400';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<span>{value ?? ''}</span>
</Badge>
