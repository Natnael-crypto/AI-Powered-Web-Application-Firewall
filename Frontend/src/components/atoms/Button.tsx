import clsx from 'clsx'
import {ReactNode} from 'react'

interface ButtonProps {
  size?: 's' | 'm' | 'l'
  variant?: 'primary' | 'secondary' | 'outline'
  classname?: string
  disabled?: boolean
  onClick?: () => void
  children?: ReactNode
}
function Button({
  size = 'm',
  variant,
  classname,
  disabled,
  onClick,
  children,
}: ButtonProps) {
  const baseClass = 'rounded-lg  focus:outline-none font-medium '

  const sizeClasses = {
    s: 'px-2 py-1',
    m: 'px-3 py-1',
    l: 'px-5 py-2',
  }
  const variantClasses = {
    primary: 'bg-blue-500 hover:bg-blue-600',
    secondary: 'bg-gray-500 hover:bg-gray-600',
    outline: 'bg-white text-black hover:bg-gray-200 border border-1 border-black',
  }

  const finalClassName = clsx(
    baseClass,
    sizeClasses[size],
    variant && variantClasses[variant],
    disabled && 'cursor-not-allowed opacity-50',
    classname,
  )
  return (
    <button className={finalClassName} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  )
}

export default Button
