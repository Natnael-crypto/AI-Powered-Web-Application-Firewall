import {useState} from 'react'
import Modal from './Modal'
import Table from './Table'
import {ColumnDef} from '@tanstack/react-table'
import {DropdownActions} from './DropdownAction'
import {
  useActivateUser,
  useDeactivateUser,
  useDeleteUser,
  useGetUsers,
} from '../hooks/api/useUser'
import {AdminUser} from '../lib/types'

const AdminTable = () => {
  const [selectedAdmin, setSelectedAdmin] = useState<AdminUser | null>(null)
  const [isAssignModalOpen, setIsAssignModalOpen] = useState(false)
  const {data, isLoading, isError} = useGetUsers()
  const {mutate} = useDeleteUser()
  const {mutate: deactivateUser} = useDeactivateUser()
  const {mutate: activateUser} = useActivateUser()

  const handleStatusChange = (
    admin: AdminUser,
    newStatus: 'active' | 'inactive' | 'suspended',
  ) => {
    console.log('header')
    const fn = newStatus === 'inactive' ? deactivateUser : activateUser

    fn(admin.username, {
      onSuccess: () => console.log(`made user with ${admin.username} ${newStatus}`),
    })
    console.log(`Changing status for ${admin.username} to ${newStatus}`)
  }

  const handleUpdatePassword = (admin: AdminUser) => {
    console.log(`Updating password for ${admin.username}`)
  }

  const handleDeleteAdmin = async (admin: AdminUser) => {
    mutate(admin.username, {
      onSuccess: () => console.log(`deleteUser with ${admin.username} username`),
    })
  }

  const handleAssign = (admin: AdminUser) => {
    setSelectedAdmin(admin)
    setIsAssignModalOpen(true)
  }

  const columns: ColumnDef<AdminUser>[] = [
    {
      accessorKey: 'username',
      header: 'Username',
    },
    {
      accessorKey: 'role',
      header: 'Role',
      cell: info => <span className="capitalize">{String(info.getValue())}</span>,
    },
    {
      accessorKey: 'status',
      header: 'Status',
      cell: info => {
        const value = info.getValue() as AdminUser['status']
        const colorClass =
          value === 'active'
            ? 'bg-green-100 text-green-800'
            : value === 'inactive'
              ? 'bg-red-100 text-red-800'
              : 'bg-yellow-100 text-yellow-800'

        return (
          <span className={`px-2 py-1 rounded-full text-xs font-medium ${colorClass}`}>
            {value.charAt(0).toUpperCase() + value.slice(1)}
          </span>
        )
      },
    },
    {
      accessorKey: 'last_login',
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
              label: 'Set Active',
              onClick: user => handleStatusChange(user, 'active'),
              show: (user: AdminUser) => user.status !== 'active',
            },
            {
              label: 'Set Inactive',
              onClick: user => handleStatusChange(user, 'inactive'),
              show: (user: AdminUser) => user.status !== 'inactive',
            },
            {
              label: 'Update Password',
              onClick: handleUpdatePassword,
              show: () => true,
            },
            {
              label: 'Delete',
              onClick: handleDeleteAdmin,
              show: () => true,
            },
            {
              label: 'Assign',
              onClick: handleAssign,
              show: () => true,
            },
          ]}
        />
      ),
    },
  ]

  return (
    <div className="space-y-4">
      {isLoading && <p>Loading admins...</p>}
      {isError && <p className="text-red-500">Failed to load admin data.</p>}
      {!isLoading && !isError && <Table columns={columns} data={data?.admins || []} />}

      <Modal
        isOpen={isAssignModalOpen}
        onClose={() => setIsAssignModalOpen(false)}
        title="Assign Admin"
      >
        {selectedAdmin && (
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Detail label="Username" value={selectedAdmin.username} />
              <Detail label="Role" value={selectedAdmin.role} capitalize />
              <Detail label="Status" value={selectedAdmin.status} capitalize />
              <Detail label="User ID" value={selectedAdmin.user_id} mono />
              <Detail
                label="Created At"
                value={new Date(selectedAdmin.created_at).toLocaleString()}
              />
              <Detail
                label="Last Login"
                value={new Date(selectedAdmin.last_login).toLocaleString()}
              />
              <Detail
                label="Profile Image"
                value={selectedAdmin.profile_image_url || 'None'}
              />
            </div>

            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                onClick={() => setIsAssignModalOpen(false)}
              >
                Cancel
              </button>
              <button
                type="button"
                className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                onClick={() => {
                  console.log(`Assigning admin ${selectedAdmin.username}`)
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

const Detail = ({
  label,
  value,
  mono = false,
  capitalize = false,
}: {
  label: string
  value: string
  mono?: boolean
  capitalize?: boolean
}) => (
  <div>
    <p className="text-sm font-medium text-gray-500">{label}</p>
    <p
      className={`mt-1 text-sm text-gray-900 ${
        mono ? 'font-mono' : ''
      } ${capitalize ? 'capitalize' : ''}`}
    >
      {value}
    </p>
  </div>
)

export default AdminTable
