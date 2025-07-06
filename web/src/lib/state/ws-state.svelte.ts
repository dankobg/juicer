import { ReconnectingWs } from '$lib/ws/reconnecting-ws';

export const ws = $state(new ReconnectingWs());
