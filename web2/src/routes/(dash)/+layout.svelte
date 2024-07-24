<script lang="ts">
	import BookOpen from 'lucide-svelte/icons/book-open';
	import Package from 'lucide-svelte/icons/package';
	import House from 'lucide-svelte/icons/house';
	import ShoppingCart from 'lucide-svelte/icons/shopping-cart';
	import Bell from 'lucide-svelte/icons/bell';
	import Menu from 'lucide-svelte/icons/menu';
	import Users from 'lucide-svelte/icons/users';
	import Logout from 'lucide-svelte/icons/log-out';
	import ChevronDown from 'lucide-svelte/icons/chevron-down';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Popover from '$lib/components/ui/popover';
	import Check from 'lucide-svelte/icons/check';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import * as Drawer from '$lib/components/ui/drawer/index.js';
	import { mediaQuery } from 'svelte-legos';
	import { page } from '$app/stores';
	import type { LayoutData } from '../$types';
	import { getInitials } from '$lib/helpers/initials';

	export let data: LayoutData;

	const sidebarItems = [
		{ href: '/dashboard', label: 'Dashboard', icon: House },
		{ href: '/dashboard/identities', label: 'Identities', icon: Users },
		{ href: '/dashboard/schemas', label: 'Schemas', icon: ShoppingCart },
		{ href: '/dashboard/sessions', label: 'Sessions', icon: Package },
		{
			label: 'Nested',
			pattern: '/dashboard/nested',
			items: [
				{ href: '/dashboard/nested/first', label: 'First', icon: BookOpen },
				{ href: '/dashboard/nested/second', label: 'Second', icon: BookOpen },
				{ href: '/dashboard/nested/third', label: 'Third', icon: BookOpen }
			]
		}
	];

	const notifications = [
		{
			title: 'Your call has been confirmed.',
			description: '1 hour ago'
		},
		{
			title: 'You have a new message!',
			description: '1 hour ago'
		},
		{
			title: 'Your subscription is expiring soon!',
			description: '2 hours ago'
		}
	];

	const ITEMS_TO_DISPLAY = 5;
	const isDesktop = mediaQuery('(min-width: 768px)');
	let open = false;
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
				href: tokenPath
			};
		});

		crumbs.unshift({ label: 'Home', href: '/' });
	}
</script>

