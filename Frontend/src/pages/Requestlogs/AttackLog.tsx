import AttackLogFilter from '../../components/AttackLogFilter'
import AttackLogTable from '../../components/Request_Logs/AttackLogTable'
import {useLogFilter} from '../../store/LogFilter'

const AttackLog = () => {
  const {filterType, filterOperation} = useLogFilter()
  console.log(filterOperation, filterType)

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
