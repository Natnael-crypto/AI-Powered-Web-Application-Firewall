import {useQuery} from '@tanstack/react-query'
import {getDeviceStat, getRequests} from '../../services/requestApi'

export function useGetRequests(application: string) {
  return useQuery({
    queryKey: ['requests', application],
    queryFn: () => getRequests(application),
  })
}

export function useGetDeviceStat() {
  return useQuery({
    queryKey: ['devices'],
    queryFn: getDeviceStat,
  })
}
