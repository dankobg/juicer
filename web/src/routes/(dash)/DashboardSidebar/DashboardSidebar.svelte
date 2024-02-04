<script lang="ts">
	import Sidebar from 'flowbite-svelte/Sidebar.svelte';
	import SidebarGroup from 'flowbite-svelte/SidebarGroup.svelte';
	import SidebarItem from 'flowbite-svelte/SidebarItem.svelte';
	import SidebarWrapper from 'flowbite-svelte/SidebarWrapper.svelte';
	import SidebarDropdownItem from 'flowbite-svelte/SidebarDropdownItem.svelte';
	import SidebarDropdownWrapper from 'flowbite-svelte/SidebarDropdownWrapper.svelte';
	import ChartPieSolid from 'flowbite-svelte-icons/ChartPieSolid.svelte';
	import ShoppingCartSolid from 'flowbite-svelte-icons/ShoppingCartSolid.svelte';
	import UserSolid from 'flowbite-svelte-icons/UserSolid.svelte';
</script>

<Sidebar class="hidden xl:block flex-shrink-0">
	<SidebarWrapper class="bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700">
		<SidebarGroup>
			<SidebarItem label="Dashboard" href="/dashboard">
				<svelte:fragment slot="icon">
					<ChartPieSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>

			<SidebarItem label="Account" href="/dashboard/account">
				<svelte:fragment slot="icon">
					<UserSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>

			<SidebarDropdownWrapper label="Identities">
				<svelte:fragment slot="icon">
					<ShoppingCartSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
				<SidebarDropdownItem label="Users" />
				<SidebarDropdownItem label="Schemas" />
				<SidebarDropdownItem label="Sessions" />
			</SidebarDropdownWrapper>
		</SidebarGroup>
	</SidebarWrapper>
</Sidebar>
