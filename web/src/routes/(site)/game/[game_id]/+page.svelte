<script lang="ts">
	import type { PageProps } from './$types';
	import Game from './game.svelte';

	let { data, params }: PageProps = $props();
</script>

<svelte:head>
	<title>Juicer Game</title>
</svelte:head>

<Game {data} {params} />
