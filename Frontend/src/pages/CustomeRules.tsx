import React, {useState} from 'react'
import {BsPersonX} from 'react-icons/bs'
import {ColumnDef} from '@tanstack/react-table'
import Table from '../components/Table'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import CreateRuleModal from '../components/CreateRuleModal'

interface CustomRulesData {
  status: string
  type: string
  name: string
  detail: string
  hitsToday: string
  updatedAt: string
}

const columns: ColumnDef<CustomRulesData>[] = [
  {
    header: 'Status',
    accessorKey: 'status',
    cell: ({row}) => (
      <span className="px-3 py-2 bg-green-400 text-white rounded">
        {row.original.status}
      </span>
    ),
  },
  {
    header: 'Name',
    accessorKey: 'name',
  },
  {
    header: 'Type',
    accessorKey: 'type',
    cell: ({row}) => (
      <div className="flex gap-3 items-center">
        {row.original.type.toLowerCase() === 'allow' ? (
          <BsPersonX size={20} />
        ) : (
          <BsPersonX />
        )}
        <span className="px-3 py-2 bg-green-400 text-white rounded">
          {row.original.type}
        </span>
      </div>
    ),
  },
  {
    header: 'Detail',
    accessorKey: 'detail',
  },
  {
    header: 'Hits Today',
    accessorKey: 'hitsToday',
  },
  {
    header: 'Updated At',
    accessorKey: 'updatedAt',
  },
]

const mockData: CustomRulesData[] = [
  {
    status: 'Enabled',
    type: 'Allow',
    name: 'Search Engine Spider',
    detail: 'Search Engine Spider',
    hitsToday: '10',
    updatedAt: '2024-11-14 12:47:35',
  },
  {
    status: 'Disabled',
    type: 'Block',
    name: 'Suspicious User',
    detail: 'Blocked due to unusual activity',
    hitsToday: '3',
    updatedAt: '2024-11-14 10:15:22',
  },
]

function CustomRules() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)

  return (
    <div>
      <CreateRuleModal isModalOpen={isModalOpen} onClose={toggleModal} />
      <Card className="flex justify-between items-center py-3 px-2">
        <p>Custom Rule</p>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          Add Applications
        </Button>
      </Card>
      <Card>
        <Table columns={columns} data={mockData} />
      </Card>
    </div>
  )
}

export default CustomRules
