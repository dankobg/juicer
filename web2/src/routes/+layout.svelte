<script lang="ts">
	import { onMount } from 'svelte';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import '../app.css';

	onMount(async () => {
		await import('@dankop/juicer-board');
	});
</script>

<ModeWatcher />
<Toaster position="top-center" richColors />
<slot />
