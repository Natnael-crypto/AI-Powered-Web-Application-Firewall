import React from 'react'
import {VectorMap} from '@react-jvectormap/core'
import {worldMill} from '@react-jvectormap/world'
import {requestData} from '../lib/Constants'
import RequestStat from './RequestStat'

const GlobeMap: React.FC = () => {
  return (
    <div className="w-full  flex flex-col bg-white shadow-lg rounded-xl p-6 gap-6">
      <div className="w-full  h-96 rounded-lg shadow-xl overflow-hidden">
        <VectorMap
          backgroundColor="#f4f4f4"
          className="h-full w-full"
          zoomOnScroll={false}
          regionStyle={{
            initial: {fill: '#D3D3D3', stroke: '#fff', strokeWidth: 1},
            hover: {fill: '#228B22', cursor: 'pointer'},
          }}
          series={{
            regions: [
              {
                values: requestData,
                scale: ['#D3D3D3', '#006400'],
                normalizeFunction: 'linear',
                attribute: 'fill',
              },
            ],
          }}
          onRegionTipShow={(_, el: any, code) => {
            const requests = requestData[code] || 0
            el.html(
              `<strong>${el.html()}</strong><br>
              Requests: <b>${requests.toLocaleString()}</b>`,
            )
          }}
          map={worldMill}
        />
      </div>
      <div className="w-full flex flex-col gap-6">
        <RequestStat className="w-full" />
      </div>
    </div>
  )
}

export default GlobeMap
