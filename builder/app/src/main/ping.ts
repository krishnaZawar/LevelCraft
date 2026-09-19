import { createServer, Server } from 'http'

// Lets the orchestrator health-check this app like every other process:
// GET /ping -> 200, on the port it handed down via LEVELCRAFT_PING_PORT.
let server: Server | null = null

export function startPingServer(): Promise<void> {
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
      server = s
      resolve()
    })
  })
}

export function stopPingServer(): void {
  server?.close()
  server = null
}
