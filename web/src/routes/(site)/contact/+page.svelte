<script>
	import { config } from '$lib/kratos/config';
	import { Button } from 'flowbite-svelte';
</script>

<h1>contact</h1>

<Button href="{config.routes.login.path}?return_to={encodeURIComponent(window.location.href)}" outline size="sm"
	>Login</Button
>
<Button href="{config.routes.registration.path}?return_to={encodeURIComponent(window.location.href)}" size="sm"
	>Sign up</Button
>
