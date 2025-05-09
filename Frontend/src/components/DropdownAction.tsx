import {Ellipsis} from 'lucide-react'
import {useEffect, useRef, useState} from 'react'

interface DropdownActionsProps<T> {
  item: T
  options: {
    label: string
    onClick: (item: T) => void
  }[]
}

export function DropdownActions<T>({item, options}: DropdownActionsProps<T>) {
  const [open, setOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  return (
    <div className="relative inline-block text-left" ref={menuRef}>
      <button
        onClick={() => setOpen(prev => !prev)}
        className="inline-flex items-center justify-center p-1.5 rounded hover:bg-gray-100 transition-colors duration-150"
        aria-label="Actions"
      >
        <Ellipsis className="text-gray-600" size={20} />
      </button>

      {open && (
        <div className="origin-top-right absolute right-0 mt-1 w-48 rounded-md shadow-md bg-white border border-gray-100 z-10">
          <div className="py-1">
            {options.map((option, index) => (
              <button
                key={index}
                onClick={() => {
                  option.onClick(item)
                  setOpen(false)
                }}
                className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-100"
              >
                {option.label}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
