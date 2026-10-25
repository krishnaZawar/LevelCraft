// The widgets a card is built from. One per `input:` type in the config.
//
// Each takes the value it edits and a commit callback, and knows nothing about
// components or the config — the generator decides which to render.

import { useState } from 'react'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Slider } from '@/components/ui/slider'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import type { FieldConfig, SelectOption } from './types'

interface FieldProps<T> {
  config: FieldConfig
  value: T
  onCommit: (value: T) => void
}

function FieldLabel({ children }: { children: React.ReactNode }): React.JSX.Element {
  return (
    <span className="text-muted-foreground text-[10px] tracking-wide uppercase">{children}</span>
  )
}

/**
 * A number input that commits on blur and on Enter, reverts on Escape, and
 * clamps rather than rejecting — matching how most inspectors behave instead
 * of sending an out-of-range value and finding out the backend refused it.
 */
export function NumberField({ config, value, onCommit }: FieldProps<number>): React.JSX.Element {
  const [text, setText] = useState(String(value))
  const isInvalid = text.trim() !== '' && !Number.isFinite(Number(text))

  function commit(): void {
    const parsed = Number(text)
    if (!Number.isFinite(parsed)) {
      setText(String(value))
      return
    }
    let next = parsed
    if (config.min !== undefined) next = Math.max(config.min, next)
    if (config.max !== undefined) next = Math.min(config.max, next)
    if (next !== value) onCommit(next)
    setText(String(next))
  }

  return (
    <label className="flex items-center gap-2" title={config.help}>
      <FieldLabel>{config.label}</FieldLabel>
      <Input
        value={text}
        aria-label={config.label}
        aria-invalid={isInvalid}
        onChange={(e) => setText(e.target.value)}
        onBlur={commit}
        onKeyDown={(e) => {
          if (e.key === 'Enter') e.currentTarget.blur()
          if (e.key === 'Escape') {
            setText(String(value))
            e.currentTarget.blur()
          }
        }}
        className="h-7 flex-1 text-xs"
      />
    </label>
  )
}

export function TextField({ config, value, onCommit }: FieldProps<string>): React.JSX.Element {
  const [text, setText] = useState(value)

  return (
    <label className="flex items-center gap-2" title={config.help}>
      <FieldLabel>{config.label}</FieldLabel>
      <Input
        value={text}
        aria-label={config.label}
        onChange={(e) => setText(e.target.value)}
        onBlur={() => text !== value && onCommit(text)}
        onKeyDown={(e) => {
          if (e.key === 'Enter') e.currentTarget.blur()
          if (e.key === 'Escape') {
            setText(value)
            e.currentTarget.blur()
          }
        }}
        className="h-7 flex-1 text-xs"
      />
    </label>
  )
}

export function CheckboxField({ config, value, onCommit }: FieldProps<boolean>): React.JSX.Element {
  return (
    <label className="flex items-center gap-2" title={config.help}>
      <Checkbox
        checked={value}
        aria-label={config.label}
        onCheckedChange={(checked) => onCommit(checked === true)}
      />
      <span className="text-xs">{config.label}</span>
    </label>
  )
}

/**
 * Commits only when the drag ends. Committing per pixel would send a request
 * for every intermediate value the thumb passes through.
 */
export function SliderField({ config, value, onCommit }: FieldProps<number>): React.JSX.Element {
  const [draft, setDraft] = useState(value)
  // Synced from the prop during render rather than by remounting on each
  // change: a remount would drop focus, so holding an arrow key would commit
  // one step and then stop.
  const [lastValue, setLastValue] = useState(value)
  if (value !== lastValue) {
    setLastValue(value)
    setDraft(value)
  }

  const min = config.min ?? 0
  const max = config.max ?? 100

  return (
    <label className="flex items-center gap-2" title={config.help}>
      <FieldLabel>{config.label}</FieldLabel>
      <Slider
        value={[draft]}
        min={min}
        max={max}
        step={config.step ?? 1}
        aria-label={config.label}
        onValueChange={([next]) => setDraft(next)}
        onValueCommit={([next]) => next !== value && onCommit(next)}
        className="flex-1"
      />
      <span className="text-muted-foreground w-8 text-right text-xs tabular-nums">{draft}</span>
    </label>
  )
}

export function SelectField({
  config,
  value,
  options,
  onCommit
}: FieldProps<string> & { options: SelectOption[] }): React.JSX.Element {
  return (
    <label className="flex items-center gap-2" title={config.help}>
      <FieldLabel>{config.label}</FieldLabel>
      <Select value={value || undefined} onValueChange={onCommit}>
        <SelectTrigger size="sm" className="h-7 flex-1 text-xs" aria-label={config.label}>
          <SelectValue placeholder="None" />
        </SelectTrigger>
        <SelectContent>
          {options.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </label>
  )
}
