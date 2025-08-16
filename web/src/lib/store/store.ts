import type { Game, Room, User } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { persistentAtom } from "@nanostores/persistent"
import { map } from "nanostores"

export const $room = map<Partial<Room>>({})

export const $game = map<Partial<Game>>({})

export const $user = persistentAtom<Partial<User>>(
  "user",
  {},
  { encode: JSON.stringify, decode: JSON.parse }
)

export function initUser() {
  if ($user.get().id) {
    request(`/api/user/init?id=${$user.get().id}`, {})
      .then((data) => $user.set(data))
      .catch(requestErr)
    return
  }
  request(`/api/user/init?id=${crypto.randomUUID().replaceAll("-", "")}`, {})
    .then((data) => $user.set(data))
    .catch(requestErr)
}
