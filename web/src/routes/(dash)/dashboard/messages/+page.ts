import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const messages = await juicer.listCourierMessages({
			pageSize: 1_000
		});
		return {
			messages
		};
	} catch (error) {
		console.log('err', error);
	}
};
