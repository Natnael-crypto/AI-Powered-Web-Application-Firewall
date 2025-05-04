import React from 'react'
import {LineChart, Line, ResponsiveContainer, Tooltip, CartesianGrid, XAxis, YAxis} from 'recharts'
import Card from './Card'

type ChartData = {
  name: string
  uv: number
}

type LineChartProps = {
  data: ChartData[]
  title?: string
  max?: number
}

const CustomLineChart: React.FC<LineChartProps> = ({data, title = 'Line Chart', max}) => {
  return (
    <Card className="flex flex-col w-full h-1/2 p-4 space-y-3 bg-white rounded-md shadow-sm border-2">
      <div className="flex justify-between items-center">
        <h3 className="text-base font-medium text-gray-700">{title}</h3>
        {max !== undefined && (
          <p className="text-xs text-gray-500 flex items-center gap-1">
            Max: <span className="font-semibold text-gray-700">{max.toLocaleString()}</span>
          </p>
        )}
      </div>

      <div className="flex-1 min-h-[120px]">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data} margin={{top: 5, right: 5, left: -20, bottom: 5}}>
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
            <XAxis 
              dataKey="name" 
              axisLine={false} 
              tickLine={false} 
              tick={{fontSize: 10, fill: '#6B7280'}}
            />
            <YAxis 
              axisLine={false} 
              tickLine={false} 
              tick={{fontSize: 10, fill: '#6B7280'}}
              width={30}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: '#fff',
                borderRadius: '0.375rem',
                borderColor: '#e5e7eb',
                boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
                padding: '8px 12px',
              }}
              labelStyle={{color: '#374151', fontSize: '0.75rem', fontWeight: 600, marginBottom: '4px'}}
              itemStyle={{color: '#4B5563', fontSize: '0.75rem'}}
              formatter={(value: number) => [value.toLocaleString(), 'Value']}
            />
            <Line
              type="monotone"
              dataKey="uv"
              stroke="#3B82F6"
              strokeWidth={2}
              dot={false}
              activeDot={{r: 4, strokeWidth: 2, fill: '#3B82F6', stroke: '#fff'}}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  )
}

export default CustomLineChart
