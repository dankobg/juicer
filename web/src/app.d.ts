// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		interface PageState {
			reseek?: {
				clockMs: number;
				incrementMs: number;
			};
		}
		// interface Platform {}
	}
}

export {};

export {};
