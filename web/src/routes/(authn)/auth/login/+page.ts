import type { PageLoad } from './$types';
import type { LoginFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken, isAxiosError } from '$lib/kratos/helpers';
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

			console.log('load getLoginFlow success', flowResponse);

			flow = { ...flowResponse.data };
		} catch (error) {
			if (isAxiosError(error)) {
				const flowData = error.response?.data as LoginFlow;
				console.log('load getLoginFlow err:', flowData);
				flow = flowData;
			}
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

			if (flowResponse.status !== 200) {
				console.log('load createBrowserLoginFlow status not 200');

				if ([403, 404, 410].includes(flowResponse.status)) {
					console.log('load createBrowserLoginFlow status [403, 404, 410]');
				}
			}

			console.log('load createBrowserLoginFlow success', flowResponse);

			flow = { ...flowResponse.data };
		} catch (error) {
			if (isAxiosError(error)) {
				const flowData = error.response?.data as LoginFlow;
				console.log('load createBrowserLoginFlow err:', flowData);
				flow = flowData;
			}
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
