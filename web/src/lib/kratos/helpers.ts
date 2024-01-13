import type { LoginFlow, RecoveryFlow, RegistrationFlow, SettingsFlow, VerificationFlow } from '@ory/client';

type KratosFlow = LoginFlow | RegistrationFlow | RecoveryFlow | VerificationFlow | SettingsFlow;

export function extractCSRFToken(flow: KratosFlow | null): string {
	if (!flow) {
		return '';
	}

	const csrfAttributes = flow?.ui?.nodes?.find(node => {
		return node.attributes.node_type === 'input' && node.group === 'default' && node.attributes.name === 'csrf_token';
	})?.attributes;

	return csrfAttributes?.node_type === 'input' ? csrfAttributes.value : '';
}
