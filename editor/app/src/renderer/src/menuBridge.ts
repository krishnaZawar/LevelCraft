import { MenuAction } from '../../shared/project'
import { saveGame } from './api/gameApi'
import { useProjectStore } from './store/projectStore'

// Shared by the native menu (macOS menu bar + keyboard accelerators) and
// AppMenu.tsx's in-app fallback (Windows/Linux, where the native menu bar
// is hidden by default) — one place owning what each action actually does.
export function dispatchMenuAction(action: MenuAction): void {
  const store = useProjectStore.getState()

  switch (action) {
    case 'new-project':
      store.setNewProjectDialogOpen(true)
      break

    case 'open-project':
      window.api.project.browseForFolder().then((path) => {
        if (path) useProjectStore.getState().openProjectFromPath(path)
      })
      break

    case 'close-project':
      store.closeProject()
      break

    case 'save-project':
      if (store.activeProject) {
        saveGame(store.activeProject.scenePath).catch((err: Error) => {
          useProjectStore.setState({ error: err.message })
        })
      }
      break
  }
}

// Dispatches native File-menu actions (and their keyboard accelerators)
// into the project store. Lives outside React so the menu can reach the
// store directly via getState()/setState() rather than needing a
// component in the tree to be listening.
export function initMenuBridge(): void {
  window.api.menu.onAction(dispatchMenuAction)
}
