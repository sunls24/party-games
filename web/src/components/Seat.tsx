import Mounted from "@/components/Mounted"
import RoomUser from "@/components/RoomUser"
import { Avatar, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { $room, $user } from "@/lib/store/store"
import type { Game } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import { UserRoundPen } from "lucide-react"
import { useMemo } from "react"

function Seat({ game, startGame }: { game: Game; startGame: () => void }) {
  const user = useStore($user)
  const room = useStore($room)

  const inSeat = useMemo(
    () => room.users?.some((u) => u.id == user.id),
    [room.users]
  )

  const lastSeat = useMemo(
    () => room.users?.findLastIndex((u) => u.id) ?? game.min,
    [room.users]
  )

  const ready = useMemo(() => room.users?.every((u) => u.id), [room.users])

  function onJoinSeat(index: number) {
    request(`/api/room/${room.id}?type=${game.path.slice(1)}`, {
      index,
      userId: user.id,
    }).catch(requestErr)
  }

  function onSeatCount(count: number) {
    request(`/api/room/${room.id}?type=${game.path.slice(1)}`, {
      count,
    }).catch(requestErr)
  }

  function onLeaveSeat(userId?: string) {
    request(`/api/room/${room.id}?type=${game.path.slice(1)}`, {
      index: -1,
      userId: userId ?? user.id,
    }).catch(requestErr)
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="grid grid-cols-4 place-items-center gap-4 p-4">
        {room.users &&
          room.users.map((v, i) => (
            <RoomUser
              key={i}
              user={v}
              index={i}
              myself={user.id}
              roomOwner={room.owner}
              onLeaveSeat={onLeaveSeat}
              onJoinSeat={onJoinSeat}
            />
          ))}
      </div>
      <Mounted>
        <div className="flex items-center justify-center gap-2">
          <Avatar className="border">
            <AvatarImage src={`/avatar/${user.icon}.png`} />
          </Avatar>
          <div className="text-sm">{user.name}</div>
          {inSeat ? (
            <Button
              size="sm"
              variant="secondary"
              className="border border-red-300"
              onClick={() => onLeaveSeat()}
            >
              离开座位
            </Button>
          ) : (
            <div className="text-muted-foreground text-xs">观战中</div>
          )}
        </div>
        {user.id === room.owner && (
          <div className="flex items-center justify-center gap-2">
            <div className="text-muted-foreground flex items-center gap-1 font-normal">
              <UserRoundPen size={18} />
              玩家数量
            </div>
            <Select
              value={String(room.users?.length)}
              onValueChange={(s) => onSeatCount(Number(s))}
            >
              <SelectTrigger size="sm">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {Array.from({ length: game.max - game.min + 1 }, (_, i) => (
                  <SelectItem
                    key={i}
                    value={String(i + game.min)}
                    disabled={i + game.min <= lastSeat}
                  >
                    {i + game.min}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        )}

        {user.id === room.owner && ready && (
          <div className="flex justify-center" onClick={startGame}>
            <Button size="sm">开始游戏</Button>
          </div>
        )}
      </Mounted>
    </div>
  )
}

export default Seat
