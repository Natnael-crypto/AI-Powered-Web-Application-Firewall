import {useState, useEffect} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'

interface RuleDefinitionItem {
  rule_type: string
  rule_method: string
  rule_definition: string
}

export interface Rule {
  rule_id: string
  rule_type: string
  rule_method: string
  rule_definition: string | RuleDefinitionItem[]
  action: string
  application_id: string
  rule_string: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  category: string
}

interface RuleFormModalProps {
  rule?: Rule
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: Partial<Rule>) => void
}

function RuleFormModal({rule, isOpen, onClose, onSubmit}: RuleFormModalProps) {
  const [form, setForm] = useState<Partial<Rule>>({
    rule_type: '',
    rule_method: '',
    rule_definition: '',
    action: '',
    rule_string: '',
    is_active: true,
    category: '',
    application_id: '',
  })

  useEffect(() => {
    if (rule) {
      setForm(rule)
    } else {
      setForm({
        rule_type: '',
        rule_method: '',
        rule_definition: '',
        action: '',
        rule_string: '',
        is_active: true,
        category: '',
        application_id: '',
      })
    }
  }, [rule])

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>,
  ) => {
    const {name, value, type} = e.target as HTMLInputElement
    const checked = (e.target as HTMLInputElement).checked
    setForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  const handleSubmit = () => {
    onSubmit(form)
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={rule ? 'Edit Rule' : 'Add New Rule'}>
      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <input
            name="rule_type"
            placeholder="Rule Type"
            value={form.rule_type || ''}
            onChange={handleChange}
            className="border rounded p-2"
          />
          <input
            name="rule_method"
            placeholder="Rule Method"
            value={form.rule_method || ''}
            onChange={handleChange}
            className="border rounded p-2"
          />
          <input
            name="action"
            placeholder="Action"
            value={form.action || ''}
            onChange={handleChange}
            className="border rounded p-2"
          />
          <input
            name="category"
            placeholder="Category"
            value={form.category || ''}
            onChange={handleChange}
            className="border rounded p-2"
          />
          <input
            name="application_id"
            placeholder="Application ID"
            value={form.application_id || ''}
            onChange={handleChange}
            className="border rounded p-2"
          />
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              name="is_active"
              checked={form.is_active ?? true}
              onChange={handleChange}
            />
            <label htmlFor="is_active">Active</label>
          </div>
        </div>

        <div>
          <label className="text-sm text-gray-500">Rule Definition</label>
          <textarea
            name="rule_definition"
            placeholder="Rule Definition"
            value={
              typeof form.rule_definition === 'string'
                ? form.rule_definition
                : JSON.stringify(form.rule_definition, null, 2)
            }
            onChange={handleChange}
            rows={4}
            className="w-full border rounded p-2"
          />
        </div>

        <div>
          <label className="text-sm text-gray-500">Rule String</label>
          <textarea
            name="rule_string"
            placeholder="Rule String"
            value={form.rule_string || ''}
            onChange={handleChange}
            rows={4}
            className="w-full border rounded p-2"
          />
        </div>

        <div className="flex justify-end space-x-2">
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSubmit}>
            {rule ? 'Update Rule' : 'Create Rule'}
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default RuleFormModal
