import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const identity = await juicer.getIdentity({ id: params.identity_id });
		return {
			identity
		};
	} catch (error) {
		console.log('err', error);
	}
};
