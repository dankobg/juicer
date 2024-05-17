<script lang="ts">
	import type { User } from '$lib/kratos/service';
	import Navbar from 'flowbite-svelte/Navbar.svelte';
	import NavBrand from 'flowbite-svelte/NavBrand.svelte';
	import Button from 'flowbite-svelte/Button.svelte';
	import Avatar from 'flowbite-svelte/Avatar.svelte';
	import Dropdown from 'flowbite-svelte/Dropdown.svelte';
	import DropdownItem from 'flowbite-svelte/DropdownItem.svelte';
	import DropdownHeader from 'flowbite-svelte/DropdownHeader.svelte';
	import DropdownDivider from 'flowbite-svelte/DropdownDivider.svelte';
	import Drawer from 'flowbite-svelte/Drawer.svelte';
	import CloseButton from 'flowbite-svelte/CloseButton.svelte';
	import Sidebar from 'flowbite-svelte/Sidebar.svelte';
	import SidebarDropdownItem from 'flowbite-svelte/SidebarDropdownItem.svelte';
	import SidebarDropdownWrapper from 'flowbite-svelte/SidebarDropdownWrapper.svelte';
	import SidebarGroup from 'flowbite-svelte/SidebarGroup.svelte';
	import SidebarItem from 'flowbite-svelte/SidebarItem.svelte';
	import SidebarWrapper from 'flowbite-svelte/SidebarWrapper.svelte';
	import ChartPieSolid from 'flowbite-svelte-icons/ChartPieSolid.svelte';
	import ShoppingBagSolid from 'flowbite-svelte-icons/ShoppingBagSolid.svelte';
	import UserSolid from 'flowbite-svelte-icons/UserSolid.svelte';
	import { sineIn } from 'svelte/easing';

	export let logoutUrl: string;
	export let user: User | null;

	let sidebarHidden = true;
	let transitionParams = {
		x: -320,
		duration: 200,
		easing: sineIn,
	};
</script>

<Navbar class="border-b border-gray-200 dark:border-gray-700">
	<div class="flex gap-4 xl:gap-0">
		<Button class="p-2 flex xl:hidden" size="lg" color="alternative" on:click={() => (sidebarHidden = false)}>
			<svg
				class="w-6 h-6 text-gray-800 dark:text-white"
				aria-hidden="true"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
			>
				<path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M5 7h14M5 12h14M5 17h10" />
			</svg>
		</Button>

		<NavBrand href="/">
			<img src="/images/logo.svg" class="me-3 h-6 sm:h-9" alt="Juicer Logo" />
			<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Juicer</span>
		</NavBrand>
	</div>

	{#if user !== null}
		<div class="flex items-center">
			<Avatar id="avatar-menu" src="/images/logo.svg" class="cursor-pointer" />
		</div>

		<Dropdown placement="bottom" triggeredBy="#avatar-menu">
			{#if !(user.firstName === '' && user.email === '')}
				<DropdownHeader>
					{#if user.firstName}
						<span class="block text-sm">{user.firstName}</span>
					{/if}
					{#if user.email}
						<span class="block truncate text-sm font-medium">{user.email}</span>
					{/if}
				</DropdownHeader>
			{/if}

			<DropdownItem href="/dashboard">Dashboard</DropdownItem>
			<DropdownItem href="/dashboard/account">Account</DropdownItem>
			<DropdownDivider />
			<DropdownItem href={logoutUrl}>Log out</DropdownItem>
		</Dropdown>
	{/if}
</Navbar>

<Drawer transitionType="fly" {transitionParams} bind:hidden={sidebarHidden} id="sidebar2">
	<div class="flex items-center">
		<CloseButton on:click={() => (sidebarHidden = true)} class="dark:text-white" />
	</div>

	<Sidebar class="w-full">
		<SidebarWrapper class="bg-white dark:bg-gray-800" divClass="overflow-y-auto py-4 px-3 rounded dark:bg-gray-800">
			<SidebarGroup>
				<SidebarItem label="Dashboard" href="/dashboard" on:click={() => (sidebarHidden = true)}>
					<svelte:fragment slot="icon">
						<ChartPieSolid
							class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
						/>
					</svelte:fragment>
				</SidebarItem>

				<SidebarItem label="Account" href="/dashboard/account" on:click={() => (sidebarHidden = true)}>
					<svelte:fragment slot="icon">
						<UserSolid
							class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
						/>
					</svelte:fragment>
				</SidebarItem>

				<SidebarDropdownWrapper label="Identities">
					<svelte:fragment slot="icon">
						<ShoppingBagSolid
							class="w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white"
						/>
					</svelte:fragment>
					<SidebarDropdownItem label="Users" on:click={() => (sidebarHidden = true)} />
					<SidebarDropdownItem label="Schemas" on:click={() => (sidebarHidden = true)} />
					<SidebarDropdownItem label="Sessions" on:click={() => (sidebarHidden = true)} />
				</SidebarDropdownWrapper>
			</SidebarGroup>
		</SidebarWrapper>
	</Sidebar>
</Drawer>
