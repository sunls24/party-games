import { clsx, type ClassValue } from "clsx"
import { toast } from "sonner"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getRoomId() {
  return window.location.pathname.slice(1)
}

export function roomNo(id: string) {
  return id.split("/")[1]
}

export async function request(path: string, body?: any, signal?: AbortSignal) {
  const resp = await requestResp(path, body, signal)
  const json = await resp.json()
  if (json.code !== 0) {
    return Promise.reject(json.message)
  }
  return json.data
}

async function requestResp(path: string, body?: any, signal?: AbortSignal) {
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
  return resp
}

export async function requestStream(
  input: string,
  init: { onChunk: (chunk: string) => void; signal?: AbortSignal; body?: any }
) {
  const resp = await requestResp(input, init.body, init.signal)
  if (!resp.body) {
    return Promise.reject("response body is null")
  }
  const decoder = new TextDecoder()
  const reader = resp.body.getReader()

  while (true) {
    const { value, done } = await reader.read()
    if (done) {
      return Promise.resolve()
    }
    init.onChunk(decoder.decode(value, { stream: true }))
  }
}

function isAbort(err: any) {
  return err.name === "AbortError"
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
