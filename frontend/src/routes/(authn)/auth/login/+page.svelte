<script lang="ts">
  import type { PageData } from './$types';
  import LoginPage from './LoginPage.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>Juicer - Log in</title>
</svelte:head>

<LoginPage {data} />
