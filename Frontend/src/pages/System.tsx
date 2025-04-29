import {ColumnDef} from '@tanstack/react-table'
import ManageUser from '../components/ManageUser'
import Table from '../components/Table'
import CleanDataSettings from '../components/CleanDataSettings'
import AttackAlertSettings from '../components/AttackAlertSetting'
import SyslogSettings from '../components/SyslogSetting'

type UserAccount = {
  name: string
  role: string
  '2FA': boolean
  lastLogin: string
}

const mockUserData: UserAccount[] = [
  {
    name: 'Alice Johnson',
    role: 'Admin',
    '2FA': true,
    lastLogin: '2023-10-01T10:00:00Z',
  },
  {
    name: 'Bob Smith',
    role: 'User',
    '2FA': false,
    lastLogin: '2023-10-02T12:30:00Z',
  },
  {
    name: 'Charlie Brown',
    role: 'Moderator',
    '2FA': true,
    lastLogin: '2023-10-03T14:45:00Z',
  },
]

function System() {
  const columns: ColumnDef<UserAccount>[] = [
    {
      accessorKey: 'name',
      header: 'Name',
    },
    {
      accessorKey: 'role',
      header: 'Role',
    },
    {
      accessorKey: '2FA',
      header: '2FA Enabled',
      cell: info => (info.getValue() ? 'Enabled' : 'Disabled'),
    },
    {
      accessorKey: 'lastLogin',
      header: 'Last Login',
      cell: info => new Date(info.getValue() as string).toLocaleString(),
    },
  ]

  return (
    <div className="flex flex-col gap-8 px-6 py-10 bg-gradient-to-br from-slate-100 to-white min-h-screen">
      <div className="max-w-7xl w-full mx-auto space-y-8">
        <section className="bg-white border border-slate-200 rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">Manage Users</h2>
          <ManageUser />
        </section>

        <section className="bg-white border border-slate-200 rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">User Table</h2>
          <Table columns={columns} data={mockUserData} />
        </section>

        <section className="bg-white border border-slate-200 rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">
            Clean Data Settings
          </h2>
          <CleanDataSettings />
        </section>

        <section className="bg-white border border-slate-200 rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">
            Attack Alert Settings
          </h2>
          <AttackAlertSettings />
        </section>

        <section className="bg-white border border-slate-200 rounded-2xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">Syslog Settings</h2>
          <SyslogSettings />
        </section>
      </div>
    </div>
  )
}

export default System
