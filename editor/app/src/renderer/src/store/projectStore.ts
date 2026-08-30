import { create } from 'zustand'
import { ProjectSummary } from '../../../shared/project'
import { isBackendReachable, loadGame } from '../api/gameApi'

interface ProjectStoreState {
  activeProject: ProjectSummary | null
  projects: ProjectSummary[]
  isLoading: boolean
  /** Set only by refreshProjects — "the list itself couldn't load." */
  listError: string | null
  /** Set by createProject/openProject — scoped to that one action. */
  error: string | null
  backendReachable: boolean | null
  // Lifted out of the dialog component so the native File menu (and its
  // Cmd+N accelerator) can open it too, not just its own trigger button.
  isNewProjectDialogOpen: boolean

  checkBackend: () => Promise<void>
  refreshProjects: () => Promise<void>
  createProject: (name: string) => Promise<boolean>
  openProject: (project: ProjectSummary) => Promise<boolean>
  openProjectFromPath: (path: string) => Promise<boolean>
  closeProject: () => void
  setNewProjectDialogOpen: (open: boolean) => void
}

export const useProjectStore = create<ProjectStoreState>((set, get) => ({
  activeProject: null,
  projects: [],
  isLoading: false,
  listError: null,
  error: null,
  backendReachable: null,
  isNewProjectDialogOpen: false,

  checkBackend: async () => {
    const reachable = await isBackendReachable()
    set({ backendReachable: reachable })
  },

  refreshProjects: async () => {
    set({ isLoading: true, listError: null })
    try {
      const projects = await window.api.project.list()
      set({ projects, isLoading: false })
    } catch (err) {
      set({ listError: (err as Error).message, isLoading: false })
    }
  },

  createProject: async (name) => {
    set({ isLoading: true, error: null })
    const result = await window.api.project.create(name)
    if (!result.ok) {
      set({ error: result.message, isLoading: false })
      return false
    }
    try {
      await loadGame(result.project.scenePath)
    } catch (err) {
      // The project folder now genuinely exists on disk even though we
      // couldn't load it as the active session — refresh the list so the
      // user can open it (retrying "Create" with this name would
      // otherwise dead-end on "already exists").
      await get().refreshProjects()
      await get().checkBackend()
      set({ error: (err as Error).message, isLoading: false })
      return false
    }
    set({ activeProject: result.project, isLoading: false })
    return true
  },

  openProject: async (project) => {
    return get().openProjectFromPath(project.path)
  },

  openProjectFromPath: async (path) => {
    set({ isLoading: true, error: null })
    const result = await window.api.project.openFromPath(path)
    if (!result.ok) {
      set({ error: result.message, isLoading: false })
      return false
    }
    try {
      await loadGame(result.project.scenePath)
    } catch (err) {
      await get().checkBackend()
      set({ error: (err as Error).message, isLoading: false })
      return false
    }
    set({ activeProject: result.project, isLoading: false })
    return true
  },

  closeProject: () => set({ activeProject: null, error: null }),

  setNewProjectDialogOpen: (open) => set({ isNewProjectDialogOpen: open })
}))
