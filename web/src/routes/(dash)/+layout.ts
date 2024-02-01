import { goto } from '$app/navigation';
import { config } from '$lib/kratos/config';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent, url }) => {
	const data = await parent();
	const { auth } = data;

	if (!auth.session?.active) {
		goto(config.routes.login.path + `?return_to=${encodeURIComponent(url.toString())}`);
	}

	return data;
};
