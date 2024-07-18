<script lang="ts">
	import type { PageData } from './$types';
	import { Button } from '$lib/components/ui/button/index.js';

	export let data: PageData;
</script>

<main class="flex h-screen w-full flex-col items-center justify-center bg-[#1A2238]">
	<h1 class="mb-4 font-extrabold text-white">Auth error</h1>

	{#if data.flow?.id}
		<p class="font-extrabold text-white">Error id: {data.flow.id}</p>
	{/if}

	{#if data.flow?.error}
		<pre>{JSON.stringify(data.flow.error, null, 2)}</pre>
	{/if}

	<Button href="/">Back to home page</Button>
</main>
