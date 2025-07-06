import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params }) => {
	try {
		const message = await juicer.getCourierMessage({ id: params.message_id });
		return {
			message
		};
	} catch (error) {
		console.log('err', error);
	}
};
