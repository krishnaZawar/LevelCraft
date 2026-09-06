import { Toaster as Sonner, ToasterProps } from 'sonner'

// Dark-only per Phase 0's theme decision (see docs/client/personality.md) —
// no light/system branching needed, unlike shadcn's default template.
function Toaster({ ...props }: ToasterProps) {
  return (
    <Sonner
      theme="dark"
      className="toaster group"
      style={
        {
          '--normal-bg': 'var(--popover)',
          '--normal-text': 'var(--popover-foreground)',
          '--normal-border': 'var(--border)'
        } as React.CSSProperties
      }
      {...props}
    />
  )
}

export { Toaster }
