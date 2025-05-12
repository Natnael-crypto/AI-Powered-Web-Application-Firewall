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
  const mockData: Rule[] = [
    {
      rule_id: '592249559871523992',
      rule_type: 'multiple',
      rule_method: 'chained',
      rule_definition:
        '[{"rule_type":"REQUEST_URI","rule_method":"streq","rule_definition":"admin"}]',
      action: 'deny',
      application_id: '16d3f539-6c7b-45ac-b977-6a51c3582d29',
      rule_string:
        'SecRule REQUEST_URI "@streq admin" "id:592249559871523992,phase:2,deny,msg:\'blocked path\'',
      created_by: 'd1c65bfd-307a-4fd4-84f8-584881bbb60a',
      created_at: '2025-03-10T22:19:49.252583+03:00',
      updated_at: '2025-03-10T22:19:49.252583+03:00',
      is_active: true,
      category: 'blocked path',
    },
    {
      rule_id: '211430091046783650',
      rule_type: 'multiple',
      rule_method: 'chained',
      rule_definition:
        '[{"rule_type":"REQUEST_URI","rule_method":"streq","rule_definition":"admin"},{"rule_type":"REQUEST_URI","rule_method":"contains","rule_definition":"test"}]',
      action: 'deny',
      application_id: '16d3f539-6c7b-45ac-b977-6a51c3582d29',
      rule_string:
        'SecRule REQUEST_URI "@streq admin" "id:211430091046783650,phase:2,deny,msg:\'blocked path2\'\n    chain\n    SecRule REQUEST_URI "@contains test"',
      created_by: 'd1c65bfd-307a-4fd4-84f8-584881bbb60a',
      created_at: '2025-03-10T22:28:09.558079+03:00',
      updated_at: '2025-03-10T22:28:09.558079+03:00',
      is_active: true,
      category: 'blocked path2',
    },
  ]

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
