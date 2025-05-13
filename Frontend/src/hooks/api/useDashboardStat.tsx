import {useQuery} from '@tanstack/react-query'
import {
  getMostTargetedEndpoint,
  getOverAllStat,
  getRateStat,
  getTopAttackTypes,
} from '../../services/DashboardApi'

export function useGetMostTargetedEndpoint(appId: string,time:any) {
  return useQuery({
    queryKey: ['mostTargetedEndpoint'],
    queryFn: ()=>  getMostTargetedEndpoint(appId,time),
  })
}
export function useGetTopThreatTypes(appId: string,time:any) {
  return useQuery({
    queryKey: ['TopAttackTypes'],
    queryFn:()=> getTopAttackTypes(appId,time),
  })
}

export function useGetOverAllStat(appId: string,time:any) {
  return useQuery({
    queryKey: ['overAllstat'],
    queryFn: () => getOverAllStat(appId,time),
  })
}

export function useRateStat(appId: string,time:any) {
  return useQuery({
    queryKey: ['overAllstat'],
    queryFn: () => getRateStat(appId,time),
  })
}

