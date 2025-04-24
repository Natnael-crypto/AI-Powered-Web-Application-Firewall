import LogLogs from '../components/Request_Logs/AttackLogTable'
import AttackLogFilter from '../components/AttackLogFilter'
import {useLogFilter} from '../store/LogFilter'

const LogAttack = () => {
  const {filterType, filterOperation} = useLogFilter()
  console.log(filterType, filterOperation)

  return (
    <div className="px-10 overflow-y-scroll my-5">
      <AttackLogFilter />
      <LogLogs />
    </div>
  )
}

export default LogAttack
