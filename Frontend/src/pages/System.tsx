import ManageUser from '../components/ManageUser'
import CleanDataSettings from '../components/CleanDataSettings'
import AttackAlertSettings from '../components/AttackAlertSetting'
import SyslogSettings from '../components/SyslogSetting'

import UserTable from '../components/UserTable'
import {useState} from 'react'
import AddUserModal from '../components/AddUserModal'
import {useAddAdmin} from '../hooks/api/useUser'
import {QueryClient} from '@tanstack/react-query'

function System() {
  const [isAddUser, setAddUser] = useState(false)
  const toggleAddUser = () => setAddUser(prev => !prev)

  const {mutate} = useAddAdmin()

  const handleAddAdmin = (data: any) => {
    mutate(data)
  }

  return (
    <div className="flex flex-col gap-8 px-6 py-10 bg-gradient-to-br from-slate-100 to-white min-h-screen">
      <div className="max-w-7xl w-full mx-auto space-y-8">
        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">Manage Users</h2>
          <ManageUser toggleAddUser={toggleAddUser} />
          <AddUserModal
            isOpen={isAddUser}
            onClose={toggleAddUser}
            onSubmit={data => handleAddAdmin(data)}
          />
        </section>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">User Table</h2>
          <UserTable />
        </section>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">
            Clean Data Settings
          </h2>
          <CleanDataSettings />
        </section>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">
            Attack Alert Settings
          </h2>
          <AttackAlertSettings />
        </section>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">Syslog Settings</h2>
          <SyslogSettings />
        </section>
      </div>
    </div>
  )
}

export default System
