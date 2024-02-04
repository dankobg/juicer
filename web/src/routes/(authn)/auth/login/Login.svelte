<script lang="ts">
	import { type AxiosError } from 'axios';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
	import type {
		ErrorBrowserLocationChangeRequired,
		GenericError,
		LoginFlow,
		UpdateLoginFlowWithPasswordMethod,
	} from '@ory/client';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';
	import { kratos } from '$lib/kratos/client';
	import Tooltip from 'flowbite-svelte/Tooltip.svelte';
	import Button from 'flowbite-svelte/Button.svelte';
	import Section from 'flowbite-svelte-blocks/Section.svelte';
	import Register from 'flowbite-svelte-blocks/Register.svelte';
	import { superForm, type ValidationErrors } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { providers } from '$lib/kratos/helpers';
	import InputEmail from '$lib/Inputs/InputEmail.svelte';
	import InputPassword from '$lib/Inputs/InputPassword.svelte';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import { toast } from 'svelte-sonner';
	import { browser } from '$app/environment';

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
		password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
	});

	type LoginFormSchema = z.infer<typeof loginFormSchema>;

	const initialLoginForm: LoginFormSchema = {
		identifier: '',
		password: '',
		method: 'password',
		csrf_token: data.csrf ?? '',
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
						updateLoginFlowBody: body,
					});
					goto(data.flow?.return_to ?? '/');
				} catch (error) {
					const axiosErr = error as AxiosError<GenericError>;
					if (!axiosErr?.isAxiosError) {
						console.error('createBrowserLoginFlow: unknown error occurred');
						return;
					}

					if (axiosErr.response?.status === 400) {
						const axiosErr = error as AxiosError<LoginFlow>;
						const errFlow = axiosErr.response?.data;

						data.flow = errFlow;

						const nodes = errFlow?.ui?.nodes ?? [];
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

					if (axiosErr.response?.status === 422) {
						const axiosErr = error as AxiosError<ErrorBrowserLocationChangeRequired>;
						const err = axiosErr.response?.data;

						window.location.href = err?.redirect_browser_to ?? '/';
						return;
					}

					if (axiosErr.response?.status) {
						const axiosErr = error as AxiosError<GenericError>;
						const err = axiosErr.response?.data;

						if (err?.id === 'session_already_available') {
							handleFlowErrAction('/', err.message);
						}
						if (err?.id === 'security_csrf_violation' || err?.id === 'security_identity_mismatch') {
							handleFlowErrAction(config.routes.login.path, err.message);
						}
						return;
					}
				}
			}
		},
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

				{#if emailVerified}
					<SimpleAlert kind="success" title={emailVerifiedMsg} />
				{/if}

				<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0 text-center">Login</h3>

				<InputEmail form={supForm} name="identifier" label="Your email" />
				<InputPassword form={supForm} name="password" label="Your password" />

				<div class="flex items-start">
					<a href={config.routes.recovery.path} class="ml-auto text-sm text-blue-700 hover:underline dark:text-blue-500"
						>Forgot password?</a
					>
				</div>

				<Button type="submit" class="w-full1 font-bold">Sign in</Button>

				<div class="text-sm font-medium text-gray-500 dark:text-gray-300">
					Don’t have an account yet? <a
						href={config.routes.registration.path}
						class="font-medium text-primary-600 hover:underline dark:text-primary-500">Sign up</a
					>
				</div>

				<section>
					<div class="inline-flex items-center justify-center w-full">
						<hr class="w-full h-px my-8 bg-gray-200 border-0 dark:bg-gray-700" />
						<span
							class="absolute px-3 font-medium text-gray-900 -translate-x-1/2 bg-white left-1/2 dark:text-white dark:bg-gray-900"
						>
							or login with
						</span>
					</div>

					<div class="flex align-center justify-between">
						{#each providers as provider}
							<form
								action={data.flow?.ui.action}
								method="post"
								encType="application/x-www-form-urlencoded"
								data-provider={provider.name}
							>
								<input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
								<input type="hidden" name="provider" value={provider.name} readonly required />
								<button id={provider.name} type="submit" class="w-12 h-12">
									<img
										class="w-full h-full object-cover"
										src="/images/providers/{provider.name}.svg"
										alt="login with {provider.label}"
									/>
								</button>
								<Tooltip triggeredBy="#{provider.name}">{provider.label}</Tooltip>
							</form>
						{/each}
					</div>
				</section>
			</form>
		</div>
	</Register>
</Section>
