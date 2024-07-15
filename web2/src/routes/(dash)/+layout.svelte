<script lang="ts">
	import LineChart from 'lucide-svelte/icons/line-chart';
	import Package from 'lucide-svelte/icons/package';
	import House from 'lucide-svelte/icons/house';
	import ShoppingCart from 'lucide-svelte/icons/shopping-cart';
	import Bell from 'lucide-svelte/icons/bell';
	import Menu from 'lucide-svelte/icons/menu';
	import Users from 'lucide-svelte/icons/users';
	import Logout from 'lucide-svelte/icons/log-out';
	import ChevronDown from 'lucide-svelte/icons/chevron-down';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as Collapsible from '$lib/components/ui/collapsible';
</script>

<div class="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[280px_1fr]">
	<div class="bg-muted/40 hidden border-r md:block">
		<div class="flex h-full max-h-screen flex-col gap-2">
			<div class="flex h-14 items-center border-b px-4 lg:h-[60px] lg:px-6">
				<a href="/" class="flex items-center gap-2 font-semibold">
					<img src="/images/logo.svg" alt="logo" class="h-8 w-8 object-cover" />
					<span class="">Juicer</span>
				</a>
				<Button variant="outline" size="icon" class="ml-auto h-8 w-8">
					<Bell class="h-4 w-4" />
					<span class="sr-only">Toggle notifications</span>
				</Button>
			</div>
			<div class="flex-1">
				<nav class="grid items-start px-2 text-sm font-medium lg:px-4">
					<a
						href="##"
						class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
					>
						<House class="h-4 w-4" />
						Dashboard
					</a>
					<a
						href="##"
						class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
					>
						<ShoppingCart class="h-4 w-4" />
						Orders
						<Badge class="ml-auto flex h-6 w-6 shrink-0 items-center justify-center rounded-full">6</Badge>
					</a>
					<a
						href="##"
						class="bg-muted text-primary hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
					>
						<Package class="h-4 w-4" />
						Products
					</a>
					<a
						href="##"
						class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
					>
						<Users class="h-4 w-4" />
						Customers
					</a>
					<a
						href="##"
						class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
					>
						<LineChart class="h-4 w-4" />
						Analytics
					</a>

					<Collapsible.Root class="w-full">
						<Collapsible.Trigger class="w-full">
							<span
								class="text-muted-foreground hover:text-primary flex w-full items-center gap-3 rounded-lg px-3 py-2 transition-all"
							>
								<LineChart class="h-4 w-4" />
								Nested
								<ChevronDown class="ml-auto h-4 w-4" />
							</span>
						</Collapsible.Trigger>
						<Collapsible.Content>
							<a
								href="##"
								class="bg-muted text-primary hover:text-primary ml-8 flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
							>
								First
							</a>
							<a
								href="##"
								class="text-muted-foreground hover:text-primary ml-8 flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
							>
								Second
							</a>
							<a
								href="##"
								class="text-muted-foreground hover:text-primary ml-8 flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
							>
								Third
							</a>
						</Collapsible.Content>
					</Collapsible.Root>
				</nav>
			</div>
			<div class="mt-auto p-4">
				<a
					href="##"
					class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 transition-all"
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
						<a
							href="##"
							class="text-muted-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
						>
							<House class="h-5 w-5" />
							Dashboard
						</a>
						<a
							href="##"
							class="bg-muted text-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
						>
							<ShoppingCart class="h-5 w-5" />
							Orders
							<Badge class="ml-auto flex h-6 w-6 shrink-0 items-center justify-center rounded-full">6</Badge>
						</a>
						<a
							href="##"
							class="text-muted-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
						>
							<Package class="h-5 w-5" />
							Products
						</a>
						<a
							href="##"
							class="text-muted-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
						>
							<Users class="h-5 w-5" />
							Customers
						</a>
						<a
							href="##"
							class="text-muted-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
						>
							<LineChart class="h-5 w-5" />
							Analytics
						</a>

						<Collapsible.Root class="w-full">
							<Collapsible.Trigger class="w-full">
								<span
									class="text-muted-foreground hover:text-foreground mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2"
								>
									<LineChart class="h-4 w-4" />
									Nested
									<ChevronDown class="ml-auto h-4 w-4" />
								</span>
							</Collapsible.Trigger>
							<Collapsible.Content>
								<a
									href="##"
									class="text-muted-foreground hover:text-foreground mx-[-0.65rem] ml-8 flex items-center gap-4 rounded-xl px-3 py-2"
								>
									First
								</a>
								<a
									href="##"
									class="text-muted-foreground hover:text-foreground mx-[-0.65rem] ml-8 flex items-center gap-4 rounded-xl px-3 py-2"
								>
									Second
								</a>
								<a
									href="##"
									class="text-muted-foreground hover:text-foreground mx-[-0.65rem] ml-8 flex items-center gap-4 rounded-xl px-3 py-2"
								>
									Third
								</a>
							</Collapsible.Content>
						</Collapsible.Root>
					</nav>
					<div class="mt-auto">
						<a
							href="##"
							class="text-muted-foreground hover:text-primary flex items-center gap-3 rounded-lg px-3 py-2 text-lg font-medium transition-all"
						>
							<Logout class="h-4 w-4" />
							Logout
						</a>
					</div>
				</Sheet.Content>
			</Sheet.Root>
			<div class="w-full flex-1"></div>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild let:builder>
					<Button builders={[builder]} variant="secondary" size="icon" class="rounded-full">
						<Avatar.Root>
							<Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" />
							<Avatar.Fallback>CN</Avatar.Fallback>
						</Avatar.Root>
						<span class="sr-only">Toggle user menu</span>
					</Button>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end">
					<DropdownMenu.Label>My Account</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item>Settings</DropdownMenu.Item>
					<DropdownMenu.Item>Support</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item>Logout</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</header>
		<main class="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
			<div class="flex items-center">
				<h1 class="text-lg font-semibold md:text-2xl">Inventory</h1>
			</div>
			<div class="flex flex-1 items-center justify-center rounded-lg border border-dashed shadow-sm">
				<div class="flex flex-col items-center gap-1 text-center">
					<h3 class="text-2xl font-bold tracking-tight">You have no products</h3>
					<p class="text-muted-foreground text-sm">You can start selling as soon as you add a product.</p>
					<Button class="mt-4">Add Product</Button>
				</div>
			</div>
		</main>
	</div>
</div>
