import type { Session } from '@ory/client';

export interface CustomTraits {
	email?: string;
	first_name?: string;
	last_name?: string;
	avatar_url?: string;
}

export interface User {
	id?: string;
	email?: string;
	first_name?: string;
	last_name?: string;
	avatar_url?: string;
	role?: string;
	isEmployee?: boolean;
}

export interface SessionService {
	user: User | null;
	session: Session | null;
}

export class KratosService implements SessionService {
	constructor(public session: Session | null) {}

	private getRole(): string {
		switch (this.session?.identity?.schema_id) {
			case 'employee':
				return 'employee';
			case 'default':
				return 'default';
			default:
				return '';
		}
	}

	private isEmployee(): boolean {
		return this.getRole() === 'employee';
	}

	public get user(): User | null {
		if (!this.session) {
			return null;
		}

		const traits = this.session?.identity?.traits as CustomTraits;

		return {
			id: this.session?.identity?.id,
			email: traits.email,
			first_name: traits.first_name,
			last_name: traits.last_name,
			avatar_url: traits.avatar_url,
			role: this.getRole(),
			isEmployee: this.isEmployee(),
		};
	}
}
