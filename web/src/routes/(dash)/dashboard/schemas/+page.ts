import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const schemas = await juicer.listIdentitySchemas({
			pageSize: 1_000
		});
		return {
			schemas
		};
	} catch (error) {
		console.log('err', error);
	}
};
