import { useState } from 'react'
import { useCreateSecurityHeader } from '../hooks/api/useSecurityHeaders'
import { useGetApplications } from '../hooks/api/useApplication'
import { data } from 'react-router-dom'

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

  const { data: applications = [], isLoading } = useGetApplications()
  const { mutate: createHeader } = useCreateSecurityHeader()

  const handleSubmit = () => {

    createHeader(
      {
        header_name:headerName,
      header_value:headerValue,
      application_id:applicationIds
      },
      {
        onSuccess: () => {
          onClose()
          setHeaderName('')
          setHeaderValue('')
          setApplicationIds([])
          alert('Header created successfully!');
        },
      }
    )
  }

  const handleSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selected = Array.from(e.target.selectedOptions, (option) => option.value)
    console.log(selected)
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
          onChange={(e) => setHeaderName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Header Value"
          className="w-full border p-2 rounded"
          value={headerValue}
          onChange={(e) => setHeaderValue(e.target.value)}
        />

        <div>
          <label className="font-medium">Assign to Applications:</label>
          {isLoading ? (
            <div>Loading applications...</div>
          ) : (
            <select
              multiple
              value={applicationIds}
              onChange={handleSelectChange}
              className="w-full border p-2 rounded h-32"
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
          >
            Create
          </button>
        </div>
      </div>
    </div>
  )
}

export default CreateSecurityHeaderModal
