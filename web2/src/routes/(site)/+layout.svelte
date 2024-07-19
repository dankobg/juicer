<script lang="ts">
	import type { LayoutData } from './$types';

	export let data: LayoutData;
</script>

<nav>NAVBAR</nav>

<a class="text-blue-400" href="/dashboard">DASHBOARD</a> <br />
<a class="text-blue-400" href="/dashboard/account">ACCOUNT</a> <br />

{#if data.auth.session?.active}
	<pre>{JSON.stringify(data.auth, null, 2)}</pre>
{:else}
	<a class="text-blue-400" href="/auth/login">LOGIN</a> <br />
	<a class="text-blue-400" href="/auth/registration">REGISTER</a> <br />
{/if}

<slot />
