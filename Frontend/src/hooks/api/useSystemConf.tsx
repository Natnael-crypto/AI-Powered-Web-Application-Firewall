import {useMutation, useQuery} from '@tanstack/react-query'
import { getSysConf, updateSysPort, updateSysRemoteLogIp } from '../../services/configApi'

export function useGetSysConf() {
  return useQuery({
    queryKey: ['getConf'],
    queryFn: () => getSysConf(),
  })
}

export function useUpdateSysRemoteLogIp() {
  return useMutation({
    mutationFn: (ip: string) => updateSysRemoteLogIp(ip),
  })
}

export function useUpdateSysPort() {
  return useMutation({
    mutationFn: (port: string) => updateSysPort(port),
  })
}
