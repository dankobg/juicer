import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends }) => {
	depends(`data:home`);

	try {
		const gameVariantsResult = await juicer.GET('/game-variants', {
			fetch
		});

		if (gameVariantsResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		const gameTimeCategoriesResult = await juicer.GET('/game-time-categories', {
			fetch
		});

		if (gameTimeCategoriesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		const gameTimeKindsResult = await juicer.GET('/game-time-kinds', {
			fetch
		});

		if (gameTimeKindsResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		const quickGamesResult = await juicer.GET('/quick-games', {
			fetch
		});

		if (quickGamesResult.error?.status_code === 403) {
			if (browser) {
				goto('/');
			}
		}

		return {
			gameVariantsResult,
			gameTimeCategoriesResult,
			gameTimeKindsResult,
			quickGamesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
