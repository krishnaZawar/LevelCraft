import { app, shell, BrowserWindow, dialog, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'
import { registerGameIpcHandlers, registerProjectIpcHandlers } from './ipc'
import { registerMenuIpcHandlers, setApplicationMenu } from './menu'
import { clearTempScene } from './project'
import { startPingServer, stopPingServer } from './ping'

// Overrides the macOS menu bar name in dev mode (unpackaged builds
// otherwise default to "Electron" until a real app bundle exists).
app.setName('LevelCraft')

let mainWindow: BrowserWindow | null = null

// This app spawns no processes at all. The orchestrator is the application's
// entry point: it starts both backends and both clients, and addresses them
// to each other through the environment.
const editorBackendBaseUrl = process.env.LEVELCRAFT_BACKEND_URL ?? ''
const orchestratorBaseUrl = process.env.LEVELCRAFT_ORCHESTRATOR_URL ?? ''

// Launching this app on its own leaves it with no backend to talk to, which
// is a startup failure rather than something to work around.
function assertEditorBackendAddressed(): void {
  if (editorBackendBaseUrl) return
  throw new Error(
    'LEVELCRAFT_BACKEND_URL is not set. Start LevelCraft through the ' +
      'orchestrator rather than launching the editor directly.'
  )
}

function createWindow(): void {
  // Create the browser window.
  const win = new BrowserWindow({
    width: 1100,
    height: 720,
    show: false,
    title: 'LevelCraft',
    autoHideMenuBar: true,
    // Inset traffic lights over app content, matching Linear/Arc/VS Code's
    // native-feeling chrome (see docs/client/personality.md). No Windows
    // equivalent implemented yet, so Windows keeps the standard title bar.
    ...(process.platform === 'darwin' ? { titleBarStyle: 'hiddenInset' } : {}),
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })
  mainWindow = win

  win.on('ready-to-show', () => {
    win.show()
  })

  win.on('closed', () => {
    if (mainWindow === win) mainWindow = null
  })

  win.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    win.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    win.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(async () => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.levelcraft.editor')

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  registerProjectIpcHandlers()
  registerGameIpcHandlers()
  registerMenuIpcHandlers()
  setApplicationMenu(false)

  // Opening a project fills the screen (editor shell wants the room);
  // closing back to the project picker restores the normal window size
  // rather than leaving a mostly-empty maximized picker screen.
  ipcMain.on('window:maximize', () => mainWindow?.maximize())
  ipcMain.on('window:unmaximize', () => mainWindow?.unmaximize())
  // Backs the in-app menu button (see AppMenu.tsx) — the fallback UI for
  // the View/Window menu actions on platforms where the native menu bar
  // is hidden by default (autoHideMenuBar: true; see createWindow above).
  ipcMain.on('window:minimize', () => BrowserWindow.getFocusedWindow()?.minimize())
  ipcMain.on('window:reload', () => BrowserWindow.getFocusedWindow()?.webContents.reload())
  ipcMain.on('window:toggleDevTools', () =>
    BrowserWindow.getFocusedWindow()?.webContents.toggleDevTools()
  )
  ipcMain.on('window:toggleFullscreen', () => {
    const win = BrowserWindow.getFocusedWindow()
    win?.setFullScreen(!win.isFullScreen())
  })
  ipcMain.on('backend:getBaseUrlSync', (event) => {
    event.returnValue = editorBackendBaseUrl
  })
  ipcMain.on('orchestrator:getBaseUrlSync', (event) => {
    event.returnValue = orchestratorBaseUrl
  })

  await startPingServer()

  try {
    assertEditorBackendAddressed()
  } catch (err) {
    dialog.showErrorBox('LevelCraft', `Couldn't start LevelCraft.\n\n${(err as Error).message}`)
    app.quit()
    return
  }

  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('will-quit', () => {
  // Best effort: quitting mid-run would otherwise leave one behind in temp.
  void clearTempScene()
  stopPingServer()
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
