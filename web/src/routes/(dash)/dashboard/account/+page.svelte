<script lang="ts">
	import type { PageData } from './$types';
	import Password from './Password.svelte';
	import Settings from './Settings.svelte';
	import SocialProviders from './SocialProviders.svelte';

	export let data: PageData;
	export let currentFlowForm: 'settings' | 'password' | 'socials' | undefined = undefined;
</script>

<section class="flex flex-col gap-6">
	<Settings {data} bind:currentFlowForm />
	<Password {data} bind:currentFlowForm />
	<SocialProviders {data} />
</section>
