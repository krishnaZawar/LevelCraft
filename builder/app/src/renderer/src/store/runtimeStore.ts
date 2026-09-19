import { create } from 'zustand'
import type { GameState } from '@/scene/types'

// The orchestrator assigns the backend's port at spawn time, so there is no
// fixed fallback: an empty address means this app was started without it.
function resolveSocketUrl(): string | null {
  const baseUrl = window.api?.backend.getBaseUrl()
  if (!baseUrl) return null
  return baseUrl.replace(/^http/, 'ws') + '/requests'
}

type RuntimeStatus = 'connecting' | 'connected' | 'error'

interface RuntimeStoreState {
  status: RuntimeStatus
  // Latest full-scene push from builder/backend, one per game loop tick.
  snapshot: GameState | null
  error: string | null

  connect: () => void
  sendKeyDown: (keyName: string) => void
  sendKeyUp: (keyName: string) => void
}

let socket: WebSocket | null = null

function send(requestType: string, keyName: string): void {
  if (socket?.readyState !== WebSocket.OPEN) return
  socket.send(JSON.stringify({ requestType, requestDetails: { keyName } }))
}

export const useRuntimeStore = create<RuntimeStoreState>((set) => ({
  status: 'connecting',
  snapshot: null,
  error: null,

  connect: () => {
    if (socket) return

    const url = resolveSocketUrl()
    if (!url) {
      set({
        status: 'error',
        error:
          'No builder backend address was provided. Start a game from the LevelCraft editor rather than launching this app directly.'
      })
      return
    }

    set({ status: 'connecting', error: null })
    const ws = new WebSocket(url)
    socket = ws

    ws.onopen = () => {
      if (socket === ws) set({ status: 'connected' })
    }

    ws.onmessage = (event) => {
      if (socket !== ws) return
      // `msg` is always the literal "inputState" even though `data` is the
      // full scene snapshot, so branching on it would be wrong.
      try {
        const parsed = JSON.parse(event.data as string)
        if (parsed?.data) set({ snapshot: parsed.data })
      } catch {
        // Malformed frame — ignored, as the backend does with bad input.
      }
    }

    ws.onerror = () => {
      if (socket === ws) set({ status: 'error', error: 'Lost connection to the builder backend.' })
    }

    ws.onclose = () => {
      if (socket === ws) socket = null
    }
  },

  sendKeyDown: (keyName) => send('KeyDownCommand', keyName),
  sendKeyUp: (keyName) => send('KeyUpCommand', keyName)
}))
