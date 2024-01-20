<script lang="ts">
	import type { PageData } from './$types';
	import Verification from './Verification.svelte';

	export let data: PageData;
</script>

<svelte:head>
	<title>Verification</title>
</svelte:head>

<Verification {data} />
