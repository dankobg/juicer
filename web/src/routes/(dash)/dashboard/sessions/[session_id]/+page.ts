import type { GetSessionExpandEnum, GetSessionRequest } from '$lib/gen/juicer_openapi';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, url }) => {
	try {
		const req: GetSessionRequest = {
			id: params.session_id,
			expand: ['identity', 'devices']
		};
		const expand = url.searchParams.getAll('expand');
		if (expand.length > 0) {
			req.expand = expand as GetSessionExpandEnum[];
		}
		const session = await juicer.getSession(req);
		return {
			session
		};
	} catch (error) {
		console.log('err', error);
	}
};
