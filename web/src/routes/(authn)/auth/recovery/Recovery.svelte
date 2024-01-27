<script lang="ts">
	import type { PageData } from './$types';
	import type { RecoveryFlow, UiNodeInputAttributes, UpdateRecoveryFlowWithCodeMethod } from '@ory/client';
	import { goto } from '$app/navigation';
	import { kratos } from '$lib/kratos/client';
	import { Button } from 'flowbite-svelte';
	import { Section, ForgotPasswordHeader, ForgotPassword } from 'flowbite-svelte-blocks';
	import { superForm } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError } from '$lib/kratos/helpers';
	import InputEmail from '$lib/Inputs/InputEmail.svelte';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import { toast } from 'svelte-sonner';
	import InputText from '$lib/Inputs/InputText.svelte';

	export let data: PageData;

	let codeSentToEmail = false;
	let secondFlowId: string | undefined;

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

					codeSentToEmail = true;
					secondFlowId = flowResponse.data.ui.action.split('flow=')[1];

					console.log('updateRecoveryFlow', flowResponse);
					// goto('/');
				} catch (error) {
					if (isAxiosError(error)) {
						if (error?.response?.status === 422 && error?.response?.data?.redirect_browser_to) {
							console.log(error.response?.data);
							goto(error.response?.data?.redirect_browser_to);
							return;
						}

						const flowData = error?.response?.data as RecoveryFlow;
						data.flow = flowData;

						const nodes = flowData?.ui?.nodes ?? [];
						const fieldErrors = new Map<keyof RecoveryFormSchema, string[]>();

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes as UiNodeInputAttributes;

							for (const msg of node?.messages ?? []) {
								errMsgs.push(msg.text);
								const fieldName = attrs?.name as keyof RecoveryFormSchema;
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

<Section name="forgotpassword">
	<ForgotPasswordHeader src="/images/logo.svg" alt="logo" href="/">Juicer</ForgotPasswordHeader>

	<ForgotPassword>
		{#if codeSentToEmail}
			<SimpleAlert kind="success" title="Recovery code was sent to your email successfully" text="" />
		{/if}

		<h1 class="mb-1 text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
			{#if codeSentToEmail}
				Complete the recovery flow
			{:else}
				Forgot your password?
			{/if}
		</h1>

		{#each data?.flow?.ui?.messages ?? [] as msg}
			{@const err = msg.type === 'error'}
			<SimpleAlert kind={err ? 'error' : 'info'} title={err ? 'Unable to recover' : undefined} text={msg.text} />
		{/each}

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
