import Back from "@/components/Back"
import DynamicIcon from "@/components/DynamicIcon"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { $room, $user } from "@/lib/store/store"
import type { Game } from "@/lib/types"
import { copyToClipboard, longPoll, request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { ExternalLink } from "lucide-react"
import { useEffect } from "react"

function RoomHeader({ game }: { game: Game }) {
  const room = useStore($room)

  useEffect(() => {
    const roomId = window.location.pathname.split("/")[2]
    request(
      `/api/room/${roomId}/cr?type=${game.path.slice(1)}&userId=${$user.get().id}`,
      {
        seat: game.min,
      }
    )
      .then((data) => $room.set(data))
      .catch(requestErr)

    const abort = new AbortController()
    longPoll(
      abort.signal,
      () =>
        `/api/room/${roomId}/long?type=${game.path.slice(1)}&version=${$room.get().version ?? 0}&userId=${$user.get().id}`,
      (date) => $room.set(date)
    ).catch(requestErr)
    return () => abort.abort()
  }, [])

  function onInviteClick() {
    copyToClipboard(
      window.location.toString(),
      "链接地址已拷贝，发送至好友，可直接进入房间"
    )
  }

  return (
    <header className="bg-secondary flex items-center justify-around border-b px-3 py-2 shadow-xs">
      <div className="flex-1">
        <Back href={game.path} />
      </div>
      <div className="flex flex-1 flex-col items-center">
        <div className="flex items-center gap-2">
          <DynamicIcon name={game.icon} size={20} />
          <span className="font-medium">{game.name}</span>
        </div>
        <div className="flex items-center gap-1 text-sm">
          <span className="text-muted-foreground">房间号:</span>
          {room.id ? (
            <span className="font-mono font-medium">{room.id}</span>
          ) : (
            <Skeleton className="bg-border h-4 w-8.5" />
          )}
        </div>
      </div>
      <div className="flex flex-1 justify-end">
        <Button
          variant="secondary"
          size="sm"
          className="bg-sidebar border"
          onClick={onInviteClick}
        >
          <ExternalLink />
          邀请好友
        </Button>
      </div>
    </header>
  )
}

export default RoomHeader
