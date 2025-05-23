import React, {useState, useContext} from 'react'
import {ToastContext} from '../providers/ToastProvider'
import {ToastType} from '../lib/types'
import Card from '../components/Card'
import {
  CheckCircle,
  XCircle,
  AlertTriangle,
  Loader,
  X,
  ChevronsDown,
  ChevronsUp,
} from 'lucide-react'

const MAX_VISIBLE_HEIGHT = 400
const TOAST_HEIGHT = 96

interface ToastStyle {
  bg: string
  border: string
  icon: React.ReactNode
  textColor: string
}

const getToastStyles = (type: ToastType): ToastStyle => {
  const styles: Record<string, ToastStyle> = {
    success: {
      bg: 'bg-green-50',
      border: 'border-green-400',
      icon: <CheckCircle className="w-5 h-5 text-green-600" />,
      textColor: 'text-green-800',
    },
    error: {
      bg: 'bg-red-50',
      border: 'border-red-400',
      icon: <XCircle className="w-5 h-5 text-red-600" />,
      textColor: 'text-red-800',
    },
    warning: {
      bg: 'bg-yellow-50',
      border: 'border-yellow-400',
      icon: <AlertTriangle className="w-5 h-5 text-yellow-600" />,
      textColor: 'text-yellow-800',
    },
    loading: {
      bg: 'bg-blue-50',
      border: 'border-blue-400',
      icon: <Loader className="w-5 h-5 text-blue-600 animate-spin" />,
      textColor: 'text-blue-800',
    },
    default: {
      bg: 'bg-gray-50',
      border: 'border-gray-300',
      icon: <></>,
      textColor: 'text-gray-800',
    },
  }

  return styles[type] || styles.default
}

const ToastItem: React.FC<{
  toast: {
    id: string
    type: ToastType
    message: string
  }
  onDismiss: (id: string) => void
}> = ({toast, onDismiss}) => {
  const style = getToastStyles(toast.type)

  return (
    <Card
      className={`p-4 rounded-xl backdrop-blur-sm transition-all duration-300 transform ${style.bg} ${style.border} ${style.textColor} border`}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center space-x-3">
          <span className="shrink-0">{style.icon}</span>
          <p className="text-sm font-medium leading-tight">{toast.message}</p>
        </div>
        <button
          onClick={() => onDismiss(toast.id)}
          className="ml-2 text-current opacity-70 hover:opacity-100 transition-opacity"
          aria-label="Close notification"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
    </Card>
  )
}

const ToastCollapseButton: React.FC<{
  count: number
  expanded: boolean
  onClick: () => void
}> = ({count, expanded, onClick}) => (
  <button
    onClick={onClick}
    className="flex items-center justify-center gap-1 px-4 py-2.5 bg-white/80 backdrop-blur-md text-gray-700 rounded-xl shadow-md hover:bg-white/95 transition-all duration-200 text-sm font-semibold"
  >
    {expanded ? 'Show less' : `Show ${count} more notification${count !== 1 ? 's' : ''}`}
    {expanded ? <ChevronsUp className="w-4 h-4" /> : <ChevronsDown className="w-4 h-4" />}
  </button>
)

const ToastContainer: React.FC = () => {
  const {toasts, removeToast} = useContext(ToastContext)!
  const [expanded, setExpanded] = useState(false)

  const maxVisibleToasts = Math.floor(MAX_VISIBLE_HEIGHT / TOAST_HEIGHT)
  const shouldCollapse = toasts.length > maxVisibleToasts && !expanded
  const visibleToasts = shouldCollapse ? toasts.slice(0, maxVisibleToasts) : toasts
  const collapsedCount = shouldCollapse ? toasts.length - maxVisibleToasts : 0

  return (
    <div className="fixed top-4 right-4 z-50 w-1/3 max-h-[calc(100vh-2rem)] overflow-scroll">
      <div className="flex flex-col-reverse gap-3">
        {visibleToasts.map(toast => (
          <ToastItem key={toast.id} toast={toast} onDismiss={removeToast} />
        ))}

        {collapsedCount > 0 && (
          <ToastCollapseButton
            count={collapsedCount}
            expanded={expanded}
            onClick={() => setExpanded(!expanded)}
          />
        )}

        {expanded && (
          <ToastCollapseButton
            count={collapsedCount}
            expanded={expanded}
            onClick={() => setExpanded(false)}
          />
        )}
      </div>
    </div>
  )
}

export default ToastContainer
