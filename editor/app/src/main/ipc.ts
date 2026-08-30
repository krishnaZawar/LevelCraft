import { dialog, ipcMain } from 'electron'
import {
  createProject,
  getProjectsRoot,
  getRecentProjectPaths,
  listProjects,
  openProjectFromPath,
  rememberRecentProject
} from './project'

export function registerProjectIpcHandlers(): void {
  ipcMain.handle('project:getRoot', () => getProjectsRoot())

  ipcMain.handle('project:list', () => listProjects())

  ipcMain.handle('project:create', async (_event, name: string) => {
    const result = await createProject(name)
    if (result.ok) await rememberRecentProject(result.project.path)
    return result
  })

  ipcMain.handle('project:openFromPath', async (_event, projectPath: string) => {
    const result = await openProjectFromPath(projectPath)
    if (result.ok) await rememberRecentProject(result.project.path)
    return result
  })

  ipcMain.handle('project:getRecentPaths', () => getRecentProjectPaths())

  ipcMain.handle('project:browseForFolder', async () => {
    const result = await dialog.showOpenDialog({
      properties: ['openDirectory'],
      title: 'Open LevelCraft Project'
    })
    if (result.canceled || result.filePaths.length === 0) return null
    return result.filePaths[0]
  })
}
