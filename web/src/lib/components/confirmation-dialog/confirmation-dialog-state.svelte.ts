import type { Snippet } from 'svelte';

type ConfirmationDialogOptions = {
	title?: string;
	cancelText?: string;
	confirmText?: string;
	description?: Snippet;
	destructive?: boolean;
	onCancel?: VoidFunction;
	onConfirm?: VoidFunction;
};

export class Confirmation {
	open = $state<boolean>(false);
	title = $state<string>('Are you sure?');
	cancelText = $state<string>('Cancel');
	confirmText = $state<string>('Continue');
	description?: Snippet;
	destructive? = $state<boolean>(false);
	onCancel?: () => void = () => {
		this.open = false;
	};
	onConfirm?: () => void;

	openDialog(obj: ConfirmationDialogOptions): void {
		this.open = true;
		this.title = obj.title ?? 'Are you sure?';
		this.cancelText = obj.cancelText ?? 'Cancel';
		this.confirmText = obj.confirmText ?? 'Continue';
		if (obj.description) {
			this.description = obj.description;
		}
		if (obj.destructive) {
			this.destructive = obj.destructive;
		}
		if (obj.onCancel) {
			this.onCancel = obj.onCancel;
		}
		if (obj.onConfirm) {
			this.onConfirm = obj.onConfirm;
		}
	}

	closeDialog(): void {
		this.open = false;
	}
}

export const confirmation = new Confirmation();
