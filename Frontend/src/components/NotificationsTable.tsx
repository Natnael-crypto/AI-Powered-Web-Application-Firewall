import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {Notification} from '../lib/types'

export const mockNotifications: Notification[] = [
  {
    notification_id: 'e4ceb5aa-2b02-4ea8-ba2f-cf258114fd06',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'Security Alert: SQL injection attempt detected\nSource IP: 192.168.1.45\nCount: 52\nAction: Blocked automatically',
    timestamp: '2025-04-18T10:28:32.56312-07:00',
    status: false,
  },
  {
    notification_id: 'f88ace8e-2fef-4bae-81af-bdc3e0a70285',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'Threshold exceeded for rule: Brute Force Attack\nSource IP: 203.0.113.17\nCount: 127\nAction: IP temporarily banned',
    timestamp: '2025-04-20T10:41:43.698004-07:00',
    status: false,
  },
  {
    notification_id: 'a1b2c3d4-5678-90ef-ghij-klmnopqrstuv',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'System Update Available\nVersion: 2.5.0\nChanges: Security patches, performance improvements\nDeadline: 2025-05-30',
    timestamp: '2025-04-22T08:15:22.123456-07:00',
    status: true,
  },
  {
    notification_id: '09876543-21fe-dcba-9876-543210abcdef',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'New Login Detected\nLocation: Tokyo, Japan\nDevice: Chrome on Windows\nTime: 2025-04-23 03:45:12 UTC',
    timestamp: '2025-04-23T03:45:12.789012-07:00',
    status: false,
  },
  {
    notification_id: '5f4e3d2c-1b0a-9f8e-7d6c-5b4a3f2e1d0c',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'Scheduled Maintenance\nDate: 2025-05-01\nDuration: 2 hours (01:00 - 03:00 UTC)\nImpact: Minimal downtime expected',
    timestamp: '2025-04-25T14:30:00.000000-07:00',
    status: true,
  },
  {
    notification_id: 'aa11bb22-cc33-dd44-ee55-ff6677889900',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'Critical Security Alert: Zero-day vulnerability detected\nAffected Component: Authentication Service\nAction Required: Immediate update',
    timestamp: '2025-04-27T09:12:35.456789-07:00',
    status: false,
  },
  {
    notification_id: '11223344-5566-7788-9900-aabbccddeeff',
    user_id: 'b4c5f248-b917-4243-b14e-d6b85d676c88',
    message:
      'Backup Completed Successfully\nSize: 45.2 GB\nDuration: 12 minutes\nNext Backup: 2025-05-02 02:00 UTC',
    timestamp: '2025-04-28T02:15:00.333444-07:00',
    status: true,
  },
]

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
  const columns = getColumns({onMarkAsRead})
  return <Table columns={columns} data={data} />
}

export default NotificationsTable
