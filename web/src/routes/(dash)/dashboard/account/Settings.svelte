<script lang="ts">
	import type { PageData } from './$types';
	import type { SettingsFlow, UiNodeInputAttributes, UpdateSettingsFlowWithProfileMethod } from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button, Card } from 'flowbite-svelte';
	import { superForm } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import InputEmail from '$lib/Inputs/InputEmail.svelte';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import { toast } from 'svelte-sonner';
	import InputText from '$lib/Inputs/InputText.svelte';
	import SuperDebug from 'sveltekit-superforms/client/SuperDebug.svelte';

	export let data: PageData;

	const SettingsFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('profile'),
		traits: z.object({
			first_name: z.string(),
			last_name: z.string(),
			email: z.string().min(1, { message: 'E-Mail is required' }).email({ message: 'E-Mail must be a valid email' }),
			avatar_url: z.string(),
		}),
	});

	type SettingsFormSchema = z.infer<typeof SettingsFormSchema>;

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
		validators: zod(SettingsFormSchema),
		SPA: true,
		dataType: 'json',
		errorSelector: '[data-invalid]',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		async onUpdated({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}

			const url = data.flow?.ui.action;
			const body = form.data as UpdateSettingsFlowWithProfileMethod & { method: 'profile' };

			if (url) {
				try {
					const responseFlow = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body,
					});

					toast.success('Account settings have been updated');

					console.log('updateSettingsFlow success', responseFlow);

					if (responseFlow.data.continue_with) {
						for (const item of responseFlow.data.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									return;
							}
						}
					}

					// reset();
				} catch (error) {
					if (isAxiosError(error)) {
						const flowData = error?.response?.data as SettingsFlow;
						data.flow = flowData;

						console.log('updateSettingsFlow err', flowData);

						if (flowData?.error?.code === 403 && flowData?.error?.id === 'session_refresh_required') {
							toast.error(flowData?.reason ?? '');
							window.location.href = flowData?.redirect_browser_to ?? '/';
						}

						const nodes = flowData?.ui?.nodes ?? [];

						const fieldErrors = new Map<keyof SettingsFormSchema, string[]>();

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes as UiNodeInputAttributes;

							for (const msg of node?.messages ?? []) {
								errMsgs.push(msg.text);
								const fieldName = attrs?.name as keyof SettingsFormSchema;
								fieldErrors.set(fieldName, errMsgs);
							}
						}

						for (const [k, v] of fieldErrors.entries()) {
							const srvErrors = {};
							set(srvErrors, k, v.join('; '));

							$errors = {
								...$errors,
								...srvErrors,
							};
						}
					}
				}
			}
		},
	});

	const { form, enhance, errors } = supForm;
</script>

<SuperDebug data={$form} />

<Card>
	<form method="POST" use:enhance class="space-y-6" action="/">
		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Account settings</h3>

		{#each data?.flow?.ui?.messages ?? [] as msg}
			{@const err = msg.type === 'error'}
			<SimpleAlert
				kind={err ? 'error' : 'info'}
				title={err ? 'Unable to change settings' : undefined}
				text={msg.text}
			/>
		{/each}

		<InputText form={supForm} name="traits.first_name" label="First name" />
		<InputText form={supForm} name="traits.last_name" label="Last name" />
		<InputEmail form={supForm} name="traits.email" label="Your email" />

		<Button type="submit" class="w-full1 font-bold">Update settings</Button>
	</form>
</Card>
