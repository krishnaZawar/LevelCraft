import { create } from 'zustand'
import type { GameState } from '../api/gameApi'

// Fixed per docs/client/approach.md §4.4 — builder/backend always binds
// here, unlike editor/backend's dynamic port.
const BUILDER_WS_URL = 'ws://localhost:8000/requests'

type RuntimeStatus = 'connecting' | 'connected' | 'error'

interface RuntimeStoreState {
  status: RuntimeStatus
  // Latest full-scene push from builder/backend (every ~5s, see game_loop.go).
  // Same shape as editorStore's gameObjects, so renderScene() takes either.
  snapshot: GameState | null
  error: string | null

  connect: () => void
  sendKeyDown: (keyName: string) => void
  sendKeyUp: (keyName: string) => void
}

// This store only ever exists inside the Play window's own renderer
// process (see PlayWindowApp.tsx) — a separate window and JS context from
// the main editor, popped out by the main process rather than replacing
// the main viewport. connect() assumes builder/backend is already up: the
// main window's Run button spawned it via window.api.builder.launch
// before this window was even created.
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
    set({ status: 'connecting', error: null })
    const ws = new WebSocket(BUILDER_WS_URL)
    socket = ws

    ws.onopen = () => {
      if (socket === ws) set({ status: 'connected' })
    }

    ws.onmessage = (event) => {
      if (socket !== ws) return
      // builder/backend's EventResponse always sets msg to the literal
      // string "inputState" even though `data` is the full scene snapshot,
      // not input-specific — a known naming bug (see approach.md §4.2).
      // Every frame with `data` present just means "here's the latest
      // state," so branching on `msg` would be wrong.
      try {
        const parsed = JSON.parse(event.data as string)
        if (parsed?.data) set({ snapshot: parsed.data })
      } catch {
        // Malformed frame — ignore, matches the backend's own
        // ignore-and-keep-going stance on bad input.
      }
    }

    ws.onerror = () => {
      if (socket === ws) set({ status: 'error', error: 'Lost connection to builder/backend' })
    }

    ws.onclose = () => {
      if (socket === ws) socket = null
    }
  },

  sendKeyDown: (keyName) => send('KeyDownCommand', keyName),
  sendKeyUp: (keyName) => send('KeyUpCommand', keyName)
}))
