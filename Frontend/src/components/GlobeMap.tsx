import React from 'react'
import {VectorMap} from '@react-jvectormap/core'
import {worldMill} from '@react-jvectormap/world'
import {requestData} from '../lib/Constants'

const GlobeMap: React.FC = () => {
  return (
    <div className="h-full flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-gray-800">Global Traffic Map</h2>
        <div className="text-xs font-medium text-gray-500 bg-gray-100 px-3 py-1 ull">
          Live Data
        </div>
      </div>

      <div className="flex-1  overflow-hidden">
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
                values: requestData,
                scale: ['#E2E8F0', '#1E40AF'],
                normalizeFunction: 'linear',
                attribute: 'fill',
              },
            ],
          }}
          onRegionTipShow={(_, el: any, code) => {
            const requests = requestData[code] || 0
            el.html(
              `<div style="font-family: 'Oxygen', sans-serif; padding: 6px;">
                <div style="font-weight: 600; color: #1F2937; margin-bottom: 4px;">${el.html()}</div>
                <div style="font-size: 12px; color: #4B5563;">
                  Requests: <span style="font-weight: 600; color: #1E40AF;">${requests.toLocaleString()}</span>
                </div>
              </div>`,
            )
          }}
          map={worldMill}
        />
      </div>
    </div>
  )
}

export default GlobeMap
