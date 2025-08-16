import Game from "@/components/games/Game"
import type { GameConf } from "@/lib/types"
import { BadgeInfo } from "lucide-react"
import { useState, type ReactNode } from "react"

function Gyyl({ conf, children }: { conf: GameConf; children: ReactNode }) {
  const [loading, setLoading] = useState(false)

  return (
    <Game
      conf={conf}
      icon={children}
      settings={
        <div className="text-muted-foreground flex items-center gap-1 text-sm">
          <BadgeInfo size={18} />
          此游戏需要发言，推荐线下玩
        </div>
      }
      settingsBody={{}}
      loading={loading}
      setLoading={setLoading}
    >
      <div>TODO</div>
    </Game>
  )
}

export default Gyyl
