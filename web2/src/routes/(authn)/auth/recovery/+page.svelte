<script lang="ts">
	import type { PageData } from './$types';
	import Recovery from './Recovery.svelte';

	export let data: PageData;
</script>

<svelte:head>
	<title>Recovery</title>
</svelte:head>

<Recovery {data} />
