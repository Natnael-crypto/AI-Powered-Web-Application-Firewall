import {LuLogOut} from 'react-icons/lu'
import {useLocation, useNavigate} from 'react-router-dom'
import ToastExample from './ToastExample'

function Navbar() {
  const location = useLocation()
  const urls = location.pathname.split('/').filter(url => url !== '')
  const navigate = useNavigate()

  const handleLogout = () => {
    localStorage.removeItem('token')
    navigate('/login')
  }

  return (
    <header className="w-full border-b border-gray-200 bg-white shadow-sm">
      <div className="flex justify-between items-center py-4 px-6">
        {/* Breadcrumbs */}
        <nav className="flex gap-2 items-center text-gray-600 text-sm font-medium">
          {urls.map((url, index) => (
            <div key={index} className="flex items-center">
              <span className="capitalize hover:text-blue-600 transition-colors duration-200">
                {url.replace(/-/g, ' ')}
              </span>
              {index !== urls.length - 1 && (
                <span className="mx-2 text-gray-300 font-light">/</span>
              )}
            </div>
          ))}
        </nav>

        {/* Logout */}
        <button
          onClick={handleLogout}
          className="p-2 ull hover:bg-red-50 transition-colors duration-200"
          title="Logout"
        >
          <LuLogOut
            size={22}
            className="text-gray-600 hover:text-red-600 transition-colors duration-200"
          />
        </button>
      </div>
    </header>
  )
}

export default Navbar
