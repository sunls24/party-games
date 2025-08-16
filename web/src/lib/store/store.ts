import type { GameState, Room, User } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { persistentAtom } from "@nanostores/persistent"
import { map } from "nanostores"

export const $room = map<Partial<Room>>({})

export const $game = map<Partial<GameState>>({})

export const $user = persistentAtom<Partial<User>>(
  "user",
  {},
  { encode: JSON.stringify, decode: JSON.parse }
)

export function initUser() {
  const id = $user.get().id
  if (id) {
    request(`/api/user/${id}/cr`)
      .then((data) => $user.set(data))
      .catch(requestErr)
    return
  }
  request("/api/user", {})
    .then((data) => $user.set(data))
    .catch(requestErr)
}
