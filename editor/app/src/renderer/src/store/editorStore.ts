import { create } from 'zustand'
import { GameState, getGameState } from '../api/gameApi'

interface EditorStoreState {
  gameObjects: GameState
  isLoading: boolean
  error: string | null

  fetchGameState: () => Promise<void>
  reset: () => void
}

// Mirrors editor/backend's in-memory scene state for the currently open
// project. Only meaningful once a project is loaded (see projectStore) —
// there's no implicit global scene to reflect otherwise.
export const useEditorStore = create<EditorStoreState>((set) => ({
  gameObjects: {},
  isLoading: false,
  error: null,

  fetchGameState: async () => {
    set({ isLoading: true, error: null })
    try {
      const { gameState } = await getGameState()
      set({ gameObjects: gameState, isLoading: false })
    } catch (err) {
      set({ error: (err as Error).message, isLoading: false })
    }
  },

  reset: () => set({ gameObjects: {}, isLoading: false, error: null })
}))
