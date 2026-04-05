<script lang="ts">
	import type { PageProps } from './$types';
	import Game from './game.svelte';

	let { data }: PageProps = $props();
</script>

<Game authenticated={data?.auth?.session?.active ?? false} user={data?.auth?.user} />
