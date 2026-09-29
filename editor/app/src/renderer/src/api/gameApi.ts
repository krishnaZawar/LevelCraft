// Resolved by the main process from the environment the orchestrator spawned
// this app with; the fallback only covers non-Electron contexts like tests.
const EDITOR_BACKEND_BASE_URL = window.api?.backend.getBaseUrl() ?? 'http://localhost:3000'

export interface GameObjectDetails {
  id: string
  name: string
  group: string
  components: Record<string, Record<string, unknown>>
}

export type GameState = Record<string, GameObjectDetails>

// What the backend actually sends. Its models are strictly typed, so a scene
// is a list of objects and each object's components are a list of
// name/data pairs — neither is keyed the way the UI wants to look things up.
interface WireComponent {
  name: string
  data: Record<string, unknown>
}

interface WireGameObject {
  id: string
  name: string
  group: string
  components: WireComponent[] | null
}

// Keys the wire shape by name/id, which is how every panel addresses it:
// the Attributes panel looks a component up by name, and the stores look an
// object up by id.
function toGameObject(wire: WireGameObject): GameObjectDetails {
  const components: Record<string, Record<string, unknown>> = {}
  for (const component of wire.components ?? []) {
    components[component.name] = component.data ?? {}
  }
  return { id: wire.id, name: wire.name, group: wire.group, components }
}

function toGameState(wire: WireGameObject[] | null): GameState {
  const state: GameState = {}
  for (const object of wire ?? []) {
    state[object.id] = toGameObject(object)
  }
  return state
}

interface SaveGameResponse {
  success: boolean
  message: string
}

interface WireGameStateResponse {
  success: boolean
  gameState: WireGameObject[] | null
}

interface WireGameObjectResponse {
  success: boolean
  objectDetails: WireGameObject
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

// fetch() throws a plain TypeError when it can't reach the server at all. This
// is the only place we call fetch, so treating any throw that way is accurate.
export const BACKEND_UNREACHABLE_MESSAGE =
  "Couldn't reach the LevelCraft backend. It may have stopped running."

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

export async function getGameState(): Promise<GetGameStateResponse> {
  const res = await request<WireGameStateResponse>('GET', '/game/state')
  return { success: res.success, gameState: toGameState(res.gameState) }
}

export function saveGame(filepath: string): Promise<SaveGameResponse> {
  return request<SaveGameResponse>('POST', '/game/save', { filepath })
}

export async function loadGame(filepath: string): Promise<LoadGameResponse> {
  const res = await request<WireGameStateResponse>('POST', '/game/load', { filepath })
  return { success: res.success, gameState: toGameState(res.gameState) }
}

export function getComponents(): Promise<GetComponentsResponse> {
  return request<GetComponentsResponse>('GET', '/components/')
}

export async function addGameobject(): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>('POST', '/gameobjects/')
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

// Only the fields actually edited are sent: the backend leaves out anything
// omitted, so a rename can't blank the group as a side effect.
export async function updateGameobject(
  objectId: string,
  metadata: { name?: string; group?: string }
): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>(
    'PUT',
    `/gameobjects/${encodePathSegment(objectId)}`,
    metadata
  )
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

// The copy is named by the backend, which is what can see the rest of the
// scene and so what can pick a name nothing else is using.
export async function duplicateGameobject(objectId: string): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>(
    'POST',
    `/gameobjects/${encodePathSegment(objectId)}/duplicate`
  )
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

export function deleteGameobject(objectId: string): Promise<DeleteGameobjectResponse> {
  return request<DeleteGameobjectResponse>('DELETE', `/gameobjects/${encodePathSegment(objectId)}`)
}

export async function addComponent(
  objectId: string,
  componentName: string
): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>(
    'POST',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`
  )
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

export async function deleteComponent(
  objectId: string,
  componentName: string
): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>(
    'DELETE',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`
  )
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

export async function updateComponent(
  objectId: string,
  componentName: string,
  details: Record<string, unknown>
): Promise<GameObjectResponse> {
  const res = await request<WireGameObjectResponse>(
    'PUT',
    `/gameobjects/${encodePathSegment(objectId)}/components/${encodePathSegment(componentName)}`,
    { details }
  )
  return { success: res.success, objectDetails: toGameObject(res.objectDetails) }
}

export async function isBackendReachable(): Promise<boolean> {
  try {
    const res = await fetch(`${EDITOR_BACKEND_BASE_URL}/ping`)
    return res.ok
  } catch {
    return false
  }
}
