import {Column} from 'react-table'
import Table from '../components/Table'
import LogFilter from '../components/Request_Logs/LogFilter'
import {AttackLog, FilterValues} from '../lib/types'
import {useState} from 'react'
import LogLogs from '../components/Request_Logs/LogLogs'

const LogAttack = () => {
  const columns: Column<AttackLog>[] = [
    {Header: 'IP Address', accessor: 'ipAddress'},
    {Header: 'Application', accessor: 'application'},
    {Header: 'Attack Count', accessor: 'attackCount'},
    {Header: 'Duration', accessor: 'duration'},
    {Header: 'Start At', accessor: 'startAt'},
  ]

  const initialData: AttackLog[] = [
    {
      ipAddress: '192.168.1.1',
      application: 'https://webserver.example.com',
      attackCount: 12,
      duration: '2h 30m',
      startAt: '2023-10-01 14:00',
    },
    {
      ipAddress: '10.0.0.1',
      application: 'https://database.example.com',
      attackCount: 5,
      duration: '1h 15m',
      startAt: '2023-10-02 09:30',
    },
    {
      ipAddress: '172.16.0.1',
      application: 'https://apigateway.example.com',
      attackCount: 8,
      duration: '3h 45m',
      startAt: '2023-10-03 18:20',
    },
  ]

  const [filteredData, setFilteredData] = useState(initialData)
  const [logType, setLogType] = useState<'log' | 'event'>('log')

  const handleFilter = (filters: FilterValues) => {
    const filtered = initialData.filter(log => {
      return (
        (!filters.ipAddress || log.ipAddress.includes(filters.ipAddress)) &&
        (!filters.domain || log.application.includes(filters.domain)) &&
        (!filters.startAt || log.startAt >= filters.startAt) &&
        (!filters.endAt || log.startAt <= filters.endAt)
      )
    })
    setFilteredData(filtered)
  }

  const toggleType = () => setLogType(logType === 'log' ? 'event' : 'log')

  return (
    <div className="px-10 overflow-y-scroll my-5">
      <LogFilter onFilter={handleFilter} logtype={logType} onLogtypeChange={toggleType} />
      {logType === 'event' ? (
        <Table columns={columns} data={filteredData} />
      ) : (
        <LogLogs />
      )}
    </div>
  )
}

export default LogAttack
