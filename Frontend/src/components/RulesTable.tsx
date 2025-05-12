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

interface rulesTabelProps {
  onUpdate: (rule: Rule) => void
}

const RulesTable = ({onUpdate}: rulesTabelProps) => {
  const {mutate: deactivateRule} = useDeactivateRule()
  const {mutate: activateRule} = useActivateRule()
  const {mutate: deleteRule} = useDeleteRules()
  const {addToast: toast} = useToast()
  const {data: rules, isLoading, isError} = useGetRules()
  async function handleDeactivation(rule: Rule) {
    deactivateRule(rule.rule_id, {
      onSuccess: () => {
        toast('deactivated successfully')
      },
    })
  }
  async function handleActivation(rule: Rule) {
    activateRule(rule.rule_id, {
      onSuccess: () => {
        toast('activated successfully')
      },
    })
  }
  async function handleDeleteRule(rule: Rule) {
    deleteRule(rule.rule_id, {
      onSuccess: () => {
        toast(`deleted a rule with id: ${rule.rule_id} successfully`)
      },
    })
  }

  const columns: ColumnDef<Rule>[] = [
    {
      accessorKey: 'rule_id',
      header: 'Rule ID',
    },
    {
      accessorKey: 'category',
      header: 'Category',
    },
    {
      accessorKey: 'action',
      header: 'Action',
    },
    {
      accessorKey: 'is_active',
      header: 'Active',
      cell: ({getValue}) => (getValue() ? 'Yes' : 'No'),
    },
    {
      accessorKey: 'created_at',
      header: 'Created At',
      cell: ({getValue}) => new Date(getValue() as string).toLocaleString(),
    },
    {
      accessorKey: 'rule_string',
      header: 'Rule String',
      cell: ({getValue}) => (
        <div className="max-w-xs truncate" title={getValue() as string}>
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
              show: (rule: Rule) => !rule.is_active,
            },
            {
              label: 'Deactivate',
              onClick: rule => handleDeactivation(rule),
              show: (rule: Rule) => rule.is_active,
            },
            {
              label: 'update Rule',
              onClick: rule => onUpdate(rule),
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
      <h1 className="text-2xl font-bold mb-4">Security Rules</h1>
      <Table columns={columns} data={rules} />
    </div>
  )
}

export default RulesTable
