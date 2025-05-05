import { useState } from 'react'
import {
  ColumnDef,
  Row,
  getCoreRowModel,
  useReactTable,
  flexRender,
} from '@tanstack/react-table'
import clsx from 'clsx'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import CreateRuleModal from '../components/CreateRuleModal'
import {
  HiOutlineBan,
  HiOutlineCheckCircle,
  HiOutlineTrash,
  HiOutlineThumbUp,
  HiOutlineThumbDown,
  HiOutlinePencilAlt,
} from 'react-icons/hi'

export interface Rule {
  rule_id: string
  rule_type: string
  rule_method: string
  rule_definition: string
  action: string
  application_id: string
  rule_string: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  category: string
}

// Local Table component
function truncateString(value: unknown, maxLength = 60): string {
  if (typeof value === 'string' && value.length > maxLength) {
    return value.substring(0, maxLength) + '...'
  }
  return String(value)
}

interface TableProps<T> {
  columns: ColumnDef<T>[]
  data: T[]
  className?: string
  onSelectionChange?: (rows: Row<T>[]) => void
}

function Table<T extends object>({
  columns,
  data,
  className,
  onSelectionChange,
}: TableProps<T>) {
  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
  })

  // Notify parent of selected rows
  const selectedRows = table.getSelectedRowModel().rows
  if (onSelectionChange) onSelectionChange(selectedRows)

  return (
    <div className={clsx('w-full shadow-md rounded-lg', className)}>
      <table className="min-w-full table-auto border-collapse bg-white">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
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
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id} className="hover:bg-gray-50 transition-colors">
              {row.getVisibleCells().map((cell) => {
                const rendered = flexRender(cell.column.columnDef.cell, cell.getContext())

                return (
                  <td
                    key={cell.id}
                    className="px-6 py-4 text-sm text-gray-900 whitespace-nowrap max-w-xs truncate"
                    title={typeof rendered === 'string' ? rendered : undefined}
                  >
                    {typeof rendered === 'string'
                      ? truncateString(rendered)
                      : rendered}
                  </td>
                )
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

const columns: ColumnDef<Rule>[] = [
  {
    id: 'select',
    header: ({ table }) => (
      <input
        type="checkbox"
        checked={table.getIsAllPageRowsSelected()}
        onChange={table.getToggleAllPageRowsSelectedHandler()}
      />
    ),
    cell: ({ row }) => (
      <input
        type="checkbox"
        checked={row.getIsSelected()}
        onChange={row.getToggleSelectedHandler()}
      />
    ),
  },
  {
    header: 'Status',
    accessorKey: 'is_active',
    cell: ({ row }) => (
      <span
        className={`px-3 py-1 rounded-full text-sm font-medium ${
          row.original.is_active
            ? 'bg-green-100 text-green-800'
            : 'bg-red-100 text-red-800'
        }`}
      >
        {row.original.is_active ? 'Active' : 'Inactive'}
      </span>
    ),
  },
  {
    header: 'Definition',
    accessorKey: 'rule_definition',
    cell: ({ row }) => {
      const value = row.original.rule_definition
      return <span>{value.length > 60 ? `${value.slice(0, 60)}...` : value}</span>
    },
  },
  {
    header: 'Action',
    accessorKey: 'action',
    cell: ({ row }) => (
      <div className="flex items-center gap-2">
        {row.original.action.toLowerCase().includes('deny') ? (
          <HiOutlineBan className="text-red-600" size={20} />
        ) : (
          <HiOutlineCheckCircle className="text-green-600" size={20} />
        )}
        {row.original.action}
      </div>
    ),
  },
  {
    header: 'Message',
    accessorKey: 'category',
  },
  {
    header: 'Updated At',
    accessorKey: 'updated_at',
  },
  {
    header: 'Edit',
    cell: ({ row }) => (
      <button onClick={() => handleEdit(row.original)}>
        <HiOutlinePencilAlt size={20} className="text-blue-600" />
      </button>
    ),
  },
]

const mockData: Rule[] = [
  {
    rule_id: '592249559871523992',
    rule_type: 'multiple',
    rule_method: 'chained',
    rule_definition:
      '[{\"rule_type\":\"REQUEST_URI\",\"rule_method\":\"streq\",\"rule_definition\":\"admin\"},{\"rule_type\":\"REQUEST_URI\",\"rule_method\":\"contains\",\"rule_definition\":\"test\"}]',
    action: 'deny',
    application_id: '16d3f539-6c7b-45ac-b977-6a51c3582d29',
    rule_string:
      'SecRule REQUEST_URI "@streq admin" "id:592249559871523992,phase:2,deny,msg:\'blocked path\'"',
    created_by: 'user',
    created_at: '',
    updated_at: '',
    is_active: true,
    category: 'blocked path',
  },
  {
    rule_id: '592249559871523993',
    rule_type: 'multiple',
    rule_method: 'chained',
    rule_definition:
      '[{\"rule_type\":\"REQUEST_URI\",\"rule_method\":\"streq\",\"rule_definition\":\"admin\"}]',
    action: 'allow',
    application_id: '16d3f539-6c7b-45ac-b977-6a51c3582d29',
    rule_string: '',
    created_by: 'user',
    created_at: '',
    updated_at: '',
    is_active: false,
    category: 'allowed path',
  },
]

const applications = [
  { id: '16d3f539-6c7b-45ac-b977-6a51c3582d29', name: 'MyApp 1' },
  { id: 'app-2', name: 'Dashboard' },
  { id: 'app-3', name: 'API Service' },
]

function CustomRules() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedApp, setSelectedApp] = useState<string | null>(null)
  const [selectedRows, setSelectedRows] = useState<Row<Rule>[]>([])
  const [isEditMode, setIsEditMode] = useState(false)
  const [currentRule, setCurrentRule] = useState<Rule | null>(null)

  const toggleModal = () => setIsModalOpen(!isModalOpen)

  const handleEdit = (rule: Rule) => {
    setIsEditMode(true)
    setCurrentRule(rule)
    toggleModal() // open modal to edit
  }

  const handleAction = (action: 'delete' | 'activate' | 'deactivate') => {
    const selectedIds = selectedRows.map((r) => r.original.rule_id)
    console.log(`Performing ${action} on`, selectedIds)
    // Implement logic to delete or update status
  }

  return (
    <div className="space-y-4">
      <CreateRuleModal
        isModalOpen={isModalOpen}
        onClose={toggleModal}
        rule={currentRule} // Pass current rule for editing
        isEditMode={isEditMode} // Indicate whether it's edit mode
      />

      <Card className="flex justify-between items-center py-4 px-6 shadow-md bg-white rounded-lg">
        <h2 className="text-lg font-semibold">Custom Rules</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          {isEditMode ? 'Edit Rule' : 'Add Rule'}
        </Button>
      </Card>

      {selectedRows.length > 0 && (
        <Card className="flex items-center justify-start gap-4 px-4 py-2 border bg-white rounded-md shadow-sm">
          <span className="text-sm text-gray-700">{selectedRows.length} selected</span>

          <button
            onClick={() => handleAction('delete')}
            className="flex items-center gap-1 text-red-600 hover:text-red-800"
          >
            <HiOutlineTrash size={18} /> Delete
          </button>

          <button
            onClick={() => handleAction('activate')}
            className="flex items-center gap-1 text-green-600 hover:text-green-800"
          >
            <HiOutlineThumbUp size={18} /> Activate
          </button>

          <button
            onClick={() => handleAction('deactivate')}
            className="flex items-center gap-1 text-yellow-600 hover:text-yellow-800"
          >
            <HiOutlineThumbDown size={18} /> Deactivate
          </button>
        </Card>
      )}

      <Table<Rule>
        columns={columns}
        data={mockData}
        onSelectionChange={setSelectedRows}
      />
    </div>
  )
}

export default CustomRules
