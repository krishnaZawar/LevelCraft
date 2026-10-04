// The scene snapshot pushed by builder/backend, same shape editor/backend
// serves. Described separately per client: they share a contract, not code.
export interface GameObjectDetails {
  id: string
  name: string
  group: string
  components: Record<string, Record<string, unknown>>
}

export type GameState = Record<string, GameObjectDetails>

// What builder/backend actually pushes. Its models are strictly typed, so a
// scene is a list of objects and each object's components are a list of
// name/data pairs — neither is keyed the way the renderer looks things up.
export interface WireComponent {
  name: string
  data: Record<string, unknown>
}

export interface WireGameObject {
  id: string
  name: string
  group: string
  components: WireComponent[] | null
}

// Keys the snapshot by id, and each object's components by name, which is how
// renderScene addresses them.
export function toGameState(wire: WireGameObject[] | null): GameState {
  const state: GameState = {}
  for (const object of wire ?? []) {
    const components: Record<string, Record<string, unknown>> = {}
    for (const component of object.components ?? []) {
      components[component.name] = component.data ?? {}
    }
    state[object.id] = {
      id: object.id,
      name: object.name,
      group: object.group,
      components
    }
  }
  return state
}
