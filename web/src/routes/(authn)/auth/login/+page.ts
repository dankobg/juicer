import type { PageLoad } from './$types';
import type { GenericError, LoginFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken, isAxiosError } from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');
	const refreshParam = browser && url.searchParams.get('refresh');
	const aalParam = browser && url.searchParams.get('aal');
	const loginChallengeParam = browser && url.searchParams.get('login_challenge');

	let flow: LoginFlow | null = null;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		flow = null;
		goto(redirectUrl);
		return;
	}

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getLoginFlow({
				id: flowIdParam,
			});
			flow = flowResponse.data;
		} catch (error) {
			if (!isAxiosError(error)) {
				console.error('getLoginFlow: unknown error occurred');
				return;
			}

			const err: GenericError = error?.response?.data?.error;

			if (err.id === 'session_already_available') {
				handleFlowErrAction('/', err.message);
			}
			if (err.id === 'self_service_flow_expired') {
				handleFlowErrAction(config.routes.login.path, err.message);
			}
		}
	} else {
		const aal: string | undefined = aalParam ? aalParam.toString() : undefined;
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;
		const refresh: boolean | undefined = refreshParam?.toString().toLowerCase() === 'true' ? true : undefined;
		const loginChallenge: string | undefined = loginChallengeParam ? loginChallengeParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserLoginFlow({
				aal,
				returnTo,
				refresh,
				loginChallenge,
			});
			flow = flowResponse.data;
		} catch (error) {
			if (!isAxiosError(error)) {
				console.error('createBrowserLoginFlow: unknown error occurred');
				return;
			}

			const err: GenericError = error?.response?.data?.error;

			if (err.id === 'session_already_available') {
				handleFlowErrAction('/', err.message);
			}
			if (err.id === 'security_csrf_violation') {
				handleFlowErrAction(config.routes.login.path, err.message);
			}
			if (err.id === 'session_aal1_required') {
				goto('/login?aal=aal1&return_to=' + window.location.href);
				return;
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
