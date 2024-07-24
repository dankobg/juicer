export function getInitials(name?: string): string {
	if (!name) {
		return '';
	}

	const parts = name.trim().toUpperCase().split(/\s+/);

	const first = parts[0].charAt(0);
	if (parts.length === 1) {
		return first;
	}

	const last = parts.at(-1)?.charAt(0);
	return first + last;
}
