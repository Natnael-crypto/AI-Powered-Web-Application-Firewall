import {useState} from 'react'
import {motion} from 'framer-motion'
import {UserIcon, LockIcon, LogInIcon} from 'lucide-react'
import {useNavigate} from 'react-router-dom'
import {useLogin} from '../hooks/api/useUser'

const LoginPage = () => {
  const [formData, setFormData] = useState({username: '', password: ''})
  const [errors, setErrors] = useState<{username?: string; password?: string}>({})
  const [errorMessage, setErrorMessage] = useState('')
  const {mutate} = useLogin()
  const navigate = useNavigate()

  const validateForm = () => {
    const newErrors: {username?: string; password?: string} = {}
    if (!formData.username) newErrors.username = 'Username is required'
    if (!formData.password) newErrors.password = 'Password is required'
    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: {preventDefault: () => void}) => {
    e.preventDefault()
    if (!validateForm()) return
    try {
      mutate(formData, {
        onSuccess: data => {
          const token = data.token
          localStorage.setItem('token', token)
          navigate('/dashboard')
        },
        onError: () => {
          setErrorMessage('Invalid username or password')
        },
      })
    } catch (error) {
      setErrorMessage('Invalid username or password')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <motion.div
        initial={{opacity: 0, y: 20}}
        animate={{opacity: 1, y: 0}}
        className="bg-white p-12  shadow-xl w-full max-w-2xl mx-4"
      >
        <motion.div
          initial={{scale: 0.95}}
          animate={{scale: 1}}
          className="flex justify-center mb-12"
        >
          <div className="h-24 w-24 bg-gray-900 xl flex items-center justify-center shadow-lg">
            <LogInIcon className="h-12 w-12 text-white" />
          </div>
        </motion.div>

        <h1 className="text-4xl font-bold text-center text-gray-900 mb-8">
          Welcome Back
        </h1>

        {errorMessage && (
          <motion.p
            initial={{opacity: 0}}
            animate={{opacity: 1}}
            className="text-red-500 text-center mb-6 p-4 bg-red-50  text-lg"
          >
            {errorMessage}
          </motion.p>
        )}

        <form onSubmit={handleSubmit} className="space-y-8">
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <UserIcon className="h-6 w-6 text-gray-400" />
            </div>
            <input
              type="text"
              value={formData.username}
              onChange={e => setFormData({...formData, username: e.target.value})}
              className="pl-12 w-full px-6 py-4 text-lg border border-gray-200  bg-gray-50 focus:bg-white focus:outline-none focus:ring-2 focus:ring-gray-400 focus:border-transparent transition-all"
              placeholder="Username"
            />
            {errors.username && (
              <motion.p
                initial={{opacity: 0}}
                animate={{opacity: 1}}
                className="text-red-500 text-base mt-2 ml-2"
              >
                {errors.username}
              </motion.p>
            )}
          </div>

          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <LockIcon className="h-6 w-6 text-gray-400" />
            </div>
            <input
              type="password"
              value={formData.password}
              onChange={e => setFormData({...formData, password: e.target.value})}
              className="pl-12 w-full px-6 py-4 text-lg border border-gray-200  bg-gray-50 focus:bg-white focus:outline-none focus:ring-2 focus:ring-gray-400 focus:border-transparent transition-all"
              placeholder="Password"
            />
            {errors.password && (
              <motion.p
                initial={{opacity: 0}}
                animate={{opacity: 1}}
                className="text-red-500 text-base mt-2 ml-2"
              >
                {errors.password}
              </motion.p>
            )}
          </div>

          <div className="space-y-6 pt-4">
            <motion.button
              whileHover={{scale: 1.01}}
              whileTap={{scale: 0.99}}
              type="submit"
              className="w-full bg-gray-900 text-white py-4  text-lg font-medium shadow-lg hover:bg-gray-800 transition-all flex items-center justify-center gap-3"
            >
              <>
                <LogInIcon className="h-6 w-6" />
                Sign In
              </>
            </motion.button>
          </div>
        </form>
      </motion.div>
    </div>
  )
}

export default LoginPage
