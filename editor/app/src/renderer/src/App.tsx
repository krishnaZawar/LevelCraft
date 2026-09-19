import { useEffect } from 'react'
import { Toaster } from '@/components/ui/sonner'
import { TooltipProvider } from '@/components/ui/tooltip'
import { useEditorStore } from '@/store/editorStore'
import { useProjectStore } from '@/store/projectStore'
import { useRunStore } from '@/store/runStore'
import EditorShell from '@/views/EditorShell'
import Home from '@/views/Home'

function App(): React.JSX.Element {
  const activeProject = useProjectStore((state) => state.activeProject)
  const isProjectOpen = Boolean(activeProject)
  const resetEditorState = useEditorStore((state) => state.reset)
  const stopGame = useRunStore((state) => state.stop)

  useEffect(() => {
    if (isProjectOpen) {
      window.api.window.maximize()
    } else {
      window.api.window.unmaximize()
      // Clear any previous project's scene state so it can't linger and
      // briefly flash once a different project is opened next.
      resetEditorState()
      // The run belongs to the project just closed, with no way back to it.
      void stopGame()
    }
    window.api.menu.notifyProjectOpen(isProjectOpen)
  }, [isProjectOpen, resetEditorState, stopGame])

  return (
    <TooltipProvider delayDuration={200}>
      {isProjectOpen ? <EditorShell /> : <Home />}
      <Toaster position="bottom-right" />
    </TooltipProvider>
  )
}

export default App
