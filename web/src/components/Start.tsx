import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import type { GameConf } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { ArrowRight, BadgePlus } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

function Start({ conf }: { conf: GameConf }) {
  const [input, setInput] = useState("")

  function onCreateRoom() {
    request(`/api/room/no?game=${conf.path}`)
      .then((data) => open(`/${conf.path}/${data}`, "_self"))
      .catch(requestErr)
  }

  function onGoRoom() {
    const value = input.trim()
    if (value.length != 4) {
      toast.warning("房间号必须是4位哦")
      return
    }
    open(`/${conf.path}/${value}`, "_self")
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key !== "Enter" || e.nativeEvent.isComposing) {
      return
    }
    e.preventDefault()
    onGoRoom()
  }

  return (
    <div className="flex flex-col gap-6 px-8">
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
