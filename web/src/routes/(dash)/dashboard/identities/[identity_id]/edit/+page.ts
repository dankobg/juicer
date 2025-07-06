import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, depends }) => {
	depends(`data:identity-${params.identity_id}`);
	try {
		const schemas = await juicer.listIdentitySchemas({
			pageSize: 1_000
		});
		const identity = await juicer.getIdentity({
			id: params.identity_id
		});
		return {
			schemas,
			identity
		};
	} catch (error) {
		console.log('err', error);
	}
};
