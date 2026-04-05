import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, depends, parent }) => {
	depends(`data:dashboard`);

	const data = await parent();

	if (data.auth.user?.isDeveloper) {
		try {
			// const analyticsStatsResult = await juicer.GET('/analytics/stats', {
			// 	fetch
			// });

			// if (analyticsStatsResult.error?.status_code === 403) {
			// 	if (browser) {
			// 		goto('/');
			// 	}
			// }

			return {
				...data
				// analyticsStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	} else {
		try {
			// const myAnalyticsStatsResult = await juicer.GET('/me/analytics/stats', {
			// 	fetch
			// });

			// if (myAnalyticsStatsResult.error?.status_code === 403) {
			// 	if (browser) {
			// 		goto('/');
			// 	}
			// }

			return {
				...data
				// myAnalyticsStatsResult
			};
		} catch (error) {
			console.log('err', error);
		}
	}
};
