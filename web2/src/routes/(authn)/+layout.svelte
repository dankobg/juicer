<script lang="ts">
	import type { LayoutData } from './$types';
	import Navbar from '$lib/components/navbar/Navbar.svelte';

	export let data: LayoutData;
</script>

<Navbar user={data.auth.user ?? undefined} logoutUrl={data.logoutUrl} />
<slot />
