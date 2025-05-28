import {useState} from 'react'
import {useCreateSecurityHeader} from '../hooks/api/useSecurityHeaders'
import {useGetApplications} from '../hooks/api/useApplication'
import { useToast } from '../hooks/useToast'

function CreateSecurityHeaderModal({
  isOpen,
  onClose,
}: {
  isOpen: boolean
  onClose: () => void
}) {
  const [headerName, setHeaderName] = useState('')
  const [headerValue, setHeaderValue] = useState('')
  const [selectedAppId, setSelectedAppId] = useState<string>('')
  const {addToast: toast} = useToast()

  const {data: applications = [], isLoading} = useGetApplications()
  const {mutate: createHeader} = useCreateSecurityHeader()

  const handleSubmit = () => {
    if (!selectedAppId) {
      toast('Please select an application')
      return
    }

    createHeader(
      {
        header_name: headerName,
        header_value: headerValue,
        application_id: [selectedAppId],
      },
      {
        onSuccess: () => {
          onClose()
          setHeaderName('')
          setHeaderValue('')
          setSelectedAppId('')
          toast('Header created successfully!')
        },
      },
    )
  }

  const handleSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedAppId(e.target.value)
  }

  const clearSelection = () => {
    setSelectedAppId('')
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
            Assign to Application:
          </label>
          {isLoading ? (
            <div className="text-gray-500 text-sm">Loading applications...</div>
          ) : (
            <>
              <select
                value={selectedAppId}
                onChange={handleSelectChange}
                className="w-full border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 rounded-lg p-2 text-sm bg-white shadow-sm"
              >
                <option value="">Select an application</option>
                {applications.map((app: any) => (
                  <option key={app.application_id} value={app.application_id}>
                    {app.hostname}
                  </option>
                ))}
              </select>

              {selectedAppId && (
                <div className="mt-2">
                  <p className="text-sm font-medium text-gray-700 mb-1">
                    Selected Application:
                  </p>
                  <div className="flex items-center justify-between bg-gray-50 p-2 rounded">
                    <span>
                      {applications.find((a: any) => a.application_id === selectedAppId)
                        ?.hostname || 'Unknown App'}
                    </span>
                    <button
                      type="button"
                      onClick={clearSelection}
                      className="text-red-500 hover:text-red-700"
                    >
                      ×
                    </button>
                  </div>
                </div>
              )}
            </>
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
