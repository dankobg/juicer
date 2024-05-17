<script lang="ts">
	import { onMount } from 'svelte';
	import { Toaster } from 'svelte-sonner';
	import '../app.css';

	onMount(async () => {
		await import('@dankop/juicer-board');
	});
</script>

<Toaster position="top-center" richColors />
<slot />
