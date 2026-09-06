import { Menu as MenuIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { dispatchMenuAction } from './menuBridge'
import { useProjectStore } from './store/projectStore'

// Fallback for the native File/Edit/View/Window menu bar (see main/menu.ts)
// on platforms where it's hidden by default — autoHideMenuBar: true only
// reveals it via the Alt key, which isn't discoverable, and macOS is the
// only platform where it renders as an always-visible global bar. New/
// Open/Save/Close Project already have dedicated buttons elsewhere in the
// UI; the items below (Edit/View/Window) don't exist anywhere else.
function AppMenu(): React.JSX.Element | null {
  const platform = window.api.platform
  const isProjectOpen = Boolean(useProjectStore((state) => state.activeProject))

  if (platform === 'darwin') return null

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="size-7">
          <MenuIcon className="size-4" aria-hidden />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        <DropdownMenuLabel>File</DropdownMenuLabel>
        <DropdownMenuItem
          disabled={isProjectOpen}
          onClick={() => dispatchMenuAction('new-project')}
        >
          New Project...
          <DropdownMenuShortcut>Ctrl+N</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem
          disabled={isProjectOpen}
          onClick={() => dispatchMenuAction('open-project')}
        >
          Open Project...
          <DropdownMenuShortcut>Ctrl+O</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          disabled={!isProjectOpen}
          onClick={() => dispatchMenuAction('save-project')}
        >
          Save Project
          <DropdownMenuShortcut>Ctrl+S</DropdownMenuShortcut>
        </DropdownMenuItem>
        <DropdownMenuItem
          disabled={!isProjectOpen}
          onClick={() => dispatchMenuAction('close-project')}
        >
          Close Project
          <DropdownMenuShortcut>Ctrl+W</DropdownMenuShortcut>
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuLabel>Edit</DropdownMenuLabel>
        <DropdownMenuItem onClick={() => document.execCommand('undo')}>Undo</DropdownMenuItem>
        <DropdownMenuItem onClick={() => document.execCommand('redo')}>Redo</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => document.execCommand('cut')}>Cut</DropdownMenuItem>
        <DropdownMenuItem onClick={() => document.execCommand('copy')}>Copy</DropdownMenuItem>
        <DropdownMenuItem onClick={() => document.execCommand('paste')}>Paste</DropdownMenuItem>
        <DropdownMenuItem onClick={() => document.execCommand('selectAll')}>
          Select All
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuLabel>View</DropdownMenuLabel>
        <DropdownMenuItem onClick={() => window.api.window.reload()}>Reload</DropdownMenuItem>
        <DropdownMenuItem onClick={() => window.api.window.toggleDevTools()}>
          Toggle Developer Tools
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => window.api.window.toggleFullscreen()}>
          Toggle Fullscreen
        </DropdownMenuItem>

        <DropdownMenuSeparator />
        <DropdownMenuLabel>Window</DropdownMenuLabel>
        <DropdownMenuItem onClick={() => window.api.window.minimize()}>Minimize</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export default AppMenu