<div class="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
	<div class="bg-muted/40 hidden border-r md:block">
		<div class="flex h-full max-h-screen flex-col gap-2">
			<div class="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
				<a href="/" class="flex items-center gap-2 font-semibold">
					<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
					<span class="">Juicer</span>
				</a>
			</div>
			<div class="flex-1">
				<nav class="grid items-start px-2 text-sm font-medium lg:px-4">
					{#each sidebarItems as item}
						{#if item.items}
							<Collapsible.Root class="w-full">
								<Collapsible.Trigger class="w-full">
									<span
										class="text-muted-foreground hover:text-primary flex w-full items-center gap-3 rounded-lg px-3 py-2 transition-all transition-all"
										class:text-primary={$page.url.pathname.startsWith(item.pattern)}
										class:bg-muted={$page.url.pathname.startsWith(item.pattern)}
									>
										<BookOpen class="h-4 w-4" />
										{item.label}
										<ChevronDown class="ml-auto h-4 w-4" />
									</span>
								</Collapsible.Trigger>
								<Collapsible.Content>
									{#each item.items as nested}
										<a
											href={nested.href}
											class="text-muted-foreground hover:text-primary ml-8 flex items-center gap-3 rounded-lg px-3 py-2 transition-all transition-all"
											class:text-primary={$page.url.pathname === nested.href}
											class:bg-muted={$page.url.pathname === nested.href}
										>
											{nested.label}
										</a>
									{/each}
								</Collapsible.Content>
							</Collapsible.Root>
						{:else}
							<a
								href={item.href}
								class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all transition-all"
								class:text-primary={$page.url.pathname === item.href}
								class:bg-muted={$page.url.pathname === item.href}
							>
								<svelte:component this={item.icon} class="h-4 w-4" />
								{item.label}
							</a>
						{/if}
					{/each}
				</nav>
			</div>
			<div class="mt-auto p-4">
				<a
					href={data.logoutUrl}
					class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all transition-all"
				>
					<Logout class="h-4 w-4" />
					Logout
				</a>
			</div>
		</div>
	</div>
	<div class="flex flex-col">
		<header class="bg-muted/40 flex h-14 items-center gap-4 border-b px-4 lg:h-[60px] lg:px-6">
			<Sheet.Root>
				<Sheet.Trigger asChild let:builder>
					<Button variant="outline" size="icon" class="shrink-0 md:hidden" builders={[builder]}>
						<Menu class="h-5 w-5" />
						<span class="sr-only">Toggle navigation menu</span>
					</Button>
				</Sheet.Trigger>
				<Sheet.Content side="left" class="flex flex-col">
					<nav class="grid gap-2 text-lg font-medium">
						<a href="##" class="flex items-center gap-2 text-lg font-semibold">
							<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
							<span class="sr-only">Juicer</span>
						</a>

						{#each sidebarItems as item}
							{#if item.items}
								<Collapsible.Root class="w-full">
									<Collapsible.Trigger class="w-full">
										<span
											class="text-muted-foreground hover:text-primary mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 transition-all"
										>
											<BookOpen class="h-4 w-4" />
											{item.label}
											<ChevronDown class="ml-auto h-4 w-4" />
										</span>
									</Collapsible.Trigger>
									<Collapsible.Content>
										{#each item.items as nested}
											<a
												href={nested.href}
												class="text-muted-foreground hover:text-primary mx-[-0.65rem] ml-8 flex items-center gap-4 rounded-xl px-3 py-2 transition-all"
												class:text-primary={$page.url.pathname === nested.href}
											>
												{nested.label}
											</a>
										{/each}
									</Collapsible.Content>
								</Collapsible.Root>
							{:else}
								<a
									href={item.href}
									class="text-muted-foreground hover:text-primary mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 transition-all"
									class:text-primary={$page.url.pathname === item.href}
								>
									<svelte:component this={item.icon} class="h-4 w-4" />
									{item.label}
								</a>
							{/if}
						{/each}
					</nav>
					<div class="mt-auto">
						<a
							href={data.logoutUrl}
							class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 text-lg font-medium transition-all transition-all"
						>
							<Logout class="h-4 w-4" />
							Logout
						</a>
					</div>
				</Sheet.Content>
			</Sheet.Root>
			<div class="w-full flex-1"></div>

			<Popover.Root>
				<Popover.Trigger>
					<Button variant="outline" size="icon" class="ml-auto h-8 w-8">
						<Bell class="h-4 w-4" />
						<span class="sr-only">Toggle notifications</span>
					</Button>
				</Popover.Trigger>
				<Popover.Content>
					<div class="grid gap-4">
						<div class="space-y-2">
							<h4 class="font-bold leading-none">Notifications</h4>
							<p class="text-muted-foreground text-sm">You have 3 unread messages</p>
						</div>
						<div class="grid gap-2">
							<div>
								{#each notifications as notification, idx (idx)}
									<div class="mb-4 grid grid-cols-[25px_1fr] items-start pb-4 last:mb-0 last:pb-0">
										<span class="flex h-2 w-2 translate-y-1 rounded-full bg-sky-500" />
										<div class="space-y-1">
											<p class="text-sm font-medium leading-none">
												{notification.title}
											</p>
											<p class="text-muted-foreground text-sm">
												{notification.description}
											</p>
										</div>
									</div>
								{/each}
							</div>
						</div>
						<Button class="w-full">
							<Check class="mr-2 h-4 w-4" /> Mark all as read
						</Button>
					</div>
				</Popover.Content>
			</Popover.Root>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="secondary" size="icon" class="rounded-full">
						<Avatar.Root>
							<Avatar.Image src={data.auth.user?.avatarUrl} alt="user avatar" />
							<Avatar.Fallback class="text-md bg-orange-500/80 font-bold">
								{getInitials(data.auth.user?.fullName || data.auth.user?.email)}
							</Avatar.Fallback>
						</Avatar.Root>
						<span class="sr-only">Toggle user menu</span>
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					{#if data.auth.user?.email}
						<DropdownMenu.Label>{data.auth.user.email}</DropdownMenu.Label>
						<DropdownMenu.Separator />
					{/if}

					{#if data.auth.user?.fullName || data.auth.user?.email}
						<DropdownMenu.Label>{data.auth.user?.fullName || data.auth.user?.email}</DropdownMenu.Label>
						<DropdownMenu.Separator />
					{/if}

					<DropdownMenu.Item href="/dashboard/account">Account</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item href={data.logoutUrl}>Logout</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</header>
		<main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
			<Breadcrumb.Root>
				<Breadcrumb.List>
					{#if crumbs.length > ITEMS_TO_DISPLAY}
						<Breadcrumb.Item>
							{#if $isDesktop}
								<DropdownMenu.Root bind:open>
									<DropdownMenu.Trigger class="flex items-center gap-1" aria-label="Toggle menu">
										<Breadcrumb.Ellipsis class="h-4 w-4" />
									</DropdownMenu.Trigger>
									<DropdownMenu.Content align="start">
										{#each crumbs.slice(1, -2) as crumb}
											<DropdownMenu.Item href={crumb.href ? crumb.href : '#'}>
												{crumb.label}
											</DropdownMenu.Item>
										{/each}
									</DropdownMenu.Content>
								</DropdownMenu.Root>
							{:else}
								<Drawer.Root bind:open>
									<Drawer.Trigger aria-label="Toggle Menu">
										<Breadcrumb.Ellipsis class="h-4 w-4" />
									</Drawer.Trigger>
									<Drawer.Content>
										<Drawer.Header class="text-left">
											<Drawer.Title>Navigate to</Drawer.Title>
											<Drawer.Description>Select a page to navigate to.</Drawer.Description>
										</Drawer.Header>
										<div class="grid gap-1 px-4">
											{#each crumbs.slice(1, -2) as crumb}
												<a href={crumb.href ? crumb.href : '#'} class="py-1 text-sm">
													{crumb.label}
												</a>
											{/each}
										</div>
										<Drawer.Footer class="pt-4">
											<Drawer.Close asChild let:builder>
												<Button variant="outline" builders={[builder]}>Close</Button>
											</Drawer.Close>
										</Drawer.Footer>
									</Drawer.Content>
								</Drawer.Root>
							{/if}
						</Breadcrumb.Item>
						<Breadcrumb.Separator />
					{/if}

					{#each crumbs.slice(-ITEMS_TO_DISPLAY + 1) as crumb, idx}
						<Breadcrumb.Item>
							{#if crumb.href}
								<Breadcrumb.Link href={crumb.href} class="max-w-20 truncate md:max-w-none">
									{crumb.label}
								</Breadcrumb.Link>
								{#if idx !== crumbs.length - 1}
									<Breadcrumb.Separator />
								{/if}
							{:else}
								<Breadcrumb.Page class="max-w-20 truncate md:max-w-none">
									{crumb.label}
								</Breadcrumb.Page>
							{/if}
						</Breadcrumb.Item>
					{/each}
				</Breadcrumb.List>
			</Breadcrumb.Root>

			<slot />
		</main>
	</div>
</div>
