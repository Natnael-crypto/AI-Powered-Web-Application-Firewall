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

export enum logFilterType {
  CLIENTIP,
  COUNTRY,
  METHOD,
  PROTOCOL,
}

export enum filterOperations {
  EQUALS_TO,
  GREATER_THAN,
  LESS_THAN,
  NOT_EQUAL,
}
