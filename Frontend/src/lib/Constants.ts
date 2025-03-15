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
  US: 1200,
  IN: 800,
  FR: 300,
  CN: 700,
  DE: 500,
  BR: 400,
}

export const getColor = (requests: number) => {
  if (requests === 0) return '#E6F4EA' // Very faint green
  const maxRequests = Math.max(...Object.values(requestData), 1000) // Prevent division by zero
  const intensity = Math.min(1, requests / maxRequests) // Normalize intensity (0 - 1)
  const greenLevel = Math.floor(255 - intensity * 155) // Darker green as requests increase
  return `rgb(0, ${greenLevel}, 0)`
}
