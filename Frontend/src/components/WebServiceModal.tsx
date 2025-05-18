import React, {useEffect, useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'

export interface WebServiceData {
  application_id?: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
}

interface WebServiceModalProps {
  application?: WebServiceData
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: WebServiceData) => void
}

const WebServiceModal: React.FC<WebServiceModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  application,
}) => {
  const [form, setForm] = useState<WebServiceData>({
    application_name: '',
    description: '',
    hostname: '',
    ip_address: '',
    port: '',
    status: true,
    tls: false,
  })

  useEffect(() => {
    setForm({
      application_name: application?.application_name ?? '',
      description: application?.description ?? '',
      hostname: application?.hostname ?? '',
      ip_address: application?.ip_address ?? '',
      port: application?.port ?? '',
      status: application?.status ?? true,
      tls: application?.tls ?? false,
    })
  }, [application])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const {name, value, type, checked} = e.target as HTMLInputElement
    setForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  const handleSubmit = async () => {
    onSubmit(
      application?.application_id
        ? {...form, application_id: application.application_id}
        : form,
    )
  }

  if (!isOpen) return null

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={application ? 'Edit Web Service' : 'Add Web Service'}
    >
      <div className="space-y-6 px-4 py-4">
        {/* Main Form Fields */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">
              Application Name
            </label>
            <input
              type="text"
              name="application_name"
              value={form.application_name}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="Enter application name"
              required
            />
          </div>

          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">Hostname</label>
            <input
              type="text"
              name="hostname"
              value={form.hostname}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="example.com"
              required
            />
          </div>

          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">IP Address</label>
            <input
              type="text"
              name="ip_address"
              value={form.ip_address}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="192.168.1.1"
              required
            />
          </div>

          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">Port</label>
            <input
              type="text"
              name="port"
              value={form.port}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition"
              placeholder="8080"
              required
            />
          </div>
        </div>

        {/* Description Field */}
        <div className="space-y-1">
          <label className="block text-sm font-medium text-gray-700">Description</label>
          <textarea
            name="description"
            value={form.description}
            onChange={handleChange}
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition resize-none"
            placeholder="Enter description"
          />
        </div>

        {/* Toggle Options */}
        <div className="flex flex-wrap gap-6 pt-2">
          <label className="inline-flex items-center">
            <input
              type="checkbox"
              name="status"
              checked={form.status}
              onChange={handleChange}
              className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <span className="ml-2 text-sm text-gray-700">Active</span>
          </label>

          <label className="inline-flex items-center">
            <input
              type="checkbox"
              name="tls"
              checked={form.tls}
              onChange={handleChange}
              className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <span className="ml-2 text-sm text-gray-700">TLS Enabled</span>
          </label>
        </div>

        {/* TLS Certificate Fields (Conditional) */}
        {form.tls && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pt-4">
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Certificate File
              </label>
              <input
                type="file"
                className="block w-full text-sm text-gray-500
                  file:mr-4 file:py-2 file:px-4
                  file:rounded-md file:border-0
                  file:text-sm file:font-semibold
                  file:bg-blue-50 file:text-blue-700
                  hover:file:bg-blue-100"
                required
              />
            </div>
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">Key File</label>
              <input
                type="file"
                className="block w-full text-sm text-gray-500
                  file:mr-4 file:py-2 file:px-4
                  file:rounded-md file:border-0
                  file:text-sm file:font-semibold
                  file:bg-blue-50 file:text-blue-700
                  hover:file:bg-blue-100"
                required
              />
            </div>
          </div>
        )}

        {/* Form Actions */}
        <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
          <Button
            variant="secondary"
            onClick={onClose}
            classname="px-4 text-white py-2 text-sm"
          >
            Cancel
          </Button>
          <Button
            variant="primary"
            onClick={handleSubmit}
            classname="px-4 py-2 text-white text-sm"
          >
            {application ? 'Update' : 'Create'} Service
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default WebServiceModal
