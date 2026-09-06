import { ChildProcess, spawn } from 'child_process'
import { join } from 'path'
import { app } from 'electron'
import { is } from '@electron-toolkit/utils'

// Fixed per docs/client/approach.md §4.4 — builder/backend has no --port
// flag, unlike editor/backend. A leaked process from a prior run blocks
// the next one from binding, so startBuilder always kills first.
const BUILDER_BASE_URL = 'http://localhost:8000'
const HEALTH_CHECK_TIMEOUT_MS = 15_000
const HEALTH_CHECK_INTERVAL_MS = 300

let builderProcess: ChildProcess | null = null

async function pingBuilder(): Promise<boolean> {
  try {
    const res = await fetch(`${BUILDER_BASE_URL}/ping`)
    return res.ok
  } catch {
    return false
  }
}

function resolveBuilderCommand(scenePath: string): {
  command: string
  args: string[]
  cwd?: string
} {
  if (is.dev) {
    // Repo layout: <root>/editor/app (this app) and <root>/builder/backend
    const builderDir = join(app.getAppPath(), '..', '..', 'builder', 'backend')
    return { command: 'go', args: ['run', './cmd', '--file-name', scenePath], cwd: builderDir }
  }
  // Packaged: a binary pre-built by scripts/build-backend.mjs, bundled
  // alongside editor-backend via electron-builder's extraResources.
  const binName = process.platform === 'win32' ? 'builder-backend.exe' : 'builder-backend'
  return {
    command: join(process.resourcesPath, 'backend', binName),
    args: ['--file-name', scenePath]
  }
}

async function waitUntilReachable(timeoutMs: number): Promise<boolean> {
  const deadline = Date.now() + timeoutMs
  while (Date.now() < deadline) {
    if (await pingBuilder()) return true
    await new Promise((resolve) => setTimeout(resolve, HEALTH_CHECK_INTERVAL_MS))
  }
  return false
}

// Same grandchild-process problem as stopEditorBackend (see backend.ts):
// `go run` in dev spawns the real binary as a grandchild, so killing just
// the wrapper leaves it orphaned. detached + process-group kill handles it.
function killBuilder(): void {
  if (builderProcess?.pid) {
    if (process.platform === 'win32') {
      spawn('taskkill', ['/pid', String(builderProcess.pid), '/t', '/f'])
    } else {
      try {
        process.kill(-builderProcess.pid, 'SIGTERM')
      } catch {
        builderProcess.kill()
      }
    }
  }
  builderProcess = null
}

// Starts builder/backend for the given scene, waiting until /ping responds
// before resolving. Kills any previously-spawned instance first — the port
// is fixed with no configurability (see BUILDER_BASE_URL above).
export async function startBuilder(scenePath: string): Promise<{ ok: boolean; message?: string }> {
  killBuilder()

  const { command, args, cwd } = resolveBuilderCommand(scenePath)
  console.log('[builder] starting builder/backend:', command, args.join(' '))

  const child = spawn(command, args, {
    cwd,
    stdio: 'pipe',
    detached: process.platform !== 'win32'
  })
  builderProcess = child

  child.stdout?.on('data', (data) => console.log('[builder/backend]', data.toString().trimEnd()))
  child.stderr?.on('data', (data) => console.error('[builder/backend]', data.toString().trimEnd()))
  child.on('error', (err) =>
    console.error('[builder] failed to start builder/backend:', err.message)
  )
  child.on('exit', (code) => {
    console.log('[builder] builder/backend exited with code', code)
    if (builderProcess === child) builderProcess = null
  })

  const reachable = await waitUntilReachable(HEALTH_CHECK_TIMEOUT_MS)
  if (!reachable) {
    killBuilder()
    return { ok: false, message: 'builder/backend did not become reachable within 15s.' }
  }
  return { ok: true }
}

export function stopBuilder(): void {
  killBuilder()
}
