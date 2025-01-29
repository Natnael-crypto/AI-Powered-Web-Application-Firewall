import clsx from 'clsx'
import {ReactNode} from 'react'

interface CardProps {
  size?: 's' | 'm' | 'l'
  className?: string
  children?: ReactNode
}
function Card({className, children, size = 'm'}: CardProps) {
  const baseClass = 'rounded-lg shadow-lg w-full'
  const sizeClasses = {
    s: 'px-3 py-2',
    m: 'px-5 py-3',
    l: 'px-6 py-4',
  }

  const finalClassName = clsx(baseClass, sizeClasses[size], className)
  return <div className={finalClassName}>{children}</div>
}

export default Card
