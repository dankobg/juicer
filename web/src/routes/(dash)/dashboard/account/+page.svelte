<script lang="ts">
	import type { PageProps } from './$types';
	import Password from './password.svelte';
	import Settings from './settings.svelte';
	import SocialProviders from './social-providers.svelte';

	let { data }: PageProps = $props();
	let currentFlowForm: 'settings' | 'password' | 'socials' | undefined = $state();
</script>

<section class="grid grid-cols-[repeat(auto-fill,minmax(24rem,1fr))] justify-center gap-4">
	<Settings {data} bind:currentFlowForm />
	<Password {data} bind:currentFlowForm />
	<SocialProviders {data} />
</section>
