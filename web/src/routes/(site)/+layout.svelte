<script lang="ts">
	import Navbar from '$lib/Navbar/Navbar.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<Navbar user={data.svc.user} logoutUrl={data.logoutUrl ?? ''} />
<slot />
