import React from 'react'
import Modal from './Modal'
import Button from './atoms/Button'
import {BiLockAlt, BiInfoCircle, BiTag} from 'react-icons/bi' // Icons for form fields

interface AddAppModalProps {
  isModalOpen: boolean
  toggleModal: () => void
}

function AddAppModal({isModalOpen, toggleModal}: AddAppModalProps) {
  return (
    <Modal isOpen={isModalOpen} onClose={toggleModal}>
      <div className="bg-gradient-to-r from-blue-500 to-blue-600 p-6 rounded-t-lg">
        <h2 className="text-xl font-semibold text-white">Add Application</h2>
      </div>
      <form className="p-6">
        <div className="space-y-6">
          {/* Application Name */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Application Name
            </label>
            <div className="relative">
              <input
                type="text"
                className="mt-1 block p-2 w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter application name"
              />
              <BiInfoCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* Port */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Port</label>
            <div className="relative">
              <input
                type="number"
                className="mt-1 block w-full p-2 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 pl-10"
                placeholder="Enter port number"
              />
              <BiTag className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            </div>
          </div>

          {/* Run Mode */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Run Mode
            </label>
            <select className="mt-1 block w-full p-2 rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500">
              <option>Production</option>
              <option>Development</option>
              <option>Staging</option>
            </select>
          </div>

          {/* Security Level */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Security Level
            </label>
            <div className="mt-1 space-y-2">
              {['High', 'Medium', 'Low'].map(option => (
                <label key={option} className="flex items-center space-x-2">
                  <input
                    type="radio"
                    name="security"
                    value={option}
                    className="form-radio h-4 w-4 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="text-gray-700">{option}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              rows={3}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              placeholder="Enter a brief description"
            />
          </div>

          {/* Tags */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Tags</label>
            <input
              type="text"
              className="mt-1 block w-full rounded-md p-2 border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
              placeholder="Enter tags (comma separated)"
            />
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
              onClick={() => {
                toggleModal()
              }}
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
