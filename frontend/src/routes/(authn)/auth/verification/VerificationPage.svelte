<script lang="ts">
  import { config } from '$lib/kratos/config';
  import type { VerificationFlow, UiNodeInputAttributes, UpdateVerificationFlowWithCodeMethod } from '@ory/client';
  import type { PageData } from './$types';
  import { superForm } from 'sveltekit-superforms/client';
  import * as z from 'zod';
  import { kratos } from '$lib/kratos/client';
  import { goto } from '$app/navigation';
  import toast from 'svelte-french-toast';
  import { FormGroup, Link, TextInput, Tile, ToastNotification } from 'carbon-components-svelte';
  import { isAxiosError } from '$lib/kratos/helpers';
  import set from 'just-safe-set';

  export let data: PageData;

  const verificationFormSchema = z.object({
    csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
    method: z.string().min(1, { message: 'method is required' }),
    code: z.string().length(6, { message: 'Code must have exactly 6 characters' }),
  });

  type VerificationFormSchema = z.infer<typeof verificationFormSchema>;

  const initialVerificationForm: VerificationFormSchema = {
    code: '',
    method: 'code',
    csrf_token: data.csrf,
  };

  const { form, enhance, errors } = superForm(initialVerificationForm, {
    validators: verificationFormSchema,
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
      const body = form.data as UpdateVerificationFlowWithCodeMethod;

      if (url) {
        try {
          await kratos.updateVerificationFlow({
            flow: data.flow?.id ?? '',
            updateVerificationFlowBody: body,
          });

          goto('/');
        } catch (error) {
          if (isAxiosError(error)) {
            const flowData = error?.response?.data as VerificationFlow;
            data.flow = flowData;

            const nodes = flowData.ui.nodes ?? [];
            const fieldErrors = new Map<keyof VerificationFormSchema, string[]>();

            for (const node of nodes) {
              const errMsgs: string[] = [];
              const attrs = node.attributes as UiNodeInputAttributes;

              for (const msg of node?.messages ?? []) {
                errMsgs.push(msg.text);
                const fieldName = attrs?.name as keyof VerificationFormSchema;
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
</script>

<form method="POST" use:enhance>
  <Tile light>
    <Link href={config.routes.home.path}>Back to home page</Link>

    <div class="form-wrapper">
      <h1 class="title">Account verification</h1>

      {#each data?.flow?.ui?.messages ?? [] as msg}
        <ToastNotification
          kind={msg.type === 'error' ? 'error' : 'info'}
          title={msg.type === 'error' ? 'Error' : 'Info'}
          subtitle={msg.type === 'error' ? 'Unable to verify account' : undefined}
          caption={msg.text}
          lowContrast
          fullWidth
        />
      {/each}

      <FormGroup noMargin style="width: 100%; margin-bottom: 0.875rem;">
        <TextInput
          bind:value={$form.code}
          name="code"
          labelText="Verification code"
          invalid={Boolean($errors.code)}
          invalidText={$errors?.code?.[0] ?? ''}
        />
      </FormGroup>

      <div class="actions">
        <button type="submit">Verify account</button>
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

  :global(.actions > .forgot) {
    margin-left: auto;
    margin-bottom: 1rem;
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
</style>
