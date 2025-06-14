import React from 'react'
import Modal from './Modal'
import {Notification} from '../lib/types'

interface NotificationDetailsModalProps {
  isOpen: boolean
  onClose: () => void
  notification: Notification | undefined
  onMarkAsRead: (notificationId: string) => void
}

const NotificationDetail: React.FC<NotificationDetailsModalProps> = ({
  isOpen,
  onClose,
  notification,
  onMarkAsRead,
}) => {
  const formatTimestamp = (timestamp: string) => {
    try {
      const date = new Date(timestamp)
      return date.toLocaleString()
    } catch (e) {
      return timestamp
    }
  }

  const getStatusBadge = (status: boolean) => {
    const baseClasses =
      'px-3 py-1 rounded-full text-xs font-medium inline-flex items-center'
    return status
      ? `${baseClasses} bg-green-100 text-green-800`
      : `${baseClasses} bg-red-100 text-red-800`
  }

  const handleMarkAsRead = () => {
    if (notification && !notification.status) {
      onMarkAsRead(notification.notification_id)
    }
  }

  const parseMessageLines = (message: string) => {
    return message.split('\n').map((line, index) => (
      <div key={index} className="py-1">
        {line}
      </div>
    ))
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Notification Details">
      <div className="space-y-6">
        {/* Basic Info Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-1">
          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </label>
            <div className={getStatusBadge(notification?.status ?? false)}>
              {notification?.status ? 'Read' : 'Unread'}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Notification ID
            </label>
            <div className="text-sm text-gray-800 font-mono break-all">
              {notification?.notification_id}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              User ID
            </label>
            <div className="text-sm text-gray-800 font-mono break-all">
              {notification?.user_id}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created At
            </label>
            <div className="text-sm text-gray-800">
              {notification?.timestamp ? formatTimestamp(notification.timestamp) : '-'}
            </div>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-gray-200 my-2"></div>

        {/* Message Section */}
        <div className="space-y-2">
          <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
            Message Details
          </label>
          <div className="text-sm text-gray-800 p-3 bg-gray-50 rounded border border-gray-200">
            {notification?.message
              ? parseMessageLines(notification.message)
              : 'No message'}
          </div>
        </div>

        {/* Footer */}
        <div className="pt-4 border-t border-gray-200 flex justify-between">
          {notification && !notification.status && (
            <button
              onClick={handleMarkAsRead}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Mark as Read
            </button>
          )}

          <button
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 ml-auto"
          >
            Close
          </button>
        </div>
      </div>
    </Modal>
  )
}

export default NotificationDetail
