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
  isModalOpen: boolean
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
      <div className="space-y-6 animate-fade-in">
        <h2 className="text-2xl font-semibold text-gray-800 text-center">Update Rule</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          {[
            {label: 'Rule Type', name: 'rule_type'},
            {label: 'Rule Definition', name: 'rule_definition'},
            {label: 'Action', name: 'action'},
            {label: 'Application ID', name: 'application_id'},
            {label: 'Category', name: 'category'},
          ].map(({label, name}) => (
            <div key={name}>
              <label className="block text-sm font-medium text-gray-600">{label}</label>
              <input
                type="text"
                name={name}
                value={payload[name as keyof RuleUpdatePayload] as string}
                onChange={handleChange}
                className="w-full px-4 py-2 mt-1 border rounded-lg shadow-sm focus:ring-2 focus:ring-indigo-400 focus:outline-none"
              />
            </div>
          ))}

          <div className="flex items-center gap-3">
            <label className="text-sm font-medium text-gray-600">Is Active</label>
            <input
              type="checkbox"
              name="is_active"
              checked={payload.is_active}
              onChange={handleChange}
              className="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
          </div>

          <div className="flex justify-end gap-3 pt-4 border-t mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 bg-gray-200 hover:bg-gray-300 rounded-lg"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-indigo-600 text-white hover:bg-indigo-700 rounded-lg"
            >
              Update Rule
            </button>
          </div>
        </form>
      </div>
    </Modal>
  )
}

export default CreateRuleModal
