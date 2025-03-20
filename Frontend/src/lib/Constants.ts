export const SiderbarContentItems = [
  {
    title: 'Dashboard',
    href: 'dashboard',
    children: [],
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
