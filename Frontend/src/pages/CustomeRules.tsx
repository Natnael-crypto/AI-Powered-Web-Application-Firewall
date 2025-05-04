import {useState} from 'react'
import {ColumnDef} from '@tanstack/react-table'
import Table from '../components/Table'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import CreateRuleModal from '../components/CreateRuleModal'
import {HiOutlineBan, HiOutlineCheckCircle} from 'react-icons/hi'

export interface Rule {
  rule_id: string
  rule_type: string
  rule_method: string
  rule_definition: string
  action: string
  application_id: string
  rule_string: string
  created_by: string
  created_at: string
  updated_at: string
  is_active: boolean
  category: string
}

const columns: ColumnDef<Rule>[] = [
  {
    header: 'Status',
    accessorKey: 'is_active',
    cell: ({row}) => (
      <span
        className={`px-3 py-1 ull text-sm font-medium ${
          row.original.is_active
            ? 'bg-green-100 text-green-800'
            : 'bg-red-100 text-red-800'
        }`}
      >
        {row.original.is_active ? 'Active' : 'Inactive'}
      </span>
    ),
  },
  {
    header: 'Type',
    accessorKey: 'rule_type',
  },
  {
    header: 'Method',
    accessorKey: 'rule_method',
  },
  {
    header: 'Definition',
    accessorKey: 'rule_definition',
  },
  {
    header: 'Action',
    accessorKey: 'action',
    cell: ({row}) => (
      <div className="flex items-center gap-2">
        {row.original.action.toLowerCase().includes('deny') ? (
          <HiOutlineBan className="text-red-600" size={20} />
        ) : (
          <HiOutlineCheckCircle className="text-green-600" size={20} />
        )}
        {row.original.action}
      </div>
    ),
  },
  {
    header: 'Updated At',
    accessorKey: 'updated_at',
  },
]

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
      'SecRule REQUEST_URI "@streq admin" "id:592249559871523992,phase:2,deny,msg:\'blocked path\'"',
    created_by: 'd1c65bfd-307a-4fd4-84f8-584881bbb60a',
    created_at: '2025-03-10T22:19:49.252583+03:00',
    updated_at: '2025-03-10T22:19:49.252583+03:00',
    is_active: true,
    category: 'blocked path',
  },
]

function CustomRules() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)

  return (
    <div className="space-y-4">
      <CreateRuleModal isModalOpen={isModalOpen} onClose={toggleModal} />
      <Card className="flex justify-between items-center py-4 px-6 shadow-md bg-white ">
        <h2 className="text-lg font-semibold">Custom Rules</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          Add Rule
        </Button>
      </Card>
      <Card className="shadow-md p-4 bg-white ">
        <Table columns={columns} data={mockData} />
      </Card>
    </div>
  )
}

export default CustomRules
