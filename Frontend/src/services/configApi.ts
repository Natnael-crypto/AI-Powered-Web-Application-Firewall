import axios from '../lib/axios'
export const getSysEmail = async () => {
  const response = await axios.get('/api/sys-email/')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}

export const createSysEmail = async (email:string,active:boolean) => {
  const response = await axios.post(`/api/sys-email/`,{email:email,active:active})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}


export const updateSysEmail = async (email:string,active:boolean) => {
  const response = await axios.put(`/api/sys-email/`,{email:email,active:active})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}

export const getSysConf = async () => {
  const response = await axios.get('/api/config/')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.data
}

export const updateSysPort = async (port:string) => {
  const response = await axios.put(`/api/config/update/listening-port`,{listening_port:port})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.email
}


export const updateSysRemoteLogIp = async (ip:string) => {
  const response = await axios.put(`/api/config/update/remote-log-server`,{remote_logServer:ip})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.data
}


export const getAllowedIp = async () => {
  const response = await axios.get('/api/service')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data
}

export const createAllowedIp = async (service:string,ip:string) => {
  const response = await axios.post(`/api/service`,{service:service,ip:ip})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data
}


export const updateAllowedIp = async (id:string,service:string,ip:string) => {
  const response = await axios.put(`/api/service/${id}`,{service:service,ip:ip})

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data
}


export const deleteAllowedIp = async (id:string) => {
  const response = await axios.delete(`/api/service/${id}`)

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data
}