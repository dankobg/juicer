<script lang="ts">
	import type { PageData } from './$types';
	import Error from './Error.svelte';

	export let data: PageData;
</script>

<svelte:head>
	<title>Error</title>
</svelte:head>

<Error {data} />
