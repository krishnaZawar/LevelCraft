// The scene snapshot pushed by builder/backend, same shape editor/backend
// serves. Described separately per client: they share a contract, not code.
export interface GameObjectDetails {
  id: string
  name: string
  group: string
  components: Record<string, Record<string, unknown>>
}

export type GameState = Record<string, GameObjectDetails>
