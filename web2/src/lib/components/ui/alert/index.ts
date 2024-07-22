import { type VariantProps, tv } from 'tailwind-variants';

import Root from './alert.svelte';
import Description from './alert-description.svelte';
import Title from './alert-title.svelte';

export const alertVariants = tv({
	base: 'relative w-full rounded-lg border px-4 py-3 text-sm [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground [&>svg~*]:pl-7',
	variants: {
		variant: {
			default: 'bg-background text-foreground',
			destructive: 'border-rose-600/50 text-rose dark:border-rose [&>svg]:text-gree bg-rose-500/30',
			error: 'border-rose-600/50 text-rose dark:border-rose [&>svg]:text-gree bg-rose-500/30',
			warn: 'border-orange-600/50 text-orange dark:border-orange [&>svg]:text-orange bg-orange-500/30',
			info: 'border-sky-600/50 text-sky dark:border-sky [&>svg]:text-sky bg-sky-500/30',
			success: 'border-green-600/50 text-green dark:border-green [&>svg]:text-gree bg-green-500/30'
		}
	},
	defaultVariants: {
		variant: 'default'
	}
});

export type Variant = VariantProps<typeof alertVariants>['variant'];
export type HeadingLevel = 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';

export {
	Root,
	Description,
	Title,
	//
	Root as Alert,
	Description as AlertDescription,
	Title as AlertTitle
};
