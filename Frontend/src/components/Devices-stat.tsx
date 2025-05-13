import {PieChart, Pie, Cell, ResponsiveContainer, Tooltip} from 'recharts'
import {useGetDeviceStat} from '../hooks/api/useRequests'

const OS_COLORS: Record<string, string> = {
  Windows: '#3366FF',
  Linux: '#7E88FF',
  Android: '#7CDDDD',
  macOS: '#FF6B81',
  iOS: '#B6F0E2',
  Other: '#999999',
}

interface UserClientsCardProps {
  selectedApp: string
  timeRange: any
}

export default function UserClientsCard({selectedApp,timeRange}:UserClientsCardProps) {
  const {data, isLoading, isError} = useGetDeviceStat()

  if (isLoading) return <p>Loading...</p>
  if (isError || !data) return <p>Something Went Wrong</p>

  const chartData = Object.entries(data)
    .filter(([_, value]) => (value as number) > 0)
    .map(([name, value]) => ({
      name,
      value,
      color: OS_COLORS[name] || '#ccc',
    }))

  return (
    <div className="w-full bg-white xl shadow-md p-6">
      <div className="flex justify-between items-center p-3">
        <p className="text-lg">User clients</p>
      </div>
      <div className="flex items-center gap-6">
        {/* Chart */}
        <div className="w-full sm:w-60 h-60">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={chartData}
                dataKey="value"
                outerRadius={90}
                innerRadius={75}
                startAngle={90}
                endAngle={-270}
                nameKey="name"
              >
                {chartData.map((entry, index) => (
                  <Cell key={`os-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Legend */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-32 gap-y-2 text-sm text-slate-700">
          {chartData.map(({name, value, color}) => (
            <div key={name} className="flex items-center gap-2">
              <span
                className="w-2.5 h-2.5 flex-shrink-0 rounded-full"
                style={{backgroundColor: color}}
              />
              <span className="truncate max-w-[8rem]">{name}</span>
              <span className="font-semibold ml-auto">{value as number}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
