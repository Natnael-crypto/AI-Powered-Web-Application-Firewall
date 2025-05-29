import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Notification} from '../lib/types'

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
      cell: ({row}) => <div className="whitespace-pre-line">{row.original.message}</div>,
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
        !row.original.status?
        <DropdownActions
          item={row.original}
          options={[
            {
              label: 'Mark as Read',
              onClick: notification => onMarkAsRead(notification.notification_id),
            },
          ]}
        />:null
      ),
    },
  ]
}
export function NotificationsTable({data, onMarkAsRead}: NotificationsTableProps) {
  const columns = getColumns({onMarkAsRead})

  return <Table columns={columns} data={data} />
}

export default NotificationsTable
