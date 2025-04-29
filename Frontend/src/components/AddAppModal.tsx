import React, {useState} from 'react'
import {Info, Tag, Save, X} from 'lucide-react'
import {useAddApplication} from '../hooks/useApplication'
import {QueryClient} from '@tanstack/react-query'

interface AddAppModalProps {
  isModalOpen: boolean
  toggleModal: () => void
}

interface AddAppPayload {
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
}

export default function AddAppModal({isModalOpen, toggleModal}: AddAppModalProps) {
  const [payload, setPayload] = useState<AddAppPayload>({
    application_name: '',
    description: '',
    hostname: '',
    ip_address: '',
    port: '',
    status: true,
    tls: false,
  })

  const queryClient = new QueryClient()

  const {mutate} = useAddApplication()

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const {name, value, type} = e.target as HTMLInputElement
    const checked =
      type === 'checkbox' ? (e.target as HTMLInputElement).checked : undefined
    setPayload({
      ...payload,
      [name]: type === 'checkbox' ? checked : value,
    })
  }

  const [_, setErrors] = useState<Record<string, string>>({})

  const validate = () => {
    const newErrors: Record<string, string> = {}

    if (!payload.application_name.trim())
      newErrors.application_name = 'Application name is required'
    if (!payload.description.trim()) newErrors.description = 'Description is required'
    if (!payload.hostname.trim()) newErrors.hostname = 'Hostname is required'

    // Validate IP Address (Simple IPv4 Check)
    const ipRegex =
      /^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])(\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])){3}$/
    if (!ipRegex.test(payload.ip_address)) newErrors.ip_address = 'Invalid IP address'

    // Validate Port (1 - 65535)
    const portNumber = Number(payload.port)
    if (isNaN(portNumber) || portNumber < 1 || portNumber > 65535) {
      newErrors.port = 'Port must be between 1 and 65535'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    mutate(payload, {
      onSuccess: () => {
        alert('Application added successfully!')
        queryClient.invalidateQueries({queryKey: ['applications']})
        toggleModal()
      },
      onError: error => {
        alert('Failed to add application: ' + error.message)
      },
    })
  }

  if (!isModalOpen) return null

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-end justify-center px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div
          className="fixed inset-0 bg-gray-500/75 transition-opacity"
          onClick={toggleModal}
        />

        <div className="relative inline-block transform overflow-hidden rounded-lg bg-white text-left align-bottom shadow-xl transition-all sm:my-8 sm:w-full sm:max-w-2xl sm:align-middle">
          <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-6 py-4">
            <div className="flex items-center justify-between">
              <h3 className="text-xl font-medium text-white">Add New Application</h3>
              <button onClick={toggleModal} className="text-white hover:text-gray-200">
                <X className="h-6 w-6" />
              </button>
            </div>
          </div>

          <form onSubmit={handleSubmit} className="bg-white px-6 py-4">
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Application Name
                </label>
                <div className="relative mt-1 rounded-md shadow-sm">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <Info className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    name="application_name"
                    value={payload.application_name}
                    onChange={handleChange}
                    className="block w-full rounded-md border border-gray-300 pl-10 py-2 focus:border-indigo-500 focus:ring-indigo-500"
                    placeholder="Enter application name"
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700">
                  Description
                </label>
                <textarea
                  name="description"
                  value={payload.description}
                  onChange={handleChange}
                  rows={3}
                  className="mt-1 block w-full rounded-md border border-gray-300 py-2 px-3 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                  placeholder="Enter a brief description"
                  required
                />
              </div>

              <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    Hostname
                  </label>
                  <div className="relative mt-1 rounded-md shadow-sm">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <Info className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="text"
                      name="hostname"
                      value={payload.hostname}
                      onChange={handleChange}
                      className="block w-full rounded-md border border-gray-300 pl-10 py-2 focus:border-indigo-500 focus:ring-indigo-500"
                      placeholder="Enter hostname"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    IP Address
                  </label>
                  <div className="relative mt-1 rounded-md shadow-sm">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <Info className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="text"
                      name="ip_address"
                      value={payload.ip_address}
                      onChange={handleChange}
                      className="block w-full rounded-md border border-gray-300 pl-10 py-2 focus:border-indigo-500 focus:ring-indigo-500"
                      placeholder="Enter IP address"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Port</label>
                  <div className="relative mt-1 rounded-md shadow-sm">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <Tag className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="text"
                      name="port"
                      value={payload.port}
                      onChange={handleChange}
                      className="block w-full rounded-md border border-gray-300 pl-10 py-2 focus:border-indigo-500 focus:ring-indigo-500"
                      placeholder="Enter port number"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">
                    Status
                  </label>
                  <select
                    name="status"
                    value={payload.status ? 'true' : 'false'}
                    onChange={handleChange}
                    className="mt-1 block w-full rounded-md border border-gray-300 py-2 pl-3 pr-10 focus:border-indigo-500 focus:ring-indigo-500"
                  >
                    <option value="true">Active</option>
                    <option value="false">Inactive</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">TLS</label>
                  <div className="relative mt-1 rounded-md shadow-sm">
                    <input
                      type="checkbox"
                      name="tls"
                      checked={payload.tls}
                      onChange={handleChange}
                      className="h-4 w-4 border-gray-300 rounded"
                    />
                    <span className="ml-2 text-sm text-gray-700">Enable TLS</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                onClick={toggleModal}
                className="inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="inline-flex items-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              >
                <Save className="mr-2 h-4 w-4" />
                Save Application
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
