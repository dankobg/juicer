<script lang="ts">
	import type { PageProps } from './$types';
	import AuthError from './auth-error.svelte';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Auth error</title>
</svelte:head>

<AuthError {data} />
