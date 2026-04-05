import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { juicer } from '$lib/juicer/client';
import { config } from '$lib/kratos/config';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params, depends }) => {
	depends(`data:dashboard-identities-${params.identity_id}-update`);

	try {
		const schemasResult = await juicer.GET('/schemas', {
			fetch,
			params: {
				query: { page_size: 500 }
			}
		});
		const identityResult = await juicer.GET('/identities/{id}', {
			fetch,
			params: {
				path: { id: params.identity_id }
			}
		});

		if (schemasResult.error?.status_code === 403 || identityResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			schemasResult,
			identityResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
