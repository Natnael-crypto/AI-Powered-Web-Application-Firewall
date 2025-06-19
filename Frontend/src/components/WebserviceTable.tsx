import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Application} from '../lib/types'
import ApplicationConfigModal from './ApplicationConfigModal'
import {useState} from 'react'
// import {useUpdateDetectBot} from '../hooks/api/useApplication'
import {Edit, Settings, Trash2, Upload} from 'lucide-react'

interface WebserviceTableProps {
  data: Application[]
  setSelectedApp?: (app: Application) => void
  openModal: () => void
  selectedApp?: Application
  handleDelete?: (application_id: string) => void
  setIsCertModalOpen?: (bool: boolean) => void
}

function getColumns({
  setSelectedApp,
  openModal,
  setIsConfigModalOpen,
  // toggleBotDetection,
  handleDelete,
  setIsCertModalOpen,
}: {
  setSelectedApp?: (app: Application) => void
  openModal: () => void
  setIsConfigModalOpen: (bool: boolean) => void
  // toggleBotDetection: (application_id: string, detectBot: boolean) => void
  handleDelete?: (application_id: string) => void
  setIsCertModalOpen?: (bool: boolean) => void
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
      header: 'Max Data Size',
      accessorKey: 'config.max_post_data_size',
      cell: ({row}) => (
        <div className="text-center">
          {row.original.config.max_post_data_size + ' MB'}
        </div>
      ),
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
              icon: <Edit className="mr-2 h-4 w-4" />,
              onClick: app => {
                setSelectedApp?.(app)
                openModal()
              },
            },
            {
              label: 'Update Config',
              icon: <Settings className="mr-2 h-4 w-4" />,
              onClick: app => {
                setSelectedApp?.(app)
                setIsConfigModalOpen(true)
              },
            },
            {
              label: 'Upload Certificates',
              icon: <Upload className="mr-2 h-4 w-4" />,
              onClick: app => {
                setSelectedApp?.(app)
                setIsCertModalOpen?.(true)
              },
            },
            {
              label: 'Delete',
              icon: <Trash2 className="mr-2 h-4 w-4 text-red-600" />,
              onClick: app => {
                setSelectedApp?.(app)
                handleDelete?.(row.original.application_id)
              },
            },
          ]}
        />
      ),
    },
  ]
}

function WebserviceTable({
  data,
  setSelectedApp,
  openModal,
  selectedApp,
  handleDelete,
  setIsCertModalOpen,
}: WebserviceTableProps) {
  const [isConfigModalOpen, setIsConfigModalOpen] = useState(false)
  // const {mutate: updateDetectBOT} = useUpdateDetectBot()

  // const toggleBotDetection = async (appId: string, currentValue: boolean) => {
  //   return updateDetectBOT({application_id: appId, data: {detect_bot: !currentValue}})
  // }

  const columns = getColumns({
    setSelectedApp,
    openModal,
    setIsConfigModalOpen,
    // toggleBotDetection,
    handleDelete,
    setIsCertModalOpen,
  })

  return (
    <>
      <ApplicationConfigModal
        appId={selectedApp?.application_id || ''}
        isOpen={isConfigModalOpen}
        onClose={() => setIsConfigModalOpen(false)}
        data={selectedApp?.config}
      />
      <Table columns={columns} data={data} emptyMessage={'No Data'} />
    </>
  )
}

export default WebserviceTable
