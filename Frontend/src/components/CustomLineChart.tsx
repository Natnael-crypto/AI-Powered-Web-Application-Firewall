import React from 'react'
import {LineChart, Line, ResponsiveContainer, Tooltip} from 'recharts'
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
    <Card className="flex flex-col w-full h-1/2 p-4 space-y-4 bg-gray-200">
      <div className="flex justify-between items-center">
        <h3 className="text-xl font-semibold text-gray-800 flex flex-nowrap">{title}</h3>
        {max !== undefined && (
          <p className="text-sm text-gray-600 flex-nowrap flex gap-3">
            Max: <span className="font-bold text-gray-800">{max}</span>
          </p>
        )}
      </div>

      <div className="flex-1">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <Tooltip
              contentStyle={{
                backgroundColor: '#fff',
                borderRadius: '0.5rem',
                borderColor: '#e5e7eb',
              }}
              labelStyle={{color: '#4B5563', fontSize: '0.875rem'}}
              itemStyle={{color: '#6B7280', fontSize: '0.875rem'}}
            />
            <Line
              type="monotone"
              dataKey="uv"
              stroke="#6366F1"
              strokeWidth={2.5}
              dot={false}
              activeDot={{r: 4, strokeWidth: 2, fill: '#6366F1', stroke: '#fff'}}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  )
}

export default CustomLineChart
