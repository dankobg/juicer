<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	let { value, ...rest }: HTMLAttributes<HTMLSpanElement> & { value?: string } = $props();
</script>

<span {...rest}>{value || ''}</span>
