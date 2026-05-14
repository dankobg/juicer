import type { Latency } from '$lib/gen/juicer_pb';

class LatencyStats {
	latencyMs = $state<number>(0);

	onLatency(latency: Latency) {
		this.latencyMs = latency.latencyMs;
	}
}

export const latencyStats = new LatencyStats();
