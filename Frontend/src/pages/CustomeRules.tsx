import {useState} from 'react'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import RulesTable from '../components/RulesTable' // Your existing table component
import {useCreateRule, useUpdateRule} from '../hooks/api/useRules'
import RuleDetailsModal from '../components/RuleDetailModal'

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

function CustomRules() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedRule, setSelectedRule] = useState<Rule | undefined>()
  const {mutate: createRule} = useCreateRule()
  const {mutate: updateRule} = useUpdateRule()

  const toggleModal = () => setIsModalOpen(!isModalOpen)

  const handleOpenDetailsModal = (rule: Rule) => {
    setSelectedRule(rule)
    setIsModalOpen(true)
  }

  const handlCreateRule = () => {
    setSelectedRule(undefined)
    setIsModalOpen(true)
  }

  const handleFormSubmit = (formData: Partial<Rule>) => {
    const isUpdate = !!formData.rule_id

    const mutationFn = isUpdate ? updateRule : createRule

    mutationFn(formData, {
      onSuccess: () => {
        toggleModal()
      },
      onError: () => {
        console.error('Something went wrong while saving the application.')
      },
    })
  }

  return (
    <div className="space-y-4">
      <RuleDetailsModal
        rule={selectedRule}
        isOpen={isModalOpen}
        onClose={toggleModal}
        onSubmit={handleFormSubmit}
      />

      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Custom Rules</h2>
        <Button onClick={handlCreateRule}>Add Rule</Button>
      </Card>

      <Card className="shadow-md p-4 bg-white">
        <RulesTable onUpdate={handleOpenDetailsModal} />
      </Card>
    </div>
  )
}

export default CustomRules
