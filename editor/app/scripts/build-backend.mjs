// Compiles editor/backend into a binary bundled into the packaged app via
// electron-builder's extraResources (see electron-builder.yml). Runs for
// the *host* platform/arch by default — building the other platforms'
// binaries (for CI cross-compiling the mac/win/linux electron-builder
// targets) is real additional work not attempted here; override with the
// GOOS/GOARCH env vars if you need one binary for a specific target.
import { execFileSync } from 'node:child_process'
import { existsSync, mkdirSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const APP_DIR = join(__dirname, '..')
const BACKEND_DIR = join(APP_DIR, '..', 'backend')
const OUT_DIR = join(APP_DIR, 'resources', 'bin')

const goos = process.env.GOOS || (process.platform === 'win32' ? 'windows' : process.platform)
const binName = goos === 'windows' ? 'editor-backend.exe' : 'editor-backend'
const outPath = join(OUT_DIR, binName)

if (!existsSync(BACKEND_DIR)) {
  console.error(`editor/backend not found at ${BACKEND_DIR}`)
  process.exit(1)
}

mkdirSync(OUT_DIR, { recursive: true })

console.log(`Building editor/backend (${goos}) -> ${outPath}`)

try {
  execFileSync('go', ['build', '-o', outPath, './cmd'], {
    cwd: BACKEND_DIR,
    stdio: 'inherit',
    env: process.env
  })
} catch (err) {
  console.error('\nFailed to build editor/backend. Is Go installed and on PATH?\n' + err.message)
  process.exit(1)
}

console.log('editor/backend built successfully.')
