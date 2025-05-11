import {Outlet, useLocation, useNavigate} from 'react-router-dom'
import {useEffect, useState} from 'react'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'
import {useIsLoggedIn} from '../hooks/api/useUser'
import {APP_ROUTES, getDefaultRoute, isPublicRoute} from '../config/routes'
import LoadingSpinner from '../components/LoadingSpinner'
import {useToast} from '../hooks/useToast'
import {useUserInfo} from '../store/UserInfo'

// Retry configuration
const MAX_RETRIES = 3
const RETRY_DELAY = 1000 // 1 second

function RootLayout() {
  const {pathname} = useLocation()
  const navigate = useNavigate()
  const {addToast: toast} = useToast()
  const {user, setUser} = useUserInfo()
  const [retryCount, setRetryCount] = useState(0)
  const {data: userInfo, isLoading, error, refetch} = useIsLoggedIn()

  useEffect(() => {
    if (userInfo) {
      setUser(userInfo)
    }
  }, [userInfo, setUser])

  useEffect(() => {
    const token = localStorage.getItem('token')

    if (!token) {
      navigate(APP_ROUTES.PUBLIC.LOGIN, {replace: true})
      return
    }

    // If we have an error and haven't exceeded max retries
    if (error && retryCount < MAX_RETRIES) {
      const timer = setTimeout(() => {
        refetch()
        setRetryCount(prev => prev + 1)
      }, RETRY_DELAY)

      return () => clearTimeout(timer)
    }

    if (userInfo || retryCount >= MAX_RETRIES) {
      handleAuthRouting()
    }
  }, [pathname, userInfo, error, retryCount])

  const handleAuthRouting = () => {
    if (isLoading || !userInfo) return

    if (error) {
      toast('Session Error: Unable to verify your session. Please login again.')
      localStorage.removeItem('token')
      navigate(APP_ROUTES.PUBLIC.LOGIN, {replace: true})
      return
    }

    if (isPublicRoute(pathname)) {
      navigate(getDefaultRoute(userInfo.role), {replace: true})
      return
    }
  }

  // Loading states
  if (isLoading || (!userInfo && !error && retryCount < MAX_RETRIES)) {
    return <LoadingSpinner fullScreen />
  }

  if (error && retryCount >= MAX_RETRIES) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <h2 className="text-xl font-bold mb-2">Session Error</h2>
          <p className="mb-4">Unable to verify your session. Please try again later.</p>
          <button
            onClick={() => {
              localStorage.removeItem('token')
              navigate(APP_ROUTES.PUBLIC.LOGIN)
            }}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Go to Login
          </button>
        </div>
      </div>
    )
  }

  if (!userInfo) {
    return null
  }

  return (
    <div className="h-screen w-full flex bg-gradient-to-r from-slate-100 to-slate-50">
      <Sidebar />
      <div className="flex flex-col w-full overflow-hidden">
        <Navbar />
        <main
          className="flex-1 overflow-y-auto p-6 bg-[#F3F6FE]"
          style={{scrollBehavior: 'smooth'}}
        >
          <Outlet context={{user: userInfo}} />
        </main>
      </div>
    </div>
  )
}

export default RootLayout
