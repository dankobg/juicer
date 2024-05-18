import { kratos } from '$lib/kratos/client';
import { createSessionService } from '$lib/kratos/service';
import type { LayoutLoad } from './$types';

export const prerender = true;

export const load: LayoutLoad = async () => {
	try {
		const result = await kratos.toSession();
		if (result.status !== 200) {
			return {
				auth: createSessionService(result.data),
			};
		}

		const svc = createSessionService(result.data);
		const logoutResp = await kratos.createBrowserLogoutFlow();

		return {
			auth: svc,
			logoutToken: logoutResp.data.logout_token,
			logoutUrl: logoutResp.data.logout_url,
		};
	} catch (error) {
		return {
			auth: createSessionService(null),
		};
	}
};
