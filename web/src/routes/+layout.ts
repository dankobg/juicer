import { kratos } from '$lib/kratos/client';
import { KratosService } from '$lib/kratos/service';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async () => {
	try {
		const result = await kratos.toSession();
		if (result.status !== 200) {
			return {
				auth: new KratosService(null),
			};
		}

		const session = result.data;
		const svc = new KratosService(session);

		const logoutResp = await kratos.createBrowserLogoutFlow();

		return {
			auth: svc,
			logoutToken: logoutResp.data.logout_token,
			logoutUrl: logoutResp.data.logout_url,
		};
	} catch (error) {
		return {
			auth: new KratosService(null),
		};
	}
};
