import { useEffect, useState } from 'react'
import { CellContext, ColumnDef } from '@tanstack/react-table'
import Table from '../Table'
import { useGetRequests } from '../../hooks/api/useRequests'
import { useLogFilter } from '../../store/LogFilter'
import LoadingSpinner from '../LoadingSpinner'

interface RequestLog {
  request_id: string
  application_name: string
  client_ip: string
  request_method: string
  request_url: string
  headers: string
  body: string
  timestamp: string
  response_code: number
  status: string
  matched_rules: string
  threat_detected: boolean
  threat_type: string
  bot_detected: boolean
  geo_location: string
  rate_limited: boolean
  user_agent: string
  ai_result: string
  rule_detected: boolean
  ai_threat_type: string
}

function AttackLogTable() {
  const { appliedFilters } = useLogFilter()
  const [page, setPage] = useState(1)

  const { data, isLoading, error } = useGetRequests({
    ...appliedFilters,
    page: String(page),
  })

  useEffect(() => {
    setPage(1)
  }, [appliedFilters])

  const columns: ColumnDef<RequestLog>[] = [
    {
      header: 'Status',
      accessorKey: 'status',
      cell: ({ getValue }: CellContext<RequestLog, unknown>) => (
        <div
          className={`py-1 px-3 text-white text-xs font-medium inline-block ${
            (getValue() as string).toLowerCase() === 'blocked'
              ? 'bg-red-600'
              : 'bg-yellow-500 text-gray-900'
          }`}
        >
          {getValue() as string}
        </div>
      ),
    },
    { header: 'Application', accessorKey: 'application_name' },
    { header: 'Method', accessorKey: 'request_method' },
    {
      header: 'URL',
      accessorKey: 'request_url',
      cell: ({ getValue }) => (
        <div className="text-sm text-blue-600 truncate max-w-[300px]">
          {getValue() as string}
        </div>
      ),
    },
    {
      header: 'Threat Type',
      accessorKey: 'threat_type',
      cell: ({ getValue }) => (
        <span className="text-sm font-medium text-red-700">{getValue() as string}</span>
      ),
    },
    {
      header: 'IP',
      accessorKey: 'client_ip',
      cell: ({ getValue }) => (
        <code className="text-xs text-gray-500">{getValue() as string}</code>
      ),
    },
    { header: 'Location', accessorKey: 'geo_location' },
    {
      header: 'Code',
      accessorKey: 'response_code',
      cell: ({ getValue }) => (
        <div className="text-sm font-semibold text-center">{getValue() as number}</div>
      ),
    },
    {
      header: 'Time',
      accessorKey: 'timestamp',
      cell: ({ getValue }) => (
        <div className="text-xs text-gray-500">
          {new Date(getValue() as string).toLocaleString()}
        </div>
      ),
    },
  ]

  if (isLoading) return <LoadingSpinner />

  if (error)
    return (
      <div className="text-center text-red-600 font-semibold mt-4">
        Error loading data: {error.message}
      </div>
    )

  return (
    <div className="bg-white p-6 shadow-xl border border-gray-200">
      <Table data={data?.requests || []} columns={columns} />
      <div className="mt-4 flex justify-center gap-4 items-center">
        <button
          onClick={() => setPage(prev => Math.max(prev - 1, 1))}
          disabled={page === 1 || isLoading}
          className="px-3 py-1 border rounded disabled:opacity-50"
        >
          Prev
        </button>
        <span className="text-sm font-medium">
          Page {data?.current_page} of {data?.total_pages}
        </span>
        <button
          onClick={() =>
            setPage(prev => Math.min(prev + 1, data?.total_pages || prev + 1))
          }
          disabled={page === data?.total_pages || isLoading}
          className="px-3 py-1 border rounded disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  )
}

export default AttackLogTable
