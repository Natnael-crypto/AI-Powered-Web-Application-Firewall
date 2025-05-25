import React, {useEffect, useState} from 'react'
import {
  AppOption,
  Condition,
  RuleInput,
  validActions,
  validRuleMethods,
  validRuleTypes,
} from '../lib/types'
import {useGetApplications} from '../hooks/api/useApplication'
import {createRule} from '../services/rulesApi'
import Modal from './Modal'

type CreateRuleModalProps = {
  isOpen: boolean
  onClose: () => void
}

const CreateRuleModal: React.FC<CreateRuleModalProps> = ({isOpen, onClose}) => {
  const [ruleInput, setRuleInput] = useState<RuleInput>({
    ruleID: '1001',
    action: 'deny',
    category: 'Message',
    conditions: [
      {
        rule_type: 'ARGS',
        rule_method: 'regex',
        rule_definition: 'value like *select*',
      },
    ],
    applications: [],
  })
  const [preview, setPreview] = useState<string>('')
  const [availableApps, setAvailableApps] = useState<AppOption[]>([])
  const {data: applications} = useGetApplications()

  useEffect(() => {
    if (applications) {
      const apps = applications.map(app => ({
        application_id: app.application_id,
        application_name: app.application_name,
      }))
      setAvailableApps(apps)
    }
  }, [applications])

  useEffect(() => {
    generateRule(ruleInput)
  }, [ruleInput])

  const updateCondition = (index: number, field: keyof Condition, value: string) => {
    const updatedConditions = [...ruleInput.conditions]
    updatedConditions[index][field] = value
    setRuleInput({...ruleInput, conditions: updatedConditions})
  }

  const handleAppAdd = (appId: string) => {
    if (!ruleInput.applications.includes(appId)) {
      setRuleInput({...ruleInput, applications: [...ruleInput.applications, appId]})
    }
  }

  const handleAppRemove = (appId: string) => {
    setRuleInput({
      ...ruleInput,
      applications: ruleInput.applications.filter(a => a !== appId),
    })
  }

  const generateRule = (input: RuleInput) => {
    const {ruleID, action, category, conditions} = input
    if (conditions.length === 0) return

    let ruleText = ''
    conditions.forEach((cond, i) => {
      const prefix = i === 0 ? 'SecRule' : '    SecRule'
      const chain = i < conditions.length - 1 ? `"chain"` : ''
      const firstLine =
        i === 0
          ? `"id:${ruleID},phase:2,${action},msg:'${category}'${conditions.length > 1 ? ',chain' : ''}"`
          : chain
      ruleText += `${prefix} ${cond.rule_type} "@${cond.rule_method} ${cond.rule_definition}" ${firstLine}\n`
    })

    setPreview(ruleText.trim())
  }

  const addCondition = () => {
    setRuleInput({
      ...ruleInput,
      conditions: [
        ...ruleInput.conditions,
        {rule_type: '', rule_method: '', rule_definition: ''},
      ],
    })
  }

  const saveRule = async () => {
    if (ruleInput.applications.length === 0) {
      alert('Please select at least one application.')
      return
    }

    const payloadTemplate = {
      action: ruleInput.action,
      category: ruleInput.category,
      is_active: true,
      conditions: ruleInput.conditions.map(cond => ({
        rule_type: cond.rule_type,
        rule_method: cond.rule_method,
        rule_definition: cond.rule_definition,
      })),
      application_ids: ruleInput.applications,
    }

    try {
      await createRule(payloadTemplate)
      alert('Rule saved successfully for all selected applications!')
      onClose()
    } catch (error) {
      console.error('Error saving rule:', error)
      alert('An error occurred while saving the rule.')
    }
  }

  const unselectedApps = availableApps.filter(
    app => !ruleInput.applications.includes(app.application_id),
  )

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Create WAF Rule">
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs text-gray-500 mb-1">Action</label>
            <select
              className="w-full px-3 py-2 border rounded-md"
              value={ruleInput.action}
              onChange={e => setRuleInput({...ruleInput, action: e.target.value})}
            >
              {validActions.map(action => (
                <option key={action} value={action}>
                  {action}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Category</label>
            <input
              className="w-full px-3 py-2 border rounded-md"
              placeholder="Category"
              value={ruleInput.category}
              onChange={e => setRuleInput({...ruleInput, category: e.target.value})}
            />
          </div>
        </div>

        <div>
          <label className="block text-xs text-gray-500 mb-1">Select Application</label>
          <select
            className="w-full px-3 py-2 border rounded-md"
            onChange={e => {
              const value = e.target.value
              if (value) handleAppAdd(value)
              e.target.value = ''
            }}
          >
            <option value="">-- Select Application --</option>
            {unselectedApps.map(app => (
              <option key={app.application_id} value={app.application_id}>
                {app.application_name}
              </option>
            ))}
          </select>

          <div className="flex flex-wrap gap-2 mt-2">
            {ruleInput.applications.map(appId => {
              const app = availableApps.find(a => a.application_id === appId)
              return (
                <span
                  key={appId}
                  className="bg-green-200 text-green-900 px-3 py-1 rounded-full cursor-pointer hover:bg-red-200 hover:text-red-900 text-sm"
                  onClick={() => handleAppRemove(appId)}
                >
                  {app?.application_name || appId} ×
                </span>
              )
            })}
          </div>
        </div>

        {/* Conditions */}
        <div className="space-y-4">
          <div className="flex items-center">
            <div className="flex-grow border-t border-gray-200"></div>
            <span className="mx-4 text-sm font-medium text-gray-500">
              RULE CONDITIONS
            </span>
            <div className="flex-grow border-t border-gray-200"></div>
          </div>

          {ruleInput.conditions.map((cond, index) => (
            <div
              key={index}
              className="grid grid-cols-1 md:grid-cols-3 gap-4 items-center"
            >
              <div>
                <label className="block text-xs text-gray-500 mb-1">Rule Type</label>
                <select
                  className="w-full px-3 py-2 border rounded-md"
                  value={cond.rule_type}
                  onChange={e => updateCondition(index, 'rule_type', e.target.value)}
                >
                  <option value="">Select Rule Type</option>
                  {validRuleTypes.map(type => (
                    <option key={type} value={type}>
                      {type}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-xs text-gray-500 mb-1">Method</label>
                <select
                  className="w-full px-3 py-2 border rounded-md"
                  value={cond.rule_method}
                  onChange={e => updateCondition(index, 'rule_method', e.target.value)}
                >
                  <option value="">Select Method</option>
                  {validRuleMethods.map(method => (
                    <option key={method} value={method}>
                      {method}
                    </option>
                  ))}
                </select>
              </div>

              <div className="flex gap-2">
                <div className="flex-1">
                  <label className="block text-xs text-gray-500 mb-1">Definition</label>
                  <input
                    className="w-full px-3 py-2 border rounded-md"
                    placeholder="Definition"
                    value={cond.rule_definition}
                    onChange={e =>
                      updateCondition(index, 'rule_definition', e.target.value)
                    }
                  />
                </div>
                {index > 0 && (
                  <button
                    className="self-end text-red-500 hover:text-red-700 font-lg"
                    onClick={() => {
                      const updated = [...ruleInput.conditions]
                      updated.splice(index, 1)
                      setRuleInput({...ruleInput, conditions: updated})
                    }}
                  >
                    ×
                  </button>
                )}
              </div>
            </div>
          ))}

          <button
            className="text-sm text-blue-600 hover:text-blue-800"
            onClick={addCondition}
          >
            + Add Condition
          </button>
        </div>

        {/* Rule Preview */}
        <div className="mt-4">
          <div className="flex items-center mb-2">
            <div className="flex-grow border-t border-gray-200"></div>
            <span className="mx-4 text-sm font-medium text-gray-500">RULE PREVIEW</span>
            <div className="flex-grow border-t border-gray-200"></div>
          </div>
          <pre className="bg-gray-100 p-4 rounded-md border text-sm whitespace-pre-wrap overflow-x-auto">
            {preview || '// Rule preview will appear here...'}
          </pre>
        </div>

        <div className="pt-4 border-t border-gray-100 flex justify-end gap-2">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
          >
            Cancel
          </button>
          <button
            className="px-4 py-2 text-sm text-white bg-black rounded-md hover:bg-gray-800 transition-colors"
            onClick={saveRule}
          >
            Save Rule
          </button>
        </div>
      </div>
    </Modal>
  )
}

export default CreateRuleModal
