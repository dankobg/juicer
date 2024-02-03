import type { PageLoad } from './$types';
import type { GenericError, RecoveryFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken, isAxiosError } from '$lib/kratos/helpers';
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
			const flowResponse = await kratos.getRecoveryFlow({
				id: flowIdParam,
			});
			flow = flowResponse.data;
		} catch (error) {
			if (!isAxiosError(error)) {
				console.error('getRecoveryFlow: unknown error occurred');
				return;
			}

			const err: GenericError = error?.response?.data?.error;

			if (err.id === 'session_already_available') {
				handleFlowErrAction('/', err.message);
			}
			if (err.id === 'self_service_flow_expired') {
				handleFlowErrAction(config.routes.recovery.path, err.message);
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserRecoveryFlow({
				returnTo,
			});
			flow = flowResponse.data;
		} catch (error) {
			if (!isAxiosError(error)) {
				console.error('createBrowserRecoveryFlow: unknown error occurred');
				return;
			}

			const err: GenericError = error?.response?.data?.error;
			handleFlowErrAction(config.routes.recovery.path, err.message);
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
