import {useState} from 'react'
import Modal from './Modal'
import Table from './Table'
import {ColumnDef} from '@tanstack/react-table'
import {DropdownActions} from './DropdownAction'
type UserAccount = {
  name: string
  role: string
  '2FA': boolean
  lastLogin: string
  status: 'active' | 'inactive' // Added status field
  email?: string // Added optional fields for user info
  phone?: string
  department?: string
}

const mockUserData: UserAccount[] = [
  {
    name: 'Alice Johnson',
    role: 'Admin',
    '2FA': true,
    lastLogin: '2023-10-01T10:00:00Z',
    status: 'active',
    email: 'alice@example.com',
    phone: '123-456-7890',
    department: 'IT',
  },
  {
    name: 'Bob Smith',
    role: 'User',
    '2FA': false,
    lastLogin: '2023-10-02T12:30:00Z',
    status: 'inactive',
    email: 'bob@example.com',
    phone: '234-567-8901',
    department: 'Marketing',
  },
  {
    name: 'Charlie Brown',
    role: 'Moderator',
    '2FA': true,
    lastLogin: '2023-10-03T14:45:00Z',
    status: 'active',
    email: 'charlie@example.com',
    phone: '345-678-9012',
    department: 'Operations',
  },
]

const UserTable = () => {
  const [selectedUser, setSelectedUser] = useState<UserAccount | null>(null)
  const [isAssignModalOpen, setIsAssignModalOpen] = useState(false)

  const handleStatusChange = (user: UserAccount, newStatus: 'active' | 'inactive') => {
    // In a real app, you would update the user status in your state/API here
    console.log(`Changing status for ${user.name} to ${newStatus}`)
  }

  const handleUpdatePassword = (user: UserAccount) => {
    // Password update logic would go here
    console.log(`Updating password for ${user.name}`)
  }

  const handleDeleteUser = (user: UserAccount) => {
    // Delete user logic would go here
    console.log(`Deleting user ${user.name}`)
  }

  const handleAssign = (user: UserAccount) => {
    setSelectedUser(user)
    setIsAssignModalOpen(true)
  }

  const columns: ColumnDef<UserAccount>[] = [
    {
      accessorKey: 'name',
      header: 'Name',
    },
    {
      accessorKey: 'role',
      header: 'Role',
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: info => (
        <span
          className={`px-2 py-1 rounded-full text-xs font-medium ${
            info.getValue() === 'active'
              ? 'bg-green-100 text-green-800'
              : 'bg-red-100 text-red-800'
          }`}
        >
          {String(info.getValue()).charAt(0).toUpperCase() +
            String(info.getValue()).slice(1)}
        </span>
      ),
    },

    {
      accessorKey: 'lastLogin',
      header: 'Last Login',
      cell: info => new Date(info.getValue() as string).toLocaleString(),
    },
    {
      header: 'Actions',
      id: 'actions',
      cell: ({row}) => (
        <DropdownActions
          item={row.original}
          options={[
            {
              label: row.original.status === 'active' ? 'Set Inactive' : 'Set Active',
              onClick: user =>
                handleStatusChange(
                  user,
                  user.status === 'active' ? 'inactive' : 'active',
                ),
            },
            {
              label: 'Update Password',
              onClick: handleUpdatePassword,
            },
            {
              label: 'Delete',
              onClick: handleDeleteUser,
            },
            {
              label: 'Assign',
              onClick: handleAssign,
            },
          ]}
        />
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <Table columns={columns} data={mockUserData} />

      {/* Assign User Modal */}
      <Modal
        isOpen={isAssignModalOpen}
        onClose={() => setIsAssignModalOpen(false)}
        title="Assign User"
      >
        {selectedUser && (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-gray-500">Name</p>
                <p className="mt-1 text-sm text-gray-900">{selectedUser.name}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Role</p>
                <p className="mt-1 text-sm text-gray-900">{selectedUser.role}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Status</p>
                <p className="mt-1 text-sm text-gray-900 capitalize">
                  {selectedUser.status}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">2FA</p>
                <p className="mt-1 text-sm text-gray-900">
                  {selectedUser['2FA'] ? 'Enabled' : 'Disabled'}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Email</p>
                <p className="mt-1 text-sm text-gray-900">
                  {selectedUser.email || 'N/A'}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Phone</p>
                <p className="mt-1 text-sm text-gray-900">
                  {selectedUser.phone || 'N/A'}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Department</p>
                <p className="mt-1 text-sm text-gray-900">
                  {selectedUser.department || 'N/A'}
                </p>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500">Last Login</p>
                <p className="mt-1 text-sm text-gray-900">
                  {new Date(selectedUser.lastLogin).toLocaleString()}
                </p>
              </div>
            </div>

            {/* Add your assignment form or controls here */}
            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                onClick={() => setIsAssignModalOpen(false)}
              >
                Cancel
              </button>
              <button
                type="button"
                className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                onClick={() => {
                  // Handle assignment logic here
                  console.log(`Assigning user ${selectedUser.name}`)
                  setIsAssignModalOpen(false)
                }}
              >
                Confirm Assignment
              </button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default UserTable
