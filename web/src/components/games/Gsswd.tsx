import RoomHeader from "@/components/RoomHeader"
import Seat from "@/components/Seat"
import { Button } from "@/components/ui/button"
import { $game } from "@/lib/store/store"
import type { Game } from "@/lib/types"
import { longPoll, request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { useEffect } from "react"

function Gsswd({ game }: { game: Game }) {
  useEffect(() => {
    const roomId = window.location.pathname.split("/")[2]
    request(`/api/game/${roomId}/cr?type=${game.path.slice(1)}`, {})
      .then((data) => $game.set(data))
      .catch(requestErr)

    const abort = new AbortController()
    longPoll(
      abort.signal,
      () =>
        `/api/game/${roomId}/long?type=${game.path.slice(1)}&version=${$game.get().version ?? 0}`,
      (date) => $game.set(date)
    ).catch(requestErr)
    return () => abort.abort()
  }, [])

  const gameState = useStore($game)

  function startGame() {
    const roomId = window.location.pathname.split("/")[2]
    request(`/api/game/${roomId}/start?type=${game.path.slice(1)}`, {}).catch(
      requestErr
    )
  }

  function stopGame() {
    const roomId = window.location.pathname.split("/")[2]
    request(`/api/game/${roomId}/stop?type=${game.path.slice(1)}`, {}).catch(
      requestErr
    )
  }

  return (
    <>
      <RoomHeader game={game} />
      {gameState.started ? (
        <div className="flex flex-col items-center justify-center gap-6 py-6">
          <span>游戏进行中... TODO</span>
          <Button size="sm" variant="secondary" onClick={stopGame}>
            停止游戏
          </Button>
        </div>
      ) : (
        <Seat game={game} startGame={startGame} />
      )}
    </>
  )
}

export default Gsswd
