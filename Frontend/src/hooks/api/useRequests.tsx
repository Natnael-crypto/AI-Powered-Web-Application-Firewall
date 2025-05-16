import {useQuery} from '@tanstack/react-query'
import {getDeviceStat, getRequests} from '../../services/requestApi'
import { Filter } from '../../lib/types'

export function useGetRequests(filter: Filter) {
  return useQuery({
    queryKey: ['requests', filter],
    queryFn: () => getRequests(filter),
  })
}

export function useGetDeviceStat(selectedApp:string,timeRange:any) {
  return useQuery({
    queryKey: ['devices'],
    queryFn: ()=> getDeviceStat(selectedApp,timeRange),
  })
}
