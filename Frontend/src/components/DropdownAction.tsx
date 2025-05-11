import {Ellipsis} from 'lucide-react'
import {useEffect, useRef, useState} from 'react'

interface DropdownActionsProps<T> {
  item: T
  options: {
    label: string
    onClick: (item: T) => void
    show?: (item: T) => boolean
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

    if (open) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [open])

  // Filter options based on show function
  const visibleOptions = options.filter(option =>
    option.show ? option.show(item) : true,
  )

  return (
    <div className="relative text-left" ref={menuRef}>
      <button
        onClick={() => setOpen(prev => !prev)}
        className="inline-flex items-center justify-center p-1.5 rounded hover:bg-gray-100 transition-colors duration-150"
        aria-label="Actions"
      >
        <Ellipsis className="text-gray-600" size={20} />
      </button>

      {open && visibleOptions.length > 0 && (
        <div className="absolute origin-top-right right-0 mt-1 rounded-md shadow-md bg-white border border-gray-100 z-50">
          {visibleOptions.map((option, index) => (
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
      )}
    </div>
  )
}
