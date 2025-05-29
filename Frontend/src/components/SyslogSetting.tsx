import {useEffect, useState} from 'react'
import {useGetSysConf} from '../hooks/api/useSystemConf'
import {useUpdateSysPort, useUpdateSysRemoteLogIp} from '../hooks/api/useSystemConf'
import {useQueryClient} from '@tanstack/react-query'
import { useToast } from '../hooks/useToast'

export default function SyslogSettings() {
  const [serverAddress, setServerAddress] = useState('')
  const [serverPort, setServerPort] = useState('')
  const [isSaving, setIsSaving] = useState(false)

  const queryClient = useQueryClient()
  const {data, isLoading, refetch} = useGetSysConf()

  const remoteLogIpMutation = useUpdateSysRemoteLogIp()
  const portMutation = useUpdateSysPort()

  const {addToast: toast} = useToast()


  useEffect(() => {
    refetch()
  })

  useEffect(() => {
    if (data) {
      if (!serverAddress) setServerAddress(data.remote_logServer || '')
      if (!serverPort) setServerPort(data.listening_port || '')
    }
  }, [data])

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()

    setIsSaving(true)

    const promises: Promise<any>[] = []

    if (serverAddress) {
      promises.push(remoteLogIpMutation.mutateAsync(serverAddress))
    }

    if (serverPort) {
      promises.push(portMutation.mutateAsync(serverPort))
    }

    try {
      await Promise.all(promises)
      toast('Syslog configuration saved successfully!')
      queryClient.invalidateQueries({queryKey: ['getSysConf']})
    } catch (error: any) {
      toast('Failed to save configuration. ' + error.message)
      console.error(error)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <form onSubmit={handleSave} className="w-full p-6 bg-white shadow-md xl">
      <h2 className="text-xl font-semibold text-gray-800 mb-6 flex items-center gap-2">
        <span>🖥️ Syslog</span>
        <span className="tooltip text-sm text-gray-500">ℹ️</span>
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Syslog server address
          </label>
          <input
            type="text"
            value={serverAddress}
            onChange={e => setServerAddress(e.target.value)}
            placeholder="192.168.10.10"
            className="w-full p-3 border border-gray-300 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={isLoading}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Syslog server port
          </label>
          <input
            type="text"
            value={serverPort}
            onChange={e => setServerPort(e.target.value)}
            placeholder="Must be in range 1 ~ 65535"
            className="w-full p-3 border border-gray-300 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={isLoading}
          />
        </div>
      </div>

      <div className="mt-6 flex justify-end">
        <button
          type="submit"
          disabled={isSaving}
          className="bg-black text-white font-semibold px-6 py-2 hover:bg-blue-700 transition disabled:opacity-50"
        >
          {isSaving ? 'Saving...' : 'Save'}
        </button>
      </div>
    </form>
  )
}
