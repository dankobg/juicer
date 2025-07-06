import { PersistedState } from 'runed';

class Settings {
	boardThemes = [
		{ name: 'blue-marble', src: 'blue-marble.jpg' },
		{ name: 'blue-marble-orig', src: 'blue-marble.orig.jpg' },
		{ name: 'blue2', src: 'blue2.jpg' },
		{ name: 'blue3', src: 'blue3.jpg' },
		{ name: 'canvas2', src: 'canvas2.jpg' },
		{ name: 'canvas2-orig', src: 'canvas2.orig.jpg' },
		{ name: 'green-plastic', src: 'green-plastic.png' },
		{ name: 'grey', src: 'grey.jpg' },
		{ name: 'horsey', src: 'horsey.jpg' },
		{ name: 'leather', src: 'leather.jpg' },
		{ name: 'leather-orig', src: 'leather.orig.jpg' },
		{ name: 'maple', src: 'maple.jpg' },
		{ name: 'maple2', src: 'maple2.jpg' },
		{ name: 'maple2-orig', src: 'maple2.orig.jpg' },
		{ name: 'marble', src: 'marble.jpg' },
		{ name: 'metal', src: 'metal.jpg' },
		{ name: 'metal-orig', src: 'metal.orig.jpg' },
		{ name: 'ncf-board', src: 'ncf-board.png' },
		{ name: 'newspaper', src: 'newspaper.png' },
		{ name: 'olive', src: 'olive.jpg' },
		{ name: 'pink-pyramid', src: 'pink-pyramid.png' },
		{ name: 'purple-diag', src: 'purple-diag.png' },
		{ name: 'svg-blue', src: 'svg/blue.svg' },
		{ name: 'svg-brown', src: 'svg/brown.svg' },
		{ name: 'svg-green', src: 'svg/green.svg' },
		{ name: 'svg-ic', src: 'svg/ic.svg' },
		{ name: 'svg-purple', src: 'svg/purple.svg' },
		{ name: 'wood', src: 'wood.jpg' },
		{ name: 'wood2', src: 'wood2.jpg' },
		{ name: 'wood3', src: 'wood3.jpg' },
		{ name: 'wood3-orig', src: 'wood3.orig.jpg' },
		{ name: 'wood4', src: 'wood4.jpg' },
		{ name: 'wood4-orig', src: 'wood4.orig.jpg' }
	];
	pieceThemes = [
		'alpha',
		'anarcandy',
		'caliente',
		'california',
		'cardinal',
		'cburnett',
		'celtic',
		'chess7',
		'chessnut',
		'companion',
		'cooke',
		'disguised',
		'dubrovny',
		'fantasy',
		'firi',
		'fresca',
		'gioco',
		'governor',
		'horsey',
		'icpieces',
		'kiwen-suwi',
		'kosal',
		'leipzig',
		'letter',
		'maestro',
		'merida',
		'monarchy',
		'mono',
		'mpchess',
		'pirouetti',
		'pixel',
		'reillycraig',
		'rhosgfx',
		'riohacha',
		'shapes',
		'spatial',
		'staunty',
		'tatiana',
		'xkcd'
	];
	boardActiveTheme = new PersistedState<(typeof this.boardThemes)[number]>('juicer-board-theme', {
		name: 'svg-brown',
		src: 'svg/brown.svg'
	});
	pieceActiveTheme = new PersistedState<(typeof this.pieceThemes)[number]>('juicer-piece-theme', 'gioco');
	soundThemes = ['futuristic', 'lisp', 'nes', 'piano', 'robot', 'sfx', 'standard', 'woodland'] as const;
	sounds = new PersistedState<{ enabled: boolean; volume: number; theme: string }>('juicer-sounds', {
		enabled: true,
		volume: 1,
		theme: 'standard'
	});
	chat = new PersistedState<'disabled' | 'friends-only' | 'everyone'>('juicer-chat', 'everyone');
	resizer = new PersistedState<'disabled' | 'first-move' | 'always'>('juicer-resizer', 'first-move');

	dialogOpen: boolean = $state(false);

	open(): void {
		this.dialogOpen = true;
	}
	close(): void {
		this.dialogOpen = false;
	}
	toggle(): void {
		this.dialogOpen = !this.dialogOpen;
	}
}

export const settings = new Settings();
