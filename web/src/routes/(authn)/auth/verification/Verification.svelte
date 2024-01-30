<script lang="ts">
	import type { PageData } from './$types';
	import {
		type VerificationFlow,
		type UiNodeInputAttributes,
		type UpdateVerificationFlowWithCodeMethod,
		VerificationFlowState,
		type GenericError,
		type ErrorBrowserLocationChangeRequired,
	} from '@ory/client';
	import { kratos } from '$lib/kratos/client';
	import { Button } from 'flowbite-svelte';
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import InputText from '$lib/Inputs/InputText.svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';

	export let data: PageData;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data.flow = null;
		goto(redirectUrl);
		return;
	}

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
					const flowResponse = await kratos.updateVerificationFlow({
						flow: data.flow?.id ?? '',
						updateVerificationFlowBody: body,
					});

					if (flowResponse.data.state === VerificationFlowState.PassedChallenge) {
						for (const item of flowResponse.data.ui.nodes) {
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
					if (!isAxiosError(error)) {
						console.error('updateVerificationFlow: unknown error occurred');
						return;
					}

					if (error.response?.status === 400) {
						const errFlowData: VerificationFlow = error.response.data;
						data.flow = errFlowData;

						const nodes = errFlowData?.ui?.nodes ?? [];
						const fieldErrors: ValidationErrors<VerificationFormSchema> = {};

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

						if (err.id === 'session_already_available') {
							handleFlowErrAction('/', err.message);
						}
						if (err.id === 'security_csrf_violation' || err.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.verification.path, err.message);
						}
						return;
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
				{#each data?.flow?.ui?.messages ?? [] as msg}
					{@const err = msg.type === 'error'}
					<SimpleAlert kind={msg.type} title={err ? 'Unable to sign up' : ''} text={msg.text} />
				{/each}

				<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0 text-center">Verifify your account</h3>

				<InputText form={supForm} name="code" label="Verification code" />

				<Button type="submit" class="w-full1 font-bold">Verify account</Button>
			</form>
		</div>
	</Register>
</Section>
