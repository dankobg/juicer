import { browser } from '$app/env';
import { goto } from '$app/navigation';
import { juicer } from '#lib/juicer/client.ts';
import type {
	operations,
	PathsGamesGetParametersQueryEmbed,
	PathsGamesGetParametersQueryUser_game_result
} from '#lib/gen/juicer_openapi.ts';
import { config } from '#lib/kratos/config.ts';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, url, depends, parent }) => {
	depends('data:dashboard-games');

	const data = await parent();
	const { auth } = data;

	const listGamesParams: operations['listGames']['parameters'] = {
		query: {
			user_id: auth?.user?.id,
			sort: ['-start_time']
		}
	};

	const pageParam = url.searchParams.get('page');
	if (pageParam) {
		const page = Number.parseInt(pageParam);
		if (!Number.isNaN(page)) {
			listGamesParams.query!.page = page;
		}
	}
	const pageSizeParam = url.searchParams.get('page_size');
	if (pageSizeParam) {
		const pageSize = Number.parseInt(pageSizeParam);
		if (!Number.isNaN(pageSize)) {
			listGamesParams.query!.page_size = pageSize;
		}
	}

	const gameVariantIdParams = url.searchParams.getAll('game_variant_id');
	if (gameVariantIdParams.length > 0) {
		const gameVariantIds = gameVariantIdParams.map(Number).filter(Number.isFinite);
		if (gameVariantIds.length > 0) {
			listGamesParams.query!.game_variant_id = gameVariantIds;
		}
	}

	const gameTimeKindIdParams = url.searchParams.getAll('game_time_kind_id');
	if (gameTimeKindIdParams.length > 0) {
		const gameTimeKindIds = gameTimeKindIdParams.map(Number).filter(Number.isFinite);
		if (gameTimeKindIds.length > 0) {
			listGamesParams.query!.game_time_kind_id = gameTimeKindIds;
		}
	}

	const gameTimeCategoryIdParams = url.searchParams.getAll('game_time_category_id');
	if (gameTimeCategoryIdParams.length > 0) {
		const gameTimeCategoryIds = gameTimeCategoryIdParams.map(Number).filter(Number.isFinite);
		if (gameTimeCategoryIds.length > 0) {
			listGamesParams.query!.game_time_category_id = gameTimeCategoryIds;
		}
	}

	const gameResultIdParams = url.searchParams.getAll('game_result_id');
	if (gameResultIdParams.length > 0) {
		const gameResultIds = gameResultIdParams.map(Number).filter(Number.isFinite);
		if (gameResultIds.length > 0) {
			listGamesParams.query!.game_result_id = gameResultIds;
		}
	}

	const gameResultStatusIdParams = url.searchParams.getAll('game_result_status_id');
	if (gameResultStatusIdParams.length > 0) {
		const gameResultStatusIds = gameResultStatusIdParams.map(Number).filter(Number.isFinite);
		if (gameResultStatusIds.length > 0) {
			listGamesParams.query!.game_result_status_id = gameResultStatusIds;
		}
	}

	const gameStateIdParams = url.searchParams.getAll('game_state_id');
	if (gameStateIdParams.length > 0) {
		const gameStateIds = gameStateIdParams.map(Number).filter(Number.isFinite);
		if (gameStateIds.length > 0) {
			listGamesParams.query!.game_state_id = gameStateIds;
		}
	}

	const userIdParam = url.searchParams.get('user_id');
	if (userIdParam) {
		listGamesParams.query!.user_id = userIdParam;
	}

	const userGameResultParams = url.searchParams.getAll('user_game_result');
	if (userGameResultParams.length > 0) {
		listGamesParams.query!.user_game_result = userGameResultParams as PathsGamesGetParametersQueryUser_game_result[];
	}

	const ratedParam = url.searchParams.get('rated');
	if (ratedParam === 'true' || ratedParam === 'false') {
		listGamesParams.query!.rated = ratedParam === 'true';
	}

	const createdAtFromParam = url.searchParams.get('created_at_from');
	if (createdAtFromParam) {
		listGamesParams.query!.created_at_from = createdAtFromParam;
	}

	const createdAtToParam = url.searchParams.get('created_at_to');
	if (createdAtToParam) {
		listGamesParams.query!.created_at_to = createdAtToParam;
	}

	const sortParams = url.searchParams.getAll('sort');
	if (sortParams.length > 0) {
		listGamesParams.query!.sort = sortParams;
	}

	const embed = url.searchParams.getAll('embed') as PathsGamesGetParametersQueryEmbed[];
	if (embed.length > 0) {
		listGamesParams.query!.embed = embed;
	}

	try {
		const gamesResult = await juicer.GET('/games', {
			fetch,
			params: listGamesParams
		});

		if (gamesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const gameVariantsResult = await juicer.GET('/game-variants', {
			fetch
		});

		if (gameVariantsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const gameTimeKindsResult = await juicer.GET('/game-time-kinds', {
			fetch
		});

		if (gameTimeKindsResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		const gameTimeCategoriesResult = await juicer.GET('/game-time-categories', {
			fetch
		});

		if (gameTimeCategoriesResult.error?.status_code === 403) {
			if (browser) {
				goto(config.routes.dashboard.path);
			}
		}

		return {
			gamesResult,
			gameVariantsResult,
			gameTimeKindsResult,
			gameTimeCategoriesResult
		};
	} catch (error) {
		console.log('err', error);
	}
};
