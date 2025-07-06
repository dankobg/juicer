import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const schema = await juicer.getIdentitySchema({ id: params.schema_id });
		return {
			id: params.schema_id,
			schema
		};
	} catch (error) {
		console.log('err', error);
	}
};
