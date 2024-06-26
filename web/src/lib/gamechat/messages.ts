import { writable } from 'svelte/store';

export type ChatMsg = {
	text: string;
	own: boolean;
};

export const chatMessages = writable<ChatMsg[]>([]);
