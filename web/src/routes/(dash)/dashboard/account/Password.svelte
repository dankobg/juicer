<script lang="ts">
	import type { PageData } from './$types';
	import type {
		ErrorBrowserLocationChangeRequired,
		GenericError,
		SettingsFlow,
		UiNodeInputAttributes,
		UpdateSettingsFlowWithPasswordMethod,
	} from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button, Card } from 'flowbite-svelte';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import InputPassword from '$lib/Inputs/InputPassword.svelte';
	import { toast } from 'svelte-sonner';
	import { config } from '$lib/kratos/config';

	export let data: PageData;
	export let currentFlowForm: 'settings' | 'password' | 'socials' | undefined;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data.flow = null;
		goto(redirectUrl);
		return;
	}

	const passwordFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('password'),
		password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
	});

	type PasswordFormSchema = z.infer<typeof passwordFormSchema>;

	const initialPasswordForm: PasswordFormSchema = {
		password: '',
		method: 'password',
		csrf_token: data.csrf,
	};

	const supForm = superForm(initialPasswordForm, {
		id: 'account_password',
		validators: zod(passwordFormSchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: false,
		async onUpdated({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}

			const url = data.flow?.ui.action;
			const body = form.data as UpdateSettingsFlowWithPasswordMethod & { method: 'password' };

			if (url) {
				try {
					currentFlowForm = 'password';

					const flowResponse = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body,
					});
					data.flow = flowResponse.data;

					toast.success('Account password have been updated');

					if (flowResponse.data.continue_with) {
						for (const item of flowResponse.data.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									return;
							}
						}
					}
				} catch (error) {
					if (!isAxiosError(error)) {
						console.error('updateSettingsFlow: unknown error occurred');
						return;
					}

					if (error.response?.status === 400) {
						const errFlowData: SettingsFlow = error.response.data;
						data.flow = errFlowData;

						const nodes = errFlowData?.ui?.nodes ?? [];
						const fieldErrors: ValidationErrors<PasswordFormSchema> = {};

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes;

							if (attrs.node_type === 'input') {
								for (const msg of node?.messages ?? []) {
									errMsgs.push(msg.text);
									const fieldName = attrs?.name;
									set(fieldErrors, fieldName, errMsgs);
								}
							}
						}

						errors.set(fieldErrors);
						return;
					}

					if (error.response?.status === 422) {
						const err: ErrorBrowserLocationChangeRequired = error.response.data?.error;
						window.location.href = err?.redirect_browser_to ?? '/';
						return;
					}

					if (error.response?.status) {
						const err: GenericError = error.response.data?.error;

						if (err.id === 'session_refresh_required') {
							handleFlowErrAction(
								config.routes.login.path + `?refresh=true&return_to=${window.location.href}`,
								err.message
							);
						}
						if (err.id === 'session_inactive') {
							handleFlowErrAction(config.routes.login.path, err.message);
						}
						if (err.id === 'session_already_available') {
							handleFlowErrAction('/', err.message);
						}
						if (err.id === 'security_csrf_violation' || err.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.settings.path, err.message);
						}
						return;
					}
				}
			}
		},
	});

	const { form, enhance, errors } = supForm;
</script>

<Card>
	<form method="POST" use:enhance class="space-y-6" action="/">
		{#if currentFlowForm === 'password'}
			{#each data?.flow?.ui?.messages ?? [] as msg}
				{@const err = msg.type === 'error'}
				<SimpleAlert kind={msg.type} title={err ? 'Unable to change password' : ''} text={msg.text} />
			{/each}
		{/if}

		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Change password</h3>

		<InputPassword form={supForm} name="password" label="New password" />

		<Button type="submit" class="w-full1 font-bold">Update password</Button>
	</form>
</Card>
