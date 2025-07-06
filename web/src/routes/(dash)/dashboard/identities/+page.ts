import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ depends }) => {
	depends('data:identities');
	try {
		const identities = await juicer.listIdentities({
			pageSize: 1_000
		});
		return {
			identities
		};
	} catch (error) {
		console.log('err', error);
	}
};
