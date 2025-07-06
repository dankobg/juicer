import IconMailX from '@lucide/svelte/icons/mail-x';
import IconRefreshCw from '@lucide/svelte/icons/refresh-cw';
import IconLayers from '@lucide/svelte/icons/layers';
import IconMailCheck from '@lucide/svelte/icons/mail-check';
import IconMail from '@lucide/svelte/icons/mail';
import IconSmartphone from '@lucide/svelte/icons/smartphone';
import IconShieldCheck from '@lucide/svelte/icons/shield-check';
import IconShieldMinus from '@lucide/svelte/icons/shield-minus';
import IconShieldPlus from '@lucide/svelte/icons/shield-plus';
import IconShieldQuestion from '@lucide/svelte/icons/shield-question';
import IconShieldUser from '@lucide/svelte/icons/shield-user';
import IconShieldX from '@lucide/svelte/icons/shield-x';
import IconAsterisk from '@lucide/svelte/icons/asterisk';
import { CourierMessageStatus, CourierMessageTemplateType, CourierMessageType } from '$lib/gen/juicer_openapi';

export const statusIcons = new Map([
	[CourierMessageStatus.Abandoned, IconMailX],
	[CourierMessageStatus.Processing, IconRefreshCw],
	[CourierMessageStatus.Queued, IconLayers],
	[CourierMessageStatus.Sent, IconMailCheck]
]);
export const statuses = Object.values(CourierMessageStatus).map(value => ({
	label: value,
	value,
	icon: statusIcons.get(value)
}));

export const typeIcons = new Map([
	[CourierMessageType.Email, IconMail],
	[CourierMessageType.Phone, IconSmartphone]
]);
export const types = Object.values(CourierMessageType).map(value => ({
	label: value,
	value,
	icon: typeIcons.get(value)
}));

export const templateTypeIcons = new Map([
	[CourierMessageTemplateType.RecoveryInvalid, IconShieldX],
	[CourierMessageTemplateType.RecoveryValid, IconShieldCheck],
	[CourierMessageTemplateType.RecoveryCodeInvalid, IconShieldMinus],
	[CourierMessageTemplateType.RecoveryCodeValid, IconShieldPlus],
	[CourierMessageTemplateType.VerificationInvalid, IconShieldX],
	[CourierMessageTemplateType.VerificationValid, IconShieldCheck],
	[CourierMessageTemplateType.VerificationCodeInvalid, IconShieldMinus],
	[CourierMessageTemplateType.VerificationCodeValid, IconShieldPlus],
	[CourierMessageTemplateType.Stub, IconAsterisk],
	[CourierMessageTemplateType.LoginCodeValid, IconShieldQuestion],
	[CourierMessageTemplateType.RegistrationCodeValid, IconShieldUser]
]);
export const templateTypes = Object.values(CourierMessageTemplateType).map(value => ({
	label: value,
	value,
	icon: templateTypeIcons.get(value)
}));
