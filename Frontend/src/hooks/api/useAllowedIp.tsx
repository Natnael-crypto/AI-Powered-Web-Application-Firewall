import {useMutation, useQuery} from '@tanstack/react-query'
import { createAllowedIp, deleteAllowedIp, getAllowedIp, updateAllowedIp } from '../../services/configApi'

export function useAllowedIp() {
  return useQuery({
    queryKey: ['getConf'],
    queryFn: () => getAllowedIp(),
  })
}

export function useCreateAllowedIp() {
  return useMutation({
    mutationFn: ({service,ip}:{service:string;ip: string}) => createAllowedIp(service,ip),
  })
}

export function useUpdateAllowedIp() {
  return useMutation({
    mutationFn: ({id,service,ip}:{id:string,service:string;ip: string}) => updateAllowedIp(id,service,ip),
  })
}

export function useDeleteAllowedIp() {
  return useMutation({
    mutationFn: (id:string) => deleteAllowedIp(id),
  })
}
