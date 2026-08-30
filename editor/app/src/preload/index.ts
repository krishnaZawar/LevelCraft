import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
import { LevelCraftApi } from '../shared/project'

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
    unmaximize: () => ipcRenderer.send('window:unmaximize')
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
  // @ts-ignore (define in dts)
  window.electron = electronAPI
  // @ts-ignore (define in dts)
  window.api = api
}
