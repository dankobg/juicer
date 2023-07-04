<script lang="ts">
  import type { PageData } from './$types';
  import Registration from './RegistrationPage.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>Juicer - Registration</title>
</svelte:head>

<Registration {data} />
