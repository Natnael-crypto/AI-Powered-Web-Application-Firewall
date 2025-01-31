import {BsLightning} from 'react-icons/bs'
import {LuLogOut} from 'react-icons/lu'
import {useLocation} from 'react-router-dom'

function Navbar() {
  const location = useLocation()
  const urls = location.pathname.split('/').filter(url => url !== '')

  return (
    <div className="w-full h-max bg-inherit shadow-sm">
      <div className="flex justify-between items-center py-4 px-8 ">
        {/* Breadcrumbs */}
        <div className="flex gap-2 items-center">
          {urls.map((url, index) => (
            <div key={index} className="flex items-center">
              <p className="text-lg font-medium text-gray-700 capitalize hover:text-blue-600 transition-colors duration-200">
                {url}
              </p>
              {index !== urls.length - 1 && <span className="mx-2 text-gray-400">/</span>}
            </div>
          ))}
        </div>

        {/* Right Section */}
        <div className="flex items-center gap-5">
          {/* Pro Badge */}
          <div className="flex items-center overflow-hidden rounded-lg shadow-sm">
            <div className="flex gap-2 items-center text-lg bg-gradient-to-r from-blue-600 to-indigo-600 text-white px-4 py-2">
              <BsLightning className="text-yellow-400" />
              <span>Pro</span>
            </div>
            <div className="bg-gray-100 text-gray-700 px-4 py-2 text-lg font-medium">
              233 days left
            </div>
          </div>

          {/* Logout Button */}
          <button className="p-2 rounded-full hover:bg-gray-100 transition-colors duration-200">
            <LuLogOut
              size={25}
              className="text-gray-700 hover:text-red-600 transition-colors duration-200"
            />
          </button>
        </div>
      </div>
    </div>
  )
}

export default Navbar
