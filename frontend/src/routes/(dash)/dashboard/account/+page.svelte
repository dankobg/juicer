<script lang="ts">
  import type { PageData } from './$types';
  import PasswordSettings from './PasswordSettings.svelte';
  import SocialProvidersLinking from './SocialProvidersLinking.svelte';

  export let data: PageData;
</script>

<PasswordSettings {data} />
<SocialProvidersLinking {data} />
