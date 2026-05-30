<script lang="ts">
	import { Button } from '$lib/components/ui/button/index';
	import IconChevronLeft from '@lucide/svelte/icons/chevron-left';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Label } from '$lib/components/ui/label/index';
	import { uiSettings } from '$lib/components/ui-settings/ui-settings-state.svelte';
	import * as Select from '$lib/components/ui/select/index';
	import * as Tooltip from '$lib/components/ui/tooltip/index';
	import { Slider } from '$lib/components/ui/slider/index';
	import IconVolume from '@lucide/svelte/icons/volume';
	import IconVolume1 from '@lucide/svelte/icons/volume-1';
	import IconVolume2 from '@lucide/svelte/icons/volume-2';
	import IconVolumeOff from '@lucide/svelte/icons/volume-off';
	import { soundManager } from '$lib/sound/sound-manager.svelte';
	import Switch from '../ui/switch/switch.svelte';
	import type { CoordsRanksPosition, CoordsFilesPosition, CoordsPlacement } from '@dankop/juicer-board';

	let nestedOpen = $state<boolean>(false);
	let editPane = $state<'board' | 'pieces' | undefined>();
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
	const boardCoordinatesPlacementOptions = [
		{ value: 'inside', label: 'Inside' },
		{ value: 'outside', label: 'Outside' }
	];
	const boardRanksPositionOptions = [
		{ value: 'left', label: 'Left' },
		{ value: 'right', label: 'Right' }
	];
	const boardFilesPositionOptions = [
		{ value: 'top', label: 'Top' },
		{ value: 'bottom', label: 'Bottom' }
	];
	let selectedBoardCoordinatesPlacement = $derived(
		boardCoordinatesPlacementOptions.find(o => uiSettings.boardCoordinates.current.placement === o.value) ??
			boardCoordinatesPlacementOptions[0]
	);
	let selectedBoardRanksPosition = $derived(
		boardRanksPositionOptions.find(o => uiSettings.boardCoordinates.current.ranksPosition === o.value) ??
			boardRanksPositionOptions[0]
	);
	let selectedBoardFilesPosition = $derived(
		boardFilesPositionOptions.find(o => uiSettings.boardCoordinates.current.filesPosition === o.value) ??
			boardFilesPositionOptions[1]
	);

	let selectedChat = $derived(chatOptions.find(o => uiSettings.chat.current === o.value) ?? chatOptions[2]);
	let selectedResizer = $derived(resizerOptions.find(o => uiSettings.resizer.current === o.value) ?? resizerOptions[1]);
	let selectedSound = $derived(soundOptions.find(o => uiSettings.sounds.current.theme === o.value) ?? soundOptions[6]);

	let sliderVolume = $derived(uiSettings.sounds.current.enabled ? uiSettings.sounds.current.volume : 0);
	let SoundIcon = $derived.by(() => {
		if (!uiSettings.sounds.current.enabled) {
			return IconVolumeOff;
		}
		const vol = uiSettings.sounds.current.volume;
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
		if (!uiSettings.dialogOpen) {
			nestedOpen = false;
		}
	});

	let previousVolume = uiSettings.sounds.current.volume;
</script>

