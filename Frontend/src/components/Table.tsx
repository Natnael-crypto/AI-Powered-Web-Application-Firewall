import {
  ColumnDef,
  getCoreRowModel,
  useReactTable,
  flexRender,
} from '@tanstack/react-table'
import clsx from 'clsx'

interface TableProps<T> {
  columns: ColumnDef<T>[]
  data: T[]
  className?: string
  emptyMessage?: string
  onRowClick?: (row: T) => void
}

function truncateString(value: unknown, maxLength = 60): string {
  if (typeof value === 'string' && value.length > maxLength) {
    return value.substring(0, maxLength) + '...'
  }
  return String(value)
}

function Table<T extends object>({
  columns,
  data,
  className,
  emptyMessage = 'No data available',
  onRowClick,
}: TableProps<T>) {
  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
  })

  return (
    <div className={clsx('w-full shadow-md', className)}>
      <table className="min-w-full table-auto border-collapse bg-white">
        <thead>
          {table.getHeaderGroups().map(headerGroup => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map(header => (
                <th
                  key={header.id}
                  className="px-6 py-4 text-left text-sm font-semibold text-gray-900 bg-gray-50 border-b"
                >
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody className="divide-y divide-gray-200">
          {table.getRowModel().rows.length > 0 ? (
            table.getRowModel().rows.map(row => (
              <tr
                key={row.id}
                className={clsx(
                  'hover:bg-gray-50 transition-colors',
                  onRowClick && 'cursor-pointer',
                )}
                onClick={() => onRowClick?.(row.original)}
              >
                {row.getVisibleCells().map(cell => {
                  const rendered = flexRender(
                    cell.column.columnDef.cell,
                    cell.getContext(),
                  )
                  const isActionsColumn = cell.column.id === 'actions'

                  return (
                    <td
                      key={cell.id}
                      className={clsx(
                        'px-6 py-4 text-sm text-gray-900',
                        isActionsColumn
                          ? 'whitespace-nowrap relative overflow-visible'
                          : 'whitespace-nowrap max-w-xs truncate',
                      )}
                      title={typeof rendered === 'string' ? rendered : undefined}
                    >
                      {typeof rendered === 'string' ? truncateString(rendered) : rendered}
                    </td>
                  )
                })}
              </tr>
            ))
          ) : (
            <tr>
              <td
                colSpan={columns.length}
                className="px-6 py-4 text-center text-sm text-gray-500"
              >
                {emptyMessage}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}

export default Table
