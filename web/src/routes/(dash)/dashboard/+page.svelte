<script lang="ts">
	import type { PageProps } from './$types';
	import Dashboard from './dashboard.svelte';

	let { data }: PageProps = $props();
</script>

<Dashboard {data} />
