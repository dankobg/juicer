<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import Registration from './registration.svelte';

	let { data }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Registration</title>
</svelte:head>

<Registration {data} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
