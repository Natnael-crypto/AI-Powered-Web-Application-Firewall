import React, {useEffect, useState} from 'react'
import {
  Condition,
  RuleResponse,
  validActions,
  validRuleMethods,
  validRuleTypes,
} from '../lib/types'
import {useGetApplications} from '../hooks/api/useApplication'
import {updateRule} from '../services/rulesApi'
import Modal from './Modal'
import { useToast } from '../hooks/useToast'

export type AppOption = {
  application_id: string
  application_name: string
}

type EditRuleModalProps = {
  isOpen: boolean
  onClose: () => void
  rule: RuleResponse
}

const EditRuleModal: React.FC<EditRuleModalProps> = ({isOpen, onClose, rule}) => {
  const [ruleInput, setRuleInput] = useState<RuleResponse>({...rule})
  const [preview, setPreview] = useState<string>('')
  const [availableApps, setAvailableApps] = useState<AppOption[]>([])
  const {data: applications} = useGetApplications()
  const {addToast: toast} = useToast()


  useEffect(() => {
    if (!applications) return

    const fetchedApps: AppOption[] = applications.map(app => ({
      application_id: app.application_id,
      application_name: app.application_name,
    }))

    const ruleApps: AppOption[] = rule.applications
    const allAppsMap = new Map<string, AppOption>()

    fetchedApps.forEach(app => {
      allAppsMap.set(app.application_id, app)
    })

    ruleApps.forEach(app => {
      if (!allAppsMap.has(app.application_id)) {
        allAppsMap.set(app.application_id, app)
      }
    })

    setAvailableApps(Array.from(allAppsMap.values()))
  }, [applications, rule.applications])

  useEffect(() => {
    generateRule(ruleInput)
  }, [ruleInput])

  const updateCondition = (index: number, field: keyof Condition, value: string) => {
    const updatedConditions = [...ruleInput.rule_definition]
    updatedConditions[index][field] = value
    setRuleInput({...ruleInput, rule_definition: updatedConditions})
  }

  const handleAppAdd = (appId: string) => {
    if (!ruleInput.applications.find(app => app.application_id === appId)) {
      const appToAdd = availableApps.find(app => app.application_id === appId) || {
        application_id: appId,
        application_name: appId,
      }
      setRuleInput({
        ...ruleInput,
        applications: [...ruleInput.applications, appToAdd],
      })
    }
  }

  const handleAppRemove = (appId: string) => {
    setRuleInput({
      ...ruleInput,
      applications: ruleInput.applications.filter(a => a.application_id !== appId),
    })
  }

  const generateRule = (input: RuleResponse) => {
    const {rule_id, action, category, rule_definition} = input
    if (rule_definition.length === 0) return

    let ruleText = ''
    rule_definition.forEach((cond, i) => {
      const prefix = i === 0 ? 'SecRule' : '    SecRule'
      const chain = i < rule_definition.length - 1 ? `"chain"` : ''
      const firstLine =
        i === 0
          ? `"id:${rule_id},phase:2,${action},msg:'${category}'${
              rule_definition.length > 1 ? ',chain' : ''
            }"`
          : chain
      ruleText += `${prefix} ${cond.rule_type} "@${cond.rule_method} ${cond.rule_definition}" ${firstLine}\n`
    })

    setPreview(ruleText.trim())
  }

  const addCondition = () => {
    setRuleInput({
      ...ruleInput,
      rule_definition: [
        ...ruleInput.rule_definition,
        {rule_type: '', rule_method: '', rule_definition: ''},
      ],
    })
  }

  const saveRule = async () => {
    if (ruleInput.applications.length === 0) {
      toast('Please select at least one application.')
      return
    }

    const appIds = ruleInput.applications.map(appOp => appOp.application_id)

    const payloadTemplate = {
      rule_id: rule.rule_id,
      action: ruleInput.action,
      category: ruleInput.category,
      is_active: true,
      conditions: ruleInput.rule_definition.map(cond => ({
        rule_type: cond.rule_type,
        rule_method: cond.rule_method,
        rule_definition: cond.rule_definition,
      })),
      application_ids: appIds,
    }

    try {
      await updateRule(payloadTemplate)
      toast('Rule updated successfully for all selected applications!')
      onClose()
    } catch (error) {
      console.error('Error saving rule:', error)
      toast('An error occurred while saving the rule.')
    }
  }

  const selectedAppIds = ruleInput.applications.map(app => app.application_id)
  const unselectedApps = availableApps.filter(
    app => !selectedAppIds.includes(app.application_id),
  )

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Edit WAF Rule">
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
            {ruleInput.applications.map(app => (
              <span
                key={app.application_id}
                className="bg-green-200 text-green-900 px-3 py-1 rounded-full cursor-pointer hover:bg-red-200 hover:text-red-900 text-sm"
                onClick={() => handleAppRemove(app.application_id)}
              >
                {app.application_name || app.application_id} ×
              </span>
            ))}
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

          {ruleInput.rule_definition.map((cond, index) => (
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
                    className="self-end text-red-500 hover:text-red-700"
                    onClick={() => {
                      const updated = [...ruleInput.rule_definition]
                      updated.splice(index, 1)
                      setRuleInput({...ruleInput, rule_definition: updated})
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

export default EditRuleModal
