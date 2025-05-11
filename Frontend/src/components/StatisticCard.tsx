interface StatisticCardProps {
  label: string
  value: number | undefined
  className?: string
}
function StatisticCard({label, value, className}: StatisticCardProps) {
  return (
    <div className={`flex flex-col justify-center pl-5 ${className}`}>
      <p className="text-slate-400 text-md mb-2 text-center">{label}</p>
      <p className="font-normal text-2xl text-center">{value}</p>
    </div>
  )
}

export default StatisticCard
