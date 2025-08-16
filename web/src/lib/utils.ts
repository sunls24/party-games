import { ABORT_SAFE } from "@/lib/contants"
import { clsx, type ClassValue } from "clsx"
import { toast } from "sonner"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export async function request(path: string, body?: any, signal?: AbortSignal) {
  const resp = await fetch(path, {
    signal,
    method: body && "POST",
    body: body && JSON.stringify(body),
    headers: body && {
      "Content-Type": "application/json",
    },
  })
  if (!resp.ok) {
    let text = await resp.text()
    try {
      text = JSON.parse(text).message
    } catch (ignore) {}
    text = text ? text : "未知错误"
    return Promise.reject(`${resp.status} ${text}`)
  }
  const json = await resp.json()
  if (json.code !== 0) {
    return Promise.reject(json.message)
  }
  return json.data
}

function isAbort(err: any) {
  return err === ABORT_SAFE || err.name === "AbortError"
}

export function requestErr(err: any) {
  if (isAbort(err)) {
    return
  }
  toast.error(err.message ?? err)
}

export function delay(ms: number): Promise<void> {
  return new Promise<void>((resolve) => setTimeout(resolve, ms))
}

export async function longPoll(
  signal: AbortSignal,
  path: () => string,
  onData: (data: any) => void
) {
  while (!signal.aborted) {
    await delay(500)
    try {
      const resp = await fetch(path(), { signal, cache: "no-store" })
      if (resp.status === 204) {
        continue
      }
      if (!resp.ok) {
        let text = await resp.text()
        try {
          text = JSON.parse(text).message
        } catch (ignore) {}
        text = text ? text : "未知错误"
        throw Error(`${resp.status} ${text}`)
      }
      const json = await resp.json()
      if (json.code !== 0) {
        throw Error(json.message)
      }
      onData(json.data)
    } catch (err: any) {
      requestErr(err)
      await delay(1000)
    }
  }
}

export function copyToClipboard(text: string, tip: string) {
  navigator.clipboard
    .writeText(text)
    .then(() => toast.success(tip))
    .catch((e) => toast.error(e.message ?? e))
}
