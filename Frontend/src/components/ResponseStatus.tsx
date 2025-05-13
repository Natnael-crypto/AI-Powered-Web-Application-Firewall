import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts'
import { useResponseStat } from '../hooks/api/useDashboardStat'
import { JSXElementConstructor, Key, ReactElement, ReactNode, ReactPortal, useEffect } from 'react'

const COLORS = ['#3366FF', '#7E88FF', '#7CDDDD', '#FF6B81', '#B6F0E2']

interface ResponseStatusProps {
  selectedApp: string
  timeRange: any
}

export default function ResponseStatus({ selectedApp, timeRange }: ResponseStatusProps) {
  const { data = [], refetch } = useResponseStat(selectedApp, timeRange)

  useEffect(() => {
    refetch()
  }, [selectedApp, timeRange])

  // Transform backend data into chart-friendly format
  const chartData = data.slice(0, 5).map((item: { response_code: number; count: number }, index: number) => ({
    name: item.response_code.toString(),
    value: item.count,
    color: COLORS[index % COLORS.length],
  }))

  return (
    <div className="w-full bg-white xl shadow-md p-6">
      <div className="flex justify-between items-center p-3">
        <p className="text-lg">Response Status Code</p>
      </div>
      <div className="flex items-center gap-6">
        <div className="w-full sm:w-60 h-60">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={chartData}
                dataKey="value"
                outerRadius={70}
                innerRadius={45}
                startAngle={90}
                endAngle={-270}
              >
                {chartData.map((entry: { color: string | undefined }, index: any) => (
                  <Cell key={`os-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip formatter={(value: any, _: any, props: any) => [`${value}`, `Code: ${props.payload.name}`]} />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="grid grid-cols-1 gap-x-32 gap-y-2 text-sm text-slate-700">
          {chartData.map((os: { name: string | undefined; color: any; value: string | number }) => (
            <div key={os.name} className="flex items-center gap-32">
              <div className="flex items-center gap-2">
                <span className="w-2.5 h-2.5 flex-shrink-0" style={{ backgroundColor: os.color }} />
                <span className="truncate max-w-[8rem]">{os.name}</span>
              </div>
              <span className="font-semibold ml-auto">{os.value}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
