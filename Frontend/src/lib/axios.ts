import axios from 'axios'

const axiosInstance = axios.create()

const excludedEndpoints = ['/login', '/register']

axiosInstance.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')?.trim()

    if (token && !excludedEndpoints.some(endpoint => config.url?.includes(endpoint))) {
      config.headers.Authorization = `${token}`
    }

    return config
  },
  error => {
    return Promise.reject(error)
  },
)

axiosInstance.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 403) {
      console.error('Forbidden - check your authentication')
    }
    return Promise.reject(error)
  },
)

export default axiosInstance
