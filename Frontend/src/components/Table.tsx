import {useTable, useSortBy, Column} from 'react-table'

interface TableProps<T extends object> {
  columns: Column<T>[]
  data: T[]
}

const Table = <T extends object>({columns, data}: TableProps<T>) => {
  const {getTableProps, getTableBodyProps, headerGroups, rows, prepareRow} = useTable(
    {
      columns,
      data,
    },
    useSortBy,
  )

  return (
    <table {...getTableProps()} className="w-full">
      <thead>
        {headerGroups.map(headerGroup => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map(column => (
              <th
                {...column.getHeaderProps((column as any).getSortByToggleProps())}
                className="p-2 text-left bg-gray-200 uppercase text-sm"
              >
                {column.render('Header')}
                <span>
                  {(column as any).isSorted
                    ? (column as any).isSortedDesc
                      ? ' 🔽'
                      : ' 🔼'
                    : ''}
                </span>
              </th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map(row => {
          prepareRow(row)
          return (
            <tr
              {...row.getRowProps()}
              className="hover:bg-gray-50 transition-all duration-200"
            >
              {row.cells.map(cell => (
                <td
                  {...cell.getCellProps()}
                  className="p-4 border-b border-gray-200 text-gray-700"
                >
                  {cell.render('Cell')}
                </td>
              ))}
            </tr>
          )
        })}
      </tbody>
    </table>
  )
}

export default Table
