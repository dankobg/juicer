import IconLayoutDashboard from '@lucide/svelte/icons/layout-dashboard';
import IconUser from '@lucide/svelte/icons/user';
import IconUsers from '@lucide/svelte/icons/users';
import IconMail from '@lucide/svelte/icons/mail';
import IconFingerprint from '@lucide/svelte/icons/fingerprint';
import IconNewspaper from '@lucide/svelte/icons/newspaper';
import IconBuilding1 from '@lucide/svelte/icons/building';
import IconGamepad2 from '@lucide/svelte/icons/gamepad-2';

import type { Component } from 'svelte';

export type NavItem = {
	title: string;
	url?: string;
	isActive?: boolean;
	items?: NavItem[];
	icon?: Component;
};

export const mainNavItems: NavItem[] = [
	{ url: '#', title: 'About' },
	{ url: '#', title: 'Blog' },
	{ url: '#', title: 'Contact' }
];

export const customerDashboardNavItems: NavItem[] = [
	{
		title: 'App',
		url: '#',
		items: [
			{
				title: 'Dashboard',
				url: '/dashboard',
				icon: IconLayoutDashboard
			}
		]
	},
	{
		title: 'User',
		url: '#',
		items: [
			{
				title: 'Account',
				url: '/dashboard/account',
				icon: IconUser
			}
		]
	}
];

export const developerDashboardNavItems: NavItem[] = [
	{
		title: 'App',
		url: '#',
		items: [
			{
				title: 'Dashboard',
				url: '/dashboard',
				icon: IconLayoutDashboard
			}
		]
	},
	{
		title: 'User',
		url: '#',
		items: [
			{
				title: 'Account',
				url: '/dashboard/account',
				icon: IconUser
			}
		]
	},
	{
		title: 'Auth',
		url: '#',
		items: [
			{
				title: 'Schemas',
				url: '/dashboard/schemas',
				icon: IconNewspaper
			},
			{
				title: 'Identities',
				url: '/dashboard/identities',
				icon: IconUsers
			},
			{
				title: 'Sessions',
				url: '/dashboard/sessions',
				icon: IconFingerprint
			},
			{
				title: 'Courier messages',
				url: '/dashboard/messages',
				icon: IconMail
			}
		]
	},
	{
		title: 'Juicer',
		url: '#',
		items: [
			{
				title: 'Games',
				url: '/dashboard/games',
				icon: IconGamepad2
			}
		]
	}
];
