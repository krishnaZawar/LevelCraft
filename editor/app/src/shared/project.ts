// Types shared between the main process, preload bridge, and renderer.
// Kept as plain relative imports (no path alias) since main/preload/renderer
// are three separate build targets and only the renderer has an alias
// resolver configured for its own source tree.

export const PROJECT_SCHEMA_VERSION = 1

export interface ProjectManifest {
  id: string
  name: string
  schemaVersion: number
  createdAt: string
  lastOpenedAt: string
}

export interface ProjectSummary {
  id: string
  name: string
  path: string
  scenePath: string
  lastOpenedAt: string
}

export type ProjectOperationError =
  'invalid-name' | 'already-exists' | 'not-found' | 'invalid-project' | 'io-error'

export type ProjectApiResult =
  | { ok: true; project: ProjectSummary }
  | { ok: false; error: ProjectOperationError; message: string }

// Actions the native File menu (and its keyboard accelerators) can
// trigger. The menu only knows "an action happened" — the renderer owns
// what each one actually does, via menuBridge.ts.
export interface TempSceneInfo {
  path: string
  scenePath: string
}

export type MenuAction = 'new-project' | 'open-project' | 'close-project' | 'save-project'

export interface LevelCraftApi {
  platform: NodeJS.Platform
  project: {
    getRoot: () => Promise<string>
    list: () => Promise<ProjectSummary[]>
    create: (name: string) => Promise<ProjectApiResult>
    openFromPath: (path: string) => Promise<ProjectApiResult>
    getRecentPaths: () => Promise<string[]>
    browseForFolder: () => Promise<string | null>
  }
  window: {
    maximize: () => void
    unmaximize: () => void
    minimize: () => void
    reload: () => void
    toggleDevTools: () => void
    toggleFullscreen: () => void
  }
  backend: {
    // Sync IPC: main already resolved this before creating the window.
    getBaseUrl: () => string
  }
  orchestrator: {
    // Sync IPC: main resolved this before the window. Empty when the app was
    // started outside the orchestrator, so no game can be run.
    getBaseUrl: () => string
  }
  game: {
    // Creates the throwaway project a run is played from. The scene itself
    // is written by the backend's save endpoint, into the returned scenePath.
    createTempScene: (sourceName: string) => Promise<TempSceneInfo>
    // Removes that folder once the run is over.
    clearTempScene: () => Promise<void>
  }
  menu: {
    onAction: (callback: (action: MenuAction) => void) => void
    notifyProjectOpen: (isOpen: boolean) => void
  }
}
