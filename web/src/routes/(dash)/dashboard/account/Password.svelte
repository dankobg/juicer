<script lang="ts">
	import type { PageData } from './$types';
	import type { SettingsFlow, UiNodeInputAttributes, UpdateSettingsFlowWithPasswordMethod } from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button, Card } from 'flowbite-svelte';
	import { superForm } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import InputPassword from '$lib/Inputs/InputPassword.svelte';
	import { toast } from 'svelte-sonner';

	export let data: PageData;

	const PasswordFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('password'),
		password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
	});

	type PasswordFormSchema = z.infer<typeof PasswordFormSchema>;

	const initialPasswordForm: PasswordFormSchema = {
		password: '',
		method: 'password',
		csrf_token: data.csrf,
	};

	const supForm = superForm(initialPasswordForm, {
		id: 'account_password',
		validators: zod(PasswordFormSchema),
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
			const body = form.data as UpdateSettingsFlowWithPasswordMethod & { method: 'password' };

			if (url) {
				try {
					const responseFlow = await kratos.updateSettingsFlow({
						flow: data.flow?.id ?? '',
						updateSettingsFlowBody: body,
					});

					toast.success('Account password have been updated');

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

					reset();
				} catch (error) {
					if (isAxiosError(error)) {
						const flowData = error?.response?.data as SettingsFlow;
						data.flow = flowData;

						console.log('updateSettingsFlow err', flowData);

						const nodes = flowData?.ui?.nodes ?? [];
						const fieldErrors = new Map<keyof PasswordFormSchema, string[]>();

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes as UiNodeInputAttributes;

							for (const msg of node?.messages ?? []) {
								errMsgs.push(msg.text);
								const fieldName = attrs?.name as keyof PasswordFormSchema;
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

	const { form, enhance, errors, reset } = supForm;
</script>

<Card>
	<form method="POST" use:enhance class="space-y-6" action="/">
		<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0">Change password</h3>

		{#each data?.flow?.ui?.messages ?? [] as msg}
			{@const err = msg.type === 'error'}
			<SimpleAlert
				kind={err ? 'error' : 'info'}
				title={err ? 'Unable to change password' : undefined}
				text={msg.text}
			/>
		{/each}

		<InputPassword form={supForm} name="password" label="New password" />

		<Button type="submit" class="w-full1 font-bold">Update password</Button>
	</form>
</Card>
