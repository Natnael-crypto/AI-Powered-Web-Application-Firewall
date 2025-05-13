import React from 'react'
import { PencilIcon, TrashIcon } from 'lucide-react'
import { useDeleteSecurityHeader } from '../hooks/api/useSecurityHeaders'
import EditSecurityHeaderModal from './EditSecurityHeaderModal'
import { useGetApplication } from '../hooks/api/useApplication'
import { useGetUserById } from '../hooks/api/useUser'

function ApplicationHostCell({ applicationId }: { applicationId: string }) {
  const { data, isLoading, isError } = useGetApplication(applicationId)
  

  if (isLoading) return <span>Loading...</span>
  if (isError || !data) return <span>Error</span>

  return <span>{data.hostname}</span>
}

function UserCell({ userId }: { userId: string }) {
  const { data, isLoading, isError } = useGetUserById(userId)

  if (isLoading) return <span>Loading...</span>
  if (isError || !data) return <span>Error</span>

  return <span>{data.username}</span>
}

function SecurityHeaderTable({ securityHeaders }: { securityHeaders: any[] }) {
  const [editingHeader, setEditingHeader] = React.useState(null)
  const { mutate: deleteHeader } = useDeleteSecurityHeader()

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this header?')) {
      deleteHeader(id)
    }
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full text-sm text-left">
        <thead className="bg-slate-100">
          <tr>
            <th className="px-4 py-2">Header Name</th>
            <th className="px-4 py-2">Header Value</th>
            <th className="px-4 py-2">Application</th>
            <th className="px-4 py-2">CreatedBy</th>
            <th className="px-4 py-2">CreatedAt</th>
            <th className="px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {securityHeaders?.map((header) => (
            <tr key={header.id} className="border-b">
              <td className="px-4 py-2">{header.header_name}</td>
              <td className="px-4 py-2">{header.header_value}</td>
              <td className="px-4 py-2">
                <ApplicationHostCell applicationId={header.application_id} />
              </td>
              <td className="px-4 py-2"><UserCell userId={header.created_by} /></td>
              <td className="px-4 py-2">{header.created_at}</td>
              <td className="px-4 py-2 flex gap-2">
                <button
                  onClick={() => setEditingHeader(header)}
                  className="text-blue-600 hover:underline px-3"
                >
                  <PencilIcon size={16} />
                </button>
                <button
                  onClick={() => handleDelete(header.id)}
                  className="text-red-600 hover:underline"
                >
                  <TrashIcon size={16} />
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {editingHeader && (
        <EditSecurityHeaderModal
          header={editingHeader}
          onClose={() => setEditingHeader(null)}
        />
      )}
    </div>
  )
}

export default SecurityHeaderTable
