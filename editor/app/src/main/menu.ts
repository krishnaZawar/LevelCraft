import { BrowserWindow, ipcMain, Menu, MenuItemConstructorOptions } from 'electron'
import { MenuAction } from '../shared/project'

function sendAction(action: MenuAction): void {
  BrowserWindow.getFocusedWindow()?.webContents.send('menu:action', action)
}

// Rebuilt whenever the renderer tells us a project opened/closed, so
// "New Project"/"Open Project" and "Close Project"/"Save Project" enable
// and disable in step with what's actually possible right now — rather
// than being always-on and failing silently when clicked at the wrong
// time.
function buildMenu(isProjectOpen: boolean): Menu {
  const fileMenu: MenuItemConstructorOptions = {
    label: 'File',
    submenu: [
      {
        label: 'New Project...',
        accelerator: 'CmdOrCtrl+N',
        enabled: !isProjectOpen,
        click: () => sendAction('new-project')
      },
      {
        label: 'Open Project...',
        accelerator: 'CmdOrCtrl+O',
        enabled: !isProjectOpen,
        click: () => sendAction('open-project')
      },
      { type: 'separator' },
      {
        label: 'Save Project',
        accelerator: 'CmdOrCtrl+S',
        enabled: isProjectOpen,
        click: () => sendAction('save-project')
      },
      {
        label: 'Close Project',
        accelerator: 'CmdOrCtrl+W',
        enabled: isProjectOpen,
        click: () => sendAction('close-project')
      },
      { type: 'separator' },
      process.platform === 'darwin' ? { role: 'close' } : { role: 'quit' }
    ]
  }

  const editMenu: MenuItemConstructorOptions = {
    label: 'Edit',
    submenu: [
      { role: 'undo' },
      { role: 'redo' },
      { type: 'separator' },
      { role: 'cut' },
      { role: 'copy' },
      { role: 'paste' },
      { role: 'selectAll' }
    ]
  }

  const viewMenu: MenuItemConstructorOptions = {
    label: 'View',
    submenu: [
      { role: 'reload' },
      { role: 'toggleDevTools' },
      { type: 'separator' },
      { role: 'resetZoom' },
      { role: 'zoomIn' },
      { role: 'zoomOut' },
      { type: 'separator' },
      { role: 'togglefullscreen' }
    ]
  }

  const windowMenu: MenuItemConstructorOptions = {
    label: 'Window',
    submenu: [{ role: 'minimize' }, { role: 'zoom' }]
  }

  const template: MenuItemConstructorOptions[] = [
    ...(process.platform === 'darwin' ? [{ role: 'appMenu' } as MenuItemConstructorOptions] : []),
    fileMenu,
    editMenu,
    viewMenu,
    windowMenu
  ]

  return Menu.buildFromTemplate(template)
}

export function setApplicationMenu(isProjectOpen: boolean): void {
  Menu.setApplicationMenu(buildMenu(isProjectOpen))
}

export function registerMenuIpcHandlers(): void {
  ipcMain.on('menu:project-state', (_event, isOpen: boolean) => {
    setApplicationMenu(isOpen)
  })
}
