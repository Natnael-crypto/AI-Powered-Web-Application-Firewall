import axios from '../lib/axios'

export async function getNotifications() {
  const response = await axios.get('/api/notifications/all', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.model
}

export async function markNotificationAsRead(notification_id: string) {
  const response = await axios.put(`/api/notifications/update/:${notification_id}`)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function deleteNotification(notification_id: string) {
  const response = await axios.get(`/api/notifications/delete/:${notification_id}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data
}
