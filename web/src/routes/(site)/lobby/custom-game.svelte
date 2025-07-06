<script lang="ts">
	import { superForm } from 'sveltekit-superforms/client';
	import * as Card from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import * as Select from '$lib/components/ui/select';
	import { Slider } from '$lib/components/ui/slider';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import * as Form from '$lib/components/ui/form';
	import { valibot } from 'sveltekit-superforms/adapters';
	import * as v from 'valibot';
	import type { GameTimeKind, GameVariant } from '$lib/gen/juicer_openapi';

	let { gameVariants, gameTimeKinds }: { gameVariants: GameVariant[]; gameTimeKinds: GameTimeKind[] } = $props();

	const customGameSchema = v.object({
		variant: v.string(),
		// timeControl: v.string(),
		timeMode: v.string(),
		minutesPerSide: v.number(),
		incrementSeconds: v.number(),
		side: v.string()
	});

	type CustomGameFormSchema = v.InferInput<typeof customGameSchema>;

	let defaultVariant = {
		value: 'standard',
		label: 'Standard'
	};

	let defaultTimeMode = {
		value: 'realtime',
		label: 'Realtime'
	};

	const initialCustomGameForm: CustomGameFormSchema = {
		variant: defaultVariant.value,
		timeMode: defaultTimeMode.value,
		// timeControl: 'blitz',
		minutesPerSide: 5,
		incrementSeconds: 0,
		side: 'random'
	};

	const supForm = superForm(initialCustomGameForm, {
		id: 'custom_game',
		validators: valibot(customGameSchema),
		SPA: true,
		dataType: 'json',
		scrollToError: 'smooth',
		autoFocusOnError: 'detect',
		stickyNavbar: undefined,
		async onUpdated({ form }) {
			console.log('custom game form: ', form);
		}
	});

	const { form, enhance, errors } = supForm;

	let sliderStepMinutesPerSide = $derived($form.minutesPerSide > 30 ? 5 : 1);
	let sliderStepIncrement = $derived($form.incrementSeconds > 30 ? 5 : 1);
</script>

<Card.Root class="mx-auto max-w-sm">
	<Card.Header>
		<Card.Title>Custom game</Card.Title>
		<Card.Description>Create a custom game</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-2">
		<form method="POST" use:enhance class="grid gap-4" id="custom-game-form">
			<Form.Field form={supForm} name="variant">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>Variant</Form.Label>
						<Select.Root type="single" bind:value={$form.variant} name={props.name}>
							<Select.Trigger {...props}>
								{$form.variant ?? 'Choose a variant'}
							</Select.Trigger>
							<Select.Content>
								{#each gameVariants as variant}
									<Select.Item value={variant.name} label={variant.name} disabled={!variant.enabled} />
								{/each}
							</Select.Content>
						</Select.Root>
					{/snippet}
				</Form.Control>
				<Form.Description />
				<Form.FieldErrors />
			</Form.Field>

			<Form.Field form={supForm} name="timeMode">
				<Form.Control>
					{#snippet children({ props })}
						<Form.Label>Time mode</Form.Label>
						<Select.Root type="single" bind:value={$form.timeMode} name={props.name}>
							<Select.Trigger {...props}>
								{$form.timeMode ?? 'Choose a time mode'}
							</Select.Trigger>
							<Select.Content>
								{#each gameTimeKinds as timeKind}
									<Select.Item value={timeKind.name} label={timeKind.name} disabled={!timeKind.enabled} />
								{/each}
							</Select.Content>
						</Select.Root>
					{/snippet}
				</Form.Control>
				<Form.Description />
				<Form.FieldErrors />
			</Form.Field>

			<Label for="clock">Minutes per side</Label>
			<div class="flex gap-4">
				<Slider id="clock" bind:value={$form.minutesPerSide} max={180} step={sliderStepMinutesPerSide} type="single" />
				<div>{$form.minutesPerSide}</div>
			</div>

			<Label for="increment">Increment in seconds</Label>
			<div class="flex gap-4">
				<Slider id="increment" bind:value={$form.incrementSeconds} max={180} step={sliderStepIncrement} type="single" />
				<div>{$form.incrementSeconds}</div>
			</div>

			<Form.Field form={supForm} name="side" class="space-y-3">
				<Form.Legend>Choose a side</Form.Legend>
				<RadioGroup.Root bind:value={$form.side} class="flex gap-4">
					<div class="flex items-center space-x-2">
						<RadioGroup.Item value="random" id="random" />
						<Label for="random">Random side</Label>
					</div>
					<div class="flex items-center space-x-2">
						<RadioGroup.Item value="white" id="white" />
						<Label for="white">White</Label>
					</div>
					<div class="flex items-center space-x-2">
						<RadioGroup.Item value="black" id="black" />
						<Label for="black">Black</Label>
					</div>
				</RadioGroup.Root>
				<Form.FieldErrors />
			</Form.Field>
		</form>
	</Card.Content>
	<Card.Footer>
		<Form.Button form="custom-game-form">Create game</Form.Button>
	</Card.Footer>
</Card.Root>
