<script lang="ts">
	import type { PageProps } from './$types';
	import { superForm } from 'sveltekit-superforms/client';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import * as Card from '$lib/components/ui/card';
	import * as Form from '$lib/components/ui/form';
	import * as Select from '$lib/components/ui/select/index.js';
	import { juicer } from '$lib/juicer/client';
	import { IdentityStateEnum, type CreateIdentityBody, type Identity } from '$lib/gen/juicer_openapi';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import PasswordInput from '$lib/components/password-input/password-input.svelte';

	let { data }: PageProps = $props();

	let createdIdentity: Identity | null = $state(null);

	const createIdentitySchema = v.object({
		schemaId: v.pipe(v.string(), v.minLength(1, 'Schema id is required')),
		state: v.union([v.literal('active'), v.literal('inactive')]),
		verified: v.boolean(),
		traits: v.object({
			first_name: v.string(),
			last_name: v.string(),
			email: v.pipe(v.string(), v.minLength(1, 'E-Mail is required'), v.email('E-Mail must be a valid email')),
			username: v.pipe(v.string(), v.minLength(1, 'Username is required')),
			avatar_url: v.string()
		}),
		password: v.union([v.literal(''), v.pipe(v.string(), v.minLength(8, 'Password must have min. 8 characters'))]),
		hashedPassword: v.string()
	});

	type CreateIdentitySchema = v.InferInput<typeof createIdentitySchema>;

	const initialIdentitySchema: CreateIdentitySchema = {
		schemaId: 'customer',
		state: 'active',
		verified: false,
		traits: {
			first_name: '',
			last_name: '',
			email: '',
			username: '',
			avatar_url: ''
		},
		password: '',
		hashedPassword: ''
	};

	const supForm = superForm(initialIdentitySchema, {
		id: 'create_identity_form',
		validators: valibot(createIdentitySchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		resetForm: false,
		async onUpdated({ form }) {
			if (!form.valid) {
				toast.error('Invalid form, please fix errors and try again');
				return;
			}
			try {
				const body: CreateIdentityBody = {
					schemaId: form.data.schemaId,
					state: form.data.state,
					traits: form.data.traits
					// metadataAdmin?: any | null;
					// metadataPublic?: any | null;
					// recoveryAddresses?: Array<RecoveryIdentityAddress>;
				};
				if (form.data.verified) {
					body.verifiableAddresses = [
						{ verified: true, via: 'email', value: form.data.traits.email, status: 'completed' }
					];
				}
				if (form.data.password || form.data.hashedPassword) {
					body.credentials ??= {};
					body.credentials.password = {
						config: {
							...(form.data.password && { password: form.data.password }),
							...(form.data.hashedPassword && { hashedPassword: form.data.hashedPassword })
						}
					};
				}
				const identity = await juicer.createIdentity({ createIdentityBody: body });
				createdIdentity = identity;
				toast.success('identity created successfuly');
			} catch (error) {
				console.log('err', error);
			}
		}
	});

	const { form, enhance, errors } = supForm;
</script>

<Card.Root class="max-w-6xl">
	<Card.Header>
		<Card.Title>Create identity</Card.Title>
		<Card.Description>Create new identity</Card.Description>
	</Card.Header>

	<Card.Content>
		<div class="grid gap-4">
			<form method="POST" use:enhance>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div class="grid gap-2">
						<Form.Field form={supForm} name="schemaId">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Schema ID</Form.Label>
									<Select.Root type="single" bind:value={$form.schemaId} name="schemaId">
										<Select.Trigger {...props}>
											{$form.schemaId ? $form.schemaId : 'Select a schema id'}
										</Select.Trigger>
										<Select.Content>
											{#each data.schemas ?? [] as schema}
												<Select.Item value={schema?.id ?? ''} label={schema?.id ?? ''} />
											{/each}
										</Select.Content>
									</Select.Root>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="state">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>State</Form.Label>
									<Select.Root type="single" bind:value={$form.state} name="state">
										<Select.Trigger {...props}>
											{$form.state ? $form.state : 'Select a state'}
										</Select.Trigger>
										<Select.Content>
											{#each Object.values(IdentityStateEnum) as x}
												<Select.Item value={x} label={x} />
											{/each}
										</Select.Content>
									</Select.Root>
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.first_name">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>First name</Form.Label>
									<Input {...props} bind:value={$form.traits.first_name} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.last_name">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Last name</Form.Label>
									<Input {...props} bind:value={$form.traits.last_name} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.username">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Username</Form.Label>
									<Input {...props} bind:value={$form.traits.username} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.email">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>E-Mail</Form.Label>
									<Input type="email" {...props} bind:value={$form.traits.email} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="traits.avatar_url">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Avatar URL</Form.Label>
									<Input {...props} bind:value={$form.traits.avatar_url} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid items-center gap-2">
						<Form.Field form={supForm} name="verified">
							<Form.Control>
								{#snippet children({ props })}
									<div class="flex items-center gap-2">
										<Checkbox {...props} bind:checked={$form.verified} />
										<div class="space-y-1 leading-none">
											<Form.Label>Verified email</Form.Label>
											<Form.Description />
										</div>
									</div>
								{/snippet}
							</Form.Control>
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="password">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Password</Form.Label>
									<PasswordInput {...props} bind:value={$form.password} />
								{/snippet}
							</Form.Control>
							<Form.Description>
								Set this only if you want to create user with the password. It's advised that you create the user and
								and send them recovery link or they can start recovery flow themselves.
							</Form.Description>
							<Form.FieldErrors />
						</Form.Field>
					</div>
					<div class="grid gap-2">
						<Form.Field form={supForm} name="hashedPassword">
							<Form.Control>
								{#snippet children({ props })}
									<Form.Label>Hashed password</Form.Label>
									<PasswordInput {...props} bind:value={$form.hashedPassword} />
								{/snippet}
							</Form.Control>
							<Form.Description />
							<Form.FieldErrors />
						</Form.Field>
					</div>
				</div>
				<Form.Button class="mt-8 font-bold">Create identity</Form.Button>
			</form>
		</div>
	</Card.Content>
</Card.Root>

{#if createdIdentity}
	<Card.Root class="max-w-6xl">
		<Card.Header>
			<Card.Title>View created identity</Card.Title>
		</Card.Header>

		<Card.Content>
			<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">ID</span>
					<a class="font-medium text-sky-500 underline" href="/dashboard/identities/{createdIdentity.id}">
						{createdIdentity.id}
					</a>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">E-Mail</span>
					<span class="font-medium">{createdIdentity.traits['email'] ?? ''}</span>
				</div>
				{#if createdIdentity.traits['first_name']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">First name</span>
						<span class="font-medium">{createdIdentity.traits['first_name'] ?? ''}</span>
					</div>
				{/if}
				{#if createdIdentity.traits['last_name']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Last name</span>
						<span class="font-medium">{createdIdentity.traits['last_name'] ?? ''}</span>
					</div>
				{/if}
				{#if createdIdentity.traits['avatar_url']}
					<div class="flex flex-col justify-center">
						<span class="text-muted-foreground">Avatar URL</span>
						<span class="font-medium">{createdIdentity.traits['avatar_url'] ?? ''}</span>
					</div>
				{/if}
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">Schema ID</span>
					<span class="font-medium">{createdIdentity.schemaId}</span>
				</div>
				<div class="flex flex-col justify-center">
					<span class="text-muted-foreground">State</span>
					<span class="font-medium">{createdIdentity.state}</span>
				</div>
			</div>
		</Card.Content>
	</Card.Root>
{/if}
