import Game from "@/components/games/Game"
import OfflineTip from "@/components/OfflineTip"
import type { GameConf } from "@/lib/types"
import { useState, type ReactNode } from "react"

function Gyyl({ conf, children }: { conf: GameConf; children: ReactNode }) {
  const [loading, setLoading] = useState(false)

  return (
    <Game
      conf={conf}
      icon={children}
      settings={<OfflineTip />}
      settingsBody={{}}
      loading={loading}
      setLoading={setLoading}
    >
      <div>TODO</div>
    </Game>
  )
}

export default Gyyl
