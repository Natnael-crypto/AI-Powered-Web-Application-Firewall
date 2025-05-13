import React from 'react'
import { PencilIcon, TrashIcon } from 'lucide-react'
import { useDeleteSecurityHeader } from '../hooks/api/useSecurityHeaders'
import EditSecurityHeaderModal from './EditSecurityHeaderModal'

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
            <th className="px-4 py-2">Name</th>
            <th className="px-4 py-2">Value</th>
            <th className="px-4 py-2">Description</th>
            <th className="px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {securityHeaders?.map((header) => (
            <tr key={header.id} className="border-b">
              <td className="px-4 py-2">{header.name}</td>
              <td className="px-4 py-2">{header.value}</td>
              <td className="px-4 py-2">{header.description}</td>
              <td className="px-4 py-2 flex gap-2">
                <button
                  onClick={() => setEditingHeader(header)}
                  className="text-blue-600 hover:underline"
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
