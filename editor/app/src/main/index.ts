import { ChildProcess, spawn } from 'child_process'
import { createServer } from 'net'
import { app, shell, BrowserWindow, dialog, ipcMain } from 'electron'
import { join } from 'path'
import { electronApp, optimizer, is } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'
import { registerProjectIpcHandlers } from './ipc'
import { startBuilder, stopBuilder } from './builder'
import { registerMenuIpcHandlers, setApplicationMenu } from './menu'
import { startPingServer, stopPingServer } from './ping'

// Overrides the macOS menu bar name in dev mode (unpackaged builds
// otherwise default to "Electron" until a real app bundle exists).
app.setName('LevelCraft')

let mainWindow: BrowserWindow | null = null
let playWindow: BrowserWindow | null = null

let editorBackendProcess: ChildProcess | null = null
// Only kill the process on quit if *we* spawned it — packaged builds
// always spawn their own; dev never does (see ensureEditorBackendRunning).
let weSpawnedEditorBackend = false
let editorBackendBaseUrl = ''

const EDITOR_BACKEND_HEALTH_CHECK_TIMEOUT_MS = 15_000
const EDITOR_BACKEND_HEALTH_CHECK_INTERVAL_MS = 300

async function pingEditorBackend(baseUrl: string): Promise<boolean> {
  try {
    const res = await fetch(`${baseUrl}/ping`)
    return res.ok
  } catch {
    return false
  }
}

// Asks the OS for an ephemeral port, then immediately releases it. Good
// enough for our purposes: the gap between release and the backend binding
// it is a well-known, accepted TOCTOU (the orchestrator's executor has the
// same retry-on-bind-failure fallback for exactly this reason).
function getFreePort(): Promise<number> {
  return new Promise((resolve, reject) => {
    const server = createServer()
    server.unref()
    server.on('error', reject)
    server.listen(0, () => {
      const address = server.address()
      if (address === null || typeof address === 'string') {
        server.close()
        reject(new Error('failed to determine a free port'))
        return
      }
      const { port } = address
      server.close(() => resolve(port))
    })
  })
}

async function waitUntilEditorBackendReachable(
  baseUrl: string,
  timeoutMs: number
): Promise<boolean> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    if (await pingEditorBackend(baseUrl)) return true
    await new Promise((resolve) => setTimeout(resolve, EDITOR_BACKEND_HEALTH_CHECK_INTERVAL_MS))
  }
  return false
}

// In dev, this app never spawns editor/backend itself — the orchestrator
// (`cd orchestrator && go run ./cmd`) owns starting both editor/backend and
// this app together and passes the backend's URL via LEVELCRAFT_BACKEND_URL.
// Running `npm run dev` directly, without the orchestrator, is no longer a
// supported dev workflow: it throws here instead of silently doing its own
// thing on the side.
//
// Packaged builds have no orchestrator wrapping them (an end user just
// launches the app), so they still spawn the bundled editor-backend binary
// themselves, same as before.
async function ensureEditorBackendRunning(): Promise<void> {
  const orchestratedUrl = process.env.LEVELCRAFT_BACKEND_URL
  if (orchestratedUrl) {
    console.log('[backend] using orchestrator-provided editor/backend at', orchestratedUrl)
    editorBackendBaseUrl = orchestratedUrl
    return
  }

  if (is.dev) {
    throw new Error(
      'LEVELCRAFT_BACKEND_URL is not set. Run the app via the orchestrator ' +
        '(cd orchestrator && go run ./cmd) instead of `npm run dev` directly.'
    )
  }

  const port = await getFreePort()
  editorBackendBaseUrl = `http://localhost:${port}`

  const binName = process.platform === 'win32' ? 'editor-backend.exe' : 'editor-backend'
  const command = join(process.resourcesPath, 'backend', binName)
  const args = ['--port', String(port)]
  console.log('[backend] starting editor/backend:', command, args.join(' '))

  const child = spawn(command, args, {
    stdio: 'pipe',
    detached: process.platform !== 'win32'
  })
  editorBackendProcess = child
  weSpawnedEditorBackend = true

  child.stdout?.on('data', (data) => console.log('[editor/backend]', data.toString().trimEnd()))
  child.stderr?.on('data', (data) => console.error('[editor/backend]', data.toString().trimEnd()))
  child.on('error', (err) =>
    console.error('[backend] failed to start editor/backend:', err.message)
  )
  child.on('exit', (code) => {
    console.log('[backend] editor/backend exited with code', code)
    if (editorBackendProcess === child) editorBackendProcess = null
  })

  const reachable = await waitUntilEditorBackendReachable(
    editorBackendBaseUrl,
    EDITOR_BACKEND_HEALTH_CHECK_TIMEOUT_MS
  )
  if (!reachable) {
    throw new Error('editor/backend did not become reachable within 15s.')
  }
}

