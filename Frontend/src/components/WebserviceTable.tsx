import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Application} from '../lib/types'
import ApplicationConfigModal from './ApplicationConfigModal'
import {useState} from 'react'
import {useUpdateDetectBot} from '../hooks/api/useApplication'
import {Edit, Settings, Trash2} from 'lucide-react'

interface WebserviceTableProps {
  data: Application[]
  setSelectedApp?: (app: Application) => void
  openModal: () => void
  selectedApp?: Application
  handleDelete?: (application_id: string) => void
}

function getColumns({
  setSelectedApp,
  openModal,
  setIsConfigModalOpen,
  toggleBotDetection,
  handleDelete,
}: {
  setSelectedApp?: (app: Application) => void
  openModal: () => void
  setIsConfigModalOpen: (bool: boolean) => void
  toggleBotDetection: (application_id: string, detectBot: boolean) => void
  handleDelete?: (application_id: string) => void
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
        const [isLoading, setIsLoading] = useState(false)
        const detectBot = row.original.config?.detect_bot ?? false

        const handleToggle = async () => {
          if (isLoading) return

          setIsLoading(true)
          try {
            const newValue = await toggleBotDetection(
              row.original.application_id,
              detectBot,
            )
            console.log(
              `Bot detection updated to ${newValue} for app ${row.original.application_id}`,
            )
          } catch (error) {
            console.error('Failed to update bot detection:', error)
          } finally {
            setIsLoading(false)
          }
        }

        return (
          <div className="flex items-center space-x-2">
            <button
              onClick={handleToggle}
              disabled={isLoading}
              className={`relative inline-flex h-4 w-8 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 ${
                detectBot ? 'bg-purple-600' : 'bg-gray-200'
              }`}
            >
              <span
                className={`inline-block h-3 w-3 transform rounded-full bg-white transition-transform ${
                  detectBot ? 'translate-x-4' : 'translate-x-1'
                } ${isLoading ? 'opacity-50' : ''}`}
              />
            </button>
            <span className="text-sm text-gray-900">
              {isLoading ? 'Updating...' : detectBot ? 'Enabled' : 'Disabled'}
            </span>
          </div>
        )
      },
    },
    {
      header: 'Max Data Size',
      accessorKey: 'config.max_post_data_size',
      cell: ({row}) => (
        <div className=" text-center">
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
}: WebserviceTableProps) {
  const [isConfigModalOpen, setIsConfigModalOpen] = useState(false)
  const {mutate: updateDetectBOT} = useUpdateDetectBot()

  const toggleBotDetection = async (appId: string, currentValue: boolean) => {
    return updateDetectBOT({application_id: appId, data: {detect_bot: !currentValue}})
  }

  const columns = getColumns({
    setSelectedApp,
    openModal,
    setIsConfigModalOpen,
    toggleBotDetection,
    handleDelete,
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
