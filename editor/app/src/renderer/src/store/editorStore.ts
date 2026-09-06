import { create } from 'zustand'
import { toast } from 'sonner'
import {
  addComponent as apiAddComponent,
  addGameobject as apiAddGameobject,
  deleteComponent as apiDeleteComponent,
  deleteGameobject as apiDeleteGameobject,
  updateComponent as apiUpdateComponent,
  getComponents,
  getGameState,
  GameState
} from '../api/gameApi'

// Keeps the Hierarchy's manual drag-reorder in sync with whatever ids
// actually exist right now — editor/backend has no persisted order field,
// so this is session-only and just needs to never go stale after an
// add/delete/reload: drop ids that no longer exist, append any new ones.
function syncOrder(order: string[], gameObjects: GameState): string[] {
  const ids = new Set(Object.keys(gameObjects))
  const kept = order.filter((id) => ids.has(id))
  const missing = Object.keys(gameObjects).filter((id) => !kept.includes(id))
  return [...kept, ...missing]
}

interface EditorStoreState {
  gameObjects: GameState
  objectOrder: string[]
  isLoading: boolean
  // Set only by fetchGameState — "the scene itself failed to load," shown
  // as a persistent inline banner with a retry action. Mutation failures
  // (add/delete/update) are transient instead, surfaced as toasts (see
  // below) so they don't blow away the Hierarchy/Workspace with a
  // full load-failure state over what's usually a one-off blip.
  loadError: string | null
  selectedObjectId: string | null
  availableComponents: string[]

  fetchGameState: () => Promise<void>
  fetchAvailableComponents: () => Promise<void>
  selectObject: (id: string | null) => void
  reorderObjects: (draggedId: string, targetId: string) => void
  addGameobject: () => Promise<void>
  deleteGameobject: (objectId: string) => Promise<void>
  addComponent: (objectId: string, componentName: string) => Promise<void>
  deleteComponent: (objectId: string, componentName: string) => Promise<void>
  updateComponent: (
    objectId: string,
    componentName: string,
    details: Record<string, unknown>
  ) => Promise<void>
  reset: () => void
}

// Mirrors editor/backend's in-memory scene state for the currently open
// project. Only meaningful once a project is loaded (see projectStore) —
// there's no implicit global scene to reflect otherwise.
export const useEditorStore = create<EditorStoreState>((set, get) => ({
  gameObjects: {},
  objectOrder: [],
  isLoading: false,
  loadError: null,
  selectedObjectId: null,
  availableComponents: [],

  fetchGameState: async () => {
    set({ isLoading: true, loadError: null })
    try {
      const { gameState } = await getGameState()
      set((state) => ({
        gameObjects: gameState,
        objectOrder: syncOrder(state.objectOrder, gameState),
        isLoading: false
      }))
    } catch (err) {
      set({ loadError: (err as Error).message, isLoading: false })
    }
  },

  fetchAvailableComponents: async () => {
    try {
      const { components } = await getComponents()
      set({ availableComponents: components })
    } catch {
      // Non-critical: the add-component dropdown just stays empty.
    }
  },

  selectObject: (id) => set({ selectedObjectId: id }),

  reorderObjects: (draggedId, targetId) => {
    if (draggedId === targetId) return
    set((state) => {
      const draggedFrom = state.objectOrder.indexOf(draggedId)
      const targetFrom = state.objectOrder.indexOf(targetId)
      if (draggedFrom === -1 || targetFrom === -1) return {}

      const order = state.objectOrder.filter((id) => id !== draggedId)
      let insertAt = order.indexOf(targetId)
      // Dragging forward drops after the target, dragging backward drops
      // before it — otherwise dropping onto the very next sibling recomputes
      // to the same position it started at and looks like nothing happened.
      if (draggedFrom < targetFrom) insertAt += 1
      order.splice(insertAt, 0, draggedId)
      return { objectOrder: order }
    })
  },

  addGameobject: async () => {
    try {
      const { objectDetails } = await apiAddGameobject()
      set((state) => ({
        gameObjects: { ...state.gameObjects, [objectDetails.id]: objectDetails },
        objectOrder: [...state.objectOrder, objectDetails.id],
        selectedObjectId: objectDetails.id
      }))
    } catch (err) {
      toast.error("Couldn't add object", { description: (err as Error).message })
    }
  },

  deleteGameobject: async (objectId) => {
    try {
      await apiDeleteGameobject(objectId)
      set((state) => {
        const gameObjects = { ...state.gameObjects }
        delete gameObjects[objectId]
        return {
          gameObjects,
          objectOrder: state.objectOrder.filter((id) => id !== objectId),
          selectedObjectId: state.selectedObjectId === objectId ? null : state.selectedObjectId
        }
      })
    } catch (err) {
      toast.error("Couldn't delete object", { description: (err as Error).message })
    }
  },

  addComponent: async (objectId, componentName) => {
    try {
      const { objectDetails } = await apiAddComponent(objectId, componentName)
      set((state) => ({ gameObjects: { ...state.gameObjects, [objectId]: objectDetails } }))
    } catch (err) {
      toast.error(`Couldn't add ${componentName}`, { description: (err as Error).message })
    }
  },

  deleteComponent: async (objectId, componentName) => {
    try {
      const { objectDetails } = await apiDeleteComponent(objectId, componentName)
      set((state) => ({ gameObjects: { ...state.gameObjects, [objectId]: objectDetails } }))
    } catch (err) {
      toast.error(`Couldn't remove ${componentName}`, { description: (err as Error).message })
    }
  },

  updateComponent: async (objectId, componentName, details) => {
    // Genuinely optimistic: applied to local state immediately, before the
    // network round-trip resolves — not just before/after bookkeeping.
    // Waiting for the server's response first (even if fast) leaves a gap
    // where the UI still shows the old value, which for canvas transforms
    // shows up as a visible flash back to the pre-edit state. Reconciled
    // with the server's response once it arrives, or reverted on failure.
    const previous = get().gameObjects[objectId]
    if (previous) {
      set((state) => ({
        gameObjects: {
          ...state.gameObjects,
          [objectId]: {
            ...previous,
            components: { ...previous.components, [componentName]: details }
          }
        }
      }))
    }
    try {
      const { objectDetails } = await apiUpdateComponent(objectId, componentName, details)
      set((state) => ({ gameObjects: { ...state.gameObjects, [objectId]: objectDetails } }))
    } catch (err) {
      set((state) => ({
        gameObjects: previous ? { ...state.gameObjects, [objectId]: previous } : state.gameObjects
      }))
      toast.error(`Couldn't update ${componentName}`, { description: (err as Error).message })
    }
  },

  reset: () =>
    set({
      gameObjects: {},
      objectOrder: [],
      isLoading: false,
      loadError: null,
      selectedObjectId: null,
      availableComponents: []
    })
}))
