import {MdDashboard, MdNotifications, MdWeb} from 'react-icons/md'
import {FaRegFileAlt} from 'react-icons/fa'
import {AiOutlineTool, AiOutlineSetting} from 'react-icons/ai'
import {HiOutlineShieldCheck} from 'react-icons/hi'

export const SiderbarContentItems = [
  {
    title: 'Dashboard',
    href: 'dashboard',
    icon: MdDashboard,
    children: [],
  },
  {
    title: 'Log',
    href: 'attacks-logs',
    icon: FaRegFileAlt,
    children: [],
  },
  {
    title: 'Custom Rules',
    href: 'custom-rules',
    icon: AiOutlineTool,
    children: [],
  },
  {
    title: 'Web Services',
    href: 'web-services',
    icon: MdWeb,
    children: [],
  },
  {
    title: 'Security Headers',
    href: 'security-headers',
    icon: HiOutlineShieldCheck,
    children: [],
  },
  {
    title: 'Notifications',
    href: 'notifications',
    icon: MdNotifications,
    children: [],
  },
  {
    title: 'System',
    href: '/system',
    icon: AiOutlineSetting,
    children: [],
  },
]

export const requestData: Record<string, number> = {
  US: 120000,
  IN: 80000,
  FR: 300,
  CN: 70000,
  DE: 50000,
  BR: 40000,
  AU: 150000,
  RU: 200000,
}

export const blockedRequestData: Record<string, number> = {
  US: 1200,
  IN: 8000,
  FR: 300,
  CN: 7000,
  DE: 50000,
  BR: 4000,
  AU: 150000,
  RU: 2000,
}

export const API_BASE_URL = 'https://waf-backend-latest.onrender.com'
