<script lang="ts">
  import type { PageData } from './$types';
  import Recovery from './RecoveryPage.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>Juicer - Recovery</title>
</svelte:head>

<Recovery {data} />
