<script lang="ts">
	import type { PageData } from './$types';
	import {
		instanceOfErrorBrowserLocationChangeRequired,
		instanceOfGenericError,
		instanceOfRecoveryFlow,
		type UpdateRecoveryFlowWithCodeMethod
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

	let codeSentToEmail = false;
	let secondFlowId: string | undefined;

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

	const recoveryFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.string().min(1, { message: 'method is required' }),
		email: z.string(),
		code: z.string()
	});

	type RecoveryFormSchema = z.infer<typeof recoveryFormSchema>;

	const initialRecoveryForm: RecoveryFormSchema = {
		email: '',
		code: '',
		method: 'code',
		csrf_token: data.csrf ?? ''
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
					const theFlowId = codeSentToEmail ? (secondFlowId ?? '') : (data.flow?.id ?? '');

					const recoveryFlow = await kratos.updateRecoveryFlow({
						flow: theFlowId,
						updateRecoveryFlowBody: body
					});

					data.flow = recoveryFlow;

					codeSentToEmail = true;
					secondFlowId = recoveryFlow.ui.action.split('flow=')[1];
					// goto('/');
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}

					if (instanceOfRecoveryFlow(error)) {
						data.flow = error;

						const nodes = error?.ui?.nodes ?? [];
						const fieldErrors: ValidationErrors<RecoveryFormSchema> = {};

						for (const node of nodes) {
							const errMsgs: string[] = [];
							if (node.attributes.node_type === 'input') {
								for (const msg of node?.messages ?? []) {
									errMsgs.push(msg.text);
									const fieldName = node.attributes.name;
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
						if (error.id === 'session_already_available') {
							handleFlowErrAction('/', error.message);
						} else if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.recovery.path, error.message);
						}
					}
				}
			}
		}
	});

	const { form, enhance, errors } = supForm;
</script>

<section class="grid h-screen place-content-center gap-4">
	<a href="/" class="justify-self-center">
		<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
	</a>

	<Card.Root class="mx-auto max-w-sm">
		<Card.Header>
			{#if codeSentToEmail}
				<Card.Title class="text-center text-2xl">Complete the recovery flow</Card.Title>
				<Card.Description>Enter the recovery code that was sent via email</Card.Description>
			{:else}
				<Card.Title class="text-center text-2xl">Forgot your password?</Card.Title>
				<Card.Description>Enter your email and you will receive a code to reset your password</Card.Description>
			{/if}
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg}
					<Alert.Root variant="{msg.type}">
						<Alert.Title>{msg.type === 'error' ? 'Unable to recover account' : ''}</Alert.Title>
						<Alert.Description>{msg.text}</Alert.Description>
					</Alert.Root>
					{/each}

					<div class="grid gap-2">
						{#if codeSentToEmail}
							<Form.Field form={supForm} name="code">
								<Form.Control let:attrs>
									<Form.Label>E-Mail</Form.Label>
									<Input type="Recovery code" {...attrs} bind:value={$form.code} />
								</Form.Control>
								<Form.FieldErrors />
							</Form.Field>
						{:else}
							<Form.Field form={supForm} name="email">
								<Form.Control let:attrs>
									<Form.Label>E-Mail</Form.Label>
									<Input type="email" {...attrs} bind:value={$form.email} />
								</Form.Control>
								<Form.FieldErrors />
							</Form.Field>
						{/if}
					</div>
					<Form.Button class="w-full font-bold">
						{#if codeSentToEmail}
							Recover account
						{:else}
							Send recovery code
						{/if}
					</Form.Button>
				</form>
			</div>
		</Card.Content>
	</Card.Root>
</section>
