<script lang="ts">
	import type { PageProps, Snapshot } from './$types';
	import Login from './login.svelte';

	let { data }: PageProps = $props();
	export const snapshot: Snapshot = { capture: () => {}, restore: () => {} };
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<Login {data} bind:capture={snapshot.capture} bind:restore={snapshot.restore} />
