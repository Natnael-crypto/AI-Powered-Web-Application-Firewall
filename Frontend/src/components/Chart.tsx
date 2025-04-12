import Card from './Card'
interface detailType {
  label: String | undefined
  value: number
}
interface chartProp {
  detailStats: detailType[]
}

function Chart({detailStats}: chartProp) {
  return (
    <Card className="flex">
      {/* pie chart */}
      <div></div>

      {/* detail stats */}
      <div className="flex flex-col flex-wrap">
        {detailStats.map(detail => (
          <div className="flex  justify-between gap-12 items-center">
            <div>{detail.label}</div>
            <div>{detail.value}</div>
          </div>
        ))}
      </div>
    </Card>
  )
}

export default Chart
