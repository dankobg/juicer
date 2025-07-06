import IconCheck from '@lucide/svelte/icons/check';
import IconX from '@lucide/svelte/icons/x';
import IconShieldUser from '@lucide/svelte/icons/shield-user';
import IconUser from '@lucide/svelte/icons/user';
import { IdentityStateEnum } from '$lib/gen/juicer_openapi';

export const stateIcons = new Map([
	[IdentityStateEnum.Active, IconCheck],
	[IdentityStateEnum.Inactive, IconX]
]);
export const states = Object.values(IdentityStateEnum).map(value => ({
	label: value,
	value,
	icon: stateIcons.get(value)
}));

export const schemaIdIcons = new Map([
	['customer', IconUser],
	['employee', IconShieldUser]
]);
export const schemaIds = ['customer', 'employee'].map(value => ({
	label: value,
	value,
	icon: schemaIdIcons.get(value)
}));
