<script lang="ts">
	import { IdentityStateEnum } from '$lib/gen/juicer_openapi';
	import { Badge } from '$lib/components/ui/badge/index';
	import { stateIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = $derived(stateIcons.get(value as IdentityStateEnum));
	let color = $derived.by(() => {
		switch (value as IdentityStateEnum) {
			case IdentityStateEnum.Active:
				return 'text-green-400';
			case IdentityStateEnum.Inactive:
				return 'text-red-400';
		}
	});
</script>

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	{#if Icon}
		<Icon />
	{/if}
	<span>{value ?? ''}</span>
</Badge>
