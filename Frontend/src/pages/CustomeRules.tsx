import React, {useState} from 'react'
import {BsPersonX} from 'react-icons/bs'
import Table from '../components/Table'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import CreateRuleModal from '../components/CreateRuleModal'

interface CustomeRulesData {
  status: string
  type: string
  name: string
  detail: string
  hitsToday: string
  updatedAt: string
}

function CustomeRules() {
  const columns: Array<{
    Header: string
    accessor: keyof CustomeRulesData
    Cell?: ({value}: {value: string}) => JSX.Element
  }> = React.useMemo(
    () => [
      {
        Header: 'Status',
        accessor: 'status',
        Cell: ({value}: {value: string}) => (
          <span className="px-3 py-2 bg-green-400 text-white">{value}</span>
        ),
      },
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'Type',
        accessor: 'type',
        Cell: ({value}: {value: string}) => (
          <>
            <div className="flex gap-3 items-center ">
              {value === 'allow' ? <BsPersonX size={20} /> : <BsPersonX />}
              <span className="px-3 py-2 bg-green-400 text-white">{value}</span>
            </div>
          </>
        ),
      },
      {
        Header: 'Detail',
        accessor: 'detail',
      },
      {
        Header: 'Hits Today',
        accessor: 'hitsToday',
      },
      {
        Header: 'Updated At',
        accessor: 'updatedAt',
      },
    ],
    [],
  )

  const data: CustomeRulesData[] = [
    {
      status: 'Enabled',
      type: 'Allow',
      name: 'Search Engine Spider',
      detail: 'Search Engine Spider',
      hitsToday: '10',
      updatedAt: '2024-11-1412:47:35',
    },
  ]

  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)

  return (
    <div className="">
      <CreateRuleModal isModalOpen={isModalOpen} onClose={toggleModal} />
      <Card className="flex justify-between items-center py3 px-2">
        <p> Custom Rule</p>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          Add Applicatins
        </Button>
      </Card>
      <Card>
        <Table columns={columns} data={data} />
      </Card>
    </div>
  )
}

export default CustomeRules
