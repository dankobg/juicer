<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from '$lib/components/ui/sonner/index';
	import ConfirmationDialog from '$lib/components/confirmation-dialog/confirmation-dialog.svelte';
	import { onMount } from 'svelte';
	import { soundManager } from '$lib/sound/sound-manager.svelte';
	import '../app.css';

	onMount(() => {
		if (!customElements.get('juicer-board')) {
			import('@dankop/juicer-board');
		}

		soundManager.preloadSounds();
	});

	let { children } = $props();
</script>

<ConfirmationDialog />
<ModeWatcher />
<Toaster position="top-center" richColors />
{@render children()}
