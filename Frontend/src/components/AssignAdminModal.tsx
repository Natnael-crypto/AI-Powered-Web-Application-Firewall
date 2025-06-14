import {
  useAssignApplication,
  useDeleteAssignment,
  useGetApplicationAssignments,
  useGetApplications,
} from '../hooks/api/useApplication'
import {AdminUser, Application, Assignment} from '../lib/types'
import Modal from './Modal'
import {Plus, X} from 'lucide-react'
import {useEffect} from 'react'

interface DetailProps {
  label: string
  value: string
  mono?: boolean
  capitalize?: boolean
}

const Detail = ({label, value, mono = false, capitalize = false}: DetailProps) => (
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

interface AssignAdminModalProps {
  isOpen: boolean
  admin: AdminUser | null
  onClose: () => void
  onConfirm: (admin: AdminUser) => void
}

export const AssignAdminModal = ({
  isOpen,
  admin,
  onClose,
  onConfirm,
}: AssignAdminModalProps) => {
  const {data: assignedData, refetch: refetchAssignments} = useGetApplicationAssignments()
  const {data: allApplications, refetch: refetchApplications} = useGetApplications()

  const {mutate: assignApplication} = useAssignApplication({
    onSuccess: () => {
      refetchAssignments()
      refetchApplications()
    },
  })

  const {mutate: deleteAssignment} = useDeleteAssignment({
    onSuccess: () => {
      refetchAssignments()
      refetchApplications()
    },
  })

  useEffect(() => {
    if (isOpen) {
      refetchAssignments()
      refetchApplications()
    }
  }, [isOpen, refetchAssignments, refetchApplications])

  if (!admin) return null

  const assignments = assignedData?.assignments ?? []
  const applications = allApplications ?? []

  const assignedAppsToThisAdmin = assignments.filter(a => a.user_id === admin.user_id)
  const unassignedApps = applications.filter(
    app =>
      !assignedAppsToThisAdmin
        .map(a => a.application_name)
        .includes(app.application_name),
  )
  const handleAssignment = (application: Application) => {
    if (!admin) return
    assignApplication({
      user_id: admin.user_id,
      application_name: application.application_name,
    })
  }

  const handleDeleteAssignment = (assignment: Assignment) => {
    if (!admin) return
    deleteAssignment(assignment.id)
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Assign Admin">
      <div className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          <Detail label="Username" value={admin.username} />
          <Detail label="Role" value={admin.role} capitalize />
          <Detail label="Status" value={admin.status} capitalize />
          <Detail label="User ID" value={admin.user_id} mono />
          <Detail
            label="Created At"
            value={new Date(admin.created_at).toLocaleString()}
          />
          <Detail
            label="Last Login"
            value={new Date(admin.last_login).toLocaleString()}
          />
          <Detail label="Profile Image" value={admin.profile_image_url || 'None'} />
        </div>

        <div>
          <h3 className="text-md font-semibold mb-2">Assigned Applications</h3>
          <ul className="space-y-2">
            {assignedAppsToThisAdmin.length === 0 ? (
              <p className="text-sm text-gray-500">No applications assigned.</p>
            ) : (
              assignedAppsToThisAdmin.map(app => (
                <li
                  key={app.id}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <span>{app.application_name}</span>
                  <button
                    type="button"
                    className="text-red-600 hover:text-red-800"
                    onClick={() => handleDeleteAssignment(app)}
                  >
                    <X size={16} />
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>

        <div>
          <h3 className="text-md font-semibold mb-2">Unassigned Applications</h3>
          <ul className="space-y-2">
            {unassignedApps.length === 0 ? (
              <p className="text-sm text-gray-500">No unassigned applications.</p>
            ) : (
              unassignedApps.map(app => (
                <li
                  key={app.application_name}
                  className="flex items-center justify-between p-2 border rounded"
                >
                  <span>{app.application_name}</span>
                  <button
                    type="button"
                    className="text-green-600 hover:text-green-800"
                    onClick={() => handleAssignment(app)}
                  >
                    <Plus size={16} />
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>

        <div className="mt-6 flex justify-end space-x-3">
          <button
            type="button"
            className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            onClick={onClose}
          >
            Cancel
          </button>
          <button
            type="button"
            className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
            onClick={() => onConfirm(admin)}
          >
            Confirm Assignment
          </button>
        </div>
      </div>
    </Modal>
  )
}
