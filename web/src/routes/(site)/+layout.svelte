<script lang="ts">
	import type { LayoutData } from '../$types';
	import { onMount, type Snippet } from 'svelte';
	import Header from '$lib/components/header/header.svelte';
	import { soundManager } from '$lib/state/sound-manager.svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	onMount(() => {
		soundManager.preloadSounds();
	});
</script>

<Header user={data.auth.user ?? undefined} logoutUrl={data.logoutUrl} />
{@render children()}
