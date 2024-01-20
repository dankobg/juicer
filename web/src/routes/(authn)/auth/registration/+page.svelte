<script lang="ts">
	import type { PageData } from './$types';
	import Registration from './Registration.svelte';

	export let data: PageData;
</script>

<svelte:head>
	<title>Registration</title>
</svelte:head>

<Registration {data} />
