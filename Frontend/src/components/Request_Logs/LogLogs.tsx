import {useEffect, useState} from 'react'
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

function LogLogs() {
  const {data, isLoading, error} = useGetRequests('waf.local')

  const columns: ColumnDef<RequestLog>[] = [
    {
      header: 'Status',
      accessorKey: 'status',
      cell: ({getValue}: CellContext<RequestLog, unknown>) => (
        <div
          className={`rounded-md py-1 px-2 text-white text-center font-semibold text-sm shadow-md ${
            (getValue() as string).toLowerCase() === 'blocked'
              ? 'bg-red-700'
              : 'bg-yellow-400'
          }`}
        >
          {getValue() as string}
        </div>
      ),
    },
    {header: 'Application', accessorKey: 'application_name'},
    {header: 'Request Method', accessorKey: 'request_method'},
    {header: 'Request URL', accessorKey: 'request_url'},
    {header: 'Threat Type', accessorKey: 'threat_type'},
    {header: 'Client IP', accessorKey: 'client_ip'},
    {header: 'Geo Location', accessorKey: 'geo_location'},
    {
      header: 'Response Code',
      accessorKey: 'response_code',
      cell: ({getValue}: CellContext<RequestLog, unknown>) => (
        <div className="text-center font-semibold text-sm">{getValue() as number}</div>
      ),
    },
    {header: 'Timestamp', accessorKey: 'timestamp'},
  ]

  if (isLoading) return <div className="text-center text-lg font-bold">Loading...</div>
  if (error)
    return (
      <div className="text-center text-red-600 font-semibold">
        Error loading data: {error.message}
      </div>
    )

  return (
    <div className="p-4 bg-gray-100 rounded-lg shadow-lg">
      <h2 className="text-xl font-bold mb-4">Request Logs</h2>
      <Table data={data || []} columns={columns} />
    </div>
  )
}

export default LogLogs
