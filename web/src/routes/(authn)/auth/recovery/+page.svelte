<script lang="ts">
	import type { PageProps } from './$types';
	import Recovery from './recovery.svelte';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Recovery</title>
</svelte:head>

<Recovery {data} />
