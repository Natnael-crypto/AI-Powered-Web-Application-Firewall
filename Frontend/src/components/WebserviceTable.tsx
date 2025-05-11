import {ColumnDef} from '@tanstack/react-table'
import {Application} from '../pages/WebServices'
import Table from './Table'
import {DropdownActions} from './DropdownAction'

interface webserviceTableProps {
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
          className={`px-3 py-1 ull text-sm font-medium ${
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
          className={`px-3 py-1 ull text-sm font-medium ${
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
      header: 'Updated At',
      accessorKey: 'updated_at',
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

const mockData: Application[] = [
  {
    application_id: '16d3f539-6c7b-45ac-b977-6a51c3582d29',
    application_name: 'waf',
    description: 'this is for tls check',
    hostname: 'waf.local',
    ip_address: '127.0.0.1',
    port: '5500',
    status: true,
    tls: false,
    created_at: '2025-03-10T19:57:54.553735+03:00',
    updated_at: '2025-03-10T19:57:54.553735+03:00',
  },
  {
    application_id: 'e9c6df27-7ab7-43d5-9309-cb967c3b54a4',
    application_name: 'auth-service',
    description: 'handles user authentication',
    hostname: 'auth.internal',
    ip_address: '192.168.0.10',
    port: '8080',
    status: true,
    tls: true,
    created_at: '2025-03-05T10:15:20.123456+03:00',
    updated_at: '2025-04-01T08:45:00.000000+03:00',
  },
  {
    application_id: 'f1a8b40d-acc9-4014-bd80-928bb6e23af3',
    application_name: 'analytics',
    description: 'collects usage metrics',
    hostname: 'analytics.service',
    ip_address: '10.10.10.5',
    port: '3000',
    status: false,
    tls: false,
    created_at: '2025-02-12T14:30:00.000000+03:00',
    updated_at: '2025-03-28T17:25:00.000000+03:00',
  },
  {
    application_id: 'a8e284a0-2d9d-4bc5-8c89-fbde1b0c6fc0',
    application_name: 'payment-gateway',
    description: 'handles transactions and payments',
    hostname: 'payments.local',
    ip_address: '10.0.0.2',
    port: '443',
    status: true,
    tls: true,
    created_at: '2025-01-20T11:11:11.111111+03:00',
    updated_at: '2025-04-10T09:00:00.000000+03:00',
  },
]

function WebserviceTable({data, setSelectedApp, openModal}: webserviceTableProps) {
  const columns = getColumns({setSelectedApp, openModal})
  return <Table columns={columns} data={data.length === 0 ? mockData : data} />
}

export default WebserviceTable
