import Back from "@/components/Back"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { $room, $user } from "@/lib/store/store"
import type { GameConf } from "@/lib/types"
import {
  copyToClipboard,
  getRoomId,
  longPoll,
  request,
  requestErr,
  roomNo,
} from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { ExternalLink, OctagonX } from "lucide-react"
import { useEffect, type ReactNode } from "react"

function RoomHeader({
  conf,
  icon,
  stopGame,
}: {
  conf: GameConf
  icon: ReactNode
  stopGame?: () => void
}) {
  const room = useStore($room)

  useEffect(() => {
    const roomId = getRoomId()
    request(`/api/room/init?id=${roomId}&uid=${$user.get().id}`, {
      seat: conf.min,
    })
      .then((data) => $room.set(data))
      .catch(requestErr)

    const signal = new AbortController()
    longPoll(
      signal.signal,
      () =>
        `/api/room/long?id=${roomId}&uid=${$user.get().id}&version=${$room.get().version ?? 0}`,
      (data) => $room.set(data)
    ).catch(requestErr)
    return () => signal.abort()
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
        <Back href={"/" + conf.path} />
      </div>
      <div className="flex flex-1 flex-col items-center">
        <div className="flex items-center gap-2">
          {icon}
          <span className="text-sm font-medium">{conf.name}</span>
        </div>
        <div className="flex items-center gap-1 text-sm">
          <span className="text-muted-foreground">房间号:</span>
          {room.id ? (
            <span className="font-mono font-medium">{roomNo(room.id)}</span>
          ) : (
            <Skeleton className="bg-border h-4 w-8.5" />
          )}
        </div>
      </div>
      <div className="flex flex-1 justify-end">
        {stopGame ? (
          <Button
            size="sm"
            variant="secondary"
            className="bg-sidebar border-destructive/50 border"
            onClick={stopGame}
          >
            <OctagonX />
            结束游戏
          </Button>
        ) : (
          <Button
            size="sm"
            variant="secondary"
            className="bg-sidebar border"
            onClick={onInviteClick}
          >
            <ExternalLink />
            邀请好友
          </Button>
        )}
      </div>
    </header>
  )
}

export default RoomHeader
