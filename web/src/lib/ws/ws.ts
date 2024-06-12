export class JuicerWS {
	#url = 'wss://juicer-dev.xyz/ws';
	#ws: WebSocket | null = null;
	#reconnectAttempts: number = 0;
	#maxReconnectAttempts: number = 5;
	#baseReconnectDelay: number = 350;
	#maxReconnectDelay: number = 10_000;
	#timerid: NodeJS.Timeout | null = null;
	onmessage: (event: MessageEvent) => void = () => {};

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
		console.debug(`reconnecting in ${delay / 1000} seconds...`);
		this.#timerid = setTimeout(() => {
			this.connect();
		}, delay);
	}

	connect() {
		this.#ws = new WebSocket(this.#url);

		this.#ws.onopen = () => {
			console.debug('ws connected');
			this.#reconnectAttempts = 0;
		};

		this.#ws.onmessage = (event: MessageEvent) => {
			console.debug('ws recv:', event.data);
			this.onmessage(event);
		};

		this.#ws.onclose = (event: CloseEvent) => {
			console.debug(`ws closed: ${event.code} ${event.reason}`);
			if (this.#timerid) {
				clearTimeout(this.#timerid);
			}
			this.#reconnect();
		};

		this.#ws.onerror = error => {
			console.debug('ws error:', error);
			this.close();
		};
	}

	close() {
		this.#ws?.close();
	}

	send(data: string) {
		if (!this.#ws) {
			return;
		}
		if (this.#ws.readyState === WebSocket.OPEN) {
			this.#ws.send(data);
		} else {
			console.debug('ws is not open, cannot send message');
		}
	}
}
