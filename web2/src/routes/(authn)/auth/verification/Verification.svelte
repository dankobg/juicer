<script lang="ts">
	import type { PageData } from './$types';
	import {
		type UpdateVerificationFlowWithCodeMethod,
		VerificationFlowState,
		instanceOfVerificationFlow,
		instanceOfErrorBrowserLocationChangeRequired,
		instanceOfGenericError
	} from '@ory/client-fetch';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';
	import { browser } from '$app/environment';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Form from '$lib/components/ui/form';

	export let data: PageData;

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

	const verificationFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('code'),
		code: z.string().length(6, { message: 'Code must have exactly 6 characters (check leading or trailing space)' })
	});

	type VerificationFormSchema = z.infer<typeof verificationFormSchema>;

	const initialVerificationForm: VerificationFormSchema = {
		code: '',
		method: 'code',
		csrf_token: data.csrf ?? ''
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
					const verificationFlow = await kratos.updateVerificationFlow({
						flow: data.flow?.id ?? '',
						updateVerificationFlowBody: body
					});

					if (verificationFlow.state === VerificationFlowState.PassedChallenge) {
						for (const node of verificationFlow.ui.nodes) {
							if (node.group === 'code') {
								if (node.attributes.node_type === 'a') {
									window.sessionStorage.setItem(
										'juicer_email_verified',
										'Your E-Mail has been verified! You can now log in'
									);

									window.location.href = node.attributes.href;
									return;
								}
							}
						}
					}
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}

					if (instanceOfVerificationFlow(error)) {
						data.flow = error;

						const nodes = error?.ui?.nodes ?? [];
						const fieldErrors: ValidationErrors<VerificationFormSchema> = {};

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
							handleFlowErrAction(config.routes.verification.path, error.message);
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
			<Card.Title class="text-center text-2xl">Verify account</Card.Title>
			<Card.Description>Verify your account</Card.Description>
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg}
						<Alert.Root variant={msg.type} variantIcon>
							<Alert.Title>{msg.type === 'error' ? 'Unable to verify account' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}

					<div class="grid gap-2">
						<Form.Field form={supForm} name="code">
							<Form.Control let:attrs>
								<Form.Label>Verification code</Form.Label>
								<Input {...attrs} bind:value={$form.code} />
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<Form.Button class="w-full font-bold">Verify account</Form.Button>
				</form>
			</div>
		</Card.Content>
	</Card.Root>
</section>
