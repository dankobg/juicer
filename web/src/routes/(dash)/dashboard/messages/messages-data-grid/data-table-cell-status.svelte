<script lang="ts">
	import { CourierMessageStatus } from '$lib/gen/juicer_openapi';
	import { Badge } from '$lib/components/ui/badge/index';
	import { statusIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(statusIcons.get(value as CourierMessageStatus));
	let color = $derived.by(() => {
		switch (value as CourierMessageStatus) {
			case CourierMessageStatus.Abandoned:
				return 'text-red-400';
			case CourierMessageStatus.Processing:
				return 'text-purple-400';
			case CourierMessageStatus.Queued:
				return 'text-yellow-400';
			case CourierMessageStatus.Sent:
				return 'text-green-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
