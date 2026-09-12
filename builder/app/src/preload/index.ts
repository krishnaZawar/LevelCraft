import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
import { LevelCraftBuilderApi } from '../shared/api'

const api: LevelCraftBuilderApi = {
  platform: process.platform,
  backend: {
    getBaseUrl: () => ipcRenderer.sendSync('backend:getBaseUrlSync')
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
  // @ts-ignore (define in dts)
  window.electron = electronAPI
  // @ts-ignore (define in dts)
  window.api = api
}
