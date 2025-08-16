import { icons } from "lucide-react"

const DynamicIcon = ({ name, size = 24, className = "" }) => {
  const LucideIcon = icons[name]
  return <LucideIcon size={size} className={className} />
}

export default DynamicIcon
