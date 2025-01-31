interface StatisticCardProps {
  label: string
  value: number
  className?: string
}
function StatisticCard({label, value, className}: StatisticCardProps) {
  return (
    <div className={`gap-y-5 flex flex-col justify-center pl-5 ${className}`}>
      <p className="text-slate-400 text-lg">{label}</p>
      <p className="font-bold text-4xl">{value}</p>
    </div>
  )
}

export default StatisticCard
