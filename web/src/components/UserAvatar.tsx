import Mounted from "@/components/Mounted"
import { Avatar, AvatarImage } from "@/components/ui/avatar"
import { $user } from "@/lib/store/store"
import { useStore } from "@nanostores/react"

function UserAvatar() {
  const user = useStore($user)
  return (
    <Mounted>
      <Avatar className="border">
        <AvatarImage src={`/avatar/${user.icon}.png`} />
      </Avatar>
    </Mounted>
  )
}

export default UserAvatar
