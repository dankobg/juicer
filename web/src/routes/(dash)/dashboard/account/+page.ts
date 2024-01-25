import type { PageLoad } from './$types';
import type { SettingsFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken, isAxiosError } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: SettingsFlow | null = null;

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getSettingsFlow({
				id: flowIdParam,
			});

			console.log('load getSettingsFlow success', flowResponse);

			flow = { ...flowResponse.data };
		} catch (error) {
			if (isAxiosError(error)) {
				const flowData = error.response?.data as SettingsFlow;
				console.log('load getSettingsFlow err:', flowData);
				flow = flowData;
			}
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserSettingsFlow({
				returnTo,
			});

			if (flowResponse.status !== 200) {
				console.log('load createBrowsersettingsFlow status not 200');

				if ([403, 404, 410].includes(flowResponse.status)) {
					console.log('load createBrowsersettingsFlow status [403, 404, 410]');
				}
			}

			console.log('load createBrowsersettingsFlow success', flowResponse);

			flow = { ...flowResponse.data };
		} catch (error) {
			if (isAxiosError(error)) {
				const flowData = error.response?.data as SettingsFlow;
				console.log('load createBrowsersettingsFlow err:', flowData);
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
