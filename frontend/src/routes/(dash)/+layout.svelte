<script lang="ts">
  import { config } from '$lib/kratos/config';
  import {
    Header,
    HeaderUtilities,
    HeaderAction,
    HeaderPanelLinks,
    HeaderPanelDivider,
    HeaderPanelLink,
    SideNav,
    SideNavItems,
    SideNavMenu,
    SideNavLink,
    SkipToContent,
    Content,
    SideNavDivider,
  } from 'carbon-components-svelte';
  import UserAvatarFilledAlt from 'carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte';
  import type { SvelteComponent } from 'svelte';
  import GroupIcon from 'carbon-icons-svelte/lib/Group.svelte';
  import DashboardIcon from 'carbon-icons-svelte/lib/Dashboard.svelte';
  import UserSettingsIcon from 'carbon-icons-svelte/lib/UserSettings.svelte';
  import LogoutIcon from 'carbon-icons-svelte/lib/Logout.svelte';
  import type { LayoutData } from './$types';

  export let data: LayoutData;

  let isSideNavOpen = false;
  let isOpen1 = false;

  let navItems: DashboardNavItem[] = [];

  $: {
    if (data.session?.identity.schema_id === 'customer') {
      navItems = customerDashboardNavItems;
    } else {
      navItems = employeeDashboardNavItems;
    }
  }

  type DashboardNavItem = {
    title: string;
    href?: string;
    dropdown?: DashboardNavItem[];
    icon?: typeof SvelteComponent;
  };

  const customerDashboardNavItems: DashboardNavItem[] = [
    { title: 'Dashboard', href: `/dashboard`, icon: DashboardIcon },
    { title: 'Account', href: `/dashboard/account`, icon: UserSettingsIcon },
    { title: 'Logout', icon: LogoutIcon, href: data.logoutUrl ?? '' },
  ];

  const employeeDashboardNavItems: DashboardNavItem[] = [
    { title: 'Dashboard', href: `/dashboard`, icon: DashboardIcon },
    { title: 'Account', href: `/dashboard/account`, icon: UserSettingsIcon },
    { title: 'Identities', href: `/dashboard/identities`, icon: GroupIcon },
    { title: 'Logout', icon: LogoutIcon, href: data.logoutUrl ?? '' },
  ];
</script>

<Header company="Juicer" platformName="Chess" bind:isSideNavOpen href={config.routes.home.path}>
  <svelte:fragment slot="skip-to-content">
    <SkipToContent />
  </svelte:fragment>
  <HeaderUtilities>
    <HeaderAction bind:isOpen={isOpen1} icon={UserAvatarFilledAlt} closeIcon={UserAvatarFilledAlt} title="Profile">
      <HeaderPanelLinks>
        <HeaderPanelDivider>{data.session?.identity.traits['email']}</HeaderPanelDivider>
        <HeaderPanelLink>Dashboard</HeaderPanelLink>
        <HeaderPanelLink>Account</HeaderPanelLink>
        <HeaderPanelLink>Logout</HeaderPanelLink>
        <HeaderPanelDivider>Other</HeaderPanelDivider>
        <HeaderPanelLink>Back to site</HeaderPanelLink>
      </HeaderPanelLinks>
    </HeaderAction>
  </HeaderUtilities>
</Header>

<SideNav bind:isOpen={isSideNavOpen} rail>
  <SideNavItems>
    {#each navItems as item}
      {#if item.dropdown}
        <SideNavMenu icon={item.icon} text={item.title}>
          {#each item.dropdown as nested}
            <SideNavLink href={nested.href} icon={nested.icon} text={nested.title} />
          {/each}
        </SideNavMenu>
      {:else}
        {#if item.title === 'Logout'}
          <SideNavDivider />
        {/if}
        <SideNavLink href={item.href} icon={item.icon} text={item.title} />
      {/if}
    {/each}
  </SideNavItems>
</SideNav>

<Content>
  <slot />
</Content>

<style>
  :global(main#main-content.bx--content) {
    background-color: #f7f7f7;
  }

  :global(main#main-content.bx--content > div.bx--grid) {
    background-color: #fff;
    padding-top: 2rem;
    padding-bottom: 2rem;
  }

  :global(nav.bx--side-nav__navigation.bx--side-nav) {
    border: 1px solid #f7f7f7;
  }
</style>
