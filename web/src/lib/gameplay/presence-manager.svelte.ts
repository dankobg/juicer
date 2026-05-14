import type { Echo, Presence, PresenceDiff, PresenceState } from '$lib/gen/juicer_pb';
import { SvelteMap, SvelteSet } from 'svelte/reactivity';

class PresenceManager {
	userPresences = $state<SvelteMap<string, Presence>>(new SvelteMap());
	channelPresences = $state<SvelteMap<string, SvelteSet<string>>>(new SvelteMap());

	lobbyPresence = $derived(this.getPresenceInChannel('lobby'));
	lobbyChatPresence = $derived(this.getPresenceInChannel('lobby.chat'));

	getPresenceInChannel(channel: string): Presence[] {
		const lobbyPresence = this.channelPresences.get(channel);
		if (lobbyPresence?.size === 0) {
			return [];
		}
		return [...(lobbyPresence?.values() ?? [])].reduce((presenceList, userId) => {
			const presence = this.userPresences.get(userId);
			if (presence) {
				presenceList.push(presence);
			}
			return presenceList;
		}, [] as Presence[]);
	}

	onEcho(echoMsg: Echo): void {
		console.log('got echo: ', echoMsg.message);
	}

	onPresenceState(presenceState: PresenceState): void {
		for (const presence of presenceState.presences) {
			this.userPresences.set(presence.userId, presence);
			const channelPresence = this.channelPresences.get(presence.channel) ?? new SvelteSet<string>();
			channelPresence.add(presence.userId);
			this.channelPresences.set(presence.channel, channelPresence);
		}
	}

	onPresenceDiff(presenceDiff: PresenceDiff): void {
		for (const presence of presenceDiff.joined) {
			let channelUsers = this.channelPresences.get(presence.channel);

			if (!channelUsers) {
				channelUsers = new SvelteSet<string>();
				this.channelPresences.set(presence.channel, channelUsers);
			}

			channelUsers.add(presence.userId);
			this.userPresences.set(presence.userId, presence);
		}

		for (const presence of presenceDiff.left) {
			const channelUsers = this.channelPresences.get(presence.channel);
			channelUsers?.delete(presence.userId);

			if (channelUsers?.size === 0) {
				this.channelPresences.delete(presence.channel);
			}

			let stillInAnyChannel = false;

			for (const userId of this.channelPresences.values()) {
				if (userId.has(presence.userId)) {
					stillInAnyChannel = true;
					break;
				}
			}

			if (!stillInAnyChannel) {
				this.userPresences.delete(presence.userId);
			}
		}
	}
}

export const presenceManager = new PresenceManager();