<Dialog.Root bind:open={uiSettings.dialogOpen}>
	<Dialog.Content class="max-w-2lg z-990 max-h-[80dvh] overflow-y-auto" showCloseButton={!nestedOpen}>
		{#if nestedOpen}
			<Button
				variant="ghost"
				size="icon"
				aria-label="back"
				class="absolute top-4 right-4"
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
					UI Settings
				{/if}
			</Dialog.Title>
			{#if !nestedOpen}
				<Dialog.Description>Change your chess ui settings</Dialog.Description>
			{/if}
		</Dialog.Header>
		{#if nestedOpen}
			<div class="grid gap-4">
				{#if editPane === 'board'}
					<div class="grid grid-cols-[repeat(auto-fit,minmax(6rem,1fr))] gap-8">
						{#each uiSettings.boardThemes as boardTheme (boardTheme.name)}
							<button
								onclick={() => (uiSettings.boardActiveTheme.current = boardTheme)}
								class={[
									'grid gap-1 rounded-sm p-1',
									uiSettings.boardActiveTheme.current.name === boardTheme.name && 'border-4 border-orange-700'
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
						{#each uiSettings.pieceThemes as pieceTheme (pieceTheme)}
							<button
								onclick={() => (uiSettings.pieceActiveTheme.current = pieceTheme)}
								class={[
									'grid gap-1 rounded-sm border bg-stone-600/10 p-1 hover:bg-stone-600/20',
									uiSettings.pieceActiveTheme.current === pieceTheme ? 'border-4 border-orange-700' : 'border-secondary'
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
										src="/images/board/{uiSettings.boardActiveTheme.current.src}"
										alt="current board theme {uiSettings.boardActiveTheme.current.name}"
									/>
								</Tooltip.Trigger>
								<Tooltip.Content class="z-995">
									<p>{uiSettings.boardActiveTheme.current.name}</p>
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
											src="/images/piece/{uiSettings.pieceActiveTheme.current}/wQ.svg"
											alt="current piece theme {uiSettings.pieceActiveTheme.current}"
										/>
									</Tooltip.Trigger>
									<Tooltip.Content class="z-995">
										<p>{uiSettings.pieceActiveTheme.current}</p>
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
							onclick={() => {
								uiSettings.sounds.current.enabled = !uiSettings.sounds.current.enabled;
								if (uiSettings.sounds.current.volume > 0) {
									soundManager.play('Move');
								}
							}}
						>
							<SoundIcon />
						</Button>
						<Slider
							id="settings-sounds"
							type="single"
							value={sliderVolume}
							onValueChange={val => {
								if (val > 0) {
									uiSettings.sounds.current.enabled = true;
								}
								uiSettings.sounds.current.volume = val;
								if (previousVolume !== val) {
									soundManager.play('Move');
									previousVolume = val;
								}
							}}
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
							uiSettings.sounds.current.theme = val;
							soundManager.preloadSounds();
							soundManager.play('Capture');
						}}
					>
						<Select.Trigger class="justify-self-end" id="settings-sound-theme">{selectedSound!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each soundOptions as opt (opt.value)}
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
						onValueChange={val => (uiSettings.chat.current = val as 'disabled' | 'friends-only' | 'everyone')}
					>
						<Select.Trigger class="justify-self-end" id="settings-chat">{selectedChat!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each chatOptions as opt (opt.value)}
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
						onValueChange={val => (uiSettings.resizer.current = val as 'disabled' | 'first-move' | 'always')}
					>
						<Select.Trigger class="justify-self-end" id="settings-resizer">{selectedResizer!.label}</Select.Trigger>
						<Select.Content class="z-990">
							{#each resizerOptions as opt (opt.value)}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>

				<div class="grid grid-cols-2">
					<Label for="settings-chat" class="col-span-1">Board coordinates placement</Label>
					<Select.Root
						type="single"
						value={selectedBoardCoordinatesPlacement!.value}
						onValueChange={val => (uiSettings.boardCoordinates.current.placement = val as CoordsPlacement)}
					>
						<Select.Trigger class="justify-self-end" id="settings-chat"
							>{selectedBoardCoordinatesPlacement!.label}</Select.Trigger
						>
						<Select.Content class="z-990">
							{#each boardCoordinatesPlacementOptions as opt (opt.value)}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-chat" class="col-span-1">Board ranks position</Label>
					<Select.Root
						type="single"
						value={selectedBoardRanksPosition!.value}
						onValueChange={val => (uiSettings.boardCoordinates.current.ranksPosition = val as CoordsRanksPosition)}
					>
						<Select.Trigger class="justify-self-end" id="settings-chat">
							{selectedBoardRanksPosition!.label}
						</Select.Trigger>
						<Select.Content class="z-990">
							{#each boardRanksPositionOptions as opt (opt.value)}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-chat" class="col-span-1">Board files position</Label>
					<Select.Root
						type="single"
						value={selectedBoardFilesPosition!.value}
						onValueChange={val => (uiSettings.boardCoordinates.current.filesPosition = val as CoordsFilesPosition)}
					>
						<Select.Trigger class="justify-self-end" id="settings-chat">
							{selectedBoardFilesPosition!.label}
						</Select.Trigger>
						<Select.Content class="z-990">
							{#each boardFilesPositionOptions as opt (opt.value)}
								<Select.Item value={opt.value}>{opt.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
				<div class="grid grid-cols-2">
					<Label for="settings-chat" class="col-span-1">Show ghost piece</Label>
					<Switch
						class="justify-self-end"
						checked={uiSettings.showGhost.current}
						onCheckedChange={val => (uiSettings.showGhost.current = val)}
					/>
				</div>
			</div>
		{/if}
	</Dialog.Content>
</Dialog.Root>
