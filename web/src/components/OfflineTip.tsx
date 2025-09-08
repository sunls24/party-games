import { BadgeInfo } from "lucide-react"

function OfflineTip() {
  return (
    <div className="text-muted-foreground flex items-center gap-1 text-sm">
      <BadgeInfo size={18} />
      此游戏需要发言，推荐线下玩
    </div>
  )
}

export default OfflineTip
