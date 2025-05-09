import {useState, useEffect} from 'react'

type ToastProps = {
  message: string
  type?: 'success' | 'error' | 'info' | 'warning'
  duration?: number
  onClose?: () => void
  position?:
    | 'top-left'
    | 'top-center'
    | 'top-right'
    | 'bottom-left'
    | 'bottom-center'
    | 'bottom-right'
}

export const Toast = ({
  message,
  type = 'info',
  duration = 3000,
  onClose,
  position = 'top-right',
}: ToastProps) => {
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false)
      onClose?.()
    }, duration)

    return () => clearTimeout(timer)
  }, [duration, onClose])

  if (!isVisible) return null

  const typeClasses = {
    success: 'bg-green-100 text-green-700 border-green-300',
    error: 'bg-red-100 text-red-700 border-red-300',
    info: 'bg-blue-100 text-blue-700 border-blue-300',
    warning: 'bg-yellow-100 text-yellow-700 border-yellow-300',
  }

  const positionClasses = {
    'top-left': 'top-4 left-4',
    'top-center': 'top-4 left-1/2 transform -translate-x-1/2',
    'top-right': 'top-4 right-4',
    'bottom-left': 'bottom-4 left-4',
    'bottom-center': 'bottom-4 left-1/2 transform -translate-x-1/2',
    'bottom-right': 'bottom-4 right-4',
  }

  return (
    <div
      className={`fixed right-3 top-3  z-1000 p-4 rounded-md border ${typeClasses[type]} shadow-lg animate-fade-in`}
    >
      {message}
    </div>
  )
}
