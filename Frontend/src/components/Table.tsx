import clsx from 'clsx'

interface Column {
  key: string
  header: string
}

interface TableProps<T extends Record<string, any>> {
  columns: Column[]
  data: T[]
  className?: string
}

const Table = <T extends Record<string, any>>({
  columns,
  data,
  className,
}: TableProps<T>) => {
  return (
    <div className={clsx('rounded-lg overflow-y-scroll ', className)}>
      <table className="min-w-full">
        <thead className="bg-gray-100 ">
          <tr>
            {columns.map(column => (
              <th
                key={column.key}
                className="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
              >
                {column.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100">
          {data.map((row, rowIndex) => (
            <tr key={rowIndex} className="hover:bg-gray-50 transition-colors">
              {columns.map(column => (
                <td key={column.key} className="px-6 py-4 text-xl text-gray-700">
                  <div className="flex flex-col gap-2">
                    {row[column.key]}
                    {column.key === 'ipAddress' && (
                      <span className="text-gray-400 text-lg">China</span>
                    )}
                  </div>
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Table
