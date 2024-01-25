<script lang="ts">
	import type { PageData } from './$types';
	import type { VerificationFlow, UiNodeInputAttributes, UpdateVerificationFlowWithCodeMethod } from '@ory/client';
	import { kratos } from '$lib/kratos/client';
	import { Button } from 'flowbite-svelte';
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { superForm } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import InputText from '$lib/Inputs/InputText.svelte';
	import { toast } from 'svelte-sonner';

	export let data: PageData;

	const verificationFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('code'),
		code: z.string().length(6, { message: 'Code must have exactly 6 characters (check leading or trailing space)' }),
	});

	type VerificationFormSchema = z.infer<typeof verificationFormSchema>;

	const initialVerificationForm: VerificationFormSchema = {
		code: '',
		method: 'code',
		csrf_token: data.csrf,
	};

	const supForm = superForm(initialVerificationForm, {
		id: 'auth_verification',
		validators: zod(verificationFormSchema),
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
			const body = form.data as UpdateVerificationFlowWithCodeMethod & { method: 'code' };

			if (url) {
				try {
					const responseFlow = await kratos.updateVerificationFlow({
						flow: data.flow?.id ?? '',
						updateVerificationFlowBody: body,
					});

					console.log('updateVerificationFlow success', responseFlow);

					if (responseFlow.data.state === 'passed_challenge') {
						for (const item of responseFlow.data.ui.nodes) {
							if (item.group === 'code') {
								const attrs = item.attributes;

								if (attrs.node_type === 'a') {
									window.sessionStorage.setItem(
										'juicer_email_verified',
										'Your E-Mail has been verified! You can now log in'
									);
									window.location.href = attrs.href;
									return;
								}
							}
						}
					}
				} catch (error) {
					if (isAxiosError(error)) {
						const flowData = error?.response?.data as VerificationFlow;
						data.flow = flowData;

						console.log('updateVerificationFlow err', flowData);

						const nodes = flowData?.ui?.nodes ?? [];
						const fieldErrors = new Map<keyof VerificationFormSchema, string[]>();

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes as UiNodeInputAttributes;

							for (const msg of node?.messages ?? []) {
								errMsgs.push(msg.text);
								const fieldName = attrs?.name as keyof VerificationFormSchema;
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

<Section name="login">
	<Register href="/">
		<svelte:fragment slot="top">
			<img class="w-8 h-8 mr-2" src="/images/logo.svg" alt="logo" />
			Juicer
		</svelte:fragment>

		<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
			<form method="POST" use:enhance class="flex flex-col space-y-6" action="/">
				<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0 text-center">Verifify your account</h3>

				{#each data?.flow?.ui?.messages ?? [] as msg}
					{@const err = msg.type === 'error'}
					<SimpleAlert kind={err ? 'error' : 'info'} title={err ? 'Unable to sign up' : ''} text={msg.text} />
				{/each}

				<InputText form={supForm} name="code" label="Verification code" />

				<Button type="submit" class="w-full1 font-bold">Verify account</Button>
			</form>
		</div>
	</Register>
</Section>
