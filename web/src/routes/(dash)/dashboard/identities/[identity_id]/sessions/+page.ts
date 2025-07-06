import type { ListIdentitySessionsRequest } from '$lib/gen/juicer_openapi';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, params, depends }) => {
	depends(`data:identity-sessions-${params.identity_id}`);
	try {
		const req: ListIdentitySessionsRequest = {
			id: params.identity_id,
			pageSize: 1_000
		};
		const active = url.searchParams.get('active');
		if (active) {
			req.active = active === 'true' ? true : active === 'false' ? false : false;
		}
		const pageToken = url.searchParams.get('page_token');
		if (pageToken) {
			req.pageToken = pageToken;
		}
		const sessions = await juicer.listIdentitySessions(req);
		const identity = await juicer.getIdentity({ id: params.identity_id });
		return {
			identity,
			sessions
		};
	} catch (error) {
		console.log('err', error);
	}
};
