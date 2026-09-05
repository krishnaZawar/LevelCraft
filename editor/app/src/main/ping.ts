import { createServer, Server } from 'http'

// Lets external process managers (the orchestrator's `editorFrontend`
// process type, see orchestrator/internal/base/const.go) health-check this
// Electron app the same way editor/backend is checked: GET /ping -> 200.
let server: Server | null = null
let baseUrl = ''

export function startPingServer(): Promise<string> {
  return new Promise((resolve, reject) => {
    const s = createServer((req, res) => {
      if (req.url === '/ping') {
        res.writeHead(200, { 'Content-Type': 'application/json' })
        res.end(JSON.stringify({ success: true }))
        return
      }
      res.writeHead(404)
      res.end()
    })
    s.on('error', reject)
    s.listen(0, () => {
      const address = s.address()
      if (address === null || typeof address === 'string') {
        s.close()
        reject(new Error('failed to determine a free port for the ping server'))
        return
      }
      server = s
      baseUrl = `http://localhost:${address.port}`
      resolve(baseUrl)
    })
  })
}

export function getPingServerBaseUrl(): string {
  return baseUrl
}

export function stopPingServer(): void {
  server?.close()
  server = null
}
