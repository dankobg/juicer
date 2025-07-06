import type { ListSessionsExpandEnum, ListSessionsRequest } from '$lib/gen/juicer_openapi';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, depends }) => {
	depends('data:sessions');
	try {
		const req: ListSessionsRequest = {
			pageSize: 1_000
		};
		const active = url.searchParams.get('active');
		if (active) {
			req.active = active === 'true' ? true : active === 'false' ? false : false;
		}
		const expand = url.searchParams.getAll('expand');
		if (expand.length > 0) {
			req.expand = expand as ListSessionsExpandEnum[];
		}
		const pageToken = url.searchParams.get('page_token');
		if (pageToken) {
			req.pageToken = pageToken;
		}
		const sessions = await juicer.listSessions(req);
		return {
			sessions
		};
	} catch (error) {
		console.log('err', error);
	}
};
