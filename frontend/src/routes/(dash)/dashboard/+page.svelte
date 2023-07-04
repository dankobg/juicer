<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

DASHBOARD

{#if data.session?.active}
  <pre>{JSON.stringify(
      {
        email: data.session.identity.traits['email'],
        active: data.session.active,
        recovery: data.session.identity.recovery_addresses?.[0]?.value,
        verif: data.session.identity.verifiable_addresses?.[0]?.value,
      },
      null,
      2
    )}</pre>
{:else}{/if}
