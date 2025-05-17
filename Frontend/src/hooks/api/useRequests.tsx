import {useQuery} from '@tanstack/react-query'
import {generateRequest, getDeviceStat, getRequests} from '../../services/requestApi'
import { Filter } from '../../lib/types'

export function useGetRequests(filters: Partial<Filter>) {
  return useQuery({
    queryKey: ['requests', filters],
    queryFn: () => getRequests(filters),
  })
}

export function useGetDeviceStat(selectedApp:string,timeRange:any) {
  return useQuery({
    queryKey: ['devices'],
    queryFn: ()=> getDeviceStat(selectedApp,timeRange),
  })
}

export function useGenerateRequests(filters: Partial<Filter>) {
  return useQuery({
    queryKey: ['generate'],
    queryFn: ()=> generateRequest(filters),
  })
}
