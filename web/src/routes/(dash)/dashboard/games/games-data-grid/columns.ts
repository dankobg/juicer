import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/juicer_openapi.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { renderComponent, renderSnippet } from '@tanstack/svelte-table';
import {
	gameResultNames,
	gameResultStatusNames,
	gameStateNames,
	getGameTimeCategoryNames,
	getGameTimeKindNames,
	getGameVariantNames
} from './data';
import DataTableCellResult from './data-table-cell-result.svelte';
import { Color } from '#lib/gen/juicer_pb.ts';
import DataTableCellColor from './data-table-cell-color.svelte';
import { createRawSnippet } from 'svelte';

const columnHelper = createAppColumnHelper<components['schemas']['Game']>();

export function makeColumns(
	userId: string | undefined,
	gameVariants: components['schemas']['GameVariant'][],
	gameTimeKinds: components['schemas']['GameTimeKind'][],
	gameTimeCategories: components['schemas']['GameTimeCategory'][]
) {
	return columnHelper.columns([
		columnHelper.display({
			id: 'select',
			header: ({ header }) => renderComponent(header.SelectHeader),
			cell: ({ cell }) => renderComponent(cell.SelectCell),
			enableSorting: false,
			enableHiding: false,
			enableResizing: false,
			maxSize: 40
		}),
		columnHelper.accessor('id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'ID'
				});
			},
			cell: ({ cell, row }) => {
				return renderComponent(cell.DataTableCellId, {
					value: row.original.id,
					href: `/dashboard/games/${row.original.id}`
				});
			}
		}),
		columnHelper.accessor(row => (userId === row.white_id ? Color.WHITE : Color.BLACK), {
			id: 'color',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Color'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(DataTableCellColor, {
					color: cell.getValue()
				});
			}
		}),
		columnHelper.accessor(row => (userId === row.white_id ? row.black_id : row.white_id), {
			id: 'opponent',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Opponent'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor('game_variant_id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Variant'
				});
			},
			cell: ({ row }) => {
				const gameVariantSnippet = createRawSnippet<[{ value: string | undefined }]>(getValue => {
					const { value } = getValue();
					return {
						render: () => `<span>${value}</span>`
					};
				});

				return renderSnippet(gameVariantSnippet, {
					value: getGameVariantNames(gameVariants).get(row.original.game_variant_id)
				});
			}
		}),
		columnHelper.accessor('game_time_kind_id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Time kind'
				});
			},
			cell: ({ row }) => {
				const gameTimeKindSnippet = createRawSnippet<[{ value: string | undefined }]>(getValue => {
					const { value } = getValue();
					return {
						render: () => `<span>${value}</span>`
					};
				});

				return renderSnippet(gameTimeKindSnippet, {
					value: getGameTimeKindNames(gameTimeKinds).get(row.original.game_time_kind_id)
				});
			}
		}),
		columnHelper.accessor('game_time_category_id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Time category'
				});
			},
			cell: ({ row }) => {
				const gameTimeCategorySnippet = createRawSnippet<[{ value: string | undefined }]>(getValue => {
					const { value } = getValue();
					return {
						render: () => `<span>${value}</span>`
					};
				});

				return renderSnippet(gameTimeCategorySnippet, {
					value: getGameTimeCategoryNames(gameTimeCategories).get(row.original.game_time_category_id)
				});
			}
		}),
		columnHelper.accessor(
			row => {
				return `${row.time_control_clock_ms / 1000 / 60}+${row.time_control_increment_ms / 1000}`;
			},
			{
				id: 'time-control',
				header: ({ header }) => {
					return renderComponent(header.DataTableColumnHeader, {
						title: 'Time control'
					});
				},
				cell: ({ cell }) => {
					return renderComponent(cell.TextCell);
				}
			}
		),
		// columnHelper.accessor('game_result_id', {
		// 	header: ({ header }) => {
		// 		return renderComponent(header.DataTableColumnHeader, {
		// 			title: 'Result'
		// 		});
		// 	},
		// 	cell: ({ row }) => {
		// 		const getRes = (): 'win' | 'loss' | 'draw' | 'interrupted' => {
		// 			const resultName = gameResultNames.get(row.original.game_result_id!);
		// 			switch (resultName) {
		// 				case 'white-won':
		// 					if (userId === row.original.white_id) return 'win';
		// 					return 'loss';
		// 				case 'black-won':
		// 					if (userId === row.original.white_id) return 'loss';
		// 					return 'win';
		// 				default:
		// 					return resultName as 'win' | 'loss' | 'draw' | 'interrupted';
		// 			}
		// 		};

		// 		return renderComponent(DataTableCellResult, {
		// 			value: getRes()
		// 		});
		// 	}
		// }),
		columnHelper.accessor(
			row => {
				const getRes = (): 'win' | 'loss' | 'draw' | 'interrupted' => {
					const resultName = gameResultNames.get(row.game_result_id!);
					switch (resultName) {
						case 'white-won':
							if (userId === row.white_id) return 'win';
							return 'loss';
						case 'black-won':
							if (userId === row.white_id) return 'loss';
							return 'win';
						default:
							return resultName as 'win' | 'loss' | 'draw' | 'interrupted';
					}
				};
				return getRes();
			},
			{
				id: 'user_game_result',
				header: ({ header }) => {
					return renderComponent(header.DataTableColumnHeader, {
						title: 'Result'
					});
				},
				cell: ({ cell }) => {
					return renderComponent(DataTableCellResult, {
						value: cell.getValue()
					});
				}
			}
		),
		columnHelper.accessor('game_result_status_id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Result status'
				});
			},
			cell: ({ cell, row }) => {
				return renderComponent(cell.TextValue, {
					value: gameResultStatusNames.get(row.original.game_result_status_id!)
				});
			}
		}),
		columnHelper.accessor('game_state_id', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'State'
				});
			},
			cell: ({ cell, row }) => {
				return renderComponent(cell.TextValue, {
					value: gameStateNames.get(row.original.game_state_id!)
				});
			}
		}),
		columnHelper.accessor('rated', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Rated'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor('fen', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'FEN'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor('pgn', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'PGN'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor('start_time', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Start time'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.DateCell);
			}
		}),
		columnHelper.accessor('end_time', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'End time'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.DateCell);
			}
		}),
		columnHelper.accessor('last_move', {
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Last move'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.DateCell);
			}
		}),
		columnHelper.display({
			id: 'actions',
			cell: () => renderComponent(DataTableRowActions),
			enableSorting: false,
			enableHiding: false,
			enableResizing: false,
			maxSize: 40
		})
	]);
}
