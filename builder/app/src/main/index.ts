import { app, BrowserWindow, ipcMain, shell } from 'electron'
import { join } from 'path'
import { electronApp, is, optimizer } from '@electron-toolkit/utils'
import { startPingServer, stopPingServer } from './ping'

// Overrides the macOS menu bar name in dev mode (unpackaged builds
// otherwise default to "Electron" until a real app bundle exists).
app.setName('LevelCraft Player')

// The builder backend this app renders, from the environment the orchestrator
// spawned it with. No fallback: this app never starts a backend of its own.
const builderBackendBaseUrl = process.env.LEVELCRAFT_BUILDER_BACKEND_URL ?? ''

let mainWindow: BrowserWindow | null = null

function createWindow(): void {
  const win = new BrowserWindow({
    width: 960,
    height: 640,
    minWidth: 480,
    minHeight: 320,
    show: false,
    title: 'LevelCraft — Play',
    autoHideMenuBar: true,
    backgroundColor: '#121212',
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })
  mainWindow = win

  win.on('ready-to-show', () => win.show())

  win.on('closed', () => {
    if (mainWindow === win) mainWindow = null
  })

  win.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    win.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    win.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

app.whenReady().then(async () => {
  electronApp.setAppUserModelId('com.levelcraft.builder')

  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  ipcMain.on('backend:getBaseUrlSync', (event) => {
    event.returnValue = builderBackendBaseUrl
  })
  ipcMain.on('window:close', () => BrowserWindow.getFocusedWindow()?.close())

  await startPingServer()

  createWindow()

  app.on('activate', function () {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Closing the window ends the run on every platform, macOS included. The
// orchestrator notices this process is gone and reaps the backend with it.
app.on('window-all-closed', () => {
  app.quit()
})

app.on('will-quit', () => {
  stopPingServer()
})
