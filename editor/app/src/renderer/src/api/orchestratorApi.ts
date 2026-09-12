// The orchestrator's control API — the whole of the editor's involvement in a
// game's lifecycle, since it never spawns a builder process itself.
const ORCHESTRATOR_BASE_URL = window.api?.orchestrator.getBaseUrl() ?? ''

export const ORCHESTRATOR_UNAVAILABLE_MESSAGE =
  'Running a game needs the LevelCraft orchestrator. Start LevelCraft through it rather than launching the editor on its own.'

export const ORCHESTRATOR_UNREACHABLE_MESSAGE =
  "Couldn't reach the LevelCraft orchestrator. It may have stopped running."

interface GameActionResponse {
  success: boolean
  message: string
}

interface GameStatusResponse {
  success: boolean
  running: boolean
}

export function isOrchestratorAvailable(): boolean {
  return ORCHESTRATOR_BASE_URL !== ''
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  if (!isOrchestratorAvailable()) {
    throw new Error(ORCHESTRATOR_UNAVAILABLE_MESSAGE)
  }

  let res: Response
  try {
    res = await fetch(`${ORCHESTRATOR_BASE_URL}${path}`, {
      method,
      headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : JSON.stringify(body)
    })
  } catch {
    throw new Error(ORCHESTRATOR_UNREACHABLE_MESSAGE)
  }

  const data = await res.json()
  if (!res.ok) {
    throw new Error((data as GameActionResponse).message ?? 'Request failed')
  }
  return data as T
}

export function runGame(filepath: string): Promise<GameActionResponse> {
  return request<GameActionResponse>('POST', '/game/run', { filepath })
}

export function stopGame(): Promise<GameActionResponse> {
  return request<GameActionResponse>('POST', '/game/stop')
}

export function getGameStatus(): Promise<GameStatusResponse> {
  return request<GameStatusResponse>('GET', '/game/status')
}
