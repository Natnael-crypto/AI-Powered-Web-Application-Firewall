import { useState, useEffect } from 'react'
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
import axios from 'axios'
import EditRuleModal from '../components/EditRuleModal'

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
const backendUrl = import.meta.env.VITE_BACKEND_URL;
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDY1NjU2MTgsInJvbGUiOiJzdXBlcl9hZG1pbiIsInVzZXJfaWQiOiJiNGM1ZjI0OC1iOTE3LTQyNDMtYjE0ZS1kNmI4NWQ2NzZjODgifQ.bmGqOlhKhxD4IsMKsomGpa04uExS6l_q5YvrPa2dMCc";

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



// 🟢 Main Component
function CustomRules() {
  const [rules, setRules] = useState<Rule[]>([])
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedRows, setSelectedRows] = useState<Row<Rule>[]>([])
  const [isEditMode, setIsEditMode] = useState(false)
  const [currentRuleID, setCurrentRuleID] = useState("")

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
        <button onClick={() => {
          setCurrentRuleID(row.original.rule_id)
          setIsEditMode(true)
          console.log(isEditMode)
          // setIsModalOpen(true)
        }}>
          <HiOutlinePencilAlt size={20} className="text-blue-600" />
        </button>
      ),
    }
    
  ]

  const fetchRules = async () => {
    try {
      const res = await axios.get(`${backendUrl}/rule`, {
        headers: {
          Authorization: token,
        },
      })
      setRules(res.data.rules)
    } catch (error) {
      console.error('Failed to fetch rules:', error)
    }
  }

  useEffect(() => {
   fetchRules()
  }, [])

  const toggleModal = () => {
    setIsModalOpen(!isModalOpen)
  }

  const handleAction = async (action: 'delete' | 'activate' | 'deactivate') => {
    const selectedIds = selectedRows.map((r) => r.original.rule_id)
  
    try {
      for (const ruleId of selectedIds) {
        if (action === 'delete') {
          await axios.delete(`${backendUrl}/rule/delete/${ruleId}`, {
            headers: {
              Authorization: token,
            },
          })
        } else if (action === 'activate') {
          await axios.get(`${backendUrl}/rule/activate/${ruleId}`, {
            headers: {
              Authorization: token,
            },
          })
        } else if (action === 'deactivate') {
          await axios.get(`${backendUrl}/rule/deactivate/${ruleId}`, {
            headers: {
              Authorization: token,
            },
          })
        }
      }
  
      // Refresh data after performing actions
      await fetchRules()
      setSelectedRows([]) // clear selection
    } catch (error) {
      console.error(`Failed to ${action} rules:`, error)
    }
  }
  

  return (
    <div className="space-y-4">
      {isModalOpen ? (
        <CreateRuleModal/>
      ) : null}

      {isEditMode ? 
        (<EditRuleModal
          ruleID={currentRuleID}
          onClose={() => {
            setCurrentRuleID("")
            setIsEditMode(false)
          }}
          onSuccess={fetchRules}
        />) : null
      }

      <Card className="flex justify-between items-center py-4 px-6 shadow-md bg-white rounded-lg">
        <h2 className="text-lg font-semibold">Custom Rules</h2>
        <Button
          classname={`text-white uppercase ${isModalOpen ? 'bg-red-500' : ''}`}
          size="l"
          variant="primary"
          onClick={() => {
            setIsEditMode(false)
            setCurrentRule(null)
            toggleModal()
          }}
        >
          {isModalOpen ? 'Back' : 'Add Rule'}
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
        data={rules}
        onSelectionChange={setSelectedRows}
      />
    </div>
  )
}

export default CustomRules


