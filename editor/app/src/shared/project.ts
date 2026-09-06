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
  }
  backend: {
    // Synchronous IPC on purpose: main resolves editor/backend's base URL
    // (see main/backend.ts) before creating the window, so it's always
    // available by the time the renderer's module graph evaluates —
    // no async bootstrap needed before the API client can be used.
    getBaseUrl: () => string
  }
  menu: {
    onAction: (callback: (action: MenuAction) => void) => void
    notifyProjectOpen: (isOpen: boolean) => void
  }
}
