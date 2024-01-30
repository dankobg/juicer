<script lang="ts">
	import type { PageData } from './$types';
	import type {
		ErrorBrowserLocationChangeRequired,
		GenericError,
		RecoveryFlow,
		UiNodeInputAttributes,
		UpdateRecoveryFlowWithCodeMethod,
	} from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button } from 'flowbite-svelte';
	import { Section, ForgotPasswordHeader, ForgotPassword } from 'flowbite-svelte-blocks';
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

	let codeSentToEmail = false;
	let secondFlowId: string | undefined;

	function handleFlowErrAction(redirectUrl: string, errMsg?: string) {
		if (errMsg) {
			toast.error(errMsg);
		}
		data.flow = null;
		goto(redirectUrl);
		return;
	}

	const recoveryFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.string().min(1, { message: 'method is required' }),
		email: z.string(),
		code: z.string(),
	});

	type RecoveryFormSchema = z.infer<typeof recoveryFormSchema>;

	const initialRecoveryForm: RecoveryFormSchema = {
		email: '',
		code: '',
		method: 'code',
		csrf_token: data.csrf,
	};

	const supForm = superForm(initialRecoveryForm, {
		id: 'auth_recovery',
		validators: zod(recoveryFormSchema),
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
			const body = form.data as UpdateRecoveryFlowWithCodeMethod & { method: 'code' };

			if (!codeSentToEmail) {
				delete body.code;
			} else {
				delete body.email;
			}

			if (url) {
				try {
					const theFlowId = codeSentToEmail ? secondFlowId ?? '' : data.flow?.id ?? '';

					const flowResponse = await kratos.updateRecoveryFlow({
						flow: theFlowId,
						updateRecoveryFlowBody: body,
					});

					data.flow = flowResponse.data;

					codeSentToEmail = true;
					secondFlowId = flowResponse.data.ui.action.split('flow=')[1];
					// goto('/');
				} catch (error) {
					if (!isAxiosError(error)) {
						console.error('updateRecoveryFlow: unknown error occurred');
						return;
					}

					if (error.response?.status === 400) {
						const errFlowData: RecoveryFlow = error.response.data;
						data.flow = errFlowData;

						const nodes = errFlowData?.ui?.nodes ?? [];
						const fieldErrors: ValidationErrors<RecoveryFormSchema> = {};

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
						const err: ErrorBrowserLocationChangeRequired = error.response.data;
						window.location.href = err?.redirect_browser_to ?? '/';
					}

					if (error.response?.status) {
						const err: GenericError = error.response.data?.error;

						if (err.id === 'session_already_available') {
							handleFlowErrAction('/', err.message);
						}
						if (err.id === 'security_csrf_violation' || err.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.recovery.path, err.message);
						}
						return;
					}
				}
			}
		},
	});

	const { form, enhance, errors } = supForm;
</script>

<Section name="forgotpassword">
	<ForgotPasswordHeader src="/images/logo.svg" alt="logo" href="/">Juicer</ForgotPasswordHeader>

	<ForgotPassword>
		{#each data?.flow?.ui?.messages ?? [] as msg}
			{@const err = msg.type === 'error'}
			<SimpleAlert kind={msg.type} title={err ? 'Unable to recover' : ''} text={msg.text} />
		{/each}

		<h1 class="mb-1 text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
			{#if codeSentToEmail}
				Complete the recovery flow
			{:else}
				Forgot your password?
			{/if}
		</h1>

		<p class="font-light text-gray-500 dark:text-gray-400">
			{#if codeSentToEmail}
				Enter the recovery code that was sent via email
			{:else}
				Don't fret! Just type in your email and we will send you a code to reset your password!
			{/if}
		</p>

		<form method="POST" use:enhance class="mt-4 space-y-4 lg:mt-5 md:space-y-5">
			{#if codeSentToEmail}
				<InputText form={supForm} name="code" label="Recovery code" />
			{:else}
				<InputEmail form={supForm} name="email" label="Your email" />
			{/if}

			<Button type="submit" color="red" class="font-bold">
				{#if codeSentToEmail}
					Recover account
				{:else}
					Send recovery code
				{/if}
			</Button>
		</form>
	</ForgotPassword>
</Section>
