<script lang="ts">
  import { config } from '$lib/kratos/config';
  import { onMount } from 'svelte';
  import type { PageData } from './$types';

  export let data: PageData;

  let ws: WebSocket | null = null;
  const wsEndpoint = 'wss://juicer-dev.xyz/ws';

  const onMessage = (event: MessageEvent) => {
    const { data } = event;
    console.log('received', data);
  };

  onMount(() => {
    ws = new WebSocket(wsEndpoint);

    ws.addEventListener('open', event => {
      console.log('open', event);
    });

    ws.addEventListener('error', event => {
      console.log('error wtf', event);
    });

    ws.addEventListener('message', onMessage);

    return () => {
      ws = null;
    };
  });

  function onSeekGame() {
    if (ws) {
      const payload = { type: 'loby:seek_game', data: { id: data.session?.identity.id, mode: 'blitz', format: '5m' } };
      ws.send(JSON.stringify(payload));
    }
  }

  function onCancelSeekGame() {
    if (ws) {
      const data = { type: 'loby:cancel_seek_game', data: {} };
      ws.send(JSON.stringify(data));
    }
  }
</script>

<svelte:head>
  <title>Juicer</title>
  <meta name="description" content="Juicer lets you play chess on the web" />
</svelte:head>

<p>Home</p>

{#if data.session?.active}
  <h2>LOGGED IN</h2>

  <button on:click={onSeekGame}>PLAY</button>

  <button on:click={onCancelSeekGame}>CANCEL SEARCH</button>

  <pre>{JSON.stringify(
      {
        id: data.session.identity.id,
        email: data.session.identity.traits['email'],
        active: data.session.active,
        recovery: data.session.identity.recovery_addresses?.[0]?.value,
        verif: data.session.identity.verifiable_addresses?.[0]?.value,
      },
      null,
      2
    )}</pre>

  <a href={config.routes.dashboard.path}>DASHBOARD</a>

  <a href={data.logoutUrl}>LOGOUT</a>
{:else}
  <h2>NO SESSION</h2>

  <a href={config.routes.login.path}>LOGIN</a>
  <a href={config.routes.registration.path}>REGISTER</a>
{/if}

<style>
</style>
