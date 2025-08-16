export type Game = {
  name: string
  icon: string
  path: string
  min: number
  max: number
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
}

export type GameState = {
  id: string
  version: number
  started: boolean
}
