<script lang="ts">
	import { CourierMessageType } from '$lib/gen/juicer_openapi';
	import { Badge } from '$lib/components/ui/badge/index';
	import { typeIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = typeIcons.get(value as CourierMessageType);
	let color = $derived.by(() => {
		switch (value as CourierMessageType) {
			case CourierMessageType.Email:
				return 'text-blue-400';
			case CourierMessageType.Phone:
				return 'text-blue-400';
			default:
				return '';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<Icon class="" />
	<span>{value ?? ''}</span>
</Badge>
