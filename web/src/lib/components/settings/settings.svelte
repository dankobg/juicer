<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import IconChevronLeft from '@lucide/svelte/icons/chevron-left';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Label } from '$lib/components/ui/label/index';
	import { settings } from './settings-state.svelte';
	import * as Select from '$lib/components/ui/select/index';
	import * as Tooltip from '$lib/components/ui/tooltip/index';
	import { Slider } from '$lib/components/ui/slider/index';
	import IconVolume from '@lucide/svelte/icons/volume';
	import IconVolume1 from '@lucide/svelte/icons/volume-1';
	import IconVolume2 from '@lucide/svelte/icons/volume-2';
	import IconVolumeOff from '@lucide/svelte/icons/volume-off';
	import { soundManager } from '$lib/state/sound-manager.svelte';

	let nestedOpen: boolean = $state(false);
	let editPane: 'board' | 'pieces' | undefined = $state();
	const chatOptions = [
		{ value: 'disabled', label: 'Disabled' },
		{ value: 'friends-only', label: 'Friends only' },
		{ value: 'everyone', label: 'Everyone' }
	];
	const resizerOptions = [
		{ value: 'disabled', label: 'Disabled' },
		{ value: 'first-move', label: 'Show before first move' },
		{ value: 'always', label: 'Show always' }
	];
	const soundOptions = [
		{ value: 'futuristic', label: 'Futuristic' },
		{ value: 'lisp', label: 'Lisp' },
		{ value: 'nes', label: 'Nes' },
		{ value: 'piano', label: 'Piano' },
		{ value: 'robot', label: 'Robot' },
		{ value: 'sfx', label: 'Sfx' },
		{ value: 'standard', label: 'Standard' },
		{ value: 'woodland', label: 'Woodland' }
	];
	let selectedChat = $derived(chatOptions.find(o => settings.chat.current === o.value) ?? chatOptions[2]);
	let selectedResizer = $derived(resizerOptions.find(o => settings.resizer.current === o.value) ?? resizerOptions[1]);
	let selectedSound = $derived(soundOptions.find(o => settings.sounds.current.theme === o.value) ?? soundOptions[6]);

	let sliderVolume = $derived(settings.sounds.current.enabled ? settings.sounds.current.volume : 0);
	let SoundIcon = $derived.by(() => {
		const vol = settings.sounds.current.volume;
		if (vol === 0) {
			return IconVolumeOff;
		} else if (vol <= 0.33) {
			return IconVolume;
		} else if (vol <= 0.66) {
			return IconVolume1;
		} else {
			return IconVolume2;
		}
	});

	$effect(() => {
		if (!settings.dialogOpen) {
			nestedOpen = false;
		}
	});
</script>

