import  { useState } from 'react'
import { useUpdateSecurityHeader } from '../hooks/api/useSecurityHeaders'
import { useToast } from '../hooks/useToast'

function EditSecurityHeaderModal({
  header,
  onClose,
}: {
  header: any
  onClose: () => void
}) {
  const [headerName, setHeaderName] = useState(header.header_name)
  const [headerValue, setHeaderValue] = useState(header.header_value)
  const {addToast: toast} = useToast()

  const { mutate: updateHeader } = useUpdateSecurityHeader()

  const handleUpdate = () => {

    var data={
        "header_name":headerName,
        "header_value":headerValue
    }
    updateHeader(
      {
          headerId: header.id,
          data: data
      },
      {
        onSuccess: () => {
          onClose()
          toast('Header Updated successfully!');
        },
      }
    )
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md space-y-4">
        <h2 className="text-lg font-semibold">Edit Security Header</h2>

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

        <div className="flex justify-end gap-2">
          <button onClick={onClose} className="px-4 py-2 bg-gray-300 rounded">
            Cancel
          </button>
          <button
            onClick={handleUpdate}
            className="px-4 py-2 bg-green-500 text-white rounded"
            style={{backgroundColor: '#1F263E'}}
          >
            Update
          </button>
        </div>
      </div>
    </div>
  )
}

export default EditSecurityHeaderModal
