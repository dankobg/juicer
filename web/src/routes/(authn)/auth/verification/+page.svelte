<script lang="ts">
	import type { PageProps } from './$types';
	import Verification from './verification.svelte';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Verification</title>
</svelte:head>

<Verification {data} />
