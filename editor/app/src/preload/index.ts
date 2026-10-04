import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
import { type LevelCraftApi, type MenuAction } from '../shared/project'

const api: LevelCraftApi = {
  platform: process.platform,
  project: {
    getRoot: () => ipcRenderer.invoke('project:getRoot'),
    list: () => ipcRenderer.invoke('project:list'),
    create: (name) => ipcRenderer.invoke('project:create', name),
    openFromPath: (path) => ipcRenderer.invoke('project:openFromPath', path),
    getRecentPaths: () => ipcRenderer.invoke('project:getRecentPaths'),
    browseForFolder: () => ipcRenderer.invoke('project:browseForFolder')
  },
  window: {
    maximize: () => ipcRenderer.send('window:maximize'),
    unmaximize: () => ipcRenderer.send('window:unmaximize'),
    minimize: () => ipcRenderer.send('window:minimize'),
    reload: () => ipcRenderer.send('window:reload'),
    toggleDevTools: () => ipcRenderer.send('window:toggleDevTools'),
    toggleFullscreen: () => ipcRenderer.send('window:toggleFullscreen')
  },
  backend: {
    getBaseUrl: (): string => ipcRenderer.sendSync('backend:getBaseUrlSync') as string
  },
  orchestrator: {
    getBaseUrl: (): string => ipcRenderer.sendSync('orchestrator:getBaseUrlSync') as string
  },
  game: {
    createTempScene: (sourceName) => ipcRenderer.invoke('game:createTempScene', sourceName),
    clearTempScene: () => ipcRenderer.invoke('game:clearTempScene')
  },
  menu: {
    onAction: (callback) => {
      ipcRenderer.on('menu:action', (_event, action: MenuAction) => callback(action))
    },
    notifyProjectOpen: (isOpen) => ipcRenderer.send('menu:project-state', isOpen)
  }
}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error(error)
  }
} else {
  // @ts-expect-error — defined in index.d.ts, not on the real window type
  window.electron = electronAPI
  // @ts-expect-error — defined in index.d.ts, not on the real window type
  window.api = api
}
