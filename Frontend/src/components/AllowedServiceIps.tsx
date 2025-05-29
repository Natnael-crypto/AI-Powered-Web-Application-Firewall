import {useState} from 'react'
import {
  useAllowedIp,
  useCreateAllowedIp,
  useUpdateAllowedIp,
  useDeleteAllowedIp,
} from '../hooks/api/useAllowedIp'
import {useQueryClient} from '@tanstack/react-query'
import {Pencil, Trash2} from 'lucide-react'
import { useToast } from '../hooks/useToast'

export default function AllowedServiceIps() {
  const [ip, setIp] = useState('')
  const [service, setService] = useState('M')
  const [editingId, setEditingId] = useState<string | null>(null)

  const queryClient = useQueryClient()

  const {data, isLoading} = useAllowedIp()
  const createIpMutation = useCreateAllowedIp()
  const updateIpMutation = useUpdateAllowedIp()
  const deleteIpMutation = useDeleteAllowedIp()
  const {addToast: toast} = useToast()
  

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    try {
      if (editingId) {
        await updateIpMutation.mutateAsync({id: editingId, ip, service})
        setEditingId(null)
      } else {
        await createIpMutation.mutateAsync({ip, service})
      }

      setIp('')
      setService('M')
      queryClient.invalidateQueries({queryKey: ['getAllowedIp']})
    } catch (err) {
      toast('Something went wrong while saving.')
    }
  }

  const handleEdit = (entry: {id: string; ip: string; service: string}) => {
    setEditingId(entry.id)
    setIp(entry.ip)
    setService(entry.service)
  }

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this IP?')) {
      await deleteIpMutation.mutateAsync(id)
      queryClient.invalidateQueries({queryKey: ['getConf']})
    }
  }

  return (
    <div className="p-6 bg-white shadow-md rounded-sm">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Allowed Service IPs</h2>

      <form
        onSubmit={handleSubmit}
        className="flex flex-col md:flex-row gap-4 items-center mb-6"
      >
        <input
          type="text"
          placeholder="Enter IP Address"
          value={ip}
          onChange={e => setIp(e.target.value)}
          required
          className="flex-1 p-3 border border-gray-300 rounded-md shadow-sm"
        />
        <select
          value={service}
          onChange={e => setService(e.target.value)}
          className="p-3 border border-gray-300 rounded-md shadow-sm"
        >
          <option value="M">ML Service</option>
          <option value="I">Interceptor Service</option>
        </select>
        <button
          type="submit"
          className="bg-black text-white font-semibold px-6 py-2 rounded-md hover:bg-slate-800 transition"
        >
          {editingId ? 'Update' : 'Add'}
        </button>
      </form>

      <table className="w-full border text-sm text-left">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2">IP Address</th>
            <th className="p-2">Service</th>
            <th className="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {!isLoading && data?.length > 0 ? (
            data.map((entry: {id: string; ip: string; service: string}) => (
              <tr key={entry.id} className="border-t">
                <td className="p-2">{entry.ip}</td>
                <td className="p-2">
                  {entry.service === 'I' ? 'Interceptor Service' : 'ML Service'}
                </td>
                <td className="p-2 flex gap-3">
                  <button onClick={() => handleEdit(entry)} className="text-blue-600">
                    <Pencil size={16} className="text-blue-600" />
                  </button>
                  <button onClick={() => handleDelete(entry.id)} className="text-red-600">
                    <Trash2 size={16} className="text-red-600" />
                  </button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td className="p-3 border text-center" colSpan={3}>
                {isLoading ? 'Loading...' : 'No allowed IPs found.'}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}
