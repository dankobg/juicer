import type { PageLoad } from './$types';
import type { VerificationFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load: PageLoad = (async ({ url }) => {
	const returnToParam = browser && url.searchParams.get('return_to');
	const flowIdParam = browser && url.searchParams.get('flow');

	let flow: VerificationFlow | null = null;

	if (flowIdParam) {
		try {
			const flowResponse = await kratos.getVerificationFlow({
				id: flowIdParam,
			});

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('verification error:', error);

			// switch (err.response?.status) {
			//   case 410:
			//   // Status code 410 means the request has expired - so let's load a fresh flow!
			//   case 403:
			//     // Status code 403 implies some other issue (e.g. CSRF) - let's reload!
			//     return router.push("/verification")
			// }

			flow = null;
		}
	} else {
		const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

		try {
			const flowResponse = await kratos.createBrowserVerificationFlow({
				returnTo,
			});

			// case 400:
			// already signed in
			// goto('/')

			if ([403, 404, 410].includes(flowResponse.status)) {
				console.log('[403, 404, 410].includes(flowResponse.status)');
			}
			if (flowResponse.status !== 200) {
				console.log('flowResponse.status !== 200');
			}

			flow = { ...flowResponse.data };
		} catch (error) {
			console.log('verification error2:', error);
		}
	}

	const csrf = extractCSRFToken(flow);

	return {
		flow,
		csrf,
	};
}) satisfies PageLoad;
