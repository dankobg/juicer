<script lang="ts">
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { config } from '$lib/kratos/config';
	import type { User } from '$lib/kratos/service';
	import Navbar from 'flowbite-svelte/Navbar.svelte';
	import NavBrand from 'flowbite-svelte/NavBrand.svelte';
	import NavLi from 'flowbite-svelte/NavLi.svelte';
	import NavUl from 'flowbite-svelte/NavUl.svelte';
	import NavHamburger from 'flowbite-svelte/NavHamburger.svelte';
	import Button from 'flowbite-svelte/Button.svelte';
	import Avatar from 'flowbite-svelte/Avatar.svelte';
	import Dropdown from 'flowbite-svelte/Dropdown.svelte';
	import DropdownItem from 'flowbite-svelte/DropdownItem.svelte';
	import DropdownHeader from 'flowbite-svelte/DropdownHeader.svelte';
	import DropdownDivider from 'flowbite-svelte/DropdownDivider.svelte';

	export let logoutUrl: string;
	export let user: User | null;

	$: activeUrl = browser ? $page.url.pathname : '/';
</script>

<Navbar>
	<NavBrand href="/">
		<img src="/images/logo.svg" class="me-3 h-6 sm:h-9" alt="Juicer Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Juicer</span>
	</NavBrand>

	<NavUl class="order-1" {activeUrl}>
		<NavLi href="/">Home</NavLi>
		<NavLi href="/contact">Contact</NavLi>
		<NavLi href="/dashboard">Dashboard</NavLi>
	</NavUl>

	{#if user !== null}
		<div class="flex items-center md:order-2">
			<Avatar id="avatar-menu" src={user.avatarUrl || '/images/logo.svg'} class="cursor-pointer" />
			<NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
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
	{:else}
		<div class="flex order-2 gap-x-2">
			<Button href={config.routes.login.path} outline size="sm">Login</Button>
			<Button href={config.routes.registration.path} size="sm">Sign up</Button>
			<NavHamburger />
		</div>
	{/if}
</Navbar>
