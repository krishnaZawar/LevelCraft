import { createServer, Server } from 'http'

// Lets the orchestrator health-check this app: GET /ping -> 200.
let server: Server | null = null
let baseUrl = ''

export function startPingServer(): Promise<string> {
  // Bind to the orchestrator-assigned port if given, else pick a free one.
  const requestedPort = Number(process.env.LEVELCRAFT_PING_PORT) || 0

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
    s.listen(requestedPort, () => {
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
