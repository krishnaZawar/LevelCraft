const EDITOR_BACKEND_BASE_URL = 'http://localhost:3000'

export interface GameObjectDetails {
  id: string
  name: string
  group: string
  components: Record<string, Record<string, unknown>>
}

export type GameState = Record<string, GameObjectDetails>

interface SaveGameResponse {
  success: boolean
  message: string
}

interface LoadGameResponse {
  success: boolean
  gameState: GameState
}

interface ErrorResponse {
  success: false
  message: string
}

// fetch() throws a plain TypeError ("Failed to fetch") when it can't reach
// the server at all (backend not running, wrong port) — indistinguishable
// from other TypeErrors by type alone, but this is the only place we call
// fetch, so treating any thrown non-HTTP error as "backend unreachable"
// is accurate here and lets callers show something actionable instead of
// the raw browser error string.
export const BACKEND_UNREACHABLE_MESSAGE =
  "Couldn't reach the LevelCraft backend. Make sure editor/backend is running, then try again."

async function postJson<T>(path: string, body: unknown): Promise<T> {
  let res: Response
  try {
    res = await fetch(`${EDITOR_BACKEND_BASE_URL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
  } catch {
    throw new Error(BACKEND_UNREACHABLE_MESSAGE)
  }
  const data = await res.json()
  if (!res.ok) {
    throw new Error((data as ErrorResponse).message ?? 'Request failed')
  }
  return data as T
}

export function saveGame(filepath: string): Promise<SaveGameResponse> {
  return postJson<SaveGameResponse>('/game/save', { filepath })
}

export function loadGame(filepath: string): Promise<LoadGameResponse> {
  return postJson<LoadGameResponse>('/game/load', { filepath })
}

export async function isBackendReachable(): Promise<boolean> {
  try {
    const res = await fetch(`${EDITOR_BACKEND_BASE_URL}/ping`)
    return res.ok
  } catch {
    return false
  }
}