function stopEditorBackend(): void {
  if (editorBackendProcess && weSpawnedEditorBackend && editorBackendProcess.pid) {
    if (process.platform === 'win32') {
      // No POSIX-style process groups on Windows — kill the whole tree by
      // PID instead (/t), otherwise a wrapper-spawned grandchild survives.
      spawn('taskkill', ['/pid', String(editorBackendProcess.pid), '/t', '/f'])
    } else {
      try {
        process.kill(-editorBackendProcess.pid, 'SIGTERM')
      } catch {
        editorBackendProcess.kill()
      }
    }
  }
  editorBackendProcess = null
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

// A small companion window for Play mode, popped out separately rather
// than replacing the main editor viewport — its own renderer process, so
// it owns its own WebSocket connection to builder/backend independently
// of the main window (see renderer/src/PlayWindowApp.tsx). Loaded from
// the same bundle as the main window via the '#play' hash, since this app
// has one renderer entry point, not a second electron-vite target.
function createPlayWindow(): void {
  if (playWindow) {
    playWindow.focus()
    return
  }

  const win = new BrowserWindow({
    width: 640,
    height: 420,
    minWidth: 480,
    minHeight: 320,
    show: false,
    title: 'LevelCraft — Play',
    autoHideMenuBar: true,
    parent: mainWindow ?? undefined,
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false
    }
  })
  playWindow = win

  win.on('ready-to-show', () => win.show())

  // Covers both the in-app Stop button (which just calls window.api.builder.stop(),
  // triggering ipcMain's 'builder:stop' -> playWindow.close() -> here) and the
  // user closing this window directly via its own close button — either way,
  // the backend process must not outlive this window.
  win.on('closed', () => {
    if (playWindow === win) playWindow = null
    stopBuilder()
    mainWindow?.webContents.send('builder:stopped')
  })

  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    win.loadURL(`${process.env['ELECTRON_RENDERER_URL']}#play`)
  } else {
    win.loadFile(join(__dirname, '../renderer/index.html'), { hash: 'play' })
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
  // Targets the focused window rather than always mainWindow, so it also
  // does the right thing from the Play window.
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
  ipcMain.handle('builder:launch', async (_event, scenePath: string) => {
    const result = await startBuilder(scenePath)
    if (result.ok) createPlayWindow()
    return result
  })
  // Closing the tracked play window (rather than calling stopBuilder()
  // directly) reuses its 'closed' handler above as the single place that
  // both kills the process and notifies the main window.
  ipcMain.on('builder:stop', () => playWindow?.close())

  await startPingServer()

  try {
    await ensureEditorBackendRunning()
  } catch (err) {
    dialog.showErrorBox(
      'LevelCraft',
      `Couldn't start the LevelCraft backend.\n\n${(err as Error).message}`
    )
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

// Only stop editor/backend if this app instance actually spawned it
// (packaged builds only — see ensureEditorBackendRunning above). Runs on
// real app quit, not just window close, so it doesn't kill the backend
// every time a macOS window closes while the app (and its menu bar) stays alive.
app.on('will-quit', () => {
  stopEditorBackend()
  stopBuilder()
  stopPingServer()
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
