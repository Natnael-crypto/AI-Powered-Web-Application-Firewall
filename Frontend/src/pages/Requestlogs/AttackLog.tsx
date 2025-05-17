import AttackLogFilter from '../../components/AttackLogFilter'
import AttackLogTable from '../../components/Request_Logs/AttackLogTable'

const AttackLog = () => {

  return (
    <div className="">
      
      <div className="">

        <h1 className="text-xl font-bold mb-4">
          Request Logs
        </h1>
        <AttackLogFilter />
        <br />
        <AttackLogTable />
      </div>
    </div>
  )
}

export default AttackLog
