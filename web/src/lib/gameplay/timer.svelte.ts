export type TimerOptions = {
	showPrecise?: 'decis' | 'centis' | 'millis';
	showPreciseFn?: (timeMs: number) => boolean;
	onTimeout?: () => void;
};

export class Timer {
	#fmt = new Intl.NumberFormat(undefined, { minimumIntegerDigits: 2 });
	#initialTimeMs: number;
	#intervalId?: NodeJS.Timeout;
	#startedAt: number = 0;
	showPrecise: Required<TimerOptions>['showPrecise'] = 'centis';
	showPreciseFn: Required<TimerOptions>['showPreciseFn'] = () => false;
	onTimeout: TimerOptions['onTimeout'];
	state: 'idle' | 'paused' | 'running' | 'timed-out' = $state('idle');
	timeMs: number = $state(0);
	time = $derived(this.#getTime(this.timeMs));

	constructor(timeMs: number, options?: TimerOptions) {
		this.#initialTimeMs = timeMs;
		this.timeMs = this.#initialTimeMs;
		if (options?.showPrecise) {
			this.showPrecise = options?.showPrecise;
		}
		if (options?.showPreciseFn) {
			this.showPreciseFn = options?.showPreciseFn;
		}
		if (options?.onTimeout) {
			this.onTimeout = options.onTimeout;
		}
	}

	#extractTimeParts(totalTimeMs: number) {
		const totalSeconds = Math.floor(totalTimeMs / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		const milliseconds = totalTimeMs % 1000;
		return { hours, minutes, seconds, milliseconds, totalTimeMs };
	}

	#formatTime(parts: { hours: number; minutes: number; seconds: number; milliseconds: number }): string {
		const mins = this.#fmt.format(parts.minutes);
		const secs = this.#fmt.format(parts.seconds);
		const ms = this.#fmt.format(parts.milliseconds);
		let precise = '';
		if (this.showPrecise === 'millis') {
			precise = ms;
		} else if (this.showPrecise === 'centis') {
			precise = ms.slice(0, 2);
		} else if (this.showPrecise === 'decis') {
			precise = ms.slice(0, 1);
		}
		return `${mins}:${secs}${this.showPreciseFn(this.timeMs) ? ':' + precise : ''}`;
	}

	#getTime(totalTimeMs: number) {
		const parts = this.#extractTimeParts(totalTimeMs);
		const formatted = this.#formatTime(parts);
		return { ...parts, formatted };
	}

	add(incrementMs: number) {
		if (this.state === 'timed-out') {
			return;
		}
		this.timeMs += incrementMs;
	}

	start() {
		if (this.state === 'running' || this.state === 'timed-out') {
			return;
		}
		this.#startedAt = performance.now();
		this.state = 'running';
		this.#intervalId = setInterval(() => {
			if (this.timeMs <= 0) {
				this.state = 'timed-out';
				clearInterval(this.#intervalId);
				this?.onTimeout?.();
				return;
			}
			const now = performance.now();
			const elapsed = now - this.#startedAt!;
			this.#startedAt = now;
			this.timeMs = Math.max(0, this.timeMs - elapsed);
		}, 10);
	}

	pause() {
		if (this.state === 'running') {
			this.state = 'paused';
			clearInterval(this.#intervalId);
			this.#intervalId = undefined;
		}
	}

	reset() {
		this.state = 'idle';
		this.timeMs = this.#initialTimeMs;
		clearInterval(this.#intervalId);
		this.#intervalId = undefined;
		this.#startedAt = 0;
	}

	restart() {
		this.reset();
		this.start();
	}

	synchronize(timeMs: number) {
		this.timeMs = timeMs;
	}
}
