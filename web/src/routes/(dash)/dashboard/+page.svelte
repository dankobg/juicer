<script lang="ts">
	import { Card } from 'flowbite-svelte';
</script>

<div class="grid gap-4 grid-cols-[repeat(auto-fit,minmax(20rem,1fr))]">
	{#each Array(10) as _}
		<Card
			class="max-w-none items-center justify-between p-4 bg-white border border-gray-200 rounded-lg shadow-sm sm:flex dark:border-gray-700 sm:p-6 dark:bg-gray-800"
		>
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
				Noteworthy technology acquisitions 2021
			</h5>
			<p class="font-normal text-gray-700 dark:text-gray-400 leading-tight">
				Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam maiores natus modi, dignissimos obcaecati
				voluptatum deserunt, alias corporis quasi libero quisquam qui quae dicta suscipit odio saepe accusantium
				nesciunt distinctio.
			</p>
		</Card>
	{/each}
</div>
