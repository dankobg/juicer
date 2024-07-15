import { Message } from '$lib/gen/juicer_pb';

export class JuicerWS {
	#url = 'wss://juicer-dev.xyz/ws';
	#ws: WebSocket | null = null;
	#reconnectAttempts: number = 0;
	#maxReconnectAttempts: number = 8;
	#baseReconnectDelay: number = 150;
	#maxReconnectDelay: number = 15_000;
	#timerid: NodeJS.Timeout | null = null;
	get readyState(): number | undefined {
		return this.#ws?.readyState;
	}
	onMessage: (event: MessageEvent) => void = () => {};
	onOpen: (event: Event) => void = () => {};
	onClose: (event: CloseEvent) => void = () => {};
	onError: (event: Event) => void = () => {};

	constructor(url: string = '') {
		if (url) {
			this.#url = url;
		}
	}

	#reconnect() {
		if (this.#reconnectAttempts > this.#maxReconnectAttempts) {
			console.debug('max reconnect attempts reached, giving up');
			return;
		}

		this.#reconnectAttempts++;
		const delay = Math.min(this.#baseReconnectDelay * Math.pow(2, this.#reconnectAttempts), this.#maxReconnectDelay);
		console.debug(`reconnecting in ${(delay / 1000).toFixed(2)} seconds...`);
		this.#timerid = setTimeout(() => {
			this.connect();
		}, delay);
	}

	connect() {
		this.#ws = new WebSocket(this.#url);

		this.#ws.onopen = (event: Event) => {
			this.onOpen(event);
			this.#reconnectAttempts = 0;
		};

		this.#ws.onmessage = (event: MessageEvent) => {
			this.onMessage(event);
		};

		this.#ws.onclose = (event: CloseEvent) => {
			this.onClose(event);
			if (this.#timerid) {
				clearTimeout(this.#timerid);
			}
			this.#reconnect();
		};

		this.#ws.onerror = (event: Event) => {
			this.onError(event);
			this.close();
		};
	}

	close() {
		this.#ws?.close();
	}

	send(msg: Message) {
		if (!this.#ws) {
			return;
		}
		if (this.readyState === WebSocket.OPEN) {
			this.#ws.send(msg.toJsonString());
		} else {
			console.debug('ws is not open, cannot send message');
		}
	}
}
