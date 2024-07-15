import type { PageLoad } from './$types';
import { instanceOfGenericError, type RecoveryFlow } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: RecoveryFlow | null = null;

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
			const recoveryFlow = await kratos.getRecoveryFlow({
				id: flowIdParam
			});
			flow = recoveryFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}

			if (instanceOfGenericError(error)) {
				if (error.id === 'session_already_available') {
					handleFlowErrAction('/', error.message);
				}
				if (error.id === 'self_service_flow_expired') {
					handleFlowErrAction(config.routes.recovery.path, error.message);
				}
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const recoveryFlow = await kratos.createBrowserRecoveryFlow({
				returnTo
			});
			flow = recoveryFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}

			if (instanceOfGenericError(error)) {
				handleFlowErrAction(config.routes.recovery.path, error.message);
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf
	};
}) satisfies PageLoad;
