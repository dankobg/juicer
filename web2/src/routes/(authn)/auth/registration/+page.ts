import { config } from './../../../../lib/kratos/config';
import type { PageLoad } from './$types';
import { instanceOfGenericError, type RegistrationFlow } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { toast } from 'svelte-sonner';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: RegistrationFlow | null = null;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		flow = null;

		if (browser) {
			goto(redirectUrl);
		}

		return;
	}

	if (flowIdParam) {
		try {
			const registrationFlow = await kratos.getRegistrationFlow({
				id: flowIdParam
			});
			flow = registrationFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (instanceOfGenericError(error)) {
				if (error.id === 'session_already_available') {
					handleFlowErrAction('/', error.message);
				} else if (error.id === 'self_service_flow_expired') {
					handleFlowErrAction(config.routes.registration.path, error.message);
				}
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const registrationFlow = await kratos.createBrowserRegistrationFlow({
				returnTo
				// afterVerificationReturnTo: ''
			});
			flow = registrationFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (instanceOfGenericError(error)) {
				if (error.id === 'session_already_available') {
					handleFlowErrAction('/', error.message);
				} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
					handleFlowErrAction(config.routes.registration.path, error.message);
				}
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf
	};
}) satisfies PageLoad;
