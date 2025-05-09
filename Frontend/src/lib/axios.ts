import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: 'https://waf-backend-latest.onrender.com',
  withCredentials: true,
  headers: {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,PATCH,OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
  },
})

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

// Add response interceptor to handle errors
axiosInstance.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 403) {
      console.error('Forbidden - check your authentication')
      // Handle token expiration or invalid tokens here
    }
    return Promise.reject(error)
  },
)

export default axiosInstance
