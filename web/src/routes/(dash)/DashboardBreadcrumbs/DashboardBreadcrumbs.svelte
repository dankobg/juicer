<script lang="ts">
	import { page } from '$app/stores';
	import Breadcrumb from 'flowbite-svelte/Breadcrumb.svelte';
	import BreadcrumbItem from 'flowbite-svelte/BreadcrumbItem.svelte';

	let crumbs: Array<{ label: string; href: string }> = [];

	$: {
		const tokens = $page.url.pathname.split('/').filter(t => t !== '');

		let tokenPath = '';
		crumbs = tokens.map((t, i) => {
			tokenPath += '/' + t;
			t = t.charAt(0).toUpperCase() + t.slice(1);

			let xxx = t;

			if (i === tokens.length - 1 && $page.data.label) {
				xxx = $page.data.label;
			}

			return {
				label: xxx,
				href: tokenPath,
			};
		});

		crumbs.unshift({ label: 'Home', href: '/' });
	}
</script>

<div class="mb-4">
	<!-- <Breadcrumb aria-label="Dashboard breadcrumbs" class="mb-4">
		<BreadcrumbItem href="/" home>Home</BreadcrumbItem>
		<BreadcrumbItem href="/">Projects</BreadcrumbItem>
		<BreadcrumbItem>Flowbite Svelte</BreadcrumbItem>
	</Breadcrumb> -->

	<Breadcrumb aria-label="Dashboard breadcrumbs" class="mb-4">
		{#each crumbs as crumb, idx}
			<BreadcrumbItem href={crumb.href} home={idx === 0}>{crumb.label}</BreadcrumbItem>
		{/each}
	</Breadcrumb>

	<!-- <h1 class="text-xl font-semibold text-gray-900 sm:text-2xl dark:text-white">Dashboard</h1> -->
</div>
