export type GameConf = {
  name: string
  icon: string
  path: string
  min: number
  max: number
  done?: boolean
}

export type User = {
  id: string
  name: string
  icon: number
}

export type RUser = User & {
  online: string
}

export type Room = {
  id: string
  version: number
  owner: string
  users: RUser[]
  started: boolean
}

export type Game = {
  id: string
  version: number
  conf: any
  state: any
}
