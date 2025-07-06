import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent, depends }) => {
	const data = await parent();
	depends('data:dashboard');
	try {
		const gameStats = await juicer.getGameStats({
			userId: data.auth.user?.id ?? ''
		});
		return {
			...data,
			gameStats
		};
	} catch (error) {
		console.log('err', error);
		return {
			...data,
			gameStats: null
		};
	}
};
