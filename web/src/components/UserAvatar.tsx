import Mounted from "@/components/Mounted"
import { Avatar, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { $user } from "@/lib/store/store"
import { request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import clsx from "clsx"
import { CloudCheck, Loader, Shuffle } from "lucide-react"
import { useState } from "react"
import { toast } from "sonner"

function UserAvatar() {
  const user = useStore($user)
  const [icon, setIcon] = useState($user.get().icon)
  const [name, setName] = useState($user.get().name)
  const [loading, setLoading] = useState(false)
  const [open, setOpen] = useState(false)

  function onOpenChange(open: boolean) {
    if (open) {
      setIcon($user.get().icon)
      setName($user.get().name)
    }
    setOpen(open)
  }

  function onRandomName() {
    if (loading) {
      return
    }
    setLoading(true)
    request("/api/user/name")
      .then((data) => setName(data))
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  function onSaveClick() {
    if (loading) {
      return
    }
    setLoading(true)
    request("/api/user/save", {
      id: user.id,
      name,
      icon,
    })
      .then((data) => {
        $user.set(data)
        toast.success("保存成功")
        setOpen(false)
      })
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  return (
    <Mounted>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogTrigger asChild>
          <Avatar className="border">
            {user.icon && <AvatarImage src={`/avatar/${user.icon}.png`} />}
          </Avatar>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>自定义头像和名称</DialogTitle>
            <DialogDescription>记得点击保存哦！</DialogDescription>
          </DialogHeader>
          <div className="border-b border-dashed pb-1 font-medium">
            系统头像
          </div>
          <div className="grid grid-cols-5 place-items-center gap-3">
            {Array.from({ length: 18 }, (_, i) => (
              <Avatar
                key={i}
                onClick={() => {
                  setIcon(i + 1)
                }}
                className={clsx(
                  "h-11 w-11 border",
                  icon === i + 1 && "outline-3 outline-blue-300"
                )}
              >
                <AvatarImage src={`/avatar/${i + 1}.png`} />
              </Avatar>
            ))}
          </div>
          <div className="border-b border-dashed pb-1 font-medium">
            随机昵称
          </div>
          <div className="flex items-center justify-between">
            <span className="text-lg font-medium">{name}</span>
            <Button
              size="sm"
              disabled={loading}
              className="bg-sidebar border"
              variant="secondary"
              onClick={onRandomName}
            >
              {loading ? <Loader className="animate-spin" /> : <Shuffle />}
              随机
            </Button>
          </div>
          <Button disabled={loading} size="sm" onClick={onSaveClick}>
            <CloudCheck />
            保存
          </Button>
        </DialogContent>
      </Dialog>
    </Mounted>
  )
}

export default UserAvatar
