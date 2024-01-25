import type { PageLoad } from './$types';
import type { RecoveryFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: RecoveryFlow | null = null;

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getRecoveryFlow({
				id: flowIdParam,
			});

			console.log('load getRecoveryFlow', flowResponse);

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('load getRecoveryFlow', error);

			flow = null;
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserRecoveryFlow({
				returnTo,
			});

			if (flowResponse.status !== 200) {
				console.log('load createBrowserRecoveryFlow status not 200');

				if ([403, 404, 410].includes(flowResponse.status)) {
					console.log('load createBrowserRecoveryFlow status [403, 404, 410]');
				}
			}

			console.log('load createBrowserRecoveryFlow', flowResponse);
			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('load createBrowserRecoveryFlow', error);
			// case: 400 -> validation err
			// setFlow(err.resp.data.flow)
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
