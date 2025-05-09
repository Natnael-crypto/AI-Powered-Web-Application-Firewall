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
  const [certFile, setCertFile] = useState<File | null>(null)
  const [keyFile, setKeyFile] = useState<File | null>(null)
  const [form, setForm] = useState<WebServiceData>({
    application_name: application?.application_name ?? '',
    description: application?.description ?? '',
    hostname: application?.hostname ?? '',
    ip_address: application?.ip_address ?? '',
    port: application?.port ?? '',
    status: application?.status ?? true,
    tls: application?.tls ?? false,
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

  console.log('application: ', form)
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

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Add Web Service">
      <div className="space-y-8 px-2 sm:px-4 py-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="flex flex-col">
            <label className="text-sm font-semibold text-gray-700 mb-1">
              Application Name
            </label>
            <input
              type="text"
              name="application_name"
              value={form.application_name}
              onChange={handleChange}
              className="w-full border border-gray-300  px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-sm transition"
              placeholder="Enter application name"
            />
          </div>

          <div className="flex flex-col">
            <label className="text-sm font-semibold text-gray-700 mb-1">Hostname</label>
            <input
              type="text"
              name="hostname"
              value={form.hostname}
              onChange={handleChange}
              className="w-full border border-gray-300  px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-sm transition"
              placeholder="example.com"
            />
          </div>

          <div className="flex flex-col">
            <label className="text-sm font-semibold text-gray-700 mb-1">IP Address</label>
            <input
              type="text"
              name="ip_address"
              value={form.ip_address}
              onChange={handleChange}
              className="w-full border border-gray-300  px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-sm transition"
              placeholder="192.168.1.1"
            />
          </div>

          <div className="flex flex-col">
            <label className="text-sm font-semibold text-gray-700 mb-1">Port</label>
            <input
              type="text"
              name="port"
              value={form.port}
              onChange={handleChange}
              className="w-full border border-gray-300  px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-sm transition"
              placeholder="8080"
            />
          </div>
        </div>

        <div className="flex flex-col">
          <label className="text-sm font-semibold text-gray-700 mb-1">Description</label>
          <textarea
            name="description"
            value={form.description}
            onChange={handleChange}
            rows={3}
            className="w-full border border-gray-300  px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:outline-none shadow-sm transition resize-none"
            placeholder="Enter description"
          />
        </div>

        <div className="flex flex-wrap gap-6">
          <label className="flex items-center gap-2 text-gray-800">
            <input
              type="checkbox"
              name="status"
              checked={form.status}
              onChange={handleChange}
              className="w-4 h-4 text-blue-600 border-gray-300  focus:ring-2 focus:ring-blue-500"
            />
            <span className="text-sm font-medium">Active</span>
          </label>

          <label className="flex items-center gap-2 text-gray-800">
            <input
              type="checkbox"
              name="tls"
              checked={form.tls}
              onChange={handleChange}
              className="w-4 h-4 text-blue-600 border-gray-300  focus:ring-2 focus:ring-blue-500"
            />
            <span className="text-sm font-medium">TLS Enabled</span>
          </label>
        </div>

        {form.tls && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="flex flex-col">
              <label className="text-sm font-semibold text-gray-700 mb-1">
                Cert File
              </label>
              <input
                type="file"
                onChange={e => setCertFile(e.target.files?.[0] || null)}
                className="block w-full text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                required
              />
            </div>
            <div className="flex flex-col">
              <label className="text-sm font-semibold text-gray-700 mb-1">Key File</label>
              <input
                type="file"
                onChange={e => setKeyFile(e.target.files?.[0] || null)}
                className="block w-full text-sm text-gray-700 file:mr-4 file:py-2 file:px-4  file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                required
              />
            </div>
          </div>
        )}

        <div className="flex justify-end gap-3 pt-6 border-t mt-6 text-white">
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSubmit}>
            Submit
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default WebServiceModal
