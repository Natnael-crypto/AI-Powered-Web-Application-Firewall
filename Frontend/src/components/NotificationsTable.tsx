import {useState} from 'react'
import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Notification} from '../lib/types'
import NotificationDetail from './NotificationDetail'

interface NotificationsTableProps {
  data: Notification[]
  onMarkAsRead: (notificationId: string) => void
}

function getColumns({
  onMarkAsRead,
}: {
  onMarkAsRead: (notificationId: string) => void
}): ColumnDef<Notification>[] {
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
          {row.original.status ? 'Read' : 'Unread'}
        </span>
      ),
    },
    {
      header: 'Message',
      accessorKey: 'message',
      cell: ({row}) => (
        <div className="whitespace-pre-line">
          {row.original.message.slice(0, 40)}
          {row.original.message.length > 40 ? '...' : ''}
        </div>
      ),
    },
    {
      header: 'Timestamp',
      accessorKey: 'timestamp',
      cell: ({row}) => new Date(row.original.timestamp).toLocaleString(),
    },
    {
      header: 'Actions',
      id: 'actions',
      cell: ({row}) => (
        <DropdownActions
          item={row.original}
          options={[
            {
              label: 'Mark as Read',
              onClick: notification => onMarkAsRead(notification.notification_id),
            },
          ]}
        />
      ),
    },
  ]
}

export function NotificationsTable({data, onMarkAsRead}: NotificationsTableProps) {
  const [selectedNotification, setSelectedNotification] = useState<Notification | null>(
    null,
  )
  const [isModalOpen, setIsModalOpen] = useState(false)

  const columns = getColumns({onMarkAsRead})

  const handleRowClick = (notification: Notification) => {
    setSelectedNotification(notification)
    setIsModalOpen(true)
  }

  const handleCloseModal = () => {
    setIsModalOpen(false)
    setSelectedNotification(null)
  }

  return (
    <>
      <Table columns={columns} data={data} onRowClick={handleRowClick} />

      <NotificationDetail
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        notification={selectedNotification ?? undefined}
        onMarkAsRead={onMarkAsRead}
      />
    </>
  )
}

export default NotificationsTable
