<script lang="ts">
	import type { PageData } from './$types';
	import {
		instanceOfErrorBrowserLocationChangeRequired,
		instanceOfGenericError,
		instanceOfSettingsFlow,
		type UpdateSettingsFlowWithProfileMethod
	} from '@ory/client-fetch';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { toast } from 'svelte-sonner';
	import { config } from '$lib/kratos/config';
	import { browser } from '$app/environment';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';

	export let data: PageData;
	export let currentFlowForm: 'settings' | 'password' | 'socials' | undefined = undefined;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data.flow = null;

		if (browser) {
			goto(redirectUrl);
		}

		return;
	}

	const settingsFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('profile'),
		traits: z.object({
			first_name: z.string(),
			last_name: z.string(),
			email: z.string().min(1, { message: 'E-Mail is required' }).email({ message: 'E-Mail must be a valid email' }),
			avatar_url: z.string()
		})
	});

	type SettingsFormSchema = z.infer<typeof settingsFormSchema>;

	const initialSettingsForm: SettingsFormSchema = {
		method: 'profile',
		csrf_token: data.csrf ?? '',
		traits: {
			first_name: data.flow?.identity.traits['first_name'] ?? '',
			last_name: data.flow?.identity.traits['last_name'] ?? '',
			email: data.flow?.identity.traits['email'] ?? '',
			avatar_url: ''
		}
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

					const settingsFlow = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body
					});
					data.flow = settingsFlow;

					toast.success('Account settings have been updated');

					if (settingsFlow.continue_with) {
						for (const item of settingsFlow.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									return;
							}
						}
					}
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}

					if (instanceOfSettingsFlow(error)) {
						data.flow = error;
						const nodes = error?.ui?.nodes ?? [];
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

					if (instanceOfErrorBrowserLocationChangeRequired(error)) {
						window.location.href = error.redirect_browser_to || '/';
						return;
					}

					if (instanceOfGenericError(error)) {
						if (error.id === 'session_refresh_required') {
							handleFlowErrAction(
								config.routes.login.path + `?refresh=true&return_to=${window.location.href}`,
								error.message
							);
						} else if (error.id === 'session_inactive') {
							handleFlowErrAction(config.routes.login.path, error.message);
						} else if (error.id === 'session_already_available') {
							handleFlowErrAction('/', error.message);
						} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.settings.path, error.message);
						}
					}
				}
			}
		}
	});

	const { form, enhance, errors } = supForm;
</script>

<Card.Root class="mx-auto max-w-sm">
	<Card.Header>
		<Card.Title class="text-center text-2xl">Account settings</Card.Title>
		<Card.Description>Change account settings</Card.Description>
	</Card.Header>

	<Card.Content>
		<div class="grid gap-4">
			<form method="POST" use:enhance class="grid gap-4">
				{#each data?.flow?.ui?.messages ?? [] as msg}
					{@const err = msg.type === 'error'}
					{@const clr = msg.type === 'error' ? 'red' : msg.type === 'success' ? 'green' : 'blue'}
					<Alert.Root class="border border-{clr}-600 bg-{clr}-50 text-{clr}-600 dark:bg-{clr}-950">
						<Alert.Title>{err ? 'Unable to change settings' : ''}</Alert.Title>
						<Alert.Description>{msg.text}</Alert.Description>
					</Alert.Root>
				{/each}

				<div class="grid gap-2">
					<Form.Field form={supForm} name="traits.first_name">
						<Form.Control let:attrs>
							<Form.Label>First name</Form.Label>
							<Input {...attrs} bind:value={$form.traits.first_name} />
						</Form.Control>
						<Form.FieldErrors />
					</Form.Field>
				</div>
				<Form.Button class="w-full font-bold">Update settings</Form.Button>
			</form>
		</div>
	</Card.Content>
</Card.Root>
