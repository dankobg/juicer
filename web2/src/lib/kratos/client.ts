import { Configuration, FrontendApi } from '@ory/client-fetch';
import { config } from './config';

const kratosConfiguration = new Configuration({
	basePath: config.kratos.publicUrl,
	credentials: 'include'
});

export const kratos = new FrontendApi(kratosConfiguration);
