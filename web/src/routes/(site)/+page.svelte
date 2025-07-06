<script lang="ts">
	import type { PageProps } from './$types';
	import Lobby from './lobby/lobby.svelte';

	let { data }: PageProps = $props();
</script>

<Lobby {data} />
