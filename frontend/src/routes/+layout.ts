import type { LayoutLoad } from './$types';
import { kratos } from '$lib/kratos/client';
import { getSession } from '$lib/kratos/service';

export const prerender = true;

export const load = (async () => {
  const session = await getSession();

  let logoutToken: string | null = null;
  let logoutUrl: string | null = null;

  if (session?.active) {
    const flowResponse = await kratos.createBrowserLogoutFlow();

    logoutToken = flowResponse.data.logout_token;
    logoutUrl = flowResponse.data.logout_url;
  }

  return {
    session,
    logoutToken,
    logoutUrl,
  };
}) satisfies LayoutLoad;
