export const SiderbarContentItems = [
  {
    title: 'Dashboard',
    href: 'dashboard',
    children: [],
  },
  {
    title: 'Log',
    href: 'log/attacks',
    children: [
      {
        title: 'Attacks',
        href: 'log/attacks',
      },
      {
        title: 'Rate Limiting',
        href: 'log/limits',
      },
      {
        title: 'Anti-Bot',
        href: 'log/captcha',
      },
    ],
  },
  {
    title: 'Custom Rules',
    href: 'custom-rules',
    children: [],
  },
  {
    title: 'Web Services',
    href: 'web-services',
    children: [],
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
