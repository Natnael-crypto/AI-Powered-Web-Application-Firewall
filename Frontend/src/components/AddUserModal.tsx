import React, {useState, useCallback} from 'react'
import Modal from './Modal'

type UserStatus = 'active' | 'inactive'
type UserRole = 'admin' | 'editor' | 'viewer' // Expanded role options

interface AdminUser {
  user_id: string
  username: string
  password_hash: string
  role: UserRole
  status: UserStatus
  profile_image_url: string
  notifications_enabled: boolean
}

interface AddUserModalProps {
  onSubmit: (data: AdminUser) => void
  isOpen: boolean
  onClose: () => void
  initialData?: Partial<AdminUser> // Allow prefilling form
}

const defaultUserData: AdminUser = {
  user_id: '',
  username: '',
  password_hash: '',
  role: 'admin',
  status: 'active',
  profile_image_url: '',
  notifications_enabled: false,
}

const AddUserModal: React.FC<AddUserModalProps> = ({
  onSubmit,
  isOpen,
  onClose,
  initialData = {},
}) => {
  const [formData, setFormData] = useState<AdminUser>({
    ...defaultUserData,
    ...initialData,
  })
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})

  const validateForm = useCallback((): boolean => {
    const newErrors: Record<string, string> = {}

    if (!formData.username.trim()) {
      newErrors.username = 'Username is required'
    }

    if (password && password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters'
    }

    if (password !== confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }, [formData.username, password, confirmPassword])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const {name, value, type} = e.target as HTMLInputElement

    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    }))

    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({...prev, [name]: ''}))
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) return

    const submitData = {...formData}
    if (password) {
      // In a real app, you would hash the password here before sending
      submitData.password_hash = password
    }

    onSubmit(submitData)
    onClose()
  }

  const handleReset = () => {
    setFormData({...defaultUserData, ...initialData})
    setPassword('')
    setConfirmPassword('')
    setErrors({})
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Admin User Settings">
      <form onSubmit={handleSubmit} className="max-w-2xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          {/* Username */}
          <div className="col-span-1">
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Username *
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              className={`w-full px-3 py-2 border ${
                errors.username ? 'border-red-500' : 'border-gray-300'
              } rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500`}
              aria-invalid={!!errors.username}
              aria-describedby={errors.username ? 'username-error' : undefined}
              required
            />
            {errors.username && (
              <p id="username-error" className="mt-1 text-sm text-red-600">
                {errors.username}
              </p>
            )}
          </div>

          {/* Role */}
          <div className="col-span-1">
            <label
              htmlFor="role"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Role
            </label>
            <select
              id="role"
              name="role"
              value={formData.role}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="admin">Admin</option>
              <option value="editor">Editor</option>
              <option value="viewer">Viewer</option>
            </select>
          </div>

          {/* Status */}
          <div className="col-span-1">
            <label
              htmlFor="status"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Status
            </label>
            <select
              id="status"
              name="status"
              value={formData.status}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          {/* Profile Image URL */}
          <div className="col-span-1 md:col-span-2">
            <label
              htmlFor="profile_image_url"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Profile Image URL
            </label>
            <input
              type="url"
              id="profile_image_url"
              name="profile_image_url"
              value={formData.profile_image_url}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="https://example.com/image.jpg"
            />
          </div>
        </div>

        {/* Password Change Section */}
        <fieldset className="border border-gray-200 rounded-lg p-4 mb-6">
          <legend className="px-2 text-lg font-medium text-gray-800">
            Change Password
          </legend>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                New Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className={`w-full px-3 py-2 border ${
                  errors.password ? 'border-red-500' : 'border-gray-300'
                } rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500`}
                aria-invalid={!!errors.password}
                aria-describedby={errors.password ? 'password-error' : undefined}
                minLength={8}
              />
              {errors.password && (
                <p id="password-error" className="mt-1 text-sm text-red-600">
                  {errors.password}
                </p>
              )}
              <p className="mt-1 text-xs text-gray-500">Minimum 8 characters</p>
            </div>

            <div>
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                Confirm Password
              </label>
              <input
                type="password"
                id="confirmPassword"
                name="confirmPassword"
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                className={`w-full px-3 py-2 border ${
                  errors.confirmPassword ? 'border-red-500' : 'border-gray-300'
                } rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500`}
                aria-invalid={!!errors.confirmPassword}
                aria-describedby={
                  errors.confirmPassword ? 'confirmPassword-error' : undefined
                }
              />
              {errors.confirmPassword && (
                <p id="confirmPassword-error" className="mt-1 text-sm text-red-600">
                  {errors.confirmPassword}
                </p>
              )}
            </div>
          </div>
        </fieldset>

        {/* Notifications */}
        <div className="flex items-center mb-6">
          <input
            type="checkbox"
            id="notifications_enabled"
            name="notifications_enabled"
            checked={formData.notifications_enabled}
            onChange={handleChange}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
          <label
            htmlFor="notifications_enabled"
            className="ml-2 block text-sm text-gray-700"
          >
            Enable Notifications
          </label>
        </div>

        <div className="flex justify-end space-x-3">
          <button
            type="button"
            onClick={handleReset}
            className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Reset
          </button>
          <button
            type="submit"
            className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={
              !formData.username.trim() || (!!password && password !== confirmPassword)
            }
          >
            Save Changes
          </button>
        </div>
      </form>
    </Modal>
  )
}

export default AddUserModal
