import { config } from './../../../../lib/kratos/config';
import type { PageLoad } from './$types';
import type { GenericError, RegistrationFlow } from '@ory/client';
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
			const flowResponse = await kratos.getRegistrationFlow({
				id: flowIdParam,
			});
			flow = flowResponse.data;
		} catch (error) {
			const axiosErr = error as AxiosError<GenericError>;
			if (!axiosErr?.isAxiosError) {
				console.error('getLoginFlow: unknown error occurred');
				return;
			}

			const err = axiosErr.response?.data;

			if (err?.id === 'session_already_available') {
				handleFlowErrAction('/', err.message);
			}
			if (err?.id === 'self_service_flow_expired') {
				handleFlowErrAction(config.routes.registration.path, err.message);
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserRegistrationFlow({
				returnTo,
				// afterVerificationReturnTo: ''
			});
			flow = flowResponse.data;
		} catch (error: unknown) {
			const axiosErr = error as AxiosError<GenericError>;
			if (!axiosErr?.isAxiosError) {
				console.error('getLoginFlow: unknown error occurred');
				return;
			}

			const err = axiosErr.response?.data;

			if (err?.id === 'session_already_available') {
				handleFlowErrAction('/', err.message);
			}
			if (err?.id === 'security_csrf_violation' || err?.id === 'security_identity_mismatch') {
				handleFlowErrAction(config.routes.registration.path, err.message);
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
