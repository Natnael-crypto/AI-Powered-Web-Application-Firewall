import {useQuery} from '@tanstack/react-query'
import { createSysEmail, getSysEmail, updateSysEmail } from '../../services/configApi'

export function useGetSysEmail() {
  return useQuery({
    queryKey: ['getEmail'],
    queryFn: () => getSysEmail(),
  })
}

export function useAddSysEmail(email:string,active:boolean) {
  return useQuery({
    queryKey: ['createEmail'],
    queryFn: ()=> createSysEmail(email,active),
  })
}

export function useUpdateSysEmail(email:string,active:boolean) {
  return useQuery({
    queryKey: ['updateEmail'],
    queryFn: ()=> updateSysEmail(email,active),
  })
}