import {CellContext, ColumnDef} from '@tanstack/react-table'
import Table from '../../components/Table'

function AntiBot() {
  const columns: ColumnDef<any>[] = [
    {
      header: 'IP Addr',
      accessorKey: 'ip',
      cell: ({getValue}: CellContext<any, unknown>) => (
        <div className=" py-1 px-2 font-semibold text-sm text-black">
          {String(getValue())}
        </div>
      ),
    },
    {
      header: 'Web Service',
      accessorKey: 'services',
      cell: ({row}) => {
        const service = row.original.services
        const sub = row.original.serviceSub || ''
        return (
          <div className="flex flex-col text-sm">
            <span className="text-blue-700 font-semibold">{service}</span>
            {sub && <span className="text-xs text-gray-500">{sub}</span>}
          </div>
        )
      },
    },
    {
      header: 'Detail',
      accessorKey: 'detail',
      cell: ({row}) => {
        const hits = row.original.hits ?? 0
        const verified = row.original.verified ?? 0
        return (
          <div className="flex flex-col text-sm">
            <span className="font-semibold">Hits {hits}</span>
            <span className="font-semibold">Verified {verified}</span>
          </div>
        )
      },
    },
    {
      header: 'Duration',
      accessorKey: 'duration',
    },
    {
      header: 'Start At',
      accessorKey: 'startAt',
    },
  ]

  const data = [
    {
      ip: '206.168.34.79',
      services: 'Anti-Replay Demo',
      serviceSub: 'Match All Host',
      detail: '',
      hits: 2,
      verified: 0,
      duration: '1 minutes',
      startAt: '2025-01-02 07:56:11',
    },
  ]

  return (
    <div className="p-4 bg-gray-100  shadow-lg">
      <h2 className="text-xl font-bold mb-4">Request Logs</h2>
      <Table data={data} columns={columns} />
    </div>
  )
}

export default AntiBot
