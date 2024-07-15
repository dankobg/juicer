import type { LoginFlow, RecoveryFlow, RegistrationFlow, SettingsFlow, VerificationFlow } from '@ory/client-fetch';

type Provider = {
	name: string;
	label: string;
};
export const providers: Provider[] = [
	{ name: 'google', label: 'Google' },
	{ name: 'github', label: 'GitHub' },
	{ name: 'facebook', label: 'Facebook' },
	{ name: 'discord', label: 'Discord' },
	{ name: 'twitch', label: 'Twitch' },
	{ name: 'slack', label: 'Slack' },
	{ name: 'spotify', label: 'Spotify' }
];

export type KratosFlow = LoginFlow | RegistrationFlow | RecoveryFlow | VerificationFlow | SettingsFlow;

export function extractCSRFToken(flow: KratosFlow | null): string {
	if (!flow) {
		return '';
	}

	const csrfAttributes = flow?.ui?.nodes?.find(node => {
		return node.attributes.node_type === 'input' && node.group === 'default' && node.attributes.name === 'csrf_token';
	})?.attributes;

	return csrfAttributes?.node_type === 'input' ? csrfAttributes.value : '';
}
