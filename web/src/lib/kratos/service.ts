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
	firstName?: string;
	lastName?: string;
	avatarUrl?: string;
	role?: string;
	isEmployee?: boolean;
	fullName?: string;
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

		const user: User = {
			id: this.session?.identity?.id,
			email: traits.email,
			firstName: traits.first_name,
			lastName: traits.last_name,
			avatarUrl: traits.avatar_url,
			role: this.getRole(),
			isEmployee: this.isEmployee(),
		};

		return user;
	}
}
