import { ElectronAPI } from '@electron-toolkit/preload'
import { LevelCraftBuilderApi } from '../shared/api'

declare global {
  interface Window {
    electron: ElectronAPI
    api: LevelCraftBuilderApi
  }
}
