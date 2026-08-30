import { ChildProcess, spawn } from 'child_process'
import { join } from 'path'
import { app } from 'electron'
import { is } from '@electron-toolkit/utils'

const BACKEND_URL = 'http://localhost:3000'
const HEALTH_CHECK_TIMEOUT_MS = 15_000
const HEALTH_CHECK_INTERVAL_MS = 300

let backendProcess: ChildProcess | null = null
// Only kill the process on quit if *we* spawned it — if it was already
// running (e.g. a developer started it by hand), leave it alone.
let weSpawnedBackend = false

async function pingBackend(): Promise<boolean> {
  try {
    const res = await fetch(`${BACKEND_URL}/ping`)
    return res.ok
  } catch {
    return false
  }
}

function resolveEditorBackendCommand(): { command: string; args: string[]; cwd?: string } {
  if (is.dev) {
    // Repo layout: <root>/editor/app (this app) and <root>/editor/backend
    const backendDir = join(app.getAppPath(), '..', 'backend')
    return { command: 'go', args: ['run', './cmd'], cwd: backendDir }
  }
  // Packaged: a binary pre-built by scripts/build-backend.mjs and bundled
  // via electron-builder's extraResources (see electron-builder.yml).
  const binName = process.platform === 'win32' ? 'editor-backend.exe' : 'editor-backend'
  return { command: join(process.resourcesPath, 'backend', binName), args: [] }
}

async function waitUntilReachable(timeoutMs: number): Promise<boolean> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    if (await pingBackend()) return true
    await new Promise((resolve) => setTimeout(resolve, HEALTH_CHECK_INTERVAL_MS))
  }
  return false
}

// Starts editor/backend if it isn't already reachable, and waits until it
// responds to /ping before resolving. Throws if it never comes up.
export async function ensureEditorBackendRunning(): Promise<void> {
  if (await pingBackend()) {
    console.log('[backend] editor/backend already reachable, not spawning a new instance')
    return
  }

  const { command, args, cwd } = resolveEditorBackendCommand()
  console.log('[backend] starting editor/backend:', command, args.join(' '))

  // `detached: true` on POSIX makes the child the leader of its own
  // process group, which matters because `go run` (used in dev) spawns
  // the actual compiled binary as a *grandchild* — killing just the `go
  // run` wrapper leaves that grandchild running, orphaned and
  // re-parented to PID 1. Killing the whole group (see stopEditorBackend)
  // takes the grandchild down too. No effect on the packaged/prebuilt
  // binary case, which has no such grandchild, but is safe there either way.
  const child = spawn(command, args, {
    cwd,
    stdio: 'pipe',
    detached: process.platform !== 'win32'
  })
  backendProcess = child
  weSpawnedBackend = true

  child.stdout?.on('data', (data) => console.log('[editor/backend]', data.toString().trimEnd()))
  child.stderr?.on('data', (data) => console.error('[editor/backend]', data.toString().trimEnd()))
  child.on('error', (err) =>
    console.error('[backend] failed to start editor/backend:', err.message)
  )
  child.on('exit', (code) => {
    console.log('[backend] editor/backend exited with code', code)
    if (backendProcess === child) backendProcess = null
  })

  const reachable = await waitUntilReachable(HEALTH_CHECK_TIMEOUT_MS)
  if (!reachable) {
    throw new Error(
      'editor/backend did not become reachable within 15s. Is Go installed and on PATH?'
    )
  }
}

export function stopEditorBackend(): void {
  if (backendProcess && weSpawnedBackend && backendProcess.pid) {
    if (process.platform === 'win32') {
      // No POSIX-style process groups on Windows — kill the whole tree by
      // PID instead (/t), otherwise a wrapper-spawned grandchild (see the
      // `go run` note above) survives the same way it did on POSIX.
      spawn('taskkill', ['/pid', String(backendProcess.pid), '/t', '/f'])
    } else {
      try {
        process.kill(-backendProcess.pid, 'SIGTERM')
      } catch {
        backendProcess.kill()
      }
    }
  }
  backendProcess = null
}
