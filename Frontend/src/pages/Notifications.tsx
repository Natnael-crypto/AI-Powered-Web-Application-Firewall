// Notifications.tsx
import {useState} from 'react'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import {NotificationsTable} from '../components/NotificationsTable'
import {
  useGetNotifications,
  useMarkNotificationAsRead,
} from '../hooks/api/useNotifications'

function Notifications() {
  const {data: notifications = [], refetch} = useGetNotifications()
  const {mutate: markAsRead} = useMarkNotificationAsRead()

  const [isRefreshing, setIsRefreshing] = useState(false)

  const handleMarkAsRead = (notificationId: string) => {
    markAsRead(notificationId, {
      onSuccess: () => {
        refetch()
      },
      onError: () => {
        console.error('Failed to mark notification as read')
      },
    })
  }

  const handleRefresh = async () => {
    setIsRefreshing(true)
    try {
      await refetch()
    } catch (error) {
      console.error('Failed to refresh notifications', error)
    } finally {
      setIsRefreshing(false)
    }
  }

  return (
    <div className="space-y-4">
      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Notifications</h2>
        <div className="flex space-x-2">
          <Button
            classname="text-white uppercase"
            size="l"
            variant="secondary"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            Refresh
          </Button>
        </div>
      </Card>

      <Card className="shadow-md p-4 bg-white">
        <NotificationsTable data={notifications} onMarkAsRead={handleMarkAsRead} />
      </Card>
    </div>
  )
}

export default Notifications
