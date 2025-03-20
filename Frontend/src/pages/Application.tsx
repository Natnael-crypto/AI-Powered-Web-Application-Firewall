import React, {useEffect, useState, useMemo} from 'react'
import Table from '../components/Table'
import Button from '../components/atoms/Button'
import Card from '../components/Card'
import {BiBot, BiKey} from 'react-icons/bi'
import {CiWavePulse1} from 'react-icons/ci'
import AddAppModal from '../components/AddAppModal'
import {CellContext, ColumnDef} from '@tanstack/react-table'
import {useGetAppliactions} from '../hooks/useApplication'

type applicationType = {
  applicationId: string
  applicationName: string
  description: string
  hostname: string
  ipAddress: string
  port: string
  status: boolean
  tls: boolean
  createdAt: string
  updatedAt: string
}

const Application: React.FC = () => {
  const [securityOptions, setSecurityOptions] = useState<Record<string, string>>({})
  const [isModalOpen, setIsModalOpen] = useState(false)

  const toggleModal = () => setIsModalOpen(prev => !prev)
  const {data: applications = [], isLoading, error} = useGetAppliactions()

  if (isLoading) return <div>Loading...</div>
  if (error) return <div>{error.message}</div>

  const columns: ColumnDef<applicationType>[] = useMemo(
    () => [
      {
        header: 'Application ID',
        accessorKey: 'applicationId',
      },
      {
        header: 'Application Name',
        accessorKey: 'applicationName',
      },
      {
        header: 'Description',
        accessorKey: 'description',
      },
      {
        header: 'Hostname',
        accessorKey: 'hostname',
      },
      {
        header: 'IP Address',
        accessorKey: 'ipAddress',
      },
      {
        header: 'Port',
        accessorKey: 'port',
      },
      {
        header: 'TLS',
        accessorKey: 'tls',
        cell: ({getValue}: CellContext<applicationType, unknown>) => {
          return getValue() ? 'Enabled' : 'Disabled'
        },
      },
      {
        header: 'Status',
        accessorKey: 'status',
        cell: ({getValue}: CellContext<applicationType, unknown>) => {
          return getValue() ? 'Active' : 'Inactive'
        },
      },
      {
        header: 'Created At',
        accessorKey: 'createdAt',
        cell: ({getValue}: CellContext<applicationType, unknown>) => {
          return new Date(getValue() as string).toLocaleString()
        },
      },
      {
        header: 'Updated At',
        accessorKey: 'updatedAt',
        cell: ({getValue}: CellContext<applicationType, string>) => {
          return new Date(getValue()).toLocaleString()
        },
      },
    ],
    [securityOptions],
  )

  return (
    <div className="px-10 py-2">
      <AddAppModal isModalOpen={isModalOpen} toggleModal={toggleModal} />
      <Card className="flex justify-between items-center py-3 px-2">
        <p>{applications.length} Applications</p>
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
        <Table columns={columns} data={applications} />
      </Card>
    </div>
  )
}

export default Application
