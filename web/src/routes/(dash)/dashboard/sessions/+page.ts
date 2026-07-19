import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';
import type { operations, PathsSessionsGetParametersQueryExpand } from '$lib/gen/juicer_openapi';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';

export const load: PageLoad = async ({ fetch, url, depends }) => {
	depends('data:dashboard-sessions');

	try {
		const listSessionsParams: operations['listSessions']['parameters'] = {
			query: { page_size: 500 }
		};
		const active = url.searchParams.get('active');
		if (active) {
			listSessionsParams.query!.active = active === 'true' ? true : active === 'false' ? false : false;
		}
		const expand = url.searchParams.getAll('expand') as PathsSessionsGetParametersQueryExpand[];
		if (expand.length > 0) {
			listSessionsParams.query!.expand = expand;
		}
		const pageToken = url.searchParams.get('page_token');
		if (pageToken) {
			listSessionsParams.query!.page_token = pageToken;
		}
		const sessionsResult = await juicer.GET('/sessions', {
			fetch,
			params: listSessionsParams
		});

		if (sessionsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			sessionsResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
