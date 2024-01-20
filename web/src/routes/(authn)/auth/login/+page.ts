import type { PageLoad } from './$types';
import type { LoginFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');
	const refreshParam = browser && url.searchParams.get('refresh');
	const aalParam = browser && url.searchParams.get('aal');

	let flow: LoginFlow | null = null;

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getLoginFlow({
				id: flowIdParam,
			});

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('login error:', error);

			flow = null;
		}
	} else {
		const aal: string | undefined = aalParam ? aalParam.toString() : undefined;
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;
		const refresh: boolean | undefined = refreshParam?.toString().toLowerCase() === 'true' ? true : undefined;

		try {
			const flowResponse = await kratos.createBrowserLoginFlow({
				aal,
				returnTo,
				refresh,
			});

			if ([403, 404, 410].includes(flowResponse.status)) {
				console.log('login createBrowserLoginFlow: [403, 404, 410]');
			}
			if (flowResponse.status !== 200) {
				console.log('login not 200');
			}

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('login:', error);
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
