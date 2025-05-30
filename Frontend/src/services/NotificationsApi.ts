import axios from '../lib/axios'
import {NotificationUpdate} from '../lib/types'

export async function getNotifications(user_id: string | undefined) {
  const response = await axios.get(`/api/notifications/all/${user_id}`, {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.notifications
}

export async function markNotificationAsRead(notification_id: string) {
  const response = await axios.put(`/api/notifications/update/${notification_id}`)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function deleteNotification(notification_id: string) {
  const response = await axios.get(`/api/notifications/delete/${notification_id}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data
}

export async function batchMarkAsRead(data: NotificationUpdate) {
  const response = await axios.put('api/notifications/update', data)
  if (!response) throw new Error('Something went wrong!')

  return response.data
}
