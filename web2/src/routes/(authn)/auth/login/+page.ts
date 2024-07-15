import type { PageLoad } from './$types';
import { instanceOfGenericError, type LoginFlow } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
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

		if (browser) {
			goto(redirectUrl);
		}

		return;
	}

	if (flowIdParam) {
		try {
			const loginFlow = await kratos.getLoginFlow({
				id: flowIdParam
			});
			flow = loginFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (instanceOfGenericError(error)) {
				if (error.id === 'session_already_available') {
					handleFlowErrAction('/', error.message);
				}
				if (error.id === 'self_service_flow_expired') {
					handleFlowErrAction(config.routes.login.path, error.message);
				}
			}
		}
	} else {
		const aal: string | undefined = aalParam ? aalParam.toString() : undefined;
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;
		const refresh: boolean | undefined = refreshParam?.toString().toLowerCase() === 'true' ? true : undefined;
		const loginChallenge: string | undefined = loginChallengeParam ? loginChallengeParam.toString() : undefined;

		try {
			const loginFlow = await kratos.createBrowserLoginFlow({
				aal,
				returnTo,
				refresh,
				loginChallenge
			});
			flow = loginFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}
			if (instanceOfGenericError(error)) {
				if (error.id === 'security_csrf_violation') {
					handleFlowErrAction(config.routes.login.path, error.message);
				}
				if (error.id === 'session_aal2_required') {
					if (browser) {
						goto('/login?aal=aal2&return_to=' + window.location.href);
					}
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
