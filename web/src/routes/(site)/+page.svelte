<script lang="ts">
	import type { PageProps } from './$types';
	import Lobby from './lobby/lobby.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Juicer home</title>
</svelte:head>

<Lobby {data} {params} />
