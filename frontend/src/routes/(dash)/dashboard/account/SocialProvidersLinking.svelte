<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;

  let providers: string[] = [];

  let linkedAccounts = data?.flow?.ui?.nodes
    .filter(x => x.group === 'oidc' && (x.attributes?.name as string).endsWith('link'))
    .reduce((acc, cur) => {
      providers.push(cur.attributes?.value);

      if ((cur.attributes?.name as string) === 'link') {
        acc.set(cur.attributes?.value as string, false);
      } else if (cur.attributes.name === 'unlink') {
        acc.set(cur.attributes.value as string, true);
      }
      return acc;
    }, new Map<string, boolean>());
</script>

<div style="margin-top: 5rem;">
  {#each providers as prov}
    {#if !linkedAccounts?.get(prov)}
      <form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded">
        <input type="hidden" name="link" value={prov} readonly required />
        <input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
        <button type="submit">Link {prov} Account</button>
      </form>
    {/if}
  {/each}
</div>

<div style="margin-top: 5rem;">
  {#each providers as prov}
    {#if linkedAccounts?.get(prov)}
      <form action={data.flow?.ui.action} method="post" encType="application/x-www-form-urlencoded">
        <input type="hidden" name="unlink" value={prov} readonly required />
        <input type="hidden" name="csrf_token" bind:value={data.csrf} readonly required />
        <button type="submit">Unlink {prov} Account</button>
      </form>
    {/if}
  {/each}
</div>
