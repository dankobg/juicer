import { juicer } from '$lib/juicer/client';
import type { PageLoad } from './$types';
import { error, isHttpError } from '@sveltejs/kit';

export const ssr = false;
export const prerender = false;

export const load: PageLoad = async ({ fetch, depends, params }) => {
	depends(`data:game-${params.game_id}`);

	try {
		const gameResult = await juicer.GET('/games/{id}', {
			fetch,
			params: {
				path: { id: Number(params.game_id) }
			}
		});

		if (gameResult.error) {
			error(404);
		}

		return {
			gameResult
		};
	} catch (err) {
		if (isHttpError(err)) {
			error(404, {
				message: 'Game not found'
			});
		}
	}
};
