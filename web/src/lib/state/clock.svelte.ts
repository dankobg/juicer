import { Timer, type TimerOptions } from '$lib/state/timer.svelte';

export type ClockOptions = {
	whiteTimeMs: number;
	blackTimeMs: number;
	whiteIncrementMs?: number;
	blackIncrementMs?: number;
	showPrecise?: TimerOptions['showPrecise'];
	showPreciseFn?: TimerOptions['showPreciseFn'];
	onTimeout?: (color: 'w' | 'b') => void;
	initialTurn?: 'w' | 'b';
};

export class Clock {
	white!: Timer;
	black!: Timer;
	#initialWhiteClockMs: number = 0;
	#initialBlackClockMs: number = 0;
	#initialWhiteIncrementMs: number = 0;
	#initialBlackIncrementMs: number = 0;
	#initialTurn: 'w' | 'b' = 'w';
	showPrecise: Required<ClockOptions>['showPrecise'] = 'centis';
	showPreciseFn: Required<ClockOptions>['showPreciseFn'] = () => false;
	onTimeout: ClockOptions['onTimeout'];
	state: 'idle' | 'paused' | 'running' | 'timed-out' = $state('idle');
	whiteIncrementMs: number = $state(0);
	blackIncrementMs: number = $state(0);
	currentTurn: 'w' | 'b' = $state('w');
	currentTimer: Timer = $derived(this.currentTurn === 'w' ? this.white : this.black);
	currentIncrement: number = $derived(this.currentTurn === 'w' ? this.whiteIncrementMs : this.blackIncrementMs);

	constructor(options: ClockOptions) {
		this.#initialWhiteClockMs = options.whiteTimeMs;
		this.#initialBlackClockMs = options.blackTimeMs;
		if (options.whiteIncrementMs) {
			this.#initialWhiteIncrementMs = options.whiteIncrementMs;
			this.whiteIncrementMs = options.whiteIncrementMs;
		}
		if (options.blackIncrementMs) {
			this.#initialBlackIncrementMs = options.blackIncrementMs;
			this.blackIncrementMs = options.blackIncrementMs;
		}
		if (options.showPrecise) {
			this.showPrecise = options.showPrecise;
		}
		if (options.showPreciseFn) {
			this.showPreciseFn = options.showPreciseFn;
		}
		if (options?.onTimeout) {
			this.onTimeout = options.onTimeout;
		}
		if (options.initialTurn) {
			this.#initialTurn = options.initialTurn;
			this.currentTurn = options.initialTurn;
		}
		const timerOpts: TimerOptions = {
			showPrecise: options.showPrecise,
			showPreciseFn: options.showPreciseFn
		};
		const whiteTimerOpts = {
			...timerOpts,
			onTimeout: () => {
				this.black.pause();
				this.state = 'timed-out';
				this.onTimeout?.('w');
			}
		};
		const blackTimerOpts = {
			...timerOpts,
			onTimeout: () => {
				this.white.pause();
				this.state = 'timed-out';
				this.onTimeout?.('b');
			}
		};
		this.white = new Timer(this.#initialWhiteClockMs, whiteTimerOpts);
		this.black = new Timer(this.#initialBlackClockMs, blackTimerOpts);
	}

	setIncrement(increment: number): void {
		if (this.state === 'timed-out') {
			return;
		}
		this.setIncrementCustom({ white: increment, black: increment });
	}

	setIncrementCustom(increment?: { white?: number; black?: number }): void {
		if (this.state === 'timed-out') {
			return;
		}
		if (increment?.white) {
			this.whiteIncrementMs = increment.white;
		}
		if (increment?.black) {
			this.blackIncrementMs = increment.black;
		}
	}

	setCurrentTurn(color: 'w' | 'b'): void {
		if (this.state === 'timed-out') {
			return;
		}
		this.currentTurn = color;
	}

	add(color: 'w' | 'b', increment: number) {
		if (this.state === 'timed-out') {
			return;
		}
		if (color === 'w') {
			this.white.add(increment);
		} else {
			this.black.add(increment);
		}
	}

	start() {
		if (this.state === 'running' || this.state === 'timed-out') {
			return;
		}
		this.state = 'running';
		this.currentTimer.start();
	}

	pause() {
		if (this.state === 'running') {
			this.state = 'paused';
			this.currentTimer.pause();
		}
	}

	reset() {
		this.state = 'idle';
		this.whiteIncrementMs = this.#initialWhiteIncrementMs;
		this.blackIncrementMs = this.#initialBlackIncrementMs;
		this.currentTurn = this.#initialTurn;
		this.white.reset();
		this.black.reset();
	}

	restart() {
		this.reset();
		this.start();
	}

	toggle() {
		if (this.state === 'timed-out') {
			return;
		}
		if (this.state === 'idle') {
			this.start();
		} else {
			this.currentTimer.pause();
			if (this.currentIncrement > 0) {
				this.currentTimer.add(this.currentIncrement);
			}
			this.setCurrentTurn(this.currentTurn === 'w' ? 'b' : 'w');
			this.currentTimer.start();
		}
	}

	synchronize(whiteClockMs: number, blackClockMs: number) {
		this.white.synchronize(whiteClockMs);
		this.black.synchronize(blackClockMs);
	}
}
