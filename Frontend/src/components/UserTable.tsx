import {useState} from 'react'
import Table from './Table'
import {ColumnDef} from '@tanstack/react-table'
import {DropdownActions} from './DropdownAction'
import {
  useActivateUser,
  useDeactivateUser,
  useDeleteUser,
  useGetUsers,
  useUpdateUser,
} from '../hooks/api/useUser'
import {AdminUser} from '../lib/types'
import {AssignAdminModal} from './AssignAdminModal'
import {CheckCircle, PauseCircle, Key, Trash2, UserPlus} from 'lucide-react'
import UpdatePasswordModal from './UpdatePasswordModal'
import {useToast} from '../hooks/useToast'

const AdminTable = () => {
  const [selectedAdmin, setSelectedAdmin] = useState<AdminUser | null>(null)
  const [isAssignModalOpen, setIsAssignModalOpen] = useState(false)
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false)
  const [passwordAdmin, setPasswordAdmin] = useState<AdminUser | null>(null)
  const {data, isLoading, isError} = useGetUsers()
  const {addToast} = useToast()
  const {mutate} = useDeleteUser()
  const {mutate: deactivateUser} = useDeactivateUser()
  const {mutate: activateUser} = useActivateUser()
  const {mutate: updateUser} = useUpdateUser()

  const handleStatusChange = (
    admin: AdminUser,
    newStatus: 'active' | 'inactive' | 'suspended',
  ) => {
    const fn = newStatus === 'inactive' ? deactivateUser : activateUser

    fn(admin.username, {
      onSuccess: () => console.log(`made user with ${admin.username} ${newStatus}`),
    })
    console.log(`Changing status for ${admin.username} to ${newStatus}`)
  }

  const handleDeleteAdmin = async (admin: AdminUser) => {
    mutate(admin.username, {
      onSuccess: () => console.log(`deleteUser with ${admin.username} username`),
    })
  }

  const handleUpdatePassword = (admin: AdminUser) => {
    setPasswordAdmin(admin)
    setIsPasswordModalOpen(true)
  }

  const handleAssign = (admin: AdminUser) => {
    setSelectedAdmin(admin)
    setIsAssignModalOpen(true)
  }

  const handlePasswordConfirm = (data: {
    username: string
    old_password: string
    new_password: string
  }) => {
    updateUser(data, {
      onSuccess: () => addToast('updated Successfully'),
    })
    setIsPasswordModalOpen(false)
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
              icon: <CheckCircle className="w-4 h-4 mr-2" />,
              onClick: user => handleStatusChange(user, 'active'),
              show: (user: AdminUser) => user.status !== 'active',
            },
            {
              label: 'Set Inactive',
              icon: <PauseCircle className="w-4 h-4 mr-2" />,
              onClick: user => handleStatusChange(user, 'inactive'),
              show: (user: AdminUser) => user.status !== 'inactive',
            },
            {
              label: 'Update Password',
              icon: <Key className="w-4 h-4 mr-2" />,
              onClick: handleUpdatePassword,
              show: () => true,
            },
            {
              label: 'Delete',
              icon: <Trash2 className="w-4 h-4 mr-2" />,
              onClick: handleDeleteAdmin,
              show: () => true,
            },
            {
              label: 'Assign Service',
              icon: <UserPlus className="w-4 h-4 mr-2" />,
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

      <AssignAdminModal
        isOpen={isAssignModalOpen}
        admin={selectedAdmin}
        onClose={() => setIsAssignModalOpen(false)}
        onConfirm={admin => {
          setIsAssignModalOpen(false)
        }}
      />
      <UpdatePasswordModal
        isOpen={isPasswordModalOpen}
        admin={passwordAdmin}
        onClose={() => setIsPasswordModalOpen(false)}
        onConfirm={handlePasswordConfirm}
      />
    </div>
  )
}

export default AdminTable
