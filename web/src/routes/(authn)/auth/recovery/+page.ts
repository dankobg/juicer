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

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('recovery error:', error);

			flow = null;
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserRecoveryFlow({
				returnTo,
			});

			console.log('load - recovery flowResponse', flowResponse);

			if ([403, 404, 410].includes(flowResponse.status)) {
				console.log('[403, 404, 410].includes(flowResponse.status)');
			}
			if (flowResponse.status !== 200) {
				console.log('flowResponse.status !== 200');
			}

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('recovery error2:', error);
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
