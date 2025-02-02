import {Column} from 'react-table'
import {LogTable} from '../../lib/types'
import Table from '../Table'

function LogLogs() {
  const columns: Column<LogTable>[] = [
    {
      Header: 'Action',
      accessor: 'action',
      Cell: ({value}: {value: string}) => (
        <div
          className={` rounded-md py-1 text-white text-center ${value.toLowerCase() === 'blocked' ? 'bg-red-700' : 'bg-yellow-400'}`}
        >
          {value}
        </div>
      ),
    },
    {Header: 'URL', accessor: 'url'},
    {Header: 'Attack Type', accessor: 'attackType'},
    {Header: 'Ip Address', accessor: 'ipAddress'},
    {Header: 'Time', accessor: 'time'},
  ]

  const mockData: LogTable[] = [
    {
      action: 'Blocked',
      url: 'https://example.com/login',
      attackType: 'SQL Injection',
      ipAddress: '192.168.1.101',
      time: '2023-10-01T14:32:45Z',
    },
    {
      action: 'Allowed',
      url: 'https://example.com/dashboard',
      attackType: 'None',
      ipAddress: '192.168.1.102',
      time: '2023-10-02T09:15:22Z',
    },
    {
      action: 'Blocked',
      url: 'https://example.com/checkout',
      attackType: 'XSS Attack',
      ipAddress: '192.168.1.103',
      time: '2023-10-03T18:47:11Z',
    },
    {
      action: 'Blocked',
      url: 'https://example.com/api/data',
      attackType: 'Brute Force',
      ipAddress: '192.168.1.104',
      time: '2023-10-04T11:05:33Z',
    },
    {
      action: 'Allowed',
      url: 'https://example.com/about',
      attackType: 'None',
      ipAddress: '192.168.1.105',
      time: '2023-10-05T16:20:59Z',
    },
  ]
  return <Table data={mockData} columns={columns} />
}

export default LogLogs
