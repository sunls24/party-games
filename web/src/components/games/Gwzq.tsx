import Alarm from "@/components/Alarm"
import Game from "@/components/games/Game"
import RoomUser from "@/components/RoomUser"
import { Button } from "@/components/ui/button"
import { $game, $room, $user } from "@/lib/store/store"
import type { GameConf } from "@/lib/types"
import { request, requestErr } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import clsx from "clsx"
import { Loader, Repeat1, TimerReset } from "lucide-react"
import { useMemo, useState, type ReactNode } from "react"
import { toast } from "sonner"

function Gwzq({ conf, children }: { conf: GameConf; children: ReactNode }) {
  const [loading, setLoading] = useState(false)

  return (
    <Game
      conf={conf}
      icon={children}
      settings={<></>}
      settingsBody={{}}
      loading={loading}
      setLoading={setLoading}
    >
      <Body />
    </Game>
  )
}

export default Gwzq

const enum CellType {
  Normal,
  Right,
  Bottom,
  Last,
}

const ROWS = 14
const COLS = 14
const SIZE = ROWS * COLS

function lastRow(i: number) {
  return (i + 1) % ROWS === 0
}

function lastCol(i: number) {
  return i + 1 > (COLS - 1) * ROWS
}

function lastCell(i: number) {
  return i + 1 === SIZE
}

function showPoint(i: number) {
  const r = i % ROWS
  const c = Math.floor(i / ROWS)
  return (r === 3 || r === 7 || r === 11) && (c === 3 || c === 7 || c === 11)
}

const board = Array.from({ length: SIZE }, (_, i) => ({
  lastRow: lastRow(i),
  lastCol: lastCol(i),
  lastCell: lastCell(i),
  showPoint: showPoint(i),
}))

function cellNumber(i: number, ct: CellType) {
  switch (ct) {
    case CellType.Right:
      return Math.floor(i / ROWS) + i + 2
    case CellType.Bottom:
      return (i % ROWS) + (ROWS + 1) * COLS + 1
    case CellType.Last:
      return (ROWS + 1) * (COLS + 1)
    default:
      return Math.floor(i / ROWS) + i + 1
  }
}

const cellClass =
  "absolute h-6 w-6 sm:h-12 sm:w-12 flex items-center justify-center"
