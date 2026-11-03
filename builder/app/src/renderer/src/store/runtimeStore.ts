import { create } from 'zustand'
import { toGameState, type GameState, type WireGameObject } from '@/scene/types'

// The orchestrator assigns the backend's port at spawn time, so there is no
// fixed fallback: an empty address means this app was started without it.
function resolveSocketUrl(): string | null {
  const baseUrl = window.api?.backend.getBaseUrl()
  if (!baseUrl) return null
  return baseUrl.replace(/^http/, 'ws') + '/requests'
}

type RuntimeStatus = 'connecting' | 'connected' | 'error'

interface InputKey {
  name: string
}

interface InputMappingResponse {
  mapping: { keyboard: InputKey[] | null; mouse: InputKey[] | null }
}

// The engine's own list of keys, so the player never sends a command for a key
// the backend would reject. Fetched rather than restated: a key added to
// mapping.json works here without a change.
async function fetchKnownKeys(baseUrl: string): Promise<Set<string>> {
  const res = await fetch(`${baseUrl}/inputmapping`)
  if (!res.ok) throw new Error(`input mapping request failed: ${res.status}`)
  const body = (await res.json()) as InputMappingResponse
  return new Set(
    [...(body.mapping.keyboard ?? []), ...(body.mapping.mouse ?? [])].map((key) => key.name)
  )
}

interface RuntimeStoreState {
  status: RuntimeStatus
  // Latest full-scene push from builder/backend, one per game loop tick.
  snapshot: GameState | null
  error: string | null
  // Empty until the mapping arrives, which is why `knowsKey` treats an empty
  // set as "not loaded yet" rather than "no keys exist".
  knownKeys: Set<string>

  connect: () => void
  knowsKey: (keyName: string) => boolean
  sendKeyDown: (keyName: string) => void
  sendKeyUp: (keyName: string) => void
}

let socket: WebSocket | null = null

function send(requestType: string, keyName: string): void {
  if (socket?.readyState !== WebSocket.OPEN) return
  socket.send(JSON.stringify({ requestType, requestDetails: { keyName } }))
}

export const useRuntimeStore = create<RuntimeStoreState>((set, get) => ({
  status: 'connecting',
  snapshot: null,
  error: null,
  knownKeys: new Set<string>(),

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

    void fetchKnownKeys(window.api.backend.getBaseUrl())
      .then((knownKeys) => set({ knownKeys }))
      .catch(() => {
        // Non-fatal: without the mapping every key is forwarded and the
        // backend rejects what it does not know, which is how it behaved
        // before the mapping existed.
      })
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
        const parsed = JSON.parse(event.data as string) as { data?: WireGameObject[] }
        if (parsed.data) set({ snapshot: toGameState(parsed.data) })
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

  knowsKey: (keyName) => {
    const { knownKeys } = get()
    return knownKeys.size === 0 || knownKeys.has(keyName)
  },

  sendKeyDown: (keyName) => send('KeyDownCommand', keyName),
  sendKeyUp: (keyName) => send('KeyUpCommand', keyName)
}))
