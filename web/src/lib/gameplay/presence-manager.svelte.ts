import type { Echo, Presence, PresenceDiff, PresenceState } from '$lib/gen/juicer_pb';

class PresenceManager {
	userPresences = $state<Record<string, Presence>>({});
	channelPresences = $state<Record<string, string[]>>({});

	lobbyPresence = $derived(this.getPresenceInChannel('lobby'));
	lobbyChatPresence = $derived(this.getPresenceInChannel('lobby.chat'));

	getPresenceInChannel(channel: string): Record<string, Presence> {
		const userIds = this.channelPresences[channel];
		const result: Record<string, Presence> = {};
		if (!userIds) {
			return result;
		}
		for (const userId of userIds) {
			const presence = this.userPresences[userId];
			if (presence) {
				result[userId] = presence;
			}
		}
		return result;
	}

	onEcho(echoMsg: Echo): void {
		console.log('got echo: ', echoMsg.message);
	}

	onPresenceState(presenceState: PresenceState): void {
		for (const presence of presenceState.presences) {
			this.userPresences[presence.userId] = presence;
			const currentChannel = this.channelPresences[presence.channel] ?? [];
			if (!currentChannel.includes(presence.userId)) {
				currentChannel.push(presence.userId);
				this.channelPresences[presence.channel] = currentChannel;
			}
		}
	}

	onPresenceDiff(presenceDiff: PresenceDiff): void {
		for (const presence of presenceDiff.joined) {
			if (!this.channelPresences[presence.channel]) {
				this.channelPresences[presence.channel] = [];
			}

			const channelUsers = this.channelPresences[presence.channel];
			if (!channelUsers?.includes(presence.userId)) {
				channelUsers?.push(presence.userId);
			}

			this.userPresences[presence.userId] = presence;
		}

		for (const presence of presenceDiff.left) {
			const channelUsers = this.channelPresences[presence.channel];
			if (channelUsers) {
				const index = channelUsers.indexOf(presence.userId);
				if (index !== -1) {
					channelUsers.splice(index, 1);
				}

				if (channelUsers.length === 0) {
					delete this.channelPresences[presence.channel];
				}
			}

			let stillInAnyChannel = false;
			for (const activeUsers of Object.values(this.channelPresences)) {
				if (activeUsers.includes(presence.userId)) {
					stillInAnyChannel = true;
					break;
				}
			}

			if (!stillInAnyChannel) {
				delete this.userPresences[presence.userId];
			}
		}
	}
}

export const presenceManager = new PresenceManager();
