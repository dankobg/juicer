import type { PageLoad } from './$types';
import type { RegistrationFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: RegistrationFlow | null = null;

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getRegistrationFlow({
				id: flowIdParam,
			});

			flow = flowResponse.data;
		} catch (error) {
			console.log('registration getRegistrationFlow:', error);
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserRegistrationFlow({
				returnTo,
			});

			if ([403, 404, 410].includes(flowResponse.status)) {
				console.log('registration: [403, 404, 410]');
			}
			if (flowResponse.status !== 200) {
				console.log('registration not status 200');
			}

			flow = flowResponse.data;
		} catch (error) {
			console.log('registration createBrowserRegistrationFlow:', error);
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
