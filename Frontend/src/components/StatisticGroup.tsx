import React from 'react'
import StatisticCard from './StatisticCard'

interface stats {
  label: string
  value: number | undefined
}
interface StatisticGroupProps {
  stats: stats[]
  className?: string
}
function StatisticGroup({stats, className}: StatisticGroupProps) {
  return (
    <div className={`flex justify-between px-5 items-center ${className}`}>
      {stats.map((stat, index) => (
        <React.Fragment key={index}>
          <StatisticCard label={stat.label} value={stat.value} />
          {index < stats.length - 1 && <div className="h-16 w-[1px] bg-slate-600" />}
        </React.Fragment>
      ))}
    </div>
  )
}

export default StatisticGroup
