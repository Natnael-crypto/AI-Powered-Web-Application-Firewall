import {useState} from 'react'
import {ColumnDef} from '@tanstack/react-table'
import Table from './Table'
import {DropdownActions} from './DropdownAction'
import {NotificationRule} from '../lib/types'
import NotificationRuleModal from '../components/UpdateNotficationRule'

interface NotificationRulesTableProps {
  data: NotificationRule[]
}

function getColumns(
  onUpdateRule: (rule: NotificationRule) => void,
): ColumnDef<NotificationRule>[] {
  return [
    {
      header: 'Name',
      accessorKey: 'name',
      cell: ({row}) => <div>{row.original.name}</div>,
    },
    {
      header: 'Threat Type',
      accessorKey: 'threat_type',
      cell: ({row}) => <div>{row.original.threat_type}</div>,
    },
    {
      header: 'Threshold',
      accessorKey: 'threshold',
      cell: ({row}) => <div>{row.original.threshold}</div>,
    },
    {
      header: 'Time Window (min)',
      accessorKey: 'time_window',
      cell: ({row}) => <div>{row.original.time_window}</div>,
    },
    {
      header: 'Status',
      accessorKey: 'is_active',
      cell: ({row}) => (
        <span
          className={`px-3 py-1 rounded-full text-sm font-medium ${
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
      header: 'Actions',
      id: 'actions',
      cell: ({row}) => (
        <DropdownActions
          item={row.original}
          options={[
            {
              label: 'Update Rule',
              onClick: () => onUpdateRule(row.original), // Trigger the update action
            },
          ]}
        />
      ),
    },
  ]
}

export function NotificationRulesTable({data}: NotificationRulesTableProps) {
  const [isModalOpen, setIsModalOpen] = useState(false) // State to manage modal visibility
  const [selectedRule, setSelectedRule] = useState<NotificationRule | null>(null) // State to store the selected rule for update

  const handleUpdateRule = (rule: NotificationRule) => {
    setSelectedRule(rule) // Set the selected rule for the modal
    setIsModalOpen(true) // Open the modal
  }

  const columns = getColumns(handleUpdateRule) // Pass the handler to the columns

  const handleModalClose = () => {
    setIsModalOpen(false) // Close the modal
    setSelectedRule(null) // Clear the selected rule
  }

  return (
    <>
      <Table columns={columns} data={data} />

      {/* Modal for updating the rule */}
      {isModalOpen && selectedRule && (
        <NotificationRuleModal
          rule={selectedRule}
          onClose={handleModalClose}
          isOpen={isModalOpen}
        />
      )}
    </>
  )
}

export default NotificationRulesTable
