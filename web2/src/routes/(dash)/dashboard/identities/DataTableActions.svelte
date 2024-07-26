<script lang="ts">
	import Ellipsis from 'lucide-svelte/icons/ellipsis';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Eye from 'lucide-svelte/icons/eye';
	import Pencil from 'lucide-svelte/icons/pencil';
	import Trash from 'lucide-svelte/icons/trash-2';

	export let id: string;
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Button variant="ghost" size="icon" builders={[builder]}>
			<span class="sr-only">Open menu</span>
			<Ellipsis class="h-4 w-4" />
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Item on:click={() => navigator.clipboard.writeText(id)}>Copy ID</DropdownMenu.Item>
		<DropdownMenu.Item class="flex gap-4">
			<Eye class="h-4 w-4 text-purple-500" />
			View
		</DropdownMenu.Item>
		<DropdownMenu.Item class="flex gap-4">
			<Pencil class="h-4 w-4 text-sky-500" />
			Edit
		</DropdownMenu.Item>
		<DropdownMenu.Item class="flex gap-4">
			<Trash class="text-destructive h-4 w-4" />
			Delete
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
