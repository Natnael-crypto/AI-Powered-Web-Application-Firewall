import {useMutation, useQuery} from '@tanstack/react-query'
import {getNotifications, markNotificationAsRead} from '../../services/NotificationsApi'

export function useGetNotifications() {
  return useQuery({
    queryKey: ['notifications'],
    queryFn: getNotifications,
  })
}

export function useMarkNotificationAsRead() {
  return useMutation({
    mutationKey: ['markNotificationAsRead'],
    mutationFn: markNotificationAsRead,
  })
}
