function TodaySummary() {
  return (
    <div className="flex flex-col gap-8  px-5 py-2 bg-white shadow-lg min-w-[250px]">
      <p className="text-gray-600 block whitespace-nowrap">Today Summary</p>
      <div className="flex flex-col gap-3">
        <p className="text-gray-500 text-lg">Page View</p>
        <p className="text-3xl font-bold">225</p>
      </div>
      <div className="flex flex-col gap-3">
        <p className="text-gray-500 text-lg">Page View</p>
        <p className="text-3xl font-bold">225</p>
      </div>
      <div className="flex flex-col gap-3">
        <p className="text-gray-500 text-lg">Page View</p>
        <p className="text-3xl font-bold">225</p>
      </div>
    </div>
  )
}

export default TodaySummary
