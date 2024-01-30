<script lang="ts">
	import type { PageData } from './$types';
	import type {
		ErrorBrowserLocationChangeRequired,
		GenericError,
		SettingsFlow,
		UiNodeInputAttributes,
		UpdateSettingsFlowWithProfileMethod,
	} from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button, Card } from 'flowbite-svelte';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import InputEmail from '$lib/Inputs/InputEmail.svelte';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import { toast } from 'svelte-sonner';
	import InputText from '$lib/Inputs/InputText.svelte';
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

	const settingsFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('profile'),
		traits: z.object({
			first_name: z.string(),
			last_name: z.string(),
			email: z.string().min(1, { message: 'E-Mail is required' }).email({ message: 'E-Mail must be a valid email' }),
			avatar_url: z.string(),
		}),
	});

	type SettingsFormSchema = z.infer<typeof settingsFormSchema>;

	const initialSettingsForm: SettingsFormSchema = {
		method: 'profile',
		csrf_token: data.csrf,
		traits: {
			first_name: data.flow?.identity.traits['first_name'] ?? '',
			last_name: data.flow?.identity.traits['last_name'] ?? '',
			email: data.flow?.identity.traits['email'] ?? '',
			avatar_url: '',
		},
	};

	const supForm = superForm(initialSettingsForm, {
		id: 'account_settings',
		validators: zod(settingsFormSchema),
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
			const body = form.data as UpdateSettingsFlowWithProfileMethod & { method: 'profile' };

			if (url) {
				try {
					currentFlowForm = 'settings';

					const flowResponse = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body,
					});
					data.flow = flowResponse.data;

					toast.success('Account settings have been updated');

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
						const fieldErrors: ValidationErrors<SettingsFormSchema> = {};

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
		{#if currentFlowForm === 'settings'}
			{#each data?.flow?.ui?.messages ?? [] as msg}
				{@const err = msg.type === 'error'}
				<SimpleAlert kind={msg.type} title={err ? 'Unable to change settings' : ''} text={msg.text} />
			{/each}
		{/if}

		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Account settings</h3>

		<InputText form={supForm} name="traits.first_name" label="First name" />
		<InputText form={supForm} name="traits.last_name" label="Last name" />
		<InputEmail form={supForm} name="traits.email" label="E-Mail" />

		<Button type="submit" class="w-full1 font-bold">Update settings</Button>
	</form>
</Card>
