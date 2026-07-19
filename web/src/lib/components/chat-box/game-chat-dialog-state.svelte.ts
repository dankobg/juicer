export const gameChatDialog = $state<{ open: boolean }>({ open: false });

export function toggleChatDialog(): void {
	gameChatDialog.open = !gameChatDialog.open;
}
