<script lang="ts">
	import { CourierMessageTemplateType } from '$lib/gen/juicer_openapi';
	import { Badge } from '$lib/components/ui/badge/index';
	import { templateTypeIcons } from './data';

	let { value }: { value?: string } = $props();

	let Icon = templateTypeIcons.get(value as CourierMessageTemplateType);
	let color = $derived.by(() => {
		switch (value as CourierMessageTemplateType) {
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

<Badge variant="outline" class="flex gap-2 border {color} w-fit">
	<Icon class="" />
	<span>{value ?? ''}</span>
</Badge>
