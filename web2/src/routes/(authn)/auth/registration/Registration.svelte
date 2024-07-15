<script lang="ts">
	import type { PageData } from './$types';
	import {
		type UpdateRegistrationFlowWithPasswordMethod,
		instanceOfRegistrationFlow,
		instanceOfErrorBrowserLocationChangeRequired,
		instanceOfGenericError
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

	const registrationFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('password'),
		password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
		traits: z.object({
			first_name: z.string(),
			last_name: z.string(),
			email: z.string().min(1, { message: 'E-Mail is required' }).email({ message: 'E-Mail must be a valid email' }),
			avatar_url: z.string()
		}),
		transient_payload: z.object({}).optional()
	});

	type RegistrationFormSchema = z.infer<typeof registrationFormSchema>;

	const initialRegistrationForm: RegistrationFormSchema = {
		password: '',
		method: 'password',
		csrf_token: data.csrf ?? '',
		traits: {
			first_name: '',
			last_name: '',
			email: '',
			avatar_url: ''
		},
		transient_payload: {}
	};

	const supForm = superForm(initialRegistrationForm, {
		id: 'auth_registration',
		validators: zod(registrationFormSchema),
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
			const body = form.data as UpdateRegistrationFlowWithPasswordMethod & { method: 'password' };

			if (url) {
				try {
					const successfulNativeRegistration = await kratos.updateRegistrationFlow({
						flow: data.flow?.id ?? '',
						updateRegistrationFlowBody: body
					});

					if (successfulNativeRegistration.continue_with) {
						for (const item of successfulNativeRegistration.continue_with) {
							switch (item.action) {
								case 'show_verification_ui':
									if (item?.flow?.id) {
										goto(item?.flow?.url as string);
									}
									return;
							}
						}
					}
				} catch (error: unknown) {
					if (!error || typeof error !== 'object') {
						return;
					}

					if (instanceOfRegistrationFlow(error)) {
						data.flow = error;
						const nodes = error.ui.nodes ?? [];
						const fieldErrors: ValidationErrors<RegistrationFormSchema> = {};
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
							handleFlowErrAction(config.routes.registration.path, error.message);
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
			<Card.Title class="text-center text-xl">Register</Card.Title>
			<Card.Description>Create a new account to play some chess</Card.Description>
		</Card.Header>

		<Card.Content>
			<div class="grid gap-4">
				<form method="POST" use:enhance class="grid gap-4">
					{#each data?.flow?.ui?.messages ?? [] as msg}
						{@const err = msg.type === 'error'}
						{@const clr = msg.type === 'error' ? 'red' : msg.type === 'success' ? 'green' : 'blue'}
						<Alert.Root class="border border-{clr}-600 bg-{clr}-50 text-{clr}-600 dark:bg-{clr}-950">
							<Alert.Title>{err ? 'Unable to sign up' : ''}</Alert.Title>
							<Alert.Description>{msg.text}</Alert.Description>
						</Alert.Root>
					{/each}

					<div class="grid grid-cols-2 gap-4">
						<div class="grid gap-2">
							<Form.Field form={supForm} name="traits.first_name">
								<Form.Control let:attrs>
									<Form.Label>First name</Form.Label>
									<Input {...attrs} bind:value={$form.traits.first_name} />
								</Form.Control>
								<Form.FieldErrors />
							</Form.Field>
						</div>
						<div class="grid gap-2">
							<Form.Field form={supForm} name="traits.last_name">
								<Form.Control let:attrs>
									<Form.Label>Last name</Form.Label>
									<Input {...attrs} bind:value={$form.traits.last_name} />
								</Form.Control>
								<Form.FieldErrors />
							</Form.Field>
						</div>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.email">
							<Form.Control let:attrs>
								<Form.Label>E-Mail</Form.Label>
								<Input type="email" {...attrs} bind:value={$form.traits.email} />
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="password">
							<Form.Control let:attrs>
								<Form.Label>Password</Form.Label>
								<Input type="password" {...attrs} bind:value={$form.password} />
							</Form.Control>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<Form.Button class="w-full font-bold">Create an account</Form.Button>
				</form>

				<section>
					<div class="inline-flex w-full items-center justify-center">
						<hr class="my-8 h-px w-full border-0 bg-gray-200 dark:bg-gray-700" />
						<span
							class="bg-card dark:bg-card absolute left-1/2 -translate-x-1/2 px-3 font-medium text-gray-900 dark:text-white"
						>
							or signup with
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
												alt="signup with {provider.label}"
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
				Already have an account?
				<a href={config.routes.login.path} class="underline"> Sign in </a>
			</div>
		</Card.Content>
	</Card.Root>
</section>

<!-- <Section name="register">
	<Register href="/">
		<svelte:fragment slot="top">
			<img class="mr-2 h-8 w-8" src="/images/logo.svg" alt="logo" />
			Juicer
		</svelte:fragment>

		<div class="space-y-4 p-6 sm:p-8 md:space-y-6">
			<form method="POST" use:enhance class="flex flex-col space-y-6">
				{#each data?.flow?.ui?.messages ?? [] as msg}
					{@const err = msg.type === 'error'}
					<SimpleAlert kind={msg.type} title={err ? 'Unable to sign up' : ''} text={msg.text} />
				{/each}

				<h3 class="p-0 text-center text-xl font-medium text-gray-900 dark:text-white">Create new account</h3>

				<InputText form={supForm} name="traits.first_name" label="First name" />
				<InputText form={supForm} name="traits.last_name" label="Last name" />
				<InputEmail form={supForm} name="traits.email" label="Your email" />
				<InputPassword form={supForm} name="password" label="Your password" />

				<Button type="submit" class="w-full1 font-bold">Create new account</Button>

				<div class="text-sm font-medium text-gray-500 dark:text-gray-300">
					Already have an account? <a
						href={config.routes.login.path}
						class="text-primary-600 dark:text-primary-500 font-medium hover:underline">Login here</a
					>
				</div>

				<section>
					<div class="inline-flex w-full items-center justify-center">
						<hr class="my-8 h-px w-full border-0 bg-gray-200 dark:bg-gray-700" />
						<span
							class="absolute left-1/2 -translate-x-1/2 bg-white px-3 font-medium text-gray-900 dark:bg-gray-900 dark:text-white"
						>
							or sign up with
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
								<input type="hidden" name="csrf_token" bind:value={$form.csrf_token} readonly required />
								<input type="hidden" name="provider" value={provider.name} readonly required />

								<button id={provider.name} type="submit" class="h-12 w-12">
									<img
										class="h-full w-full object-cover"
										src="/images/providers/{provider.name}.svg"
										alt="sign up with {provider.label}"
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
</Section> -->
