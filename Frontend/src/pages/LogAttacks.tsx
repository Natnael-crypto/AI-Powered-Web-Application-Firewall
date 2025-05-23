import LogLogs from '../components/Request_Logs/AttackLogTable'
import AttackLogFilter from '../components/AttackLogFilter'

const LogAttack = () => {
  return (
    <div className="px-10 overflow-y-scroll my-5">
      <AttackLogFilter />
      <LogLogs />
    </div>
  )
}

export default LogAttack
