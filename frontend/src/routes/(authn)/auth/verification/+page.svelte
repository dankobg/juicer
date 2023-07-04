<script lang="ts">
  import type { PageData } from './$types';
  import Verification from './VerificationPage.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>Juicer - Verification</title>
</svelte:head>

<Verification {data} />
