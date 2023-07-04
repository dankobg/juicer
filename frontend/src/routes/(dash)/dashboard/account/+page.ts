import type { PageLoad } from './$types';
import type { SettingsFlow } from '@ory/client';
import { kratos } from '$lib/kratos/client';
import { extractCSRFToken } from '$lib/kratos/helpers';
import { browser } from '$app/environment';

export const load = (async ({ url }) => {
  const returnToParam = browser && url.searchParams.get('return_to');
  const flowIdParam = browser && url.searchParams.get('flow');

  let flow: SettingsFlow | null = null;

  if (flowIdParam) {
    try {
      const flowResponse = await kratos.getSettingsFlow({
        id: flowIdParam,
      });

      flow = { ...flowResponse.data };
    } catch (error) {
      console.log('settings error:', error);

      flow = null;
    }
  } else {
    const returnTo: string | undefined = returnToParam ? returnToParam.toString() : undefined;

    try {
      const flowResponse = await kratos.createBrowserSettingsFlow({
        returnTo,
      });

      if ([403, 404, 410].includes(flowResponse.status)) {
        console.log('[403, 404, 410].includes(flowResponse.status)');
      }
      if (flowResponse.status !== 200) {
        console.log('flowResponse.status !== 200');
      }

      flow = { ...flowResponse.data };
    } catch (error) {
      console.log('settings error2:', error);
    }
  }

  const csrf = extractCSRFToken(flow);

  return {
    flow,
    csrf,
  };
}) satisfies PageLoad;
