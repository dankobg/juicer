import { browser } from '$app/environment';
import { uiSettings } from '$lib/components/ui-settings/ui-settings-state.svelte';

export type SoundName =
	| 'Capture'
	| 'Check'
	| 'Checkmate'
	| 'Confirmation'
	| 'Defeat'
	| 'Draw'
	| 'Error'
	| 'LowTime'
	| 'Move'
	| 'NewChallenge'
	| 'NewChatMessage'
	| 'NewPM'
	| 'OutOfBound'
	| 'Victory';

class SoundManager {
	audioCtx?: AudioContext = browser ? new AudioContext() : undefined;
	audioBuffers: Record<string, AudioBuffer> = {};
	sounds = $derived.by(() => {
		const theme = uiSettings.sounds.current.theme;
		const newChatMsgUrl =
			theme === 'futuristic'
				? '/sounds/futuristic/NewChatMessage.ogg'
				: theme === 'piano'
					? '/sounds/piano/stock/NewChatMessage.ogg'
					: `/sounds/${theme}/GenericNotify.ogg`;

		const sounds: Record<SoundName, string> = {
			Capture: `/sounds/${theme}/Capture.ogg`,
			Check: `/sounds/${theme}/Check.ogg`,
			Checkmate: `/sounds/${theme}/Checkmate.ogg`,
			Confirmation: `/sounds/${theme}/Confirmation.ogg`,
			Defeat: `/sounds/${theme}/Defeat.ogg`,
			Draw: `/sounds/${theme}/Draw.ogg`,
			Error: `/sounds/${theme}/Error.ogg`,
			LowTime: `/sounds/${theme}/LowTime.ogg`,
			Move: `/sounds/${theme}/Move.ogg`,
			NewChallenge: `/sounds/${theme}/NewChallenge.ogg`,
			NewChatMessage: newChatMsgUrl,
			NewPM: `/sounds/${theme}/NewPM.ogg`,
			OutOfBound: `/sounds/${theme}/OutOfBound.ogg`,
			Victory: `/sounds/${theme}/Victory.ogg`
		};
		return sounds;
	});

	async preloadSounds(): Promise<void> {
		for (const sound of Object.values(this.sounds)) {
			this.loadSound(sound);
		}
	}

	async loadSound(url: string): Promise<AudioBuffer> {
		if (this.audioBuffers[url]) {
			return this.audioBuffers[url];
		}
		const resp = await fetch(url);
		const buf = await resp.arrayBuffer();
		const audioBuffer = await this.audioCtx!.decodeAudioData(buf);
		this.audioBuffers[url] = audioBuffer;
		return audioBuffer;
	}

	async playSound(name: SoundName, volume: number = 1) {
		const url = this.sounds[name];
		const buffer = await this.loadSound(url);
		const source = this.audioCtx!.createBufferSource();
		const gainNode = this.audioCtx!.createGain();
		source.buffer = buffer;
		gainNode.gain.value = volume;
		source.connect(gainNode).connect(this.audioCtx!.destination);
		source.start(0);
	}

	async play(name: SoundName) {
		if (!uiSettings.sounds.current.enabled) {
			return;
		}
		this.playSound(name, uiSettings.sounds.current.volume);
	}
}

export const soundManager = new SoundManager();
