import { MessageSchema, type Message } from '$lib/gen/juicer_pb';
import { toJsonString } from '@bufbuild/protobuf';
const wsEndpoint = import.meta.env['VITE_PUBLIC_WS_ENDPOINT'] as string;

export class ReconnectingWs {
	#url = wsEndpoint;
	#ws: WebSocket | null = null;
	#reconnectAttempts: number = 0;
	#maxReconnectAttempts: number = 10;
	#baseReconnectDelayMs: number = 500;
	#maxReconnectDelayMs: number = 60_000;
	#timerid: NodeJS.Timeout | null = null;
	get readyState(): number | undefined {
		return this.#ws?.readyState;
	}
	onMessage: (event: MessageEvent) => void = () => {};
	onOpen: (event: Event) => void = () => {};
	onClose: (event: CloseEvent) => void = () => {};
	onError: (event: Event) => void = () => {};

	constructor(params?: URLSearchParams) {
		if (params) {
			this.#url = `${this.#url}?${params.toString()}`;
		}
	}

	connect(params?: URLSearchParams): void {
		if (params) {
			this.#url = `${this.#url}?${params.toString()}`;
		}
		if (this.readyState === WebSocket.OPEN) {
			console.log('websocket is already open');
			return;
		}
		this.#ws = new WebSocket(this.#url);
		this.#ws.onopen = (event: Event) => {
			this.#reconnectAttempts = 0;
			if (this.#timerid) {
				clearTimeout(this.#timerid);
			}
			this.onOpen(event);
		};
		this.#ws.onclose = (event: CloseEvent) => {
			if (this.#timerid) {
				clearTimeout(this.#timerid);
			}
			this.onClose(event);
			this.#reconnect();
		};
		this.#ws.onmessage = (event: MessageEvent) => {
			this.onMessage(event);
		};
		this.#ws.onerror = (event: Event) => {
			this.onError(event);
			this.close();
		};
	}

	close(): void {
		this.#ws?.close();
	}

	send(msg: Message) {
		if (!this.#ws) {
			return;
		}
		if (this.readyState === WebSocket.OPEN) {
			this.#ws.send(toJsonString(MessageSchema, msg));
		} else {
			console.debug('ws is not open, cannot send message');
		}
	}

	#reconnect(): void {
		if (!navigator.onLine) {
			console.log('connection offline');
			window.ononline = () => {
				this.connect();
			};
			return;
		}
		if (this.readyState === WebSocket.OPEN) {
			console.log('websocket is already open');
			return;
		}
		if (this.#reconnectAttempts >= this.#maxReconnectAttempts) {
			console.debug('max reconnect attempts reached, giving up');
			return;
		}
		const delay = this.#getBackoffDelayMs(this.#reconnectAttempts);
		console.log(`reconnecting in ${(delay / 1000).toFixed(2)} seconds...`);
		this.#timerid = setTimeout(() => {
			this.#reconnectAttempts++;
			this.connect();
		}, delay);
	}

	#getBackoffDelayMs(attempt: number): number {
		const jitter = Math.random() * 1000;
		return Math.min(this.#baseReconnectDelayMs * Math.pow(2, attempt) + jitter, this.#maxReconnectDelayMs);
	}
}
