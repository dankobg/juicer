const apiEndpoint = import.meta.env['VITE_PUBLIC_API_ENDPOINT'] as string;
import { Configuration, JuicerApi } from '$lib/gen/juicer_openapi';

const conf = new Configuration({
	basePath: apiEndpoint,
	credentials: 'include',
	fetchApi: fetch,
	headers: {
		'Content-Type': 'application/json',
		Accept: 'application/json'
	}
});

export const juicer = new JuicerApi(conf);
