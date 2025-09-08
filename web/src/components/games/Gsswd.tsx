import Alarm from "@/components/Alarm"
import Game from "@/components/games/Game"
import OfflineTip from "@/components/OfflineTip"
import RoomUser from "@/components/RoomUser"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { $game, $room, $user } from "@/lib/store/store"
import type { GameConf } from "@/lib/types"
import { request, requestErr, requestStream } from "@/lib/utils"
import { useStore } from "@nanostores/react"
import clsx from "clsx"
import {
  BadgeInfo,
  BadgeQuestionMark,
  CircleArrowUp,
  HandHelping,
  ListTodo,
  Loader,
  Speech,
  TimerReset,
} from "lucide-react"
import { useEffect, useMemo, useRef, useState, type ReactNode } from "react"
import { toast } from "sonner"

function Gsswd({ conf, children }: { conf: GameConf; children: ReactNode }) {
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
      <Body />
    </Game>
  )
}

export default Gsswd

function Body() {
  const [loading, setLoading] = useState(false)
  const [loadingAI, setLoadingAI] = useState(false)
  const [helpAI, setHelpAI] = useState("")

  const user = useStore($user)
  const room = useStore($room)
  const game = useStore($game)

  const idx = useMemo(
    () => room.users?.findIndex((v) => v.id === user.id) ?? -1,
    [room.users]
  )

  const result = useMemo(() => {
    if (helpAI && game.state.stage === 2) {
      setHelpAI("")
    }
    const m = new Map<number, number>()
    let ties: number[] = []
    let alive = 0
    game.state?.players.forEach((p: any, idx: number) => {
      if (!p.out) {
        alive++
      }
      if (p.tie) {
        ties.push(idx)
      }
      if (p.vote < 0) {
        return
      }
      m.set(p.vote, (m.get(p.vote) ?? 0) + 1)
    })
    const failed = alive === 2
    return {
      m,
      ties,
      failed,
      gameOver:
        failed || (game.state && game.state.out === game.state.undercover),
      showTie:
        ties.length > 0 && (game.state.voteDone || game.state.stage === 1),
    }
  }, [game.state])

  function onNextStage() {
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/sswd/stage?id=${room.id}`, {})
      .finally(() => setLoading(false))
      .catch(requestErr)
  }

  function onVote(idxVote: number) {
    if (!room.users?.some((v) => v.id === user.id)) {
      return
    }
    if (loading) {
      return
    }
    setLoading(true)
    request(`/api/sswd/vote?id=${room.id}`, {
      idx,
      idxVote,
    })
      .then(() => toast.success(`已投票给玩家${idxVote + 1}`))
      .finally(() => setLoading(false))
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

  const signalRef = useRef(new AbortController())
  useEffect(() => signalRef.current.abort(), [])

  function onAIHelp() {
    if (loadingAI) {
      return
    }
    setHelpAI("")
    setLoadingAI(true)
    const word =
      game.state.undercover === idx ? game.state.spyWord : game.state.word
    requestStream(`/api/sswd/help?word=${word}`, {
      signal: signalRef.current.signal,
      onChunk(chunk) {
        setHelpAI((v) => v + chunk)
      },
    })
      .finally(() => setLoadingAI(false))
      .catch(requestErr)
  }

  return (
    <div className="flex flex-col gap-4 p-4">
      <div className="self-center px-2 font-medium">
        {game.state ? (
          <div className="flex items-center gap-1 border-b-2 border-blue-300">
            {game.state.stage === 1 ? (
              <>
                <Speech size={20} />
                发言阶段
              </>
            ) : (
              <>
                <ListTodo size={20} />
                投票阶段
              </>
            )}
          </div>
        ) : (
          <Skeleton className="bg-border/50 h-6.5 w-22" />
        )}
      </div>
      <div className="px-2 font-medium">玩家列表</div>
      <div className="grid w-full grid-cols-4 place-items-center gap-4">
        {room.users?.map((u, i) => (
          <div key={i} className="flex h-full flex-col items-center gap-4">
            <RoomUser
              user={u}
              index={i}
              myself={user.id}
              roomOwner={room.owner}
            >
              {game.state?.players[i].out && (
                <div className="bg-primary/50 text-primary-foreground absolute flex h-full w-full items-center justify-center rounded-full text-xs">
                  OUT
                </div>
              )}
              {game.state &&
                game.state.stage === 2 &&
                !game.state.voteDone &&
                !game.state?.players[i].out &&
                !game.state?.players[i].tie &&
                game.state?.players[i].vote < 0 && <Alarm />}
            </RoomUser>

            {game.state && (u.id === user.id || result.gameOver) ? (
              <div
                className={clsx(
                  "rounded-md px-2 text-sm font-medium outline",
                  result.gameOver && i === game.state.undercover
                    ? "bg-blue-100"
                    : "bg-amber-100"
                )}
              >
                {game.state.undercover === i
                  ? game.state.spyWord
                  : game.state.word}
              </div>
            ) : (
              <BadgeQuestionMark size={20} />
            )}
            {game.state && game.state.stage === 2 && (
              <>
                {game.state.players[idx].vote < 0 &&
                  u.id !== user.id &&
                  (result.ties.length <= 0 ||
                    (game.state.players[i].tie &&
                      !game.state.players[idx].tie)) &&
                  !game.state.voteDone &&
                  !game.state.players[i].out &&
                  !game.state.players[idx].out && (
                    <Button
                      size="sm"
                      disabled={loading}
                      variant="secondary"
                      className="bg-sidebar border"
                      onClick={() => onVote(i)}
                    >
                      {loading ? (
                        <Loader className="animate-spin" />
                      ) : (
                        <CircleArrowUp />
                      )}
                      投票
                    </Button>
                  )}
                {(game.state.players[idx].vote >= 0 || game.state.voteDone) && (
                  <div className="text-sm">
                    {(u.id === user.id || game.state.voteDone) && (
                      <div
                        className={clsx(
                          "text-center font-medium",
                          game.state.voteDone ? "border-b-2" : "border-b-4 py-1"
                        )}
                      >
                        {game.state.players[i].vote === -1
                          ? "未投票"
                          : `投给玩家${game.state.players[i].vote + 1}`}
                      </div>
                    )}
                    {game.state.voteDone && result.m.has(i) && (
                      <div className="text-center">
                        被投
                        <span className="px-0.5 font-medium">
                          {result.m.get(i)}
                        </span>
                        票
                      </div>
                    )}
                  </div>
                )}
              </>
            )}
          </div>
        ))}
      </div>
      <div className="px-2">
        <span className="font-medium">发言顺序</span>
        {game.state && (
          <div className="mt-2 border-l-3 border-blue-300 pl-1 text-sm">
            从
            <span className="px-0.5 font-medium underline underline-offset-4">{`玩家${game.state.start + 1}`}</span>
            开始，
            <span className="pr-0.5 font-medium underline underline-offset-4">
              {game.state.clockwise ? "顺时针" : "逆时针"}
            </span>
            发言
          </div>
        )}
      </div>
      {game.state && game.state.stage === 2 && !game.state.voteDone && (
        <div className="text-muted-foreground mx-auto text-sm">
          等待所有人投票完成
        </div>
      )}

      {game.state && (game.state.out >= 0 || result.showTie) && (
        <div className="flex flex-col items-center gap-2">
          {game.state.out >= 0 && (
            <div className="text-sm font-medium">
              玩家{game.state.out + 1}是
              {game.state.out === game.state.undercover
                ? "卧底，平民获胜"
                : `平民，${result.failed ? "卧底胜利" : "游戏继续"}`}
            </div>
          )}
          {result.showTie && (
            <div className="mx-auto text-sm">
              <span className="font-medium">
                {result.ties.map((idx) => `玩家${idx + 1}`).join(", ")}
              </span>
              平票
            </div>
          )}

          {user.id === room.owner &&
            game.state.stage === 2 &&
            game.state.voteDone &&
            !result.gameOver && (
              <Button
                size="sm"
                disabled={loading}
                variant="secondary"
                className="bg-sidebar border border-blue-300"
                onClick={onNextStage}
              >
                {loading ? <Loader className="animate-spin" /> : <Speech />}{" "}
                重新发言
              </Button>
            )}
        </div>
      )}

      {game.state && (
        <div className="flex flex-col items-center gap-4 border-t border-dashed pt-4">
          {user.id === room.owner &&
            (game.state.stage === 1 ? (
              <Button
                size="sm"
                disabled={loading}
                variant="secondary"
                className="bg-sidebar mx-auto w-fit border border-blue-300"
                onClick={onNextStage}
              >
                {loading ? <Loader className="animate-spin" /> : <ListTodo />}
                开始投票
              </Button>
            ) : (
              (result.ties.length <= 0 || !game.state.voteDone) && (
                <Button
                  size="sm"
                  disabled={loading}
                  variant="secondary"
                  className="bg-sidebar border-destructive/50 mx-auto w-fit border"
                  onClick={onRestart}
                >
                  {loading ? (
                    <Loader className="animate-spin" />
                  ) : (
                    <TimerReset />
                  )}
                  重开一局
                </Button>
              )
            ))}
          {game.state.stage === 1 && (
            <div className="bg-border/30 flex w-full flex-col gap-2 rounded-md border p-2 text-sm">
              <div className="flex w-full items-center gap-1">
                <BadgeInfo size={18} />
                <span>不知道怎么描述</span>
                <div className="flex-1" />
                <Button
                  size="sm"
                  disabled={loadingAI}
                  variant="secondary"
                  className="bg-sidebar border"
                  onClick={onAIHelp}
                >
                  {loadingAI ? (
                    <Loader className="animate-spin" />
                  ) : (
                    <HandHelping />
                  )}
                  问问 AI
                </Button>
              </div>
              {helpAI && (
                <div className="border-t border-dashed pt-2 whitespace-break-spaces">
                  {helpAI}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
