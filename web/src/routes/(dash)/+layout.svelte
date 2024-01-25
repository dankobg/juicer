<script lang="ts">
	import type { LayoutData } from './$types';
	import DashboardNavbar from './DashboardNavbar/DashboardNavbar.svelte';
	import DashboardSidebar from './DashboardSidebar/DashboardSidebar.svelte';
	import DashboardBreadcrumbs from './DashboardBreadcrumbs/DashboardBreadcrumbs.svelte';

	export let data: LayoutData;
</script>

<DashboardNavbar user={data.svc.user} logoutUrl={data.logoutUrl ?? ''} />

<div class="w-full overflow-hidden bg-gray-50 dark:bg-gray-900">
	<div class="mx-auto flex flex-wrap justify-between container">
		<DashboardSidebar />

		<main class="flex-1">
			<div class="px-4 py-6">
				<DashboardBreadcrumbs />
				<slot />
			</div>
		</main>
	</div>
</div>
