import {useQuery} from '@tanstack/react-query'
import {getRequests} from '../../services/requestApi'

export function useGetRequests(application: string) {
  return useQuery({
    queryKey: ['requests'],
    queryFn: () => getRequests(application),
  })
}
