import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
	try {
		const [gameVariants, gameTimeKinds, gameTimeCategories, quickGames] = await Promise.all([
			juicer.listGameVariants(),
			juicer.listGameTimeKinds(),
			juicer.listGameTimeCategories(),
			juicer.listQuickGames()
		]);

		return {
			gameVariants,
			gameTimeKinds,
			gameTimeCategories,
			quickGames
		};
	} catch (error) {
		console.log('err', error);
	}
};
