import DataTableRowActions from './data-table-row-actions.svelte';
import type { components } from '#lib/gen/juicer_openapi.ts';
import { createAppColumnHelper } from '#lib/components/data-grid-shared/app-table.ts';
import { renderComponent } from '@tanstack/svelte-table';
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
		columnHelper.accessor(row => getGameVariantNames(gameVariants).get(row.game_variant_id), {
			id: 'game_variant',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Variant'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor(row => getGameTimeKindNames(gameTimeKinds).get(row.game_time_kind_id), {
			id: 'game_time_kind',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Time kind'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor(row => getGameTimeCategoryNames(gameTimeCategories).get(row.game_time_category_id), {
			id: 'game_time_category',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Time category'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
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
		columnHelper.accessor(
			row => {
				const resultName = gameResultNames.get(row.game_result_id!);
				switch (resultName) {
					case 'white-won':
						if (userId === row.white_id) return 'win';
						return 'loss';
					case 'black-won':
						if (userId === row.white_id) return 'loss';
						return 'win';
					default:
						return resultName;
				}
			},
			{
				id: 'game_result',
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
		columnHelper.accessor(row => gameResultStatusNames.get(row.game_result_status_id!), {
			id: 'game_result_status',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'Result status'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
			}
		}),
		columnHelper.accessor(row => gameStateNames.get(row.game_state_id!), {
			id: 'game_state',
			header: ({ header }) => {
				return renderComponent(header.DataTableColumnHeader, {
					title: 'State'
				});
			},
			cell: ({ cell }) => {
				return renderComponent(cell.TextCell);
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
