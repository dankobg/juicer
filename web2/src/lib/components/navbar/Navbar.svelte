<script lang="ts">
	import Menu from 'lucide-svelte/icons/menu';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { getInitials } from '$lib/helpers/initials';
	import { config } from '$lib/kratos/config';
	import { page } from '$app/stores';
	import { navItems } from './navitems';
	import type { User } from '$lib/kratos/service';

	export let logoutUrl: string | undefined;
	export let user: User | undefined;
</script>

<header class="bg-background sticky top-0 flex h-16 max-w-[1920px] items-center gap-4 px-4 md:mx-auto md:px-6">
	<nav
		class="hidden w-full flex-col gap-6 text-lg font-medium md:flex md:flex-row md:items-center md:justify-start md:gap-5 md:text-sm lg:gap-6"
	>
		<a href="/" class="flex items-center gap-2 font-semibold md:text-base">
			<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
			<span class="text-lg">Juicer</span>
		</a>
		<div class="flex items-center gap-2">
			{#each navItems as item}
				<a
					href={item.href}
					class="text-muted-foreground hover:text-primary text-lg transition-all"
					class:text-primary={$page.url.pathname === item.href}>{item.label}</a
				>
			{/each}
		</div>
	</nav>
	<Sheet.Root>
		<Sheet.Trigger asChild let:builder>
			<Button variant="outline" size="icon" class="shrink-0 md:hidden" builders={[builder]}>
				<Menu class="h-5 w-5" />
				<span class="sr-only">Toggle navigation menu</span>
			</Button>
		</Sheet.Trigger>
		<Sheet.Content side="left">
			<nav class="grid gap-6 text-lg font-medium">
				<a href="/" class="flex items-center gap-2 text-lg font-semibold">
					<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
					<span class="text-lg">Juicer</span>
				</a>
				{#each navItems as item}
					<a
						href={item.href}
						class="text-muted-foreground hover:text-primary transition-all"
						class:text-primary={$page.url.pathname === item.href}>{item.label}</a
					>
				{/each}
			</nav>
		</Sheet.Content>
	</Sheet.Root>

	<div class="ml-auto flex items-center gap-4 md:gap-2 lg:gap-4">
		{#if user}
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="secondary" size="icon" class="rounded-full">
						<Avatar.Root>
							<Avatar.Image src={user.avatarUrl} alt="user avatar" />
							<Avatar.Fallback class="text-md bg-orange-500/80 font-bold">
								{getInitials(user.fullName ?? user.email)}
							</Avatar.Fallback>
						</Avatar.Root>
						<span class="sr-only">Toggle user menu</span>
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					{#if user.fullName || user.email}
						<DropdownMenu.Label>{user.fullName || user.email}</DropdownMenu.Label>
						<DropdownMenu.Separator />
					{/if}
					<DropdownMenu.Item href="/dashboard">Dashboard</DropdownMenu.Item>
					<DropdownMenu.Item href="/dashboard/account">Account</DropdownMenu.Item>

					{#if logoutUrl}
						<DropdownMenu.Separator />
						<DropdownMenu.Item href={logoutUrl}>Logout</DropdownMenu.Item>
					{/if}
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		{:else}
			<Button href={config.routes.login.path} variant="outline" class="hover:text-primary font-bold">Login</Button>
			<Button href={config.routes.registration.path} class="bg-primary font-bold">Register</Button>
		{/if}
	</div>
</header>
