import {CellContext, ColumnDef} from '@tanstack/react-table'
import Table from '../Table'
import {useGetRequests} from '../../hooks/useRequests'

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
  ai_analysis_result: string
}

function AttackLogTable() {
  const {data, isLoading, error} = useGetRequests('waf.local')

  const columns: ColumnDef<RequestLog>[] = [
    {
      header: 'Status',
      accessorKey: 'status',
      cell: ({getValue}: CellContext<RequestLog, unknown>) => (
        <div
          className={`rounded-full py-1 px-3 text-white text-xs font-medium shadow-sm inline-block ${
            (getValue() as string).toLowerCase() === 'blocked'
              ? 'bg-red-600'
              : 'bg-yellow-500 text-gray-900'
          }`}
        >
          {getValue() as string}
        </div>
      ),
    },
    {header: 'Application', accessorKey: 'application_name'},
    {header: 'Method', accessorKey: 'request_method'},
    {
      header: 'URL',
      accessorKey: 'request_url',
      cell: ({getValue}) => (
        <div className="text-sm text-blue-600 truncate max-w-[300px]">
          {getValue() as string}
        </div>
      ),
    },
    {
      header: 'Threat Type',
      accessorKey: 'threat_type',
      cell: ({getValue}) => (
        <span className="text-sm font-medium text-red-700">{getValue() as string}</span>
      ),
    },
    {
      header: 'IP',
      accessorKey: 'client_ip',
      cell: ({getValue}) => (
        <code className="text-xs text-gray-500">{getValue() as string}</code>
      ),
    },
    {header: 'Location', accessorKey: 'geo_location'},
    {
      header: 'Code',
      accessorKey: 'response_code',
      cell: ({getValue}) => (
        <div className="text-sm font-semibold text-center">{getValue() as number}</div>
      ),
    },
    {
      header: 'Time',
      accessorKey: 'timestamp',
      cell: ({getValue}) => (
        <div className="text-xs text-gray-500">
          {new Date(getValue() as string).toLocaleString()}
        </div>
      ),
    },
  ]

  if (isLoading)
    return (
      <div className="text-center text-lg font-bold text-green-700">
        Loading request logs...
      </div>
    )

  if (error)
    return (
      <div className="text-center text-red-600 font-semibold">
        Error loading data: {error.message}
      </div>
    )

  return (
    <div className="bg-white p-6 rounded-2xl shadow-xl border border-gray-200">
      <h2 className="text-2xl font-bold text-green-800 mb-4">Request Logs</h2>
      <Table data={data || []} columns={columns} />
    </div>
  )
}

export default AttackLogTable
