import { Avatar, AvatarImage } from "@/components/ui/avatar"
import type { RUser } from "@/lib/types"
import clsx from "clsx"
import { X } from "lucide-react"
import { toast } from "sonner"

function RoomUser({
  user,
  myself,
  roomOwner,
  index,
  onLeaveSeat,
  onJoinSeat,
}: {
  user: RUser
  index: number
  myself?: string
  roomOwner?: string
  onLeaveSeat?: (uid: string) => void
  onJoinSeat?: (index: number) => void
}) {
  return (
    <div
      className={clsx(
        "relative flex h-11 w-11 items-center justify-center rounded-full border bg-amber-100",
        myself == user.id && "outline-2 outline-blue-300"
      )}
      onClick={() => {
        if (user.id && myself !== user.id) {
          toast.info(user.name)
        }
      }}
    >
      {user.id && user.id === roomOwner && (
        <div className="bg-secondary/90 absolute -top-2 -left-2 z-10 flex items-center justify-center rounded-md border border-blue-300 px-1 text-[10px]">
          房主
        </div>
      )}
      {onLeaveSeat && user.id && myself !== user.id && myself === roomOwner && (
        <div
          className="text-foreground/90 absolute -top-1 -right-1 z-10 flex h-4 w-4 items-center justify-center rounded-full border bg-red-300/90"
          onClick={(e) => {
            e.stopPropagation()
            onLeaveSeat(user.id)
          }}
        >
          <X size={11} />
        </div>
      )}

      {onJoinSeat && (
        <div
          className="bg-secondary absolute rounded-sm border px-1 text-[11px] shadow-xs"
          onClick={() => onJoinSeat(index)}
        >
          坐下
        </div>
      )}

      {user.icon && (
        <Avatar className="h-10 w-10">
          <AvatarImage src={`/avatar/${user.icon}.png`} />
        </Avatar>
      )}
      <div className="text-primary-foreground absolute -right-1 -bottom-1 z-10 flex h-4 w-4 items-center justify-center rounded-full border bg-stone-500/90 font-mono text-[10px]">
        {index + 1}
      </div>
      {user.id && (
        <div
          className={clsx(
            "absolute -bottom-0.5 -left-0.5 z-10 h-3 w-3 rounded-full border",
            user.online ? "bg-green-400" : "bg-red-400"
          )}
        />
      )}
    </div>
  )
}

export default RoomUser
