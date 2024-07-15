<script lang="ts">
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
	import {
		type UpdateLoginFlowWithPasswordMethod,
		instanceOfErrorBrowserLocationChangeRequired,
		instanceOfGenericError,
		instanceOfLoginFlow
	} from '@ory/client-fetch';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';
	import { kratos } from '$lib/kratos/client';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { providers } from '$lib/kratos/helpers';
	import { toast } from 'svelte-sonner';
	import { browser } from '$app/environment';
	import { Input } from '$lib/components/ui/input';
	import * as Tooltip from '$lib/components/ui/tooltip';
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

	const loginFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.string().min(1, { message: 'method is required' }),
		identifier: z.string().min(1, { message: 'E-Mail is required' }).email({ message: 'E-Mail must be a valid email' }),
		password: z.string().min(8, { message: 'Password must have min. 8 characters' })
	});

	type LoginFormSchema = z.infer<typeof loginFormSchema>;

	const initialLoginForm: LoginFormSchema = {
		identifier: '',
		password: '',
		method: 'password',
		csrf_token: data.csrf ?? ''
	};

	const supForm = superForm(initialLoginForm, {
		id: 'auth_login',
		validators: zod(loginFormSchema),
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
			const body = form.data as UpdateLoginFlowWithPasswordMethod & { method: 'password' };

			if (url) {
				try {
					await kratos.updateLoginFlow({
						flow: data.flow?.id ?? '',
						updateLoginFlowBody: body
					});
					goto(data.flow?.return_to ?? '/', { invalidateAll: true });
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}

					if (instanceOfLoginFlow(error)) {
						data.flow = error;
						const nodes = error.ui.nodes ?? [];
						const fieldErrors: ValidationErrors<LoginFormSchema> = {};
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
						if (error.id === 'session_already_available') {
							handleFlowErrAction('/', error.message);
						}
						if (error.id === 'security_csrf_violation' || error.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.login.path, error.message);
						}
						return;
					}
				}
			}
		}
	});

	const { form, enhance, errors } = supForm;

	let emailVerified = false;
	let emailVerifiedMsg = '';

	onMount(() => {
		const val = window.sessionStorage.getItem('juicer_email_verified');

		if (val) {
			emailVerified = true;
			emailVerifiedMsg = val;
		}

		return () => {
			emailVerified = false;
			emailVerifiedMsg = '';
			sessionStorage.removeItem('juicer_email_verified');
		};
	});
</script>

<svelte:window
	on:beforeunload={() => {
		sessionStorage.removeItem('juicer_email_verified');
	}}
/>

<section class="grid h-screen place-content-center gap-4">
	<a href="/" class="justify-self-center">
		<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
	</a>

	<Card.Root class="mx-auto max-w-sm">
		<Card.Header>
			<Card.Title class="text-center text-2xl">Login</Card.Title>
			<Card.Description>Sign in to your account</Card.Description>
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg}
						{@const err = msg.type === 'error'}
						{@const clr = msg.type === 'error' ? 'red' : msg.type === 'success' ? 'green' : 'blue'}
						<Alert.Root class="border border-{clr}-600 bg-{clr}-50 text-{clr}-600 dark:bg-{clr}-950">
							<Alert.Title>{err ? 'Unable to log in' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}

					{#if emailVerified}
						<Alert.Root class="border border-green-600 bg-green-50 text-green-600 dark:bg-green-950">
							<Alert.Title>Success</Alert.Title>
							<Alert.Description>{emailVerifiedMsg}</Alert.Description>
						</Alert.Root>
					{/if}

					<div class="grid gap-2">
						<Form.Field form={supForm} name="identifier">
							<Form.Control let:attrs>
								<Form.Label>E-Mail</Form.Label>
								<Input type="email" {...attrs} bind:value={$form.identifier} />
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<div class="flex items-center">
							<a href={config.routes.recovery.path} class="ml-auto inline-block text-sm underline">
								Forgot your password?
							</a>
						</div>
						<Form.Field form={supForm} name="password">
							<Form.Control let:attrs>
								<Form.Label>Password</Form.Label>
								<Input type="password" {...attrs} bind:value={$form.password} />
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<Form.Button class="w-full font-bold">Login</Form.Button>
				</form>

				<section>
					<div class="inline-flex w-full items-center justify-center">
						<hr class="my-8 h-px w-full border-0 bg-gray-200 dark:bg-gray-700" />
						<span
							class="bg-card dark:bg-card absolute left-1/2 -translate-x-1/2 px-3 font-medium text-gray-900 dark:text-white"
						>
							or login with
						</span>
					</div>

					<div class="align-center flex justify-between">
						{#each providers as provider}
							<form
								action={data.flow?.ui.action}
								method="post"
								encType="application/x-www-form-urlencoded"
								data-provider={provider.name}
							>
								<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
								<input type="hidden" name="provider" value={provider.name} readonly required />
								<Tooltip.Root>
									<Tooltip.Trigger>
										<button type="submit" class="h-12 w-12">
											<img
												class="h-full w-full object-cover"
												src="/images/providers/{provider.name}.svg"
												alt="login with {provider.label}"
											/>
										</button>
									</Tooltip.Trigger>
									<Tooltip.Content>
										{provider.label}
									</Tooltip.Content>
								</Tooltip.Root>
							</form>
						{/each}
					</div>
				</section>
			</div>
			<div class="mt-4 text-center text-sm">
				Don't have an account yet?
				<a href={config.routes.registration.path} class="underline">Sign up</a>
			</div>
		</Card.Content>
	</Card.Root>
</section>
