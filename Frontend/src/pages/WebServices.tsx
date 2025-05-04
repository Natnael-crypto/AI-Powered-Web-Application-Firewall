import {useState} from 'react'
import {ColumnDef} from '@tanstack/react-table'
import Table from '../components/Table'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import WebServiceModal from '../components/WebServiceModal'

export interface Application {
  application_id: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
  created_at: string
  updated_at: string
}

const columns: ColumnDef<Application>[] = [
  {
    header: 'Status',
    accessorKey: 'status',
    cell: ({row}) => (
      <span
        className={`px-3 py-1 ull text-sm font-medium ${
          row.original.status ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
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
]

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

function WebService() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)

  return (
    <div className="space-y-4">
      <WebServiceModal
        isOpen={isModalOpen}
        onClose={toggleModal}
        onSubmit={data => {
          console.log('Sending this data:', data)
        }}
      />

      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Web Services</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          Add Service
        </Button>
      </Card>
      <Card className="shadow-md p-4 bg-white ">
        <Table columns={columns} data={mockData} />
      </Card>
    </div>
  )
}

export default WebService
