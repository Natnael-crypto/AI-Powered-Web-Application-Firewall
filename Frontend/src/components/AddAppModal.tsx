import React, {useState} from 'react'
import axios from 'axios'
import Modal from './Modal'
import Button from './atoms/Button'
import {BiInfoCircle, BiTag} from 'react-icons/bi'

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
}

function AddAppModal({isModalOpen, toggleModal}: AddAppModalProps) {
  const [payload, setPayload] = useState<AddAppPayload>({
    application_name: '',
    description: '',
    hostname: '',
    ip_address: '',
    port: '',
    status: true,
  })

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const {name, value, type} = e.target
    setPayload({
      ...payload,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await axios.post(
        'http://localhost:8080/application/add',
        payload,
        {
          headers: {
            Authorization:
              'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzgyNTA3MDksInJvbGUiOiJzdXBlcl9hZG1pbiIsInVzZXJfaWQiOiIyYTY2OGJmNy02ZGEwLTRiYjEtYTIzZS1kYzI2NTNiYjZmMmYifQ.gPOUxy-z3Xfc9jrpMD63g076SS46XZ8RMJcViIXeuvA',
          },
        },
      )
      console.log('Application added successfully:', response.data)
      toggleModal()
    } catch (error) {
      console.error('Error adding application:', error)
    }
  }

  return (
    <Modal isOpen={isModalOpen} onClose={toggleModal}>
      <div className="bg-gradient-to-r from-blue-500 to-blue-600 p-6 rounded-t-lg">
        <h2 className="text-xl font-semibold text-white">Add Application</h2>
      </div>
      <form className="p-6" onSubmit={handleSubmit}>
        <div className="space-y-6">
          {/* Application Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Application Name
            </label>
            <div className="relative">
              <input
                type="text"
                name="application_name"
                value={payload.application_name}
                onChange={handleChange}
                className="mt-1 block p-2 w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter application name"
                required
              />
              <BiInfoCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              name="description"
              value={payload.description}
              onChange={handleChange}
              rows={3}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              placeholder="Enter a brief description"
              required
            />
          </div>

          {/* Hostname */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Hostname
            </label>
            <div className="relative">
              <input
                type="text"
                name="hostname"
                value={payload.hostname}
                onChange={handleChange}
                className="mt-1 block p-2 w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter hostname"
                required
              />
              <BiInfoCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* IP Address */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              IP Address
            </label>
            <div className="relative">
              <input
                type="text"
                name="ip_address"
                value={payload.ip_address}
                onChange={handleChange}
                className="mt-1 block p-2 w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter IP address"
                required
              />
              <BiInfoCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* Port */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Port</label>
            <div className="relative">
              <input
                type="text"
                name="port"
                value={payload.port}
                onChange={handleChange}
                className="mt-1 block w-full p-2 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter port number"
                required
              />
              <BiTag className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* Status */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              name="status"
              value={payload.status ? 'true' : 'false'}
              onChange={handleChange}
              className="mt-1 block w-full p-2 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              required
            >
              <option value="true">Active</option>
              <option value="false">Inactive</option>
            </select>
          </div>

          {/* Buttons */}
          <div className="flex justify-end space-x-4">
            <Button
              classname="bg-gray-300 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-400 transition-all duration-200"
              onClick={toggleModal}
            >
              Cancel
            </Button>
            <Button
              classname="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 transition-all duration-200"
              onClick={() => handleSubmit}
            >
              Save
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  )
}

export default AddAppModal
