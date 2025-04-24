import React from 'react'
import {VectorMap} from '@react-jvectormap/core'
import {worldMill} from '@react-jvectormap/world'
import {requestData} from '../lib/Constants'
import RequestStat from './RequestStat'

const GlobeMap: React.FC = () => {
  return (
    <div className="w-full h-full flex justify-between bg-white">
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
              values: requestData, // { [countryCode]: number }
              scale: ['#D3D3D3', '#006400'], // from light gray to dark green
              normalizeFunction: 'linear', // or 'polynomial', 'log'
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
      <RequestStat />
    </div>
  )
}

export default GlobeMap
