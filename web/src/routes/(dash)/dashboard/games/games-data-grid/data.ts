import type { components } from '#lib/gen/juicer_openapi.ts';

export const gameResults: components['schemas']['GameResult'][] = [
	{ id: 1, name: 'white-won', created_at: '', updated_at: '' },
	{ id: 2, name: 'black-won', created_at: '', updated_at: '' },
	{ id: 3, name: 'draw', created_at: '', updated_at: '' },
	{ id: 4, name: 'interrupted', created_at: '', updated_at: '' }
];

export const gameResultStatuses: components['schemas']['GameResultStatus'][] = [
	{ id: 1, name: 'checkmate', created_at: '', updated_at: '' },
	{ id: 2, name: 'insufficient-material', created_at: '', updated_at: '' },
	{ id: 3, name: 'threefold-repetition', created_at: '', updated_at: '' },
	{ id: 4, name: 'fivefold-repetition', created_at: '', updated_at: '' },
	{ id: 5, name: 'fifty-move-rule', created_at: '', updated_at: '' },
	{ id: 6, name: 'seventyfive-move-rule', created_at: '', updated_at: '' },
	{ id: 7, name: 'stalemate', created_at: '', updated_at: '' },
	{ id: 8, name: 'resignation', created_at: '', updated_at: '' },
	{ id: 9, name: 'draw-agreed', created_at: '', updated_at: '' },
	{ id: 10, name: 'flagged', created_at: '', updated_at: '' },
	{ id: 11, name: 'adjudication', created_at: '', updated_at: '' },
	{ id: 12, name: 'timed-out', created_at: '', updated_at: '' },
	{ id: 13, name: 'aborted', created_at: '', updated_at: '' },
	{ id: 14, name: 'interrupted', created_at: '', updated_at: '' }
];

export const gameStates: components['schemas']['GameState'][] = [
	{ id: 1, name: 'active', created_at: '', updated_at: '' },
	{ id: 2, name: 'finished', created_at: '', updated_at: '' },
	{ id: 3, name: 'interrupted', created_at: '', updated_at: '' }
];

type NamedEntity = {
	id?: number;
	name?: string;
};

function toIdNameMap<T extends NamedEntity>(items: T[]): Map<number, string> {
	const map = new Map<number, string>();
	for (const item of items) {
		if (typeof item.id === 'number' && typeof item.name === 'string') {
			map.set(item.id, item.name);
		}
	}
	return map;
}

export function getGameVariantNames(gameVariants: components['schemas']['GameVariant'][]): Map<number, string> {
	return toIdNameMap(gameVariants);
}

export function getGameTimeKindNames(gameTimeKinds: components['schemas']['GameTimeKind'][]): Map<number, string> {
	return toIdNameMap(gameTimeKinds);
}

export function getGameTimeCategoryNames(
	gameTimeCategories: components['schemas']['GameTimeCategory'][]
): Map<number, string> {
	return toIdNameMap(gameTimeCategories);
}

export const gameResultNames = new Map<number, string>([
	[1, 'white-won'],
	[2, 'black-won'],
	[3, 'draw'],
	[4, 'interrupted']
]);

export const gameResultStatusNames = new Map<number, string>([
	[1, 'checkmate'],
	[2, 'insufficient-material'],
	[3, 'threefold-repetition'],
	[4, 'fivefold-repetition'],
	[5, 'fifty-move-rule'],
	[6, 'seventyfive-move-rule'],
	[7, 'stalemate'],
	[8, 'resignation'],
	[9, 'draw-agreed'],
	[10, 'flagged'],
	[11, 'adjudication'],
	[12, 'timed-out'],
	[13, 'aborted'],
	[14, 'interrupted']
]);

export const gameStateNames = new Map<number, string>([
	[1, 'active'],
	[2, 'finished'],
	[3, 'interrupted']
]);
