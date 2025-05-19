import {ColumnDef} from '@tanstack/react-table'
import {Rule} from '../pages/CustomeRules'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {
  useActivateRule,
  useDeactivateRule,
  useDeleteRules,
  useGetRules,
} from '../hooks/api/useRules'
import {useToast} from '../hooks/useToast'
import LoadingSpinner from './LoadingSpinner'
import {useState} from 'react'
import EditRuleModal from './EditRuleModal'
import {RuleResponse} from '../lib/types'
import RuleDetailsModal from './RuleDetailModal'

interface RulesTableProps {
  onUpdate: (rule: Rule) => void
}

const RulesTable = ({onUpdate}: RulesTableProps) => {
  const {mutate: deactivateRule} = useDeactivateRule()
  const {mutate: activateRule} = useActivateRule()
  const {mutate: deleteRule} = useDeleteRules()
  const {addToast: toast} = useToast()
  const {data: rules, isLoading, isError} = useGetRules()

  const [isEditModalOpen, setEditModalOpen] = useState(false)
  const [isDetailModalOpen, setDetailModalOpen] = useState(false)
  const [selectedRule, setSelectedRule] = useState<RuleResponse | null>(null)

  const closeEditModal = () => {
    setEditModalOpen(false)
    setSelectedRule(null)
  }

  const openEditModal = (rule: RuleResponse) => {
    setSelectedRule(rule)
    setEditModalOpen(true)
  }

  const handleRowClick = (rule: RuleResponse) => {
    setSelectedRule(rule)
    setDetailModalOpen(true)
  }

  async function handleDeactivation(rule: RuleResponse) {
    deactivateRule(rule.rule_id, {
      onSuccess: () => {
        toast('deactivated successfully')
      },
    })
  }

  async function handleActivation(rule: RuleResponse) {
    activateRule(rule.rule_id, {
      onSuccess: () => {
        toast('activated successfully')
      },
    })
  }

  async function handleDeleteRule(rule: RuleResponse) {
    deleteRule(rule.rule_id, {
      onSuccess: () => {
        toast(`deleted a rule with id: ${rule.rule_id} successfully`)
      },
    })
  }

  const columns: ColumnDef<RuleResponse>[] = [
    {
      accessorKey: 'rule_id',
      header: 'Rule ID',
      cell: ({getValue, row}) => (
        <div onClick={() => handleRowClick(row.original)} className="cursor-pointer">
          {getValue() as string}
        </div>
      ),
    },
    {
      accessorKey: 'category',
      header: 'Category',
      cell: ({getValue, row}) => (
        <div onClick={() => handleRowClick(row.original)} className="cursor-pointer">
          {getValue() as string}
        </div>
      ),
    },
    {
      accessorKey: 'action',
      header: 'Action',
      cell: ({getValue, row}) => (
        <div onClick={() => handleRowClick(row.original)} className="cursor-pointer">
          {getValue() as string}
        </div>
      ),
    },
    {
      accessorKey: 'is_active',
      header: 'Active',
      cell: ({getValue, row}) => (
        <div onClick={() => handleRowClick(row.original)} className="cursor-pointer">
          {getValue() ? 'Yes' : 'No'}
        </div>
      ),
    },
    {
      accessorKey: 'created_at',
      header: 'Created At',
      cell: ({getValue, row}) => (
        <div onClick={() => handleRowClick(row.original)} className="cursor-pointer">
          {new Date(getValue() as string).toLocaleString()}
        </div>
      ),
    },
    {
      accessorKey: 'rule_string',
      header: 'Rule String',
      cell: ({getValue, row}) => (
        <div
          className="max-w-xs truncate cursor-pointer"
          title={getValue() as string}
          onClick={() => handleRowClick(row.original)}
        >
          {getValue() as string}
        </div>
      ),
    },
    {
      header: 'Actions',
      id: 'actions',
      cell: ({row}) => (
        <DropdownActions
          item={row.original}
          options={[
            {
              label: 'Activate',
              onClick: rule => handleActivation(rule),
              show: (rule: RuleResponse) => !rule.is_active,
            },
            {
              label: 'Deactivate',
              onClick: rule => handleDeactivation(rule),
              show: (rule: RuleResponse) => rule.is_active,
            },
            {
              label: 'Update Rule',
              onClick: rule => openEditModal(rule),
            },
            {
              label: 'Delete Rule',
              onClick: rule => handleDeleteRule(rule),
            },
          ]}
        />
      ),
    },
  ]

  if (isLoading) return <LoadingSpinner />
  if (isError) return <p>Something went wrong</p>

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">WAF Rules</h1>
      <Table columns={columns} data={rules} />

      {selectedRule && (
        <>
          <EditRuleModal
            isOpen={isEditModalOpen}
            onClose={closeEditModal}
            rule={selectedRule}
          />
          <RuleDetailsModal
            isOpen={isDetailModalOpen}
            onClose={() => setDetailModalOpen(false)}
            rule={selectedRule}
          />
        </>
      )}
    </div>
  )
}

export default RulesTable
