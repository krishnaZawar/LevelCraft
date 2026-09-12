import { create } from 'zustand'
import { toast } from 'sonner'
import { ProjectSummary } from '../../../shared/project'
import { saveGame } from '../api/gameApi'
import { getGameStatus, runGame, stopGame } from '../api/orchestratorApi'

// Polled rather than pushed: the orchestrator has no channel back here, and a
// run can end unseen — the player window closing, or the builder crashing.
const STATUS_POLL_INTERVAL_MS = 1500

export type RunStatus = 'idle' | 'starting' | 'running' | 'stopping'

interface RunStoreState {
  status: RunStatus

  run: (project: ProjectSummary) => Promise<void>
  stop: () => Promise<void>
}

let pollTimer: ReturnType<typeof setInterval> | null = null

function stopPolling(): void {
  if (pollTimer !== null) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

// Never throws: this runs on paths already handling a failure, and a stuck
// temp folder shouldn't mask whatever went wrong first.
async function clearTempScene(): Promise<void> {
  try {
    await window.api.game.clearTempScene()
  } catch (err) {
    console.error('failed to clear the temporary game project:', (err as Error).message)
  }
}

export const useRunStore = create<RunStoreState>((set, get) => ({
  status: 'idle',

  run: async (project) => {
    if (get().status !== 'idle') return
    set({ status: 'starting' })

    try {
      // A throwaway copy, so running never saves over the project's scene.
      const temp = await window.api.game.createTempScene(project.name)
      await saveGame(temp.scenePath)
      await runGame(temp.scenePath)
    } catch (err) {
      await clearTempScene()
      set({ status: 'idle' })
      toast.error("Couldn't run the game", { description: (err as Error).message })
      return
    }

    set({ status: 'running' })

    stopPolling()
    pollTimer = setInterval(async () => {
      // Only 'running' is worth asking about; stop() knows the answer.
      if (get().status !== 'running') return
      try {
        const { running } = await getGameStatus()
        if (running) return
      } catch {
        // Unreachable means nothing is managing the run either, so treat it
        // the same as the run ending.
      }
      stopPolling()
      await clearTempScene()
      set({ status: 'idle' })
    }, STATUS_POLL_INTERVAL_MS)
  },

  stop: async () => {
    const status = get().status
    if (status === 'idle' || status === 'stopping') return

    stopPolling()
    set({ status: 'stopping' })

    try {
      await stopGame()
    } catch (err) {
      toast.error("Couldn't stop the game", { description: (err as Error).message })
    }

    await clearTempScene()
    set({ status: 'idle' })
  }
}))
