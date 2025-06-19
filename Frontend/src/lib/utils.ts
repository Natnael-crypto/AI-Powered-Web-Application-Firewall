import {clsx, type ClassValue} from 'clsx'
import {twMerge} from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const validateIPAddressOrDomain = (input: string): boolean => {
  const ipv4Regex =
    /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/

  const domainRegex = /^(?!:\/\/)([a-zA-Z0-9-_]+\.)+[a-zA-Z]{2,}$/ // e.g., example.com, sub.example.org

  return ipv4Regex.test(input) || domainRegex.test(input)
}

export const validateHostname = (hostname: string): boolean => {
  // Simple hostname validation
  const hostnameRegex =
    /^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$/
  return hostnameRegex.test(hostname)
}

export const validatePort = (port: string): boolean => {
  const portNum = parseInt(port, 10)
  return !isNaN(portNum) && portNum > 0 && portNum <= 65535
}
