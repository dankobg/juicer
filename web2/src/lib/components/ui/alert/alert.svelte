<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { type Variant, alertVariants } from './index.js';
	import CircleAlert from 'lucide-svelte/icons/circle-alert';
	import TriangleAlert from 'lucide-svelte/icons/triangle-alert';
	import Info from 'lucide-svelte/icons/info';
	import { cn } from '$lib/utils.js';

	type $$Props = HTMLAttributes<HTMLDivElement> & {
		variant?: Variant;
		variantIcon?: boolean;
	};

	let className: $$Props['class'] = undefined;
	export let variant: $$Props['variant'] = 'default';
	export { className as class };
</script>

<div class={cn(alertVariants({ variant }), className)} {...$$restProps} role="alert">
	{#if variant !== 'default'}
		{#if variant === 'destructive' || variant === 'error'}
			<TriangleAlert class="h-4 w-4" />
		{/if}
		{#if variant === 'warn'}
			<CircleAlert class="h-4 w-4" />
		{/if}
		{#if variant === 'info'}
			<Info class="h-4 w-4" />
		{/if}
	{/if}

	<slot />
</div>
