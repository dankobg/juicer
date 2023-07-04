<script lang="ts">
  import type { UpdateSettingsFlowWithPasswordMethod } from '@ory/client';
  import type { PageData } from './$types';
  import { superForm } from 'sveltekit-superforms/client';
  import * as z from 'zod';
  import { kratos } from '$lib/kratos/client';
  import { goto } from '$app/navigation';
  import toast from 'svelte-french-toast';

  export let data: PageData;

  const settingsFormSchema = z.object({
    csrf_token: z.string().min(1, { message: 'csrf_token is required' }),
    method: z.string().min(1, { message: 'method is required' }),
    password: z.string().min(8, { message: 'Password must have min. 8 characters' }),
  });

  type SettingsFormSchema = z.infer<typeof settingsFormSchema>;

  const initialSettingsForm: SettingsFormSchema = {
    password: '',
    method: 'password',
    csrf_token: data.csrf,
  };

  const { form, enhance, errors } = superForm(initialSettingsForm, {
    validators: settingsFormSchema,
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
      const body = form.data as UpdateSettingsFlowWithPasswordMethod;

      if (url) {
        try {
          const responseFlow = await kratos.updateSettingsFlow({
            flow: data.flow?.id ?? '',
            updateSettingsFlowBody: body,
          });

          if (responseFlow.data.continue_with) {
            for (const item of responseFlow.data.continue_with) {
              switch (item.action) {
                case 'show_verification_ui':
                  if (item?.flow?.id) {
                    // goto(`${config.routes.verification.path}?flow=${item?.flow?.id}`);
                    goto(item?.flow?.url as string);
                    // @TODO: maybe hardcode if doesnt exist
                  }
                  return;
              }
            }
          }

          goto('/');
        } catch (error) {
          const flowData = error?.response?.data;
          // flow = { ...flowData };

          const nodes = error?.response?.data?.ui?.nodes || [];
          const errGroup = new Map<keyof SettingsFormSchema, string[]>();

          nodes.forEach((node: any) => {
            const messages = node?.messages || [];
            const arr: Array<any> = [];
            messages.forEach((msg: any) => {
              if (msg.type === 'error') {
                arr.push(msg.text);
                errGroup.set(node.attributes.name, arr);
              }
            });
          });

          errGroup.forEach((value, key) => {
            errors.update(errs => {
              errs[key] = value.join(', ');
            });
          });
        }
      }
    },
  });
</script>

<p>Password settings</p>

<form method="POST" use:enhance>
  <div>
    <input name="password" placeholder="New Password" bind:value={$form.password} data-invalid={$errors.password} />
    {#if $errors.password}<span style="color: red;">{$errors.password}</span>{/if}
  </div>

  <button type="submit">update password</button>
</form>

<style>
  form > div {
    display: flex;
    flex-direction: column;
    width: 20rem;
  }

  input {
    margin-bottom: 1rem;
  }
</style>
