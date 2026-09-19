import { create } from 'zustand'
import { toast } from 'sonner'
import {
  addComponent as apiAddComponent,
  addGameobject as apiAddGameobject,
  deleteComponent as apiDeleteComponent,
  deleteGameobject as apiDeleteGameobject,
  duplicateGameobject as apiDuplicateGameobject,
  updateComponent as apiUpdateComponent,
  updateGameobject as apiUpdateGameobject,
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

// Drop-to-reorder, shared by the Hierarchy's objects and the Attributes
// panel's components. Dragging forward drops after the target and backward
// drops before it, so dropping onto the next sibling actually moves.
function moveWithin(order: string[], draggedId: string, targetId: string): string[] {
  const draggedFrom = order.indexOf(draggedId)
  const targetFrom = order.indexOf(targetId)
  if (draggedFrom === -1 || targetFrom === -1 || draggedId === targetId) return order

  const next = order.filter((id) => id !== draggedId)
  let insertAt = next.indexOf(targetId)
  if (draggedFrom < targetFrom) insertAt += 1
  next.splice(insertAt, 0, draggedId)
  return next
}

// The order components are shown in for one object, so a user can float the
// ones they touch often to the top. Session-only, like the Hierarchy's order:
// the backend stores components in a map and has no order of its own.
export function componentNamesInOrder(order: string[] | undefined, attached: string[]): string[] {
  if (!order) return attached
  const kept = order.filter((name) => attached.includes(name))
  return [...kept, ...attached.filter((name) => !kept.includes(name))]
}

interface EditorStoreState {
  gameObjects: GameState
  objectOrder: string[]
  // objectId -> the order its components are listed in
  componentOrder: Record<string, string[]>
  // Collapsed by component name rather than per object: collapsing Transform
  // is a statement about how much you care about Transform, not about one
  // object, and it should hold as you click between objects.
  collapsedComponents: string[]
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
  reorderComponents: (objectId: string, draggedName: string, targetName: string) => void
  toggleComponentCollapsed: (name: string) => void
  addGameobject: () => Promise<void>
  duplicateGameobject: (objectId: string) => Promise<void>
  deleteGameobject: (objectId: string) => Promise<void>
  updateGameobjectMetadata: (
    objectId: string,
    metadata: { name?: string; group?: string }
  ) => Promise<void>
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
  componentOrder: {},
  collapsedComponents: [],
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

  reorderObjects: (draggedId, targetId) =>
    set((state) => ({ objectOrder: moveWithin(state.objectOrder, draggedId, targetId) })),

  reorderComponents: (objectId, draggedName, targetName) =>
    set((state) => {
      const attached = Object.keys(state.gameObjects[objectId]?.components ?? {})
      const current = componentNamesInOrder(state.componentOrder[objectId], attached)
      return {
        componentOrder: {
          ...state.componentOrder,
          [objectId]: moveWithin(current, draggedName, targetName)
        }
      }
    }),

  toggleComponentCollapsed: (name) =>
    set((state) => ({
      collapsedComponents: state.collapsedComponents.includes(name)
        ? state.collapsedComponents.filter((collapsed) => collapsed !== name)
        : [...state.collapsedComponents, name]
    })),

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

  duplicateGameobject: async (objectId) => {
    try {
      const { objectDetails } = await apiDuplicateGameobject(objectId)
      set((state) => {
        // Placed directly after what it was copied from rather than at the
        // end, so the copy appears where the user is already looking.
        const order = [...state.objectOrder]
        const sourceAt = order.indexOf(objectId)
        order.splice(sourceAt === -1 ? order.length : sourceAt + 1, 0, objectDetails.id)
        return {
          gameObjects: { ...state.gameObjects, [objectDetails.id]: objectDetails },
          objectOrder: order,
          selectedObjectId: objectDetails.id
        }
      })
    } catch (err) {
      toast.error("Couldn't duplicate the object", { description: (err as Error).message })
    }
  },

  deleteGameobject: async (objectId) => {
    try {
      await apiDeleteGameobject(objectId)
      set((state) => {
        const gameObjects = { ...state.gameObjects }
        delete gameObjects[objectId]
        const componentOrder = { ...state.componentOrder }
        delete componentOrder[objectId]
        return {
          gameObjects,
          componentOrder,
          objectOrder: state.objectOrder.filter((id) => id !== objectId),
          selectedObjectId: state.selectedObjectId === objectId ? null : state.selectedObjectId
        }
      })
    } catch (err) {
      toast.error("Couldn't delete object", { description: (err as Error).message })
    }
  },

  updateGameobjectMetadata: async (objectId, metadata) => {
    // Optimistic for the same reason component edits are: a name typed into
    // the panel should appear in the Hierarchy as it is typed, not a
    // round-trip later.
    const previous = get().gameObjects[objectId]
    if (!previous) return
    set((state) => ({
      gameObjects: { ...state.gameObjects, [objectId]: { ...previous, ...metadata } }
    }))
    try {
      const { objectDetails } = await apiUpdateGameobject(objectId, metadata)
      set((state) => ({ gameObjects: { ...state.gameObjects, [objectId]: objectDetails } }))
    } catch (err) {
      set((state) => ({ gameObjects: { ...state.gameObjects, [objectId]: previous } }))
      toast.error("Couldn't update the object", { description: (err as Error).message })
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
      componentOrder: {},
      collapsedComponents: [],
      isLoading: false,
      loadError: null,
      selectedObjectId: null,
      availableComponents: []
    })
}))
