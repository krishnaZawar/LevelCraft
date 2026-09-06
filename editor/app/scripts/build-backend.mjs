// Compiles editor/backend and builder/backend into binaries bundled into
// the packaged app via electron-builder's extraResources (see
// electron-builder.yml). Runs for the *host* platform/arch by default —
// building the other platforms' binaries (for CI cross-compiling the
// mac/win/linux electron-builder targets) is real additional work not
// attempted here; override with the GOOS/GOARCH env vars if you need one
// binary for a specific target.
import { execFileSync } from 'node:child_process'
import { existsSync, mkdirSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const APP_DIR = join(__dirname, '..')
const OUT_DIR = join(APP_DIR, 'resources', 'bin')

const goos = process.env.GOOS || (process.platform === 'win32' ? 'windows' : process.platform)
const exeSuffix = goos === 'windows' ? '.exe' : ''

const targets = [
  { name: 'editor-backend', dir: join(APP_DIR, '..', 'backend') },
  { name: 'builder-backend', dir: join(APP_DIR, '..', '..', 'builder', 'backend') }
]

mkdirSync(OUT_DIR, { recursive: true })

for (const { name, dir } of targets) {
  if (!existsSync(dir)) {
    console.error(`${name} source not found at ${dir}`)
    process.exit(1)
  }

  const outPath = join(OUT_DIR, name + exeSuffix)
  console.log(`Building ${name} (${goos}) -> ${outPath}`)

  try {
    execFileSync('go', ['build', '-o', outPath, './cmd'], {
      cwd: dir,
      stdio: 'inherit',
      env: process.env
    })
  } catch (err) {
    console.error(`\nFailed to build ${name}. Is Go installed and on PATH?\n` + err.message)
    process.exit(1)
  }
}

console.log('Backends built successfully.')
