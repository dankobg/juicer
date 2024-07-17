import type { PageLoad } from './$types';
import { instanceOfGenericError, type GenericError, type SettingsFlow } from '@ory/client-fetch';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: SettingsFlow | null = null;

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
			const settingsFlow = await kratos.getSettingsFlow({
				id: flowIdParam
			});
			flow = settingsFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}

			if (instanceOfGenericError(error)) {
				if (error.id === 'session_inactive' || error.id === 'session_refresh_required') {
					handleFlowErrAction(
						config.routes.login.path + `?return_to=${encodeURIComponent(window.location.href)}`,
						error.message
					);
				} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
					handleFlowErrAction(config.routes.settings.path, error.message);
				}
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const settingsFlow = await kratos.createBrowserSettingsFlow({
				returnTo
			});
			flow = settingsFlow;
		} catch (error: unknown) {
			if (!error || typeof error !== 'object') {
				return;
			}

			if (instanceOfGenericError(error)) {
				if (error.id === 'session_inactive') {
					handleFlowErrAction(
						config.routes.login.path + `?return_to=${encodeURIComponent(window.location.href)}`,
						error.message
					);
				} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
					handleFlowErrAction(config.routes.settings.path, error.message);
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
