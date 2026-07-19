import { AuthenticatorAssuranceLevel } from '$lib/gen/juicer_openapi';

export const aals = Object.values(AuthenticatorAssuranceLevel).map(value => ({
	label: value,
	value
}));
