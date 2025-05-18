import {Outlet, useLocation, useNavigate} from 'react-router-dom'
import {useEffect, useState} from 'react'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'
import {useIsLoggedIn} from '../hooks/api/useUser'
import {APP_ROUTES, getDefaultRoute, isPublicRoute} from '../config/routes'
import LoadingSpinner from '../components/LoadingSpinner'
import {useToast} from '../hooks/useToast'
import {useUserInfo} from '../store/UserInfo'
import {useGetNotifications} from '../hooks/api/useNotifications'
import {Notification} from '../lib/types'

const MAX_RETRIES = 3
const RETRY_DELAY = 1000
const NOTIFICATION_POLL_INTERVAL = 30000

function RootLayout() {
  const {pathname} = useLocation()
  const navigate = useNavigate()
  const {addToast: toast} = useToast()
  const {setUser, user} = useUserInfo()
  const [retryCount, setRetryCount] = useState(0)
  const {data: userInfo, isLoading, error, refetch} = useIsLoggedIn()
  const [lastNotificationTime, setLastNotificationTime] = useState<string | null>(null)
  const [notificationPolling, setNotificationPolling] = useState<NodeJS.Timeout | null>(
    null,
  )

  const {data: notifications, refetch: refetchNotifications} = useGetNotifications(
    user?.user_id,
  )

  useEffect(() => {
    if (!notifications?.length || !userInfo) return

    const newNotifications = lastNotificationTime
      ? notifications.filter((n: Notification) => !n.status)
      : notifications.slice(0, 1)

    if (newNotifications.length > 0) {
      newNotifications.forEach((notification: any) => {
        toast(notification.message, notification.status ? 'success' : 'warning', 5000)
      })
      setLastNotificationTime(newNotifications[0].timestamp)
    }
  }, [notifications, lastNotificationTime, toast, userInfo])

  useEffect(() => {
    if (userInfo) {
      const interval = setInterval(() => {
        refetchNotifications()
      }, NOTIFICATION_POLL_INTERVAL)
      setNotificationPolling(interval)
    } else if (notificationPolling) {
      clearInterval(notificationPolling)
      setNotificationPolling(null)
    }

    return () => {
      if (notificationPolling) {
        clearInterval(notificationPolling)
      }
    }
  }, [userInfo, refetchNotifications])

  useEffect(() => {
    if (userInfo) {
      setUser(userInfo)
      if (notifications?.length) {
        setLastNotificationTime(notifications[0].timestamp)
      }
    }
  }, [userInfo, setUser, notifications])

  useEffect(() => {
    const token = localStorage.getItem('token')

    if (!token) {
      navigate(APP_ROUTES.PUBLIC.LOGIN, {replace: true})
      return
    }

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
      toast('Session Error: Unable to verify your session. Please login again.', 'error')
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
