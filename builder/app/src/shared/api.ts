// Types shared between the main process, preload bridge, and renderer. Plain
// relative imports: only the renderer has an alias resolver configured.

export interface LevelCraftBuilderApi {
  platform: NodeJS.Platform
  backend: {
    // Sync IPC: main resolved this from the environment before the window.
    getBaseUrl: () => string
  }
  window: {
    close: () => void
  }
}
