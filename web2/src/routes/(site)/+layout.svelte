<script lang="ts">
	import type { LayoutData } from './$types';

	export let data: LayoutData;
</script>

<nav>NAVBAR</nav>

<a class="text-blue-400" href="/dashboard">DASHBOARD</a> <br />
<a class="text-blue-400" href="/auth/login">LOGIN</a> <br />
<a class="text-blue-400" href="/auth/registration">REGISTER</a> <br />

<slot />
