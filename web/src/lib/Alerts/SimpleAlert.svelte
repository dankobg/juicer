<script lang="ts">
	import { Alert } from 'flowbite-svelte';
	import { InfoCircleSolid } from 'flowbite-svelte-icons';

	export let title = '';
	export let text = '';
	export let kind: 'error' | 'info' | 'warning' | 'success' = 'info';
	export let color:
		| 'gray'
		| 'red'
		| 'yellow'
		| 'green'
		| 'indigo'
		| 'purple'
		| 'pink'
		| 'blue'
		| 'light'
		| 'dark'
		| 'default'
		| 'dropdown'
		| 'navbar'
		| 'navbarUl'
		| 'form'
		| 'primary'
		| 'orange'
		| 'none'
		| undefined = undefined;

	$: if (kind && !color) {
		if (kind === 'error') {
			color = 'red';
		} else if (kind === 'info') {
			color = 'blue';
		} else if (kind === 'warning') {
			color = 'yellow';
		} else if (kind === 'success') {
			color = 'green';
		}
	}
</script>

<Alert {color} rounded={false} class="border-t-4">
	<InfoCircleSolid slot="icon" class="w-4 h-4" />
	<span class="font-medium">{title}</span>
	{text}
</Alert>
