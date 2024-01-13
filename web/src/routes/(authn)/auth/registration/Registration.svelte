<script lang="ts">
	import { Section, Register } from 'flowbite-svelte-blocks';
	import { Button, Label, Input } from 'flowbite-svelte';
	import { config } from '$lib/kratos/config';
</script>

<Section name="register">
	<Register href="/">
		<svelte:fragment slot="top">
			<img class="w-8 h-8 mr-2" src="/images/logo.jpeg" alt="logo" />
			Juicer
		</svelte:fragment>

		<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
			<form class="flex flex-col space-y-6">
				<h3 class="text-xl font-medium text-gray-900 dark:text-white p-0 text-center">Create new account</h3>

				<Label class="space-y-2">
					<span>Your email</span>
					<Input type="email" name="email" required />
				</Label>

				<Label class="space-y-2">
					<span>Your password</span>
					<Input type="password" name="password" required />
				</Label>

				<Button type="submit" class="w-full1 font-bold">Create new account</Button>

				<div class="text-sm font-medium text-gray-500 dark:text-gray-300">
					Already have an account? <a
						href={config.routes.login.path}
						class="font-medium text-primary-600 hover:underline dark:text-primary-500">Login here</a
					>
				</div>
			</form>
		</div>
	</Register>
</Section>
