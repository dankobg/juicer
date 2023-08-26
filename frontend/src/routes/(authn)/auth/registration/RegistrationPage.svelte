<script lang="ts">
  import { config } from '$lib/kratos/config';
  import type {
    ContinueWith,
    RegistrationFlow,
    UiNodeInputAttributes,
    UpdateRegistrationFlowWithPasswordMethod,
  } from '@ory/client';
  import type { PageData } from './$types';
  import { superForm } from 'sveltekit-superforms/client';
  import * as z from 'zod';
  import { kratos } from '$lib/kratos/client';
  import { goto } from '$app/navigation';
  import toast from 'svelte-french-toast';
  import { FormGroup, Link, PasswordInput, TextInput, Tile, ToastNotification } from 'carbon-components-svelte';
  import { isAxiosError } from '$lib/kratos/helpers';
  import set from 'just-safe-set';

  export let data: PageData;

  const registrationFormSchema = z.object({
    csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
    method: z.string().min(1, { message: 'method is required' }),
    password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
    traits: z.object({
      first_name: z.string(),
      last_name: z.string(),
      email: z.string().min(1, { message: 'E-Mail is required' }).email('E-Mail must be a valid email'),
      avatar_url: z.string(),
    }),
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
  };

  const { form, enhance, errors } = superForm(initialRegistrationForm, {
    validators: registrationFormSchema,
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
      const body = form.data as UpdateRegistrationFlowWithPasswordMethod;

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

  type Provider = {
    name: string;
    label: string;
  };
  let providers: Provider[] = [
    { name: 'google', label: 'Google' },
    { name: 'github', label: 'GitHub' },
    { name: 'facebook', label: 'Facebook' },
    { name: 'discord', label: 'Discord' },
    { name: 'twitch', label: 'Twitch' },
    { name: 'slack', label: 'Slack' },
    { name: 'spotify', label: 'Spotify' },
  ];
</script>

<form method="POST" use:enhance>
  <Tile light>
    <Link href={config.routes.home.path}>Back to home page</Link>

    <div class="form-wrapper">
      <h1 class="title">Sign up</h1>

      {#each data?.flow?.ui?.messages ?? [] as msg}
        <ToastNotification
          kind={msg.type === 'error' ? 'error' : 'info'}
          title={msg.type === 'error' ? 'Error' : 'Info'}
          subtitle={msg.type === 'error' ? 'Unable to sign up' : undefined}
          caption={msg.text}
          lowContrast
          fullWidth
        />
      {/each}

      <FormGroup noMargin style="width: 100%; margin-bottom: 0.875rem;">
        <TextInput
          bind:value={$form.traits.email}
          name="traits.email"
          labelText="E-Mail"
          invalid={Boolean($errors?.traits?.email)}
          invalidText={$errors?.traits?.email?.[0] ?? ''}
        />
      </FormGroup>

      <FormGroup noMargin style="width: 100%; margin-bottom: 0.875rem;">
        <PasswordInput
          bind:value={$form.password}
          name="password"
          labelText="Password"
          invalid={Boolean($errors.password)}
          invalidText={$errors?.password?.[0] ?? ''}
        />
      </FormGroup>

      <div class="actions">
        <button type="submit">Sign up</button>
      </div>

      <div class="separator">OR</div>

      <div class="socials">
        {#each providers as provider}
          <form
            action={data.flow?.ui.action}
            method="post"
            encType="application/x-www-form-urlencoded"
            class="provider {provider.name}"
          >
            <input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
            <input type="hidden" name="provider" value={provider.name} readonly required />
            <button type="submit">
              <img src="/images/auth/provider-{provider.name}.svg" alt="Continue with {provider.label}" />
              <span>Continue with {provider.label}</span>
            </button>
          </form>
        {/each}
      </div>

      <div class="footer">
        <span class="register">
          Already have an account? <Link href={config.routes.login.path}>Log in here!</Link>
        </span>
      </div>
    </div>
  </Tile>
</form>

<style>
  .form-wrapper {
    max-width: 26rem;
  }

  .title {
    font-size: 2.6rem;
    font-weight: 400;
    text-align: center;
  }

  .actions {
    display: flex;
    flex-direction: column;
    margin-bottom: 1rem;
  }

  .footer {
    margin-top: 1rem;
    display: flex;
    justify-content: center;
  }

  .actions > button[type='submit'] {
    width: 100%;
    padding: 0.625rem 1rem;
    margin-top: 1rem;
    background-color: darkorange;
    color: #fff;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    font-size: 18px;
    line-height: 26px;
    font-weight: bold;
  }

  .socials {
    margin: 1rem 0;
    display: grid;
    gap: 0.4rem;
  }

  .provider > button[type='submit'] {
    width: 100%;
    padding: 0.475rem 2rem;
    color: #fff;
    border-radius: 4px;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    line-height: 26px;
    font-weight: bold;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .provider > button > img {
    margin-right: auto;
    object-fit: cover;
    width: 24px;
    height: 24px;
  }

  .provider > button > span {
    margin-right: auto;
  }

  .provider.google > button[type='submit'] {
    background-color: #fff;
    color: #333;
  }
  .provider.github > button[type='submit'] {
    background-color: #211f1f;
  }
  .provider.facebook > button[type='submit'] {
    background-color: #0165e1;
  }
  .provider.discord > button[type='submit'] {
    background-color: #5865f2;
  }
  .provider.slack > button[type='submit'] {
    background-color: #4a154b;
  }
  .provider.twitch > button[type='submit'] {
    background-color: #9146ff;
  }
  .provider.spotify > button[type='submit'] {
    background-color: #1ed760;
  }

  .separator {
    display: flex;
    align-items: center;
    text-align: center;
  }
  .separator::before,
  .separator::after {
    content: '';
    flex: 1;
    border-bottom: 1px solid #d6d6d6;
    font-size: 1rem;
  }
  .separator:not(:empty)::before {
    margin-right: 1rem;
  }
  .separator:not(:empty)::after {
    margin-left: 1rem;
  }
</style>
