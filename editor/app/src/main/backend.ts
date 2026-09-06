import { ChildProcess, spawn } from 'child_process'
import { createServer } from 'net'
import { join } from 'path'
import { app } from 'electron'
import { is } from '@electron-toolkit/utils'

// The port a manually-started (e.g. by a developer) editor/backend defaults
// to, per editor/backend/cmd/main.go's `--port` flag default. Used only as
// a first guess so we don't spawn a duplicate instance in that case — any
// instance *we* spawn gets a dynamically chosen free port instead, matching
// how the orchestrator spawns editor/backend (see orchestrator/internal/executor).
const DEFAULT_BACKEND_PORT = 3000
const HEALTH_CHECK_TIMEOUT_MS = 15_000
const HEALTH_CHECK_INTERVAL_MS = 300

let backendProcess: ChildProcess | null = null
// Only kill the process on quit if *we* spawned it — if it was already
// running (e.g. a developer started it by hand), leave it alone.
let weSpawnedBackend = false
let backendBaseUrl = `http://localhost:${DEFAULT_BACKEND_PORT}`

async function pingBackend(baseUrl: string): Promise<boolean> {
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

function resolveEditorBackendCommand(port: number): { command: string; args: string[]; cwd?: string } {
  if (is.dev) {
    // Repo layout: <root>/editor/app (this app) and <root>/editor/backend
    const backendDir = join(app.getAppPath(), '..', 'backend')
    return { command: 'go', args: ['run', './cmd', '--port', String(port)], cwd: backendDir }
  }
  // Packaged: a binary pre-built by scripts/build-backend.mjs and bundled
  // via electron-builder's extraResources (see electron-builder.yml).
  const binName = process.platform === 'win32' ? 'editor-backend.exe' : 'editor-backend'
  return {
    command: join(process.resourcesPath, 'backend', binName),
    args: ['--port', String(port)]
  }
}

async function waitUntilReachable(baseUrl: string, timeoutMs: number): Promise<boolean> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    if (await pingBackend(baseUrl)) return true
    await new Promise((resolve) => setTimeout(resolve, HEALTH_CHECK_INTERVAL_MS))
  }
  return false
}

// Returns the base URL of the editor/backend instance this app is using —
// only meaningful after ensureEditorBackendRunning() has resolved.
export function getEditorBackendBaseUrl(): string {
  return backendBaseUrl
}

// Starts editor/backend if it isn't already reachable, and waits until it
// responds to /ping before resolving. Throws if it never comes up.
export async function ensureEditorBackendRunning(): Promise<void> {
  const defaultUrl = `http://localhost:${DEFAULT_BACKEND_PORT}`
  if (await pingBackend(defaultUrl)) {
    console.log('[backend] editor/backend already reachable, not spawning a new instance')
    backendBaseUrl = defaultUrl
    return
  }

  const port = await getFreePort()
  backendBaseUrl = `http://localhost:${port}`

  const { command, args, cwd } = resolveEditorBackendCommand(port)
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

  const reachable = await waitUntilReachable(backendBaseUrl, HEALTH_CHECK_TIMEOUT_MS)
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
