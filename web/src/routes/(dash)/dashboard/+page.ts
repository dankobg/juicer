import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends, parent }) => {
	depends(`data:dashboard`);

	const data = await parent();

	if (data.auth.user?.isDeveloper) {
		try {
			const gameStatsResult = await juicer.GET('/game-stats/{user_id}', {
				params: { path: { user_id: data?.auth?.user?.id } },
				fetch
			});

			if (gameStatsResult.error?.status_code === 403) {
				if (browser) {
					goto('/');
				}
			}

			return {
				...data,
				gameStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	} else {
		try {
			const gameStatsResult = await juicer.GET('/game-stats/{user_id}', {
				params: { path: { user_id: data?.auth?.user?.id ?? '' } },
				fetch
			});

			if (gameStatsResult.error?.status_code === 403) {
				if (browser) {
					goto('/');
				}
			}

			return {
				...data,
				gameStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	}
};
