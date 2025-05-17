import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Application} from '../lib/types'

interface WebserviceTableProps {
  data: Application[]
  setSelectedApp?: (app: Application) => void
  openModal: () => void
}

function getColumns({
  setSelectedApp,
  openModal,
}: {
  setSelectedApp?: (app: Application) => void
  openModal: () => void
}): ColumnDef<Application>[] {
  return [
    {
      header: 'Status',
      accessorKey: 'status',
      cell: ({row}) => (
        <span
          className={`px-3 py-1 rounded-full text-sm font-medium ${
            row.original.status
              ? 'bg-green-100 text-green-800'
              : 'bg-red-100 text-red-800'
          }`}
        >
          {row.original.status ? 'Running' : 'Stopped'}
        </span>
      ),
    },
    {
      header: 'TLS',
      accessorKey: 'tls',
      cell: ({row}) => (
        <span
          className={`px-3 py-1 rounded-full text-sm font-medium ${
            row.original.tls ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
          }`}
        >
          {row.original.tls ? 'Enabled' : 'Disabled'}
        </span>
      ),
    },
    {
      header: 'Name',
      accessorKey: 'application_name',
    },
    {
      header: 'Description',
      accessorKey: 'description',
    },
    {
      header: 'Host',
      accessorKey: 'hostname',
    },
    {
      header: 'IP Address',
      accessorKey: 'ip_address',
    },
    {
      header: 'Port',
      accessorKey: 'port',
    },
    {
      header: 'Detect Bot',
      accessorKey: 'config.detect_bot',
      cell: ({row}) => {
        const detectBot = row.original.config?.detect_bot ?? false
        return (
          <label className="flex items-center space-x-2">
            <input
              type="radio"
              checked={detectBot}
              readOnly
              className="form-radio h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300"
            />
            <span className="text-sm text-gray-900">
              {detectBot ? 'Enabled' : 'Disabled'}
            </span>
          </label>
        )
      },
    },
    {
      header: 'Updated At',
      accessorKey: 'updated_at',
      cell: ({row}) => new Date(row.original.updated_at).toLocaleString(),
    },
    {
      header: 'Actions',
      id: 'actions',
      cell: ({row}) => (
        <DropdownActions
          item={row.original}
          options={[
            {
              label: 'Update Detail',
              onClick: app => {
                setSelectedApp?.(app)
                openModal()
              },
            },
            {
              label: 'Update Config',
              onClick: app => {
                setSelectedApp?.(app)
                openModal()
              },
            },
          ]}
        />
      ),
    },
  ]
}


function WebserviceTable({data, setSelectedApp, openModal}: WebserviceTableProps) {
  const columns = getColumns({setSelectedApp, openModal})
  return (
    <Table columns={columns} data={data} />
  )
}

export default WebserviceTable
