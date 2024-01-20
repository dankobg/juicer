<script lang="ts" context="module">
	type T = Record<string, unknown>;
</script>

<script lang="ts" generics="T extends Record<string, unknown>">
	import { Helper, Input, Label } from 'flowbite-svelte';
	import { formFieldProxy, type SuperForm, type FormPathLeaves } from 'sveltekit-superforms';

	export let form: SuperForm<T, unknown>;
	export let name: FormPathLeaves<T>;
	export let label = '';
	export let labelClasses = 'space-y-2';

	const { value, errors, constraints } = formFieldProxy(form, name);
</script>

<Label class={labelClasses}>
	<span>{label}</span>

	<Input
		type="text"
		{name}
		bind:value={$value}
		{...$constraints}
		{...$$restProps}
		aria-invalid={$errors ? 'true' : undefined}
	/>
	{#if $errors}
		<Helper class="mt-2" color="red">
			<span class="font-medium">{$errors}</span>
		</Helper>
	{/if}
</Label>
