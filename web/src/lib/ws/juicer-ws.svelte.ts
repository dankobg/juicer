import { MessageSchema, type Message } from '$lib/gen/juicer_pb';
import { toJsonString } from '@bufbuild/protobuf';
import { createSubscriber } from 'svelte/reactivity';
const wsEndpoint = import.meta.env['VITE_PUBLIC_WS_ENDPOINT'] as string;

class JuicerWebSocket {
	baseUrl: string = wsEndpoint;
	#ws: WebSocket | null = null;
	#update: VoidFunction | null = null;
	subscribe: VoidFunction;

	#reconnectAttempts: number = 0;
	#maxReconnectAttempts: number = 10;
	#baseReconnectDelayMs: number = 500;
	#maxReconnectDelayMs: number = 60_000;
	#timeoutId: ReturnType<typeof setTimeout> | null = null;

	onMessage: (event: MessageEvent) => void = () => {};
	onOpen: (event: Event) => void = () => {};
	onClose: (event: CloseEvent) => void = () => {};
	onError: (event: Event) => void = () => {};

	constructor() {
		this.subscribe = createSubscriber(update => {
			this.#update = update;

			return () => {
				this.#update = null;
			};
		});
	}

	get readyState(): WebSocket['readyState'] {
		this.subscribe();
		return this.#ws ? this.#ws.readyState : WebSocket.CLOSED;
	}

	connect(params?: URLSearchParams): void {
		this.#update?.();

		if (this.#ws && (this.#ws.readyState === WebSocket.CONNECTING || this.#ws.readyState === WebSocket.OPEN)) {
			console.debug('websocket connection connecting or already open');
			return;
		}

		const url = params ? `${this.baseUrl}?${params.toString()}` : this.baseUrl;

		this.#ws = new WebSocket(url);

		this.#update?.();

		this.#ws.onopen = (event: Event): void => {
			this.#update?.();
			this.#maxReconnectAttempts = 0;
			if (this.#timeoutId !== null) {
				clearTimeout(this.#timeoutId);
			}
			this?.onOpen?.(event);
		};

		this.#ws.onclose = (event: CloseEvent): void => {
			this.#update?.();
			if (this.#timeoutId !== null) {
				clearTimeout(this.#timeoutId);
			}
			this?.onClose?.(event);
			this.#reconnect();
		};

		this.#ws.onerror = (event: Event): void => {
			this.#update?.();
			this?.onError?.(event);
			this.close();
		};

		this.#ws.onmessage = (event: MessageEvent): void => {
			this.#update?.();
			this?.onMessage?.(event);
		};
	}

	close(): void {
		this.#ws?.close();
	}

	send(msg: Message) {
		if (!this.#ws || this.#ws.readyState !== WebSocket.OPEN) {
			console.debug('ws is not open, cannot send message');
			return;
		}
		this.#ws.send(toJsonString(MessageSchema, msg));
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
		this.#timeoutId = setTimeout(() => {
			this.#reconnectAttempts++;
			this.connect();
		}, delay);
	}

	#getBackoffDelayMs(attempt: number): number {
		const jitter = Math.random() * 1000;
		return Math.min(this.#baseReconnectDelayMs * Math.pow(2, attempt) + jitter, this.#maxReconnectDelayMs);
	}
}

export const ws = new JuicerWebSocket();
