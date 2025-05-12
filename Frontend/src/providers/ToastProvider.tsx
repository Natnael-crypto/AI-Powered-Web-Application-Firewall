// components/ToastProvider.tsx
import React, {createContext, useState, useCallback} from 'react'
import {Toast, ToastContextType, ToastType} from '../lib/types'
import ToastContainer from '../components/ToastContainer'

export const ToastContext = createContext<ToastContextType | undefined>(undefined)

export const ToastProvider: React.FC<{children: React.ReactNode}> = ({children}) => {
  const [toasts, setToasts] = useState<Toast[]>([])

  const addToast = useCallback(
    (message: string, type: ToastType = 'info', duration = 8000) => {
      const id = Math.random().toString(36).substring(2, 9)
      const newToast: Toast = {id, message, type, createdAt: Date.now(), duration}

      setToasts(prev => [newToast, ...prev])

      if (duration > 0) {
        setTimeout(() => {
          setToasts(prev => prev.filter(toast => toast.id !== id))
        }, duration)
      }
    },
    [],
  )

  const removeToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id))
  }, [])

  const clearToasts = useCallback(() => {
    setToasts([])
  }, [])

  return (
    <ToastContext.Provider value={{toasts, addToast, removeToast, clearToasts}}>
      {children}
      <ToastContainer />
    </ToastContext.Provider>
  )
}
