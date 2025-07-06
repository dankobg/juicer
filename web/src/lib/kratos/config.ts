const baseUrl = import.meta.env['VITE_PUBLIC_SITE_URL'] as string;
const publicUrl = import.meta.env['VITE_PUBLIC_KRATOS_PUBLIC_URL'] as string;
const adminUrl = import.meta.env['VITE_PUBLIC_KRATOS_ADMIN_URL'] as string;

export const config = {
	kratos: {
		baseUrl,
		publicUrl,
		adminUrl
	},

	routes: {
		registration: {
			path: '/auth/registration',
			selfServiceUrl: `${publicUrl}/self-service/registration/browser`
		},
		login: {
			path: '/auth/login',
			selfServiceUrl: `${publicUrl}/self-service/login/browser`
		},
		logout: {
			path: '/auth/logout',
			selfServiceUrl: `${publicUrl}/self-service/logout/browser`
		},
		verification: {
			path: '/auth/verification',
			selfServiceUrl: `${publicUrl}/self-service/verification/browser`
		},
		recovery: {
			path: '/auth/recovery',
			selfServiceUrl: `${publicUrl}/self-service/recovery/browser`
		},
		settings: {
			path: '/dashboard/account',
			selfServiceUrl: `${publicUrl}/self-service/settings/browser`
		},
		dashboard: {
			path: '/dashboard'
		},
		refresh: {
			selfServiceUrl: (returnTo?: string): string =>
				`${publicUrl}/self-service/browser/flows/login?refresh=true&return_to=${returnTo || baseUrl}`
		}
	},

	labels: {
		to_verify: {
			label: 'Email',
			priority: 100
		},
		csrf_token: {
			label: '',
			priority: 100
		},
		'traits.email': {
			label: 'Email',
			priority: 90
		},
		email: {
			label: 'Email',
			priority: 80
		},
		identifier: {
			label: 'Email',
			priority: 80
		},
		password: {
			label: 'Password',
			priority: 80
		},
		'traits.first_name': {
			label: 'First name',
			priority: 95
		},
		'traits.last_name': {
			label: 'Last name',
			priority: 96
		}
	}
};
