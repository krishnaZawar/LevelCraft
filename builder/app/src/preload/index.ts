import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
import { type LevelCraftBuilderApi } from '../shared/api'

const api: LevelCraftBuilderApi = {
  platform: process.platform,
  backend: {
    getBaseUrl: (): string => ipcRenderer.sendSync('backend:getBaseUrlSync') as string
  },
  window: {
    close: () => ipcRenderer.send('window:close')
  }
}

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
