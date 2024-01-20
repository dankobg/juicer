<script lang="ts">
	import type { PageData } from './$types';
	import Login from './Login.svelte';

	export let data: PageData;
</script>

<svelte:head>
	<title>Login</title>
</svelte:head>

<Login {data} />
