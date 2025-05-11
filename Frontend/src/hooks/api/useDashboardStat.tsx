import {useQuery} from '@tanstack/react-query'
import {
  getMostTargetedEndpoint,
  getOverAllStat,
  getTopAttackTypes,
} from '../../services/DashboardApi'

export function useGetMostTargetedEndpoint() {
  return useQuery({
    queryKey: ['mostTargetedEndpoint'],
    queryFn: getMostTargetedEndpoint,
  })
}
export function useGetTopThreatTypes() {
  return useQuery({
    queryKey: ['TopAttackTypes'],
    queryFn: getTopAttackTypes,
  })
}

export function useGetOverAllStat(appId: string) {
  return useQuery({
    queryKey: ['overAllstat'],
    queryFn: () => getOverAllStat(appId),
  })
}
