import {useEffect, useState} from 'react'
import {VectorMap} from '@react-jvectormap/core'
import {worldMill} from '@react-jvectormap/world'
import {useMapState} from '../hooks/api/useDashboardStat'
import {whereCountry} from 'iso-3166-1' // Import the iso-3166-1 package

interface GlobeMapProps {
  selectedApp: string
  timeRange: any
}

const GlobeMap = ({selectedApp, timeRange}: GlobeMapProps) => {
  const [filter, setFilter] = useState<'all' | 'blocked'>('all')
  const threat = filter === 'all' ? '' : 'blocked'
  const {data, refetch} = useMapState(selectedApp, timeRange, threat)
  const [country, setCountry] = useState<any>({})

  useEffect(() => {
    console.log(threat, '--:--', filter)
    refetch()
  }, [filter])

  useEffect(() => {
    const updatedCountryData: any = {}
    if (data) {
      Object.keys(data).forEach(countryName => {
        const countryCode = whereCountry(countryName)
        if (countryCode) {
          updatedCountryData[countryCode.alpha2] = data[countryName]
        }
      })
    }
    setCountry(updatedCountryData)
  }, [data])

  return (
    <div className="h-full flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg text-gray-800">Global Traffic Map</h2>
        <div className="flex gap-2 items-center">
          <select
            value={filter}
            onChange={e => setFilter(e.target.value as 'all' | 'blocked')}
            className="text-xs font-medium text-gray-700 bg-gray-100 px-2 py-1 rounded border border-gray-300 focus:outline-none"
          >
            <option value="all">All Requests</option>
            <option value="blocked">Blocked Requests</option>
          </select>
          <div className="text-xs font-medium text-gray-500 bg-gray-100 px-3 py-1 rounded-full">
            Live Data
          </div>
        </div>
      </div>

      <div className="flex-1 rounded-lg overflow-hidden">
        <VectorMap
          backgroundColor="#f8fafc"
          className="h-full w-full"
          zoomOnScroll={false}
          regionStyle={{
            initial: {fill: '#E2E8F0', stroke: '#fff', strokeWidth: 1},
            hover: {fill: '#3B82F6', cursor: 'pointer'},
          }}
          series={{
            regions: [
              {
                values: country,
                scale: ['#E2E8F0', '#1E40AF'],
                normalizeFunction: 'linear',
                attribute: 'fill',
              },
            ],
          }}
          onRegionTipShow={(_, el: any, code) => {
            const requests = country[code] || 0
            el.html(`
              <div style="
                font-family: 'Inter', sans-serif;
                padding: 12px;
                background: white;
                color: #1F2937;
                border-radius: 8px;
                box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
                border: 1px solid #E5E7EB;
                min-width: 180px;
              ">
                <div style="
                  display: flex;
                  align-items: center;
                  margin-bottom: 8px;
                  padding-bottom: 8px;
                  border-bottom: 1px solid #F3F4F6;
                ">
                  <div style="
                    width: 24px;
                    height: 24px;
                    margin-right: 8px;
                    background-color: #3B82F6;
                    border-radius: 4px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    color: white;
                    font-weight: bold;
                    font-size: 12px;
                  ">
                    ${code}
                  </div>
                  <div style="font-weight: 600; font-size: 14px;">${el.html()}</div>
                </div>
                <div style="display: flex; justify-content: space-between; font-size: 13px;">
                  <span style="color: #6B7280;">Requests:</span>
                  <span style="font-weight: 600; color: ${filter === 'blocked' ? '#EF4444' : '#10B981'}">
                    ${requests.toLocaleString()}
                  </span>
                </div>
                ${
                  filter === 'blocked'
                    ? `
                <div style="
                  margin-top: 8px;
                  padding-top: 8px;
                  border-top: 1px solid #F3F4F6;
                  font-size: 12px;
                  color: #6B7280;
                  display: flex;
                  align-items: center;
                ">
                  <span style="
                    display: inline-block;
                    width: 8px;
                    height: 8px;
                    background-color: #EF4444;
                    border-radius: 50%;
                    margin-right: 6px;
                  "></span>
                  Blocked requests
                </div>
                `
                    : ''
                }
              </div>
            `)
          }}
          map={worldMill}
        />
      </div>
    </div>
  )
}

export default GlobeMap
