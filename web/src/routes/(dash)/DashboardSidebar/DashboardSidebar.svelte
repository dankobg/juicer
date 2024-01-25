<script lang="ts">
	import {
		Sidebar,
		SidebarGroup,
		SidebarItem,
		SidebarWrapper,
		SidebarDropdownItem,
		SidebarDropdownWrapper,
	} from 'flowbite-svelte';
	import {
		ChartPieSolid,
		GridSolid,
		ShoppingCartSolid,
		MailBoxSolid,
		UserSolid,
		ArrowRightToBracketSolid,
		FileEditSolid,
	} from 'flowbite-svelte-icons';

	let spanClass = 'flex-1 ms-3 whitespace-nowrap';
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

			<!-- <SidebarDropdownWrapper label="E-commerce">
				<svelte:fragment slot="icon">
					<ShoppingCartSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
				<SidebarDropdownItem label="Products" />
				<SidebarDropdownItem label="Billing" />
				<SidebarDropdownItem label="Invoice" />
			</SidebarDropdownWrapper>

			<SidebarItem label="Sign In">
				<svelte:fragment slot="icon">
					<ArrowRightToBracketSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem>

			<SidebarItem label="Sign Up">
				<svelte:fragment slot="icon">
					<FileEditSolid
						class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
					/>
				</svelte:fragment>
			</SidebarItem> -->
		</SidebarGroup>
	</SidebarWrapper>
</Sidebar>
