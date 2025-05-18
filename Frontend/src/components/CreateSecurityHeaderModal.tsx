import {useState} from 'react'
import {useCreateSecurityHeader} from '../hooks/api/useSecurityHeaders'
import {useGetApplications} from '../hooks/api/useApplication'

function CreateSecurityHeaderModal({
  isOpen,
  onClose,
}: {
  isOpen: boolean
  onClose: () => void
}) {
  const [headerName, setHeaderName] = useState('')
  const [headerValue, setHeaderValue] = useState('')
  const [applicationIds, setApplicationIds] = useState<string[]>([])

  const {data: applications = [], isLoading} = useGetApplications()
  const {mutate: createHeader} = useCreateSecurityHeader()

  const handleSubmit = () => {
    createHeader(
      {
        header_name: headerName,
        header_value: headerValue,
        application_id: applicationIds,
      },
      {
        onSuccess: () => {
          onClose()
          setHeaderName('')
          setHeaderValue('')
          setApplicationIds([])
          alert('Header created successfully!')
        },
      },
    )
  }

  const handleSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selected = Array.from(e.target.selectedOptions, option => option.value)
    setApplicationIds(selected)
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md space-y-4">
        <h2 className="text-lg font-semibold">Create Security Header</h2>

        <input
          type="text"
          placeholder="Header Name"
          className="w-full border p-2 rounded"
          value={headerName}
          onChange={e => setHeaderName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Header Value"
          className="w-full border p-2 rounded"
          value={headerValue}
          onChange={e => setHeaderValue(e.target.value)}
        />

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Assign to Applications:
          </label>
          {isLoading ? (
            <div className="text-gray-500 text-sm">Loading applications...</div>
          ) : (
            <select
              value={applicationIds}
              onChange={handleSelectChange}
              className="w-full border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 rounded-lg p-2 text-sm bg-white shadow-sm"
            >
              {applications.map((app: any) => (
                <option key={app.application_id} value={app.application_id}>
                  {app.hostname}
                </option>
              ))}
            </select>
          )}
        </div>

        <div className="flex justify-end gap-2">
          <button onClick={onClose} className="px-4 py-2 bg-gray-300 rounded">
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="px-4 py-2 bg-blue-500 text-white rounded"
            style={{backgroundColor: '#1F263E'}}
          >
            Create
          </button>
        </div>
      </div>
    </div>
  )
}

export default CreateSecurityHeaderModal
