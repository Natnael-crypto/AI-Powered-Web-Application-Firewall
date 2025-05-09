import React, {useEffect, useCallback} from 'react'
import {MdCancel} from 'react-icons/md'

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

  console.log('modal: ', isOpen)

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 transition-opacity"
      onClick={handleBackdropClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby={title ? 'modal-title' : undefined}
    >
      <div className="bg-white rounded-lg shadow-2xl w-11/12 max-w-2xl max-h-[90vh] overflow-y-auto p-6 animate-scaleFadeIn">
        <div className="flex items-start justify-between">
          {title && (
            <h3 id="modal-title" className="text-lg font-semibold text-gray-800">
              {title}
            </h3>
          )}
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            aria-label="Close modal"
          >
            <MdCancel size={22} />
          </button>
        </div>
        <div className="mt-4">{children}</div>
      </div>
    </div>
  )
}

export default Modal
