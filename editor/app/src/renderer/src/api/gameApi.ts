// Resolved by the main process (see main/backend.ts) since editor/backend
// no longer binds to a fixed port — the orchestrator (and this app, when
// self-spawning it in dev) assigns it a dynamic one. Falls back to the
// backend's own default port for non-Electron contexts (e.g. unit tests).
const EDITOR_BACKEND_BASE_URL = window.api?.backend.getBaseUrl() ?? 'http://localhost:3000'

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

interface GetGameStateResponse {
  success: boolean
  gameState: GameState
}

interface GetComponentsResponse {
  components: string[]
}

interface GameObjectResponse {
  success: boolean
  objectDetails: GameObjectDetails
}

interface DeleteGameobjectResponse {
  success: boolean
  message: string
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

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  let res: Response
  try {
    res = await fetch(`${EDITOR_BACKEND_BASE_URL}${path}`, {
      method,
      headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : JSON.stringify(body)
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

function encodePathSegment(value: string): string {
  return encodeURIComponent(value)
}

export function getGameState(): Promise<GetGameStateResponse> {
  return request<GetGameStateResponse>('GET', '/game/state')
}

export function saveGame(filepath: string): Promise<SaveGameResponse> {
  return request<SaveGameResponse>('POST', '/game/save', { filepath })
}

export function loadGame(filepath: string): Promise<LoadGameResponse> {
  return request<LoadGameResponse>('POST', '/game/load', { filepath })
}

export function getComponents(): Promise<GetComponentsResponse> {
  return request<GetComponentsResponse>('GET', '/components/')
}

export function addGameobject(): Promise<GameObjectResponse> {
  return request<GameObjectResponse>('POST', '/gameobjects/')
}

export function deleteGameobject(objectId: string): Promise<DeleteGameobjectResponse> {
  return request<DeleteGameobjectResponse>('DELETE', `/gameobjects/${encodePathSegment(objectId)}`)
}

export function addComponent(objectId: string, componentName: string): Promise<GameObjectResponse> {
  return request<GameObjectResponse>(
    'POST',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`
  )
}

export function deleteComponent(
  objectId: string,
  componentName: string
): Promise<GameObjectResponse> {
  return request<GameObjectResponse>(
    'DELETE',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`
  )
}

export function updateComponent(
  objectId: string,
  componentName: string,
  details: Record<string, unknown>
): Promise<GameObjectResponse> {
  return request<GameObjectResponse>(
    'PUT',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`,
    { details }
  )
}

export async function isBackendReachable(): Promise<boolean> {
  try {
    const res = await fetch(`${EDITOR_BACKEND_BASE_URL}/ping`)
    return res.ok
  } catch {
    return false
  }
}
