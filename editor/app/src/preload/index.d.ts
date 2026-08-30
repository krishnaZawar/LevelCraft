import { ElectronAPI } from '@electron-toolkit/preload'
import { LevelCraftApi } from '../shared/project'

declare global {
  interface Window {
    electron: ElectronAPI
    api: LevelCraftApi
  }
}
