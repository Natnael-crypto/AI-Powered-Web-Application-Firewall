import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from 'react-router-dom'
import RootLayout from './layout/RootLayout'
import Dashboard from './pages/Dashboard'
import LoginPage from './pages/Login'
import PageNotFound from './pages/PageNotFound'
import CustomeRules from './pages/CustomeRules'
import System from './pages/System'
import AttackLog from './pages/Requestlogs/AttackLog'
import RateLimiting from './pages/Requestlogs/RateLimiting'
import AntiBot from './pages/Requestlogs/Antibot'
import WebServices from './pages/WebServices'
import {ToastProvider} from './providers/ToastProvider'

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<RootLayout />}>
        <Route index element={<Dashboard />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="log">
          <Route path="attacks" element={<AttackLog />} />
          <Route path="limits" element={<RateLimiting />} />
          <Route path="captcha" element={<AntiBot />} />
        </Route>
        <Route path="custom-rules" element={<CustomeRules />} />
        <Route path="web-services" element={<WebServices />} />
        <Route path="system" element={<System />} />
        <Route path="*" element={<PageNotFound />} />
      </Route>
    </Route>,
  ),
)

function App() {
  return (
    <ToastProvider>
      <RouterProvider router={router} />
    </ToastProvider>
  )
}

export default App
