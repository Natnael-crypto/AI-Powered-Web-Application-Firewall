import axios from '../lib/axios'
export const getSysEmail = async () => {
  const response = await axios.get('/api/sys-email')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}

export const createSysEmail = async (email:string,active:boolean) => {
  const response = await axios.post(`/api/sys-email`,{email:email,active:active})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}


export const updateSysEmail = async (email:string,active:boolean) => {
  const response = await axios.put(`/api/sys-email`,{email:email,active:active})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}

export const getSysConf = async () => {
  const response = await axios.get('/api/config')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.data
}

export const updateSysPort = async (port:string) => {
  const response = await axios.put(`/api/sys-email`,{listening_port:port})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}


export const updateSysRemoteLogIp = async (ip:string) => {
  const response = await axios.put(`/api/sys-email`,{remote_logServer:ip})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.data
}
