import React, {useState} from 'react'
import Table from '../components/Table'
import Button from '../components/atoms/Button'
import Card from '../components/Card'
import {BiBot, BiKey} from 'react-icons/bi'
import {CiWavePulse1} from 'react-icons/ci'
import AddAppModal from '../components/AddAppModal'

interface ApplicationData {
  application: string
  port: number
  runMode: string | any
  security: string
  today: {
    totalRequests: number
    blockedRequests: number
  }
}

function Application() {
  const [securityOptions, setSecurityOptions] = useState<Record<number, string>>({
    0: 'High',
    1: 'Medium',
    2: 'Low',
  })

  const columns = React.useMemo(
    () => [
      {
        Header: 'Application',
        accessor: 'application' as const,
      },
      {
        Header: 'Port',
        accessor: 'port' as const,
      },
      {
        Header: 'Run Mode',
        accessor: 'runMode' as const,
      },
      {
        Header: 'Security',
        accessor: 'security' as const,
        Cell: ({row}: {row: {index: number; original: ApplicationData}}) => {
          const {index, original} = row
          const selectedOption = securityOptions[index]

          return (
            <div className="flex justify-between gap-5">
              {[
                {label: 'Http Flood', icon: <CiWavePulse1 size={20} />},
                {label: 'Bot Protection', icon: <BiBot />},
                {label: 'Auth', icon: <BiKey />},
              ].map(option => (
                <button
                  key={option.label}
                  onClick={() =>
                    setSecurityOptions(prev => ({
                      ...prev,
                      [index]: option.label,
                    }))
                  }
                  className={`px-4 py-2 rounded flex items-center gap-2 flex-1 justify-center ${
                    selectedOption === option.label
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-200 text-gray-700'
                  }`}
                >
                  {option.icon}
                  {option.label}
                </button>
              ))}
            </div>
          )
        },
      },
      {
        Header: 'Today',
        accessor: 'today' as const,
        Cell: ({value}: {value: {totalRequests: number; blockedRequests: number}}) => (
          <div>
            <p>Total Requests: {value.totalRequests}</p>
            <p>Blocked Requests: {value.blockedRequests}</p>
          </div>
        ),
      },
    ],
    [securityOptions],
  )

  const data: ApplicationData[] = React.useMemo(
    () => [
      {
        application: 'Web Server',
        port: 8080,
        runMode: 'Production',
        security: 'High',
        today: {
          totalRequests: 1200,
          blockedRequests: 200,
        },
      },
      {
        application: 'Database',
        port: 3306,
        runMode: 'Development',
        security: 'Medium',
        today: {
          totalRequests: 800,
          blockedRequests: 50,
        },
      },
      {
        application: 'API Gateway',
        port: 443,
        runMode: <Button classname="bg-green-600 px-3 py-2 text-white">Defense</Button>,
        security: 'Low',
        today: {
          totalRequests: 1500,
          blockedRequests: 300,
        },
      },
    ],
    [],
  )
  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)

  return (
    <div className="px-10 py-2">
      <AddAppModal isModalOpen={isModalOpen} toggleModal={toggleModal} />
      <Card className="flex justify-between items-center py3 px-2">
        <p>3 Applications</p>
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

export default Application
