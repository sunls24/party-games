import { Button } from "@/components/ui/button"
import { Undo2 } from "lucide-react"

function Back({ href }: { href: string }) {
  return (
    <Button
      size="sm"
      variant="secondary"
      className="bg-sidebar border"
      onClick={() => open(href, "_self")}
    >
      <Undo2 />
    </Button>
  )
}

export default Back
