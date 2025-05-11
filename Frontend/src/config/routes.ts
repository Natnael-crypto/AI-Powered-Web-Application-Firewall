// src/config/routes.ts
import {UserRole} from '../lib/types'

export const APP_ROUTES = {
  PUBLIC: {
    LOGIN: '/login',
    REGISTER: '/register',
    FORGOT_PASSWORD: '/forgot-password',
    RESET_PASSWORD: '/reset-password',
    ABOUT: '/about',
  },
  PRIVATE: {
    DASHBOARD: {
      path: '/dashboard',
      roles: ['super_admin', 'admin', 'user'] as UserRole[],
    },
    ADMIN_PANEL: {
      path: '/admin',
      roles: ['super_admin', 'admin'] as UserRole[],
    },
    USER_PROFILE: {
      path: '/profile',
      roles: ['super_admin', 'admin', 'user'] as UserRole[],
    },
    SYSTEM_SETTINGS: {
      path: '/settings',
      roles: ['super_admin'] as UserRole[],
    },
  },
  MISC: {
    HOME: '/',
    NOT_FOUND: '/404',
  },
} as const

// Type helpers
type ValueOf<T> = T[keyof T]
type RoutePath =
  | ValueOf<typeof APP_ROUTES.PUBLIC>
  | ValueOf<typeof APP_ROUTES.PRIVATE>['path']

// Navigation guard helpers
export const isPublicRoute = (path: string): boolean => {
  return Object.values(APP_ROUTES.PUBLIC).includes(path as any)
}

// export const hasRouteAccess = (path: string, userRole: UserRole): boolean => {
//   console.log('details: ', userRole)
//   if (isPublicRoute(path)) return true

//   const privateRoute = Object.values(APP_ROUTES.PRIVATE).find(
//     route => route.path === path,
//   )

//   return privateRoute ? privateRoute.roles.includes(userRole) : false
// }

export const getDefaultRoute = (userRole: UserRole): string => {
  if (userRole === 'super_admin') return APP_ROUTES.PRIVATE.SYSTEM_SETTINGS.path
  if (userRole === 'admin') return APP_ROUTES.PRIVATE.ADMIN_PANEL.path
  return APP_ROUTES.PRIVATE.DASHBOARD.path
}