function Body() {
  const [loading, setLoading] = useState(false)

  const user = useStore($user)
  const game = useStore($game)
  const room = useStore($room)

  const idx = useMemo(
    () => room.users?.findIndex((v) => v.id === user.id) ?? -1,
    [room.users]
  )

  function onCellClick(i: number, ct: CellType = CellType.Normal) {
    if (
      idx < 0 ||
      game.state.current !== idx ||
      game.state.board[i][ct] ||
      loading
    ) {
      return
    }
    const val = game.state.black === idx ? -1 : 1
    request(`/api/wzq/chess?id=${game.id}`, {
      idx: i,
      ct,
      val,
    })
      .then(() => checkWin(game.state.board, i, ct, val))
      .catch(requestErr)
  }

  function onRestart() {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/game/restart?id=${room.id}`, {})
      .then(() => toast.success("已经重新开了一局"))
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  function onRegret() {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/wzq/regret?id=${room.id}`, { idx: idx === 0 ? 1 : 0 })
      .then(() => toast.success("已向对方请求悔棋"))
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  function onRegretAction(action: number) {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/wzq/regret/action?id=${room.id}`, { action })
      .then(() => toast.success(`已${action === 1 ? "同意" : "拒绝"}对方悔棋`))
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  return (
    <div className="flex flex-col gap-6 px-4 py-6 sm:p-8">
      <div className="flex justify-around">
        {room.users?.map((u, i) => (
          <div key={i} className="flex w-20 flex-col items-center gap-4">
            <RoomUser
              user={u}
              index={i}
              myself={user.id}
              roomOwner={room.owner}
            >
              {game.state?.current === i && <Alarm />}
            </RoomUser>
            <div className="flex items-center gap-1 text-sm underline underline-offset-4">
              <Chess value={game.state?.black === i ? -1 : 1} inBoard={false} />
              {game.state?.current === i && <span>思考中</span>}
            </div>
          </div>
        ))}
      </div>
      {game.state && (
        <>
          <Board
            state={game.state.board}
            onCellClick={onCellClick}
            lastIdx={game.state.lastIdx}
            lastCt={game.state.lastCt}
          />
          <div className="flex flex-col gap-4 border-t border-dashed pt-4">
            {game.state.regret === idx && (
              <div className="flex items-center justify-center gap-2">
                <span className="text-sm font-medium">
                  是否同意对方的悔棋请求
                </span>
                <Button
                  size="sm"
                  disabled={loading}
                  variant="secondary"
                  className="bg-sidebar border border-blue-300"
                  onClick={() => onRegretAction(1)}
                >
                  同意
                </Button>
                <Button
                  size="sm"
                  disabled={loading}
                  variant="secondary"
                  className="bg-sidebar border-destructive border"
                  onClick={() => onRegretAction(0)}
                >
                  拒绝
                </Button>
              </div>
            )}
            <div className="flex justify-center gap-4">
              {idx >= 0 &&
                game.state.regret < 0 &&
                game.state.lastCt >= 0 &&
                game.state.current !== idx && (
                  <Button
                    size="sm"
                    disabled={loading}
                    variant="secondary"
                    className="bg-sidebar border"
                    onClick={onRegret}
                  >
                    <Repeat1 />
                    悔棋
                  </Button>
                )}
              {user.id === room.owner && (
                <Button
                  size="sm"
                  disabled={loading}
                  variant="secondary"
                  className="bg-sidebar border"
                  onClick={onRestart}
                >
                  {loading ? (
                    <Loader className="animate-spin" />
                  ) : (
                    <TimerReset />
                  )}
                  重开一局
                </Button>
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}

function Board({
  state,
  lastIdx,
  lastCt,
  onCellClick,
}: {
  state: number[][]
  lastIdx: number
  lastCt: number
  onCellClick: (i: number, ct?: CellType) => void
}) {
  return (
    <div className="grid grid-cols-14 border-b border-l border-stone-600 bg-orange-100 outline outline-stone-600">
      {board.map((c, i) => (
        <div
          key={i}
          className="relative aspect-square border-t border-r border-stone-600 text-sm"
        >
          <div
            className={clsx(
              "-top-[12.5px] -left-[12.5px] sm:-top-[24.5px] sm:-left-[24.5px]",
              cellClass
            )}
            onClick={() => onCellClick(i)}
          >
            {c.showPoint && (
              <div className="h-1.5 w-1.5 rounded-full bg-stone-600" />
            )}
            <Chess
              value={state[i][CellType.Normal]}
              last={i === lastIdx && lastCt === CellType.Normal}
            />
          </div>
          {c.lastRow && (
            <div
              className={clsx(
                "-top-[12.5px] -right-[12.5px] sm:-top-[24.5px] sm:-right-[24.5px]",
                cellClass
              )}
              onClick={() => onCellClick(i, CellType.Right)}
            >
              <Chess
                value={state[i][CellType.Right]}
                last={i === lastIdx && lastCt === CellType.Right}
              />
            </div>
          )}
          {c.lastCol && (
            <div
              className={clsx(
                "-bottom-[12.5px] -left-[12.5px] sm:-bottom-[24.5px] sm:-left-[24.5px]",
                cellClass
              )}
              onClick={() => onCellClick(i, CellType.Bottom)}
            >
              <Chess
                value={state[i][CellType.Bottom]}
                last={i === lastIdx && lastCt === CellType.Bottom}
              />
            </div>
          )}
          {c.lastCell && (
            <div
              className={clsx(
                "-right-[12.5px] -bottom-[12.5px] sm:-right-[24.5px] sm:-bottom-[24.5px]",
                cellClass
              )}
              onClick={() => onCellClick(i, CellType.Last)}
            >
              <Chess
                value={state[i][CellType.Last]}
                last={i === lastIdx && lastCt === CellType.Last}
              />
            </div>
          )}
        </div>
      ))}
    </div>
  )
}

function Chess({
  value,
  last = false,
  inBoard = true,
}: {
  value: number
  last?: boolean
  inBoard?: boolean
}) {
  if (!value) {
    return
  }

  return (
    <div
      className={clsx(
        "z-10 h-5 w-5 rounded-full border shadow-md sm:h-10 sm:w-10",
        value > 0 ? "border-stone-700 bg-white" : "border-white bg-stone-700",
        inBoard && "absolute",
        last && "outline-2 outline-blue-300 sm:outline-4"
      )}
    />
  )
}

function countBy(
  state: number[][],
  count: number,
  val: number,
  next: (
    i: number,
    last: boolean
  ) => { idx: number; ct: CellType; ct2?: CellType }
) {
  let ret = 0
  for (let i = 1; i <= count; i++) {
    const n = next(i, i === count)
    const nv = state[n.idx][n.ct]
    if (nv !== val) {
      break
    }
    ret++
    if (n.ct2) {
      state[n.idx][n.ct2] === val && ret++
    }
  }
  return ret
}

function checkWin(state: number[][], idx: number, ct: CellType, val: number) {
  const row = idx % ROWS
  const col = Math.floor(idx / ROWS)
  console.log("checkWin", idx, ct, val, row, "*", col)
  let count = 1
  switch (ct) {
    case CellType.Normal:
      // left
      count += countBy(state, Math.min(4, row), val, (i) => ({
        idx: idx - i,
        ct,
      }))
      console.log("check left", count)
      if (count >= 5) {
        break
      }
      // right
      count += countBy(state, Math.min(4, ROWS - 1 - row), val, (i, last) => ({
        idx: idx + i,
        ct,
        ct2: last && board[idx + i].lastRow ? CellType.Right : CellType.Normal,
      }))
      console.log("check right", count)
      if (count >= 5) {
        break
      }
      // top
      count = 1
      count += countBy(state, Math.min(4, col), val, (i, last) => ({
        idx: idx - ROWS * i,
        ct,
      }))
      console.log("check top", count)
      if (count >= 5) {
        break
      }

      // bottom
      count += countBy(state, Math.min(4, COLS - 1 - row), val, (i, last) => ({
        idx: idx + ROWS * i,
        ct, // TODO: ct2
      }))
      console.log("check bottom", count)
      if (count >= 5) {
        break
      }

      break
    case CellType.Right:
      break
    case CellType.Bottom:
      break
    case CellType.Last:
      break
  }
  console.log("count", count)

  if (count < 5) {
    return
  }
}
