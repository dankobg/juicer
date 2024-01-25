<script lang="ts">
	import type { PageData } from './$types';
	import Password from './Password.svelte';
	import Settings from './Settings.svelte';
	import SocialProviders from './SocialProviders.svelte';

	export let data: PageData;
</script>

<section class="flex flex-col gap-6">
	<Settings {data} />
	<Password {data} />
	<SocialProviders {data} />
</section>
