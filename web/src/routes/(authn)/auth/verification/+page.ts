import type { PageLoad } from './$types';
import type { GenericError, VerificationFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';
import type { AxiosError } from 'axios';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: VerificationFlow | null = null;

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
			const flowResponse = await kratos.getVerificationFlow({
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
				handleFlowErrAction(config.routes.verification.path, err.message);
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserVerificationFlow({
				returnTo,
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
			if (err?.id === 'security_csrf_violation' || err?.id === 'security_identity_mismatch') {
				handleFlowErrAction(config.routes.verification.path, err.message);
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
