import RoomHeader from "@/components/RoomHeader"
import Seat from "@/components/Seat"
import type { Game } from "@/lib/types"

function Gyyl({ game }: { game: Game }) {
  return (
    <>
      <RoomHeader game={game} />
      <Seat game={game} />
    </>
  )
}

export default Gyyl
