import {useQuery} from '@tanstack/react-query'
import { getSysConf, updateSysPort, updateSysRemoteLogIp } from '../../services/configApi'

export function useGetSysConf() {
  return useQuery({
    queryKey: ['getConf'],
    queryFn: () => getSysConf(),
  })
}

export function useUpdateSysPort(port:string) {
  return useQuery({
    queryKey: ['updatePort'],
    queryFn: ()=> updateSysPort(port),
  })
}

export function useUpdateSysRemoteLogIp(ip:string) {
  return useQuery({
    queryKey: ['updateRemoteLogIp'],
    queryFn: ()=> updateSysRemoteLogIp(ip),
  })
}