export const SiderbarContentItems = [
  {
    title: 'Statistics',
    href: 'statistics/dashboard',
    children: [
      {
        title: 'Dashboard',
        href: 'statistics/dashboard',
      },
    ],
  },
  {
    title: 'Logs',
    href: 'logs/attacks',
    children: [
      {
        title: 'Attacks',
        href: 'logs/attacks',
      },
      {
        title: 'Rate Limiting',
        href: '/limits',
      },
      {
        title: 'Anti-Bot',
        href: '/captcha',
      },
      {
        title: 'Waiting Room',
        href: '/waiting_room',
      },
    ],
  },
  {
    title: 'Application',
    href: 'application/applications',
    children: [
      {
        title: 'applications',
        href: 'application/applications',
      },
      {
        title: 'SSL Cert',
        href: '/cert',
      },
      {
        title: 'Global Settings',
        href: '/config',
      },
    ],
  },
  {
    title: 'Protections',
    href: 'protection/custom_rules',
    children: [
      {
        title: 'Rate Limiting',
        href: '/limits',
      },
      {
        title: 'Custom Rules',
        href: 'protection/custom_rules',
      },
      {
        title: 'Detection Mod',
        href: '/semantics',
      },
      {
        title: 'Settings',
        href: '/settings',
      },
    ],
  },
  {
    title: 'System',
    href: '/system',
    children: [],
  },
]
