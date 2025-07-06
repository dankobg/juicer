<script lang="ts">
	import Header from '$lib/components/header/header.svelte';
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button/index';

	let { data }: PageProps = $props();
</script>

<Header />

<main class="grid h-[calc(100vh-4rem)] place-content-center px-4">
	<div class="text-center">
		<p class="text-primary-900 text-2xl font-bold tracking-tight sm:text-4xl">Authentication error occurred</p>
		{#if data.flow?.id}
			<p class="text-primary-500 mt-4">Error id: {data.flow.id}</p>
		{/if}
		<Button href="/" class="mt-4 font-bold">Go back home</Button>
	</div>
</main>
