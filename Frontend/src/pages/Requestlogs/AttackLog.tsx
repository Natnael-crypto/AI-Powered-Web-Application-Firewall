import AttackLogFilter from '../../components/AttackLogFilter'
import AttackLogTable from '../../components/Request_Logs/AttackLogTable'
import {useLogFilter} from '../../store/LogFilter'

const AttackLog = () => {
  const {filterType, filterOperation} = useLogFilter()
  console.log(filterOperation, filterType)

  return (
    <div className="p-10 space-y-6 bg-gradient-to-b from-white to-green-50 min-h-screen">
      <h1 className="text-3xl font-extrabold text-green-800 tracking-tight">
        Attack Logs
      </h1>
      <div className="bg-white xl shadow-lg p-6 border border-gray-200">
        <AttackLogFilter />
      </div>
      <AttackLogTable />
    </div>
  )
}

export default AttackLog
