import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import type { Game } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { ArrowRight, BadgePlus } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

function Start({ game }: { game: Game }) {
  const [input, setInput] = useState("")

  function onCreateRoom() {
    request(`/api/room?type=${game.path.slice(1)}`, {
      seat: game.min,
    })
      .then((data) => open(game.path + "/" + data.id, "_self"))
      .catch(requestErr)
  }

  function onGoRoom() {
    const value = input.trim()
    if (value.length != 4) {
      toast.warning("房间号必须是4位哦")
      return
    }
    open(game.path + "/" + value, "_self")
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key !== "Enter" || e.nativeEvent.isComposing) {
      return
    }
    e.preventDefault()
    onGoRoom()
  }

  return (
    <div className="flex flex-col gap-8 px-8">
      <div className="flex gap-4">
        <Input
          value={input}
          onChange={(e) => setInput(e.currentTarget.value)}
          placeholder="输入房间号，直接进入房间"
          onKeyDown={onKeyDown}
        />
        <Button variant="secondary" onClick={onGoRoom} className="border">
          <ArrowRight />
          GO
        </Button>
      </div>
      <Button onClick={onCreateRoom}>
        <BadgePlus />
        创建房间
      </Button>
    </div>
  )
}

export default Start
