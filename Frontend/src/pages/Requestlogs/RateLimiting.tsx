import {CellContext, ColumnDef} from '@tanstack/react-table'
import Table from '../../components/Table'

const mockData = [
  {
    ip: '47.104.86.168',
    country: 'China',
    serviceName: 'Block Attack Dem',
    serviceUrl: 'demo.waf.com',
    detail: 'Reason: 10 Attacks within 60 seconds\nAction: Block 30 minutes',
    blocked: 9065,
    startAt: '2024-11-21 01:27:24',
  },
  {
    ip: '47.96.26.154',
    country: 'China',
    serviceName: 'Block Attack Dem',
    serviceUrl: 'demo.waf.com',
    detail:
      'Reason: 100 Requests within 10 seconds\nAction: Anti-Bot Challenge 10 minutes',
    blocked: 0,
    startAt: '2024-11-12 18:01:00',
  },
]

function RateLimiting() {
  const columns: ColumnDef<any>[] = [
    {
      header: 'IP Addr',
      accessorKey: 'ip',
      cell: ({row}: CellContext<any, unknown>) => (
        <div>
          <div className="font-semibold text-sm">{row.original.ip}</div>
          <div className="text-xs text-gray-500">{row.original.country}</div>
        </div>
      ),
    },
    {
      header: 'Web Service',
      accessorKey: 'serviceName',
      cell: ({row}: CellContext<any, unknown>) => (
        <div>
          <div className="text-blue-600 font-medium">{row.original.serviceName}</div>
          <div className="text-xs text-gray-500">{row.original.serviceUrl}</div>
        </div>
      ),
    },
    {
      header: 'Detail',
      accessorKey: 'detail',
      cell: ({getValue}) => {
        const value = getValue() as string
        return <div className="text-xs whitespace-pre-line text-gray-700">{value}</div>
      },
    },
    {
      header: 'Blocked',
      accessorKey: 'blocked',
      cell: ({getValue}) => (
        <div className="text-center font-semibold text-sm text-gray-800">
          {getValue() as number}
        </div>
      ),
    },
    {
      header: 'Start At',
      accessorKey: 'startAt',
      cell: ({row}: CellContext<any, unknown>) => (
        <div className="text-sm text-gray-600">
          {row.original.startAt}
          <button className="ml-4 text-blue-600 text-xs font-medium hover:underline">
            Unblock
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className="p-4 bg-white  shadow-lg">
      <h2 className="text-xl font-bold mb-4">Request Logs</h2>
      <Table data={mockData} columns={columns} />
    </div>
  )
}

export default RateLimiting
