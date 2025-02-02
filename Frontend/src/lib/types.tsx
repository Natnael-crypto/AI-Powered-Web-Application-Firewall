export type FilterValues = {
  ipAddress: string
  port: string
  domain: string
  startAt: string
  endAt: string
}

export type AttackLog = {
  ipAddress: string
  application: string
  attackCount: number
  duration: string
  startAt: string
}

export type LogTable = {
  action: string
  url: string
  attackType: string
  ipAddress: string
  time: string
}
