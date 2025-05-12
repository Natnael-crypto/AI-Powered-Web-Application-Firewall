import React, {useEffect, useCallback} from 'react'
import {MdClose} from 'react-icons/md'

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  children: React.ReactNode
  title?: string
  closeOnOutsideClick?: boolean
}

const Modal: React.FC<ModalProps> = ({
  isOpen,
  onClose,
  children,
  title,
  closeOnOutsideClick = true,
}) => {
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    },
    [onClose],
  )

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
      document.addEventListener('keydown', handleKeyDown)
    } else {
      document.body.style.overflow = ''
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      document.body.style.overflow = ''
    }
  }, [isOpen, handleKeyDown])

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (closeOnOutsideClick && e.target === e.currentTarget) {
      onClose()
    }
  }

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm transition-opacity"
      onClick={handleBackdropClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? 'modal-title' : undefined}
    >
      <div className="relative bg-white dark:bg-zinc-900 rounded-2xl border border-gray-200 dark:border-zinc-700 shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 animate-fadeIn scale-95 sm:scale-100 transition-transform">
        <div className="flex items-center justify-between">
          {title && (
            <h3
              id="modal-title"
              className="text-xl font-semibold tracking-tight text-gray-900 dark:text-white"
            >
              {title}
            </h3>
          )}
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-red-500 transition-colors duration-200 focus:outline-none"
            aria-label="Close modal"
          >
            <MdClose size={24} />
          </button>
        </div>
        <div className="mt-5">{children}</div>
      </div>
    </div>
  )
}

export default Modal
