import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {
  batchMarkAsRead,
  getNotifications,
  markNotificationAsRead,
} from '../../services/NotificationsApi'

export function useGetNotifications(user_id: string | undefined) {
  return useQuery({
    queryKey: ['notifications'],
    queryFn: () => getNotifications(user_id),
  })
}

export function useMarkNotificationAsRead() {
  return useMutation({
    mutationKey: ['markNotificationAsRead'],
    mutationFn: markNotificationAsRead,
  })
}
export function useBatchMarkNotificationAsRead() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['markNotificationAsRead'],
    mutationFn: batchMarkAsRead,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['notifications']}),
  })
}
