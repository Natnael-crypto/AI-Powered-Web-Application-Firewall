import { requestData } from '../lib/Constants'
import countries from 'i18n-iso-countries'
import localeRegistration from 'i18n-iso-countries/langs/en.json'

countries.registerLocale(localeRegistration)

function RequestStat({ className }: { className?: string }) {
  const formattedData = Object.entries(requestData)
    .map(([key, value]) => ({
      label: countries.getName(key, 'en'),
      value,
    }))
    .sort((a, b) => b.value - a.value)
    .slice(0, 8)

  const maxValue = Math.max(...formattedData.map(d => d.value))
  const minValue = Math.min(...formattedData.map(d => d.value))
  const range = maxValue - minValue || 1

  const points = formattedData.map((data, index) => {
    const x = (index / (formattedData.length - 1)) * 100
    const y = 100 - ((data.value - minValue) / range) * 100
    return { x, y, ...data }
  })

  const path = points.map((p, i) =>
    `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(2)},${p.y.toFixed(2)}`
  ).join(' ')

  const gradientId = 'line-gradient'

  return (
    <div className={`w-full ${className}`}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-semibold text-gray-800">Top Countries by Request Volume</h3>
        <span className="text-sm text-gray-500 font-medium bg-gray-100 px-3 py-1 rounded-full">
          {formattedData.length} Countries
        </span>
      </div>

      <div className="relative h-[300px] w-full">
        <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="absolute inset-0 w-full h-full">
          {/* Gradient under the line */}
          <defs>
            <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="#3b82f6" stopOpacity="0.3" />
              <stop offset="100%" stopColor="#3b82f6" stopOpacity="0" />
            </linearGradient>
          </defs>

          {/* Filled area under curve */}
          <path
            d={`${path} L 100,100 L 0,100 Z`}
            fill={`url(#${gradientId})`}
          />

          {/* Line path */}
          <path
            d={path}
            fill="none"
            stroke="#3b82f6"
            strokeWidth="1.5"
            strokeLinejoin="round"
          />

          {/* Dots */}
          {points.map((point, index) => (
            <circle
              key={index}
              cx={point.x}
              cy={point.y}
              r="1.5"
              fill="#3b82f6"
              className="hover:scale-125 transition-transform"
            />
          ))}
        </svg>

        {/* Labels & tooltips */}
        {points.map((point, index) => (
          <div
            key={index}
            className="absolute flex flex-col items-center"
            style={{
              left: `${point.x}%`,
              bottom: `${100 - point.y}%`,
              transform: 'translate(-50%, 50%)',
            }}
          >
            <div className="hidden group-hover:block absolute -top-8 whitespace-nowrap bg-white px-2 py-1 text-xs rounded shadow text-gray-800">
              {point.value.toLocaleString()}
            </div>
            <div className="text-[10px] text-gray-600 mt-2 text-center max-w-[80px] truncate">
              {point.label}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-6 flex justify-between text-xs text-gray-500 border-t border-gray-100 pt-4">
        <span>{minValue.toLocaleString()}</span>
        <span>{maxValue.toLocaleString()}</span>
      </div>
    </div>
  )
}

export default RequestStat
