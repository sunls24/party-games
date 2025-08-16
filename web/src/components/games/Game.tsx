import RoomHeader from "@/components/RoomHeader"
import Seat from "@/components/Seat"
import { $game, $room, $user } from "@/lib/store/store"
import type { GameConf } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { type ReactNode } from "react"

function Game({
  conf,
  icon,
  children,
  settings,
  settingsBody,
  loading,
  setLoading,
}: {
  conf: GameConf
  icon: ReactNode
  children: ReactNode
  settings: ReactNode
  settingsBody: any
  loading: boolean
  setLoading: (b: boolean) => void
}) {
  const room = useStore($room)

  function startGame() {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/game/start?id=${room.id}`, settingsBody)
      .then((data) => $game.set(data))
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  function stopGame() {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/game/stop?id=${room.id}`, {})
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  return (
    <>
      <RoomHeader
        conf={conf}
        icon={icon}
        stopGame={
          room.started && $user.get().id === room.owner ? stopGame : undefined
        }
      />
      {room.started ? (
        children
      ) : (
        <Seat conf={conf} startGame={startGame} loading={loading}>
          <div className="flex flex-col items-center border-t border-dashed pt-4">
            {settings}
          </div>
        </Seat>
      )}
    </>
  )
}

export default Game
