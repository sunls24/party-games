import Mounted from "@/components/Mounted"
import RoomUser from "@/components/RoomUser"
import { Avatar, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { $room, $user } from "@/lib/store/store"
import type { GameConf } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { Loader, UserRoundPen } from "lucide-react"
import { useMemo, type ReactNode } from "react"

function Seat({
  conf,
  children,
  startGame,
  loading,
}: {
  conf: GameConf
  children: ReactNode
  startGame: () => void
  loading: boolean
}) {
  const user = useStore($user)
  const room = useStore($room)

  const inSeat = useMemo(
    () => room.users?.some((u) => u.id == user.id),
    [room.users]
  )

  const lastSeat = useMemo(
    () => room.users?.findLastIndex((u) => u.id) ?? conf.min,
    [room.users]
  )

  const ready = useMemo(() => room.users?.every((u) => u.id), [room.users])

  function onJoinSeat(idx: number) {
    request(`/api/room/seat?id=${room.id}`, {
      idx,
      uid: user.id,
    }).catch(requestErr)
  }

  function onSeatCount(count: number) {
    request(`/api/room/seat?id=${room.id}`, {
      count,
    }).catch(requestErr)
  }

  function onLeaveSeat(userId: string) {
    if (loading) {
      return
    }
    request(`/api/room/seat?id=${room.id}`, {
      idx: -1,
      uid: userId,
    }).catch(requestErr)
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="grid grid-cols-4 place-items-center gap-4 p-4">
        {room.users?.map((u, i) => (
          <RoomUser
            key={i}
            index={i}
            user={u}
            myself={user.id}
            roomOwner={room.owner}
            onLeaveSeat={onLeaveSeat}
            onJoinSeat={onJoinSeat}
          />
        ))}
      </div>
      <Mounted>
        <div className="flex items-center justify-center gap-2">
          <Avatar className="h-9 w-9 border">
            {user.icon && <AvatarImage src={`/avatar/${user.icon}.png`} />}
          </Avatar>
          <div className="text-sm">{user.name}</div>
          {inSeat ? (
            <Button
              size="sm"
              disabled={loading}
              variant="secondary"
              className="border-destructive/50 border"
              onClick={() => onLeaveSeat(user.id!)}
            >
              离开座位
            </Button>
          ) : (
            <div className="text-muted-foreground text-xs">观战中</div>
          )}
        </div>
        {user.id === room.owner && (
          <div className="flex items-center justify-center gap-2">
            <Label htmlFor="userCount">
              <UserRoundPen size={18} />
              玩家数量
            </Label>
            <Select
              disabled={loading}
              value={String(room.users?.length)}
              onValueChange={(v) => onSeatCount(Number(v))}
            >
              <SelectTrigger size="sm" id="userCount">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {Array.from({ length: conf.max - conf.min + 1 }, (_, i) => (
                  <SelectItem
                    key={i}
                    value={String(i + conf.min)}
                    disabled={i + conf.min <= lastSeat}
                  >
                    {i + conf.min}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {user.id === room.owner && (
          <div className="flex flex-col items-center justify-center gap-1">
            {!ready && (
              <div className="text-muted-foreground text-xs">等待人齐</div>
            )}
            <Button size="sm" onClick={startGame} disabled={!ready || loading}>
              {loading && <Loader className="animate-spin" />}
              开始游戏
            </Button>
          </div>
        )}
        {children}
      </Mounted>
    </div>
  )
}

export default Seat
