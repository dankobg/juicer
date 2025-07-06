<script lang="ts">
	import type { PageProps } from './$types';
	import * as Table from '$lib/components/ui/table/index';
	import { statusIcons, templateTypeIcons, typeIcons } from '../messages-data-grid/data';
	import {
		CourierMessageStatus,
		CourierMessageType,
		CourierMessageTemplateType,
		MessageDispatchStatusEnum
	} from '$lib/gen/juicer_openapi';
	import IconCheck from '@lucide/svelte/icons/check';
	import IconX from '@lucide/svelte/icons/x';

	let { data }: PageProps = $props();
	const fmt = new Intl.DateTimeFormat(undefined, {
		dateStyle: 'short',
		timeStyle: 'short',
		hour12: false
	});

	let StatusIcon = $derived(data.message?.status && statusIcons.get(data.message.status));
	let statusIconClasses = $derived.by(() => {
		switch (data.message?.status as CourierMessageStatus) {
			case CourierMessageStatus.Abandoned:
				return 'text-red-400';
			case CourierMessageStatus.Processing:
				return 'text-purple-400';
			case CourierMessageStatus.Queued:
				return 'text-yellow-400';
			case CourierMessageStatus.Sent:
				return 'text-green-400';
			default:
				return '';
		}
	});

	let TypeIcon = $derived(data.message?.type && typeIcons.get(data.message.type));
	let typeIconClasses = $derived.by(() => {
		switch (data.message?.type as CourierMessageType) {
			case CourierMessageType.Email:
				return 'text-blue-400';
			case CourierMessageType.Phone:
				return 'text-purple-400';
			default:
				return '';
		}
	});

	let TemplateTypeIcon = $derived(data.message?.templateType && templateTypeIcons.get(data.message.templateType));
	let templateTypeIconClasses = $derived.by(() => {
		switch (data.message?.templateType as CourierMessageTemplateType) {
			case CourierMessageTemplateType.RecoveryValid:
			case CourierMessageTemplateType.RecoveryCodeValid:
			case CourierMessageTemplateType.VerificationValid:
			case CourierMessageTemplateType.VerificationCodeValid:
			case CourierMessageTemplateType.LoginCodeValid:
			case CourierMessageTemplateType.RegistrationCodeValid:
				return 'text-green-400';
			case CourierMessageTemplateType.VerificationCodeInvalid:
			case CourierMessageTemplateType.VerificationInvalid:
			case CourierMessageTemplateType.RecoveryCodeInvalid:
			case CourierMessageTemplateType.RecoveryInvalid:
				return 'text-red-400';
			case CourierMessageTemplateType.Stub:
				return 'text-purple-400';
			default:
				return '';
		}
	});
</script>

{#if data.message}
	<h1 class="mb-6 text-2xl font-bold">Courier Message</h1>
	<div class="grid grid-cols-1 gap-x-8 gap-y-4 text-sm sm:grid-cols-2">
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">ID</span>
			<span class="font-medium">{data.message.id}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Recipient</span>
			<span class="font-medium">{data.message.recipient}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Send count</span>
			<span class="font-medium">{data.message.sendCount}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Status</span>
			<span class="flex gap-2 font-medium">{data.message.status} <StatusIcon class={statusIconClasses} /></span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Channel</span>
			<span class="font-medium">{data.message.channel}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Type</span>
			<span class="flex gap-2 font-medium">{data.message.type} <TypeIcon class={typeIconClasses} /></span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Template type</span>
			<span class="flex gap-2 font-medium">
				{data.message.templateType}
				<TemplateTypeIcon class={templateTypeIconClasses} />
			</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Subject</span>
			<span class="font-medium">{data.message.subject}</span>
		</div>
		<div class="col-span-1 flex flex-col sm:col-span-2">
			<span class="text-muted-foreground">Body</span>
			<span class="font-medium">{data.message.body}</span>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Created time</span>
			<time class="font-medium">{fmt.format(data.message.createdAt)}</time>
		</div>
		<div class="flex flex-col justify-center">
			<span class="text-muted-foreground">Updated time</span>
			<time class="font-medium">{fmt.format(data.message.updatedAt)}</time>
		</div>
	</div>

	{#if data.message.dispatches && data.message.dispatches.length > 0}
		<p class="mt-8 text-lg">Message dispatches</p>
		<Table.Root>
			<Table.Caption>A list of message dispatches</Table.Caption>
			<Table.Header>
				<Table.Row>
					<Table.Head>ID</Table.Head>
					<Table.Head>Message ID</Table.Head>
					<Table.Head>Status</Table.Head>
					<Table.Head>Created time</Table.Head>
					<Table.Head>Update time</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each data.message.dispatches as dispatch (dispatch)}
					<Table.Row>
						<Table.Cell class="font-medium">{dispatch.id}</Table.Cell>
						<Table.Cell>{dispatch.messageId}</Table.Cell>
						<Table.Cell>
							<div class="flex gap-2">
								{dispatch.status}
								{#if dispatch.status === MessageDispatchStatusEnum.Success}
									<IconCheck class="text-green-400" />
								{/if}
								{#if dispatch.status === MessageDispatchStatusEnum.Failed}
									<IconX class="text-red-400" />
								{/if}
							</div>
						</Table.Cell>
						<Table.Cell>{fmt.format(dispatch.createdAt)}</Table.Cell>
						<Table.Cell>{fmt.format(dispatch.updatedAt)}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
{/if}