<Dialog.Root bind:open={settings.dialogOpen}>
	<Dialog.Content class="max-w-2lg z-990 max-h-[80dvh] overflow-y-auto" showCloseButton={!nestedOpen}>
		{#if nestedOpen}
			<Button
				variant="ghost"
				size="icon"
				aria-label="back"
				class="absolute right-4 top-4"
				onclick={() => (nestedOpen = false)}
			>
				<IconChevronLeft />
			</Button>
		{/if}
		<Dialog.Header class="mb-8">
			<Dialog.Title>
				{#if nestedOpen}
					{#if editPane === 'board'}
						Choose board theme
					{:else}
						Choose piece set
					{/if}
				{:else}
					Settings
				{/if}
			</Dialog.Title>
			{#if !nestedOpen}
				<Dialog.Description>Change some common game settings here</Dialog.Description>
			{/if}
		</Dialog.Header>
		{#if nestedOpen}
			<div class="grid gap-4">
				{#if editPane === 'board'}
					<div class="grid grid-cols-[repeat(auto-fit,minmax(6rem,1fr))] gap-8">
						{#each settings.boardThemes as boardTheme}
							<button
								onclick={() => (settings.boardActiveTheme.current = boardTheme)}
								class={[
									'grid gap-1 rounded-sm p-1',
									settings.boardActiveTheme.current.name === boardTheme.name && 'border border-4 border-orange-700'
								]}
							>
								<p class="text-center">{boardTheme.name}</p>
								<img
									src="/images/board/{boardTheme.src}"
									alt="board theme {boardTheme.name}"
									class="max-w-full object-contain"
								/>
							</button>
						{/each}
					</div>
				{:else}
					<div class="grid grid-cols-[repeat(auto-fit,minmax(6rem,1fr))] gap-8">
						{#each settings.pieceThemes as pieceTheme}
							<button
								onclick={() => (settings.pieceActiveTheme.current = pieceTheme)}
								class={[
									'grid gap-1 rounded-sm border bg-stone-600/10 p-1 hover:bg-stone-600/20',
									settings.pieceActiveTheme.current === pieceTheme ? 'border-4 border-orange-700' : 'border-secondary'
								]}
							>
								<p class="text-center">{pieceTheme}</p>
								<img
									src="/images/piece/{pieceTheme}/wQ.svg"
									alt="piece theme {pieceTheme}"
									class="w-full object-cover"
								/>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			<div class="grid gap-4">
				<div class="grid grid-cols-2">
					<div class="col-span-1 flex flex-wrap items-center gap-2">
						<p>Board theme</p>
						<Tooltip.Provider delayDuration={150}>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<img
										class="h-12 w-12 max-w-full object-contain"
										src="/images/board/{settings.boardActiveTheme.current.src}"
										alt="current board theme {settings.boardActiveTheme.current.name}"
									/>
								</Tooltip.Trigger>
								<Tooltip.Content class="z-995">
									<p>{settings.boardActiveTheme.current.name}</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</div>
					<div class="col-span-1 justify-self-end">
						<Button
							onclick={() => {
								nestedOpen = true;
								editPane = 'board';
							}}
						>
							Edit
						</Button>
					</div>
				</div>
				<div class="grid grid-cols-2">
					<div class="col-span-1">
						<div class="col-span-1 flex flex-wrap items-center gap-2">
							<p>Piece theme</p>
							<Tooltip.Provider delayDuration={150}>
								<Tooltip.Root>
									<Tooltip.Trigger>
										<img
											class="h-12 w-12 max-w-full object-contain"
											src="/images/piece/{settings.pieceActiveTheme.current}/wQ.svg"
											alt="current piece theme {settings.pieceActiveTheme.current}"
										/>
									</Tooltip.Trigger>
									<Tooltip.Content class="z-995">
										<p>{settings.pieceActiveTheme.current}</p>
									</Tooltip.Content>
								</Tooltip.Root>
							</Tooltip.Provider>
						</div>
					</div>

					<div class="col-span-1 justify-self-end">
						<Button
							onclick={() => {
								nestedOpen = true;
								editPane = 'pieces';
							}}
						>
							Edit
						</Button>
					</div>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-sounds" class="col-span-1">Sounds</Label>
					<div class="flex gap-2">
						<Button
							variant="secondary"
							size="icon"
							class="size-8"
							onclick={() => (settings.sounds.current.enabled = !settings.sounds.current.enabled)}
						>
							<SoundIcon />
						</Button>
						<Slider
							id="settings-sounds"
							type="single"
							value={sliderVolume}
							onValueChange={val => (settings.sounds.current.volume = val)}
							max={1}
							step={0.1}
						/>
					</div>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-sound-theme" class="col-span-1">Sound theme</Label>
					<Select.Root
						type="single"
						value={selectedSound!.value}
						onValueChange={val => {
							settings.sounds.current.theme = val as any;
							soundManager.preloadSounds();
						}}
					>
						<Select.Trigger class="justify-self-end" id="settings-sound-theme">{selectedSound!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each soundOptions as opt}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-chat" class="col-span-1">Chat</Label>
					<Select.Root
						type="single"
						value={selectedChat!.value}
						onValueChange={val => (settings.chat.current = val as any)}
					>
						<Select.Trigger class="justify-self-end" id="settings-chat">{selectedChat!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each chatOptions as opt}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-resizer" class="col-span-1">Resizer</Label>
					<Select.Root
						type="single"
						value={selectedResizer!.value}
						onValueChange={val => (settings.resizer.current = val as any)}
					>
						<Select.Trigger class="justify-self-end" id="settings-resizer">{selectedResizer!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each resizerOptions as opt}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
