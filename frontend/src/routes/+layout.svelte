<script lang="ts">
  import 'carbon-components-svelte/css/all.css';
  import { Toaster } from 'svelte-french-toast';
  import { shittyFirefoxPatchForDragEventPointerCoordinates } from '$lib/util/firefox_fix';
  import { onMount } from 'svelte';

  onMount(() => {
    shittyFirefoxPatchForDragEventPointerCoordinates();
  });
</script>

<Toaster />
<slot />
