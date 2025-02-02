import React, {useState} from 'react'
import axios from 'axios'
import Modal from './Modal'

interface RuleUpdatePayload {
  rule_type: string
  rule_definition: string
  action: string
  application_id: string
  is_active: boolean
  category: string
}

interface ModalProps {
  ruleId?: string
  onClose: () => void
  isModalOpen?: boolean
}

const CreateRuleModal: React.FC<ModalProps> = ({ruleId, onClose, isModalOpen}) => {
  const [payload, setPayload] = useState<RuleUpdatePayload>({
    rule_type: 'REQUEST_URI',
    rule_definition: '@rx ^/admin',
    action: 'Block access to admin',
    application_id: 'b028fd26-fd2c-4486-8a1f-d1b510b652f0',
    is_active: false,
    category: 'Access Control',
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const {name, value, type} = e.target
    setPayload({
      ...payload,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await axios.put(
        `http://localhost:8080/rule/update/${ruleId}`,
        payload,
        {
          headers: {
            Authorization:
              'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQyNTQxMjksInJvbGUiOiJzdXBlcl9hZG1pbiIsInVzZXJfaWQiOiIyYTY2OGJmNy02ZGEwLTRiYjEtYTIzZS1kYzI2NTNiYjZmMmYifQ.uzuC7ehu9KjL4hdaxmHFamkctValYb6WEf_XCWAR2-k',
          },
        },
      )
      console.log('Rule updated successfully:', response.data)
      onClose()
    } catch (error) {
      console.error('Error updating rule:', error)
    }
  }

  return (
    <Modal isOpen={isModalOpen} onClose={onClose}>
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
        <div className="bg-white p-6 rounded-lg w-1/3">
          <h2 className="text-xl font-bold mb-4">Update Rule</h2>
          <form onSubmit={handleSubmit}>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Rule Type</label>
              <input
                type="text"
                name="rule_type"
                value={payload.rule_type}
                onChange={handleChange}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">
                Rule Definition
              </label>
              <input
                type="text"
                name="rule_definition"
                value={payload.rule_definition}
                onChange={handleChange}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Action</label>
              <input
                type="text"
                name="action"
                value={payload.action}
                onChange={handleChange}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">
                Application ID
              </label>
              <input
                type="text"
                name="application_id"
                value={payload.application_id}
                onChange={handleChange}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Is Active</label>
              <input
                type="checkbox"
                name="is_active"
                checked={payload.is_active}
                onChange={handleChange}
                className="mt-1 block"
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700">Category</label>
              <input
                type="text"
                name="category"
                value={payload.category}
                onChange={handleChange}
                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div className="flex justify-end">
              <button
                type="button"
                onClick={onClose}
                className="mr-2 px-4 py-2 bg-gray-300 rounded-md hover:bg-gray-400"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
              >
                Update
              </button>
            </div>
          </form>
        </div>
      </div>
    </Modal>
  )
}

export default CreateRuleModal
