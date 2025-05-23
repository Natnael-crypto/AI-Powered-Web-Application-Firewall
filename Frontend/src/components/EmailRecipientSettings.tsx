import { useState } from 'react'
import { Pencil, Trash2 } from 'lucide-react'
import {
  useGetUseEmail,
  useAddUserEmail,
  useUpdateUserEmail,
  useDeleteUserEmail,
} from '../hooks/api/useSystemEmail'
import { useGetAllUsers } from '../hooks/api/useUser'

const EmailRecipientSettings = () => {
  const [selectedUserId, setSelectedUserId] = useState('')
  const [email, setEmail] = useState('')
  const [editingId, setEditingId] = useState<string | null>(null)

  const { data: emails, refetch } = useGetUseEmail()
  const { data: users } = useGetAllUsers()

  const addEmailMutation = useAddUserEmail()
  const updateEmailMutation = useUpdateUserEmail()
  const deleteEmailMutation = useDeleteUserEmail()

  const handleSave = () => {
  


    if (editingId) {
      if (!email || !selectedUserId) {
          alert('Please select a user and enter an email.')
          return
      }
      const payload = { email, id: selectedUserId }

      updateEmailMutation.mutate(payload, {
        onSuccess: () => {
          alert('Email updated.')
          resetForm()
        },
        onError: () => alert('Failed to update email'),
      })
    } else {
        if (!email || !selectedUserId) {
          alert('Please select a user and enter an email.')
          return
        }
      const payload = { email, id: selectedUserId }
      addEmailMutation.mutate(payload, {
        onSuccess: () => {
          alert('Email added.')
          resetForm()
        },
        onError: () => alert('Failed to add email'),
      })
    }
  }

  const handleEdit = (emailItem: any) => {
    setEmail(emailItem.email)
    setSelectedUserId(emailItem.user_id)
    setEditingId(emailItem.user_id)
  }

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this email?')) {
      deleteEmailMutation.mutate({ id }, {
        onSuccess: () => {
          alert('Deleted successfully')
          refetch()
        },
        onError: () => alert('Delete failed'),
      })
    }
  }

  const resetForm = () => {
    setEditingId(null)
    setEmail('')
    setSelectedUserId('')
    refetch()
  }

  return (
    <div className="p-6 bg-white shadow-lg w-full">
      <div className="mb-4">
        <label className="block mb-1 font-medium text-gray-700">Select Admin User</label>
        {users ? (
          <select
            value={selectedUserId}
            onChange={e => setSelectedUserId(e.target.value)}
            className="w-full border px-3 py-2 rounded"
          >
            <option value="">-- Select User --</option>
            {users.map((user: any) => (
              <option key={user.user_id} value={user.user_id}>
                {user.username}
              </option>
            ))}
          </select>
        ) : (
          <div className="text-gray-500 text-sm">Loading users...</div>
        )}
      </div>

      <div className="mb-4">
        <label className="block mb-1 font-medium text-gray-700">Email Address</label>
        <input
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          placeholder="waf-alert-server-noreply@org.com"
          className="w-full px-4 py-2 border border-gray-300 text-sm placeholder-gray-400 rounded"
        />
      </div>

      <div className="mb-6 flex justify-end">
        <button
          onClick={handleSave}
          className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 transition"
        >
          {editingId ? 'Save Email' : 'Add Email'}
        </button>
      </div>

      <h5 className="text-md font-semibold text-gray-800 mb-2">Configured Emails</h5>
      <table className="w-full border text-sm text-left">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2">User</th>
            <th className="p-2">Email</th>
            <th className="p-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {emails?.length > 0 ? (
            emails.map((entry: any) => (
              <tr key={entry.id} className="border-t">
                <td className="p-2">
                  {
                    users? (
                    users?.find((u: any) => u.user_id === entry.user_id)?.username
                    ):
                    <div className="text-gray-500 text-sm">Loading user...</div>
                  }
                </td>

                <td className="p-2">{entry.email}</td>
                <td className="p-2 flex gap-3">
                  <button onClick={() => handleEdit(entry)}>
                    <Pencil size={16} className="text-blue-600" />
                  </button>
                  <button onClick={() => handleDelete(entry.user_id)}>
                    <Trash2 size={16} className="text-red-600" />
                  </button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={4} className="p-4 text-center text-gray-500">
                No email configurations found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}

export default EmailRecipientSettings
