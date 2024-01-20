<script lang="ts">
	import type { PageData } from './$types';
	import type { RegistrationFlow, UiNodeInputAttributes, UpdateRegistrationFlowWithPasswordMethod } from '@ory/client';
	import { goto } from '$app/navigation';
	import { config } from '$lib/kratos/config';
	import { kratos } from '$lib/kratos/client';
	import { Button, Tooltip } from 'flowbite-svelte';
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { superForm } from 'sveltekit-superforms/client';
	import set from 'just-safe-set';
	import { zod } from 'sveltekit-superforms/adapters';
	import { z } from 'zod';
	import { isAxiosError, providers } from '$lib/kratos/helpers';
	import InputEmail from '$lib/Inputs/InputEmail.svelte';
	import InputPassword from '$lib/Inputs/InputPassword.svelte';
	import SimpleAlert from '$lib/Alerts/SimpleAlert.svelte';
	import InputText from '$lib/Inputs/InputText.svelte';

	export let data: PageData;

	const registrationFormSchema = z.object({
		csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
		method: z.literal('password'),
		password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
		traits: z.object({
			first_name: z.string(),
			last_name: z.string(),
			email: z.string().min(1, { message: 'E-Mail is required' }).email('E-Mail must be a valid email'),
			avatar_url: z.string(),
		}),
		transient_payload: z.object({}).optional(),
	});

	type RegistrationFormSchema = z.infer<typeof registrationFormSchema>;

	const initialRegistrationForm: RegistrationFormSchema = {
		password: '',
		method: 'password',
		csrf_token: data.csrf,
		traits: {
			first_name: '',
			last_name: '',
			email: '',
			avatar_url: '',
		},
		transient_payload: {},
	};

	const supForm = superForm(initialRegistrationForm, {
		validators: zod(registrationFormSchema),
		SPA: true,
		dataType: 'json',
		errorSelector: '[data-invalid]',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		async onUpdated({ form }) {
			if (!form.valid) {
				// toast.error('Invalid form, please fix errors and try again');
				return;
			}

			const url = data.flow?.ui.action;
			const body = form.data as UpdateRegistrationFlowWithPasswordMethod & { method: 'password' };

			if (url) {
				try {
					const responseFlow = await kratos.updateRegistrationFlow({
						flow: data.flow?.id ?? '',
						updateRegistrationFlowBody: body,
					});

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
				} catch (error) {
					if (isAxiosError(error)) {
						const flowData = error?.response?.data as RegistrationFlow;
						data.flow = flowData;

						const nodes = flowData.ui.nodes ?? [];
						const fieldErrors = new Map<keyof RegistrationFormSchema, string[]>();

						for (const node of nodes) {
							const errMsgs: string[] = [];
							const attrs = node.attributes as UiNodeInputAttributes;

							for (const msg of node?.messages ?? []) {
								errMsgs.push(msg.text);
								const fieldName = attrs?.name as keyof RegistrationFormSchema;
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

<Section name="register">
	<Register href="/">
		<svelte:fragment slot="top">
			<img class="w-8 h-8 mr-2" src="/images/logo.jpeg" alt="logo" />
			Juicer
		</svelte:fragment>

		<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
			<form method="POST" use:enhance class="flex flex-col space-y-6">
				<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0 text-center">Create new account</h3>

				{#each data?.flow?.ui?.messages ?? [] as msg}
					{@const err = msg.type === 'error'}
					<SimpleAlert kind={err ? 'error' : 'info'} title={err ? 'Unable to sign up' : undefined} text={msg.text} />
				{/each}

				<InputText form={supForm} name="traits.first_name" label="First name" />
				<InputText form={supForm} name="traits.last_name" label="Last name" />
				<InputEmail form={supForm} name="traits.email" label="Your email" />
				<InputPassword form={supForm} name="password" label="Your password" />

				<Button type="submit" class="w-full1 font-bold">Create new account</Button>

				<div class="text-sm font-medium text-gray-500 dark:text-gray-300">
					Already have an account? <a
						href={config.routes.login.path}
						class="font-medium text-primary-600 hover:underline dark:text-primary-500">Login here</a
					>
				</div>

				<section>
					<div class="inline-flex items-center justify-center w-full">
						<hr class="w-full h-px my-8 bg-gray-200 border-0 dark:bg-gray-700" />
						<span
							class="absolute px-3 font-medium text-gray-900 -translate-x-1/2 bg-white left-1/2 dark:text-white dark:bg-gray-900"
						>
							or sign up with
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
</Section>
