// The generator: turns a component's config rows into its card body.
//
// This is the only place that decides how a field is edited, which is what
// keeps every card consistent no matter who adds the next component.

import { useEditorStore } from '@/store/editorStore'
import { resolveRows } from './config'
import { CheckboxField, NumberField, SelectField, SliderField, TextField } from './fields'
import type { Condition, FieldConfig, SelectOption } from './types'

interface CardBodyProps {
  componentName: string
  details: Record<string, unknown>
  // Takes a patch rather than one field, because switching a flag on can need
  // its dependent fields set in the same update to pass engine validation.
  onCommit: (patch: Record<string, unknown>) => void
}

function isVisible(condition: Condition | undefined, details: Record<string, unknown>): boolean {
  if (!condition) return true
  return details[condition.field] === condition.equals
}

function renderField(
  config: FieldConfig,
  details: Record<string, unknown>,
  optionsFor: (config: FieldConfig) => SelectOption[],
  onCommit: CardBodyProps['onCommit']
): React.JSX.Element {
  const value = details[config.field]
  // Text-like inputs are remounted whenever the stored value changes, so they
  // always start from what the component holds rather than a stale draft. The
  // slider syncs itself instead, because remounting it would drop focus
  // mid-adjustment.
  const key = config.input === 'slider' ? config.field : `${config.field}:${String(value)}`
  const commit = (next: unknown): void => onCommit({ [config.field]: next })

  switch (config.input) {
    case 'checkbox':
      return (
        <CheckboxField
          key={key}
          config={config}
          value={value === true}
          onCommit={(checked) =>
            onCommit(
              checked && config.enables
                ? { [config.field]: true, ...config.enables }
                : { [config.field]: checked }
            )
          }
        />
      )
    case 'select':
      return (
        <SelectField
          key={key}
          config={config}
          value={typeof value === 'string' ? value : ''}
          options={optionsFor(config)}
          onCommit={commit}
        />
      )
    case 'slider':
      return <SliderField key={key} config={config} value={Number(value ?? 0)} onCommit={commit} />
    case 'text':
      return <TextField key={key} config={config} value={String(value ?? '')} onCommit={commit} />
    case 'number':
      return <NumberField key={key} config={config} value={Number(value ?? 0)} onCommit={commit} />
  }
}

export function CardBody({
  componentName,
  details,
  onCommit
}: CardBodyProps): React.JSX.Element | null {
  const inputKeys = useEditorStore((state) => state.inputKeys)

  const { rows } = resolveRows(componentName, details)
  if (rows.length === 0) return null

  // Dynamic lists come from the engine rather than the config, so the key
  // mapping is stated once in the backend and never copied into YAML.
  const optionsFor = (config: FieldConfig): SelectOption[] =>
    config.optionsFrom === 'inputKeys' ? inputKeys : (config.options ?? [])

  return (
    <div className="space-y-2 px-2.5 pt-0.5 pb-2.5">
      {rows.map((row, rowIndex) => {
        if (!isVisible(row.visibleIf, details)) return null

        const fields = row.fields.filter((field) => isVisible(field.visibleIf, details))
        if (fields.length === 0) return null

        return (
          <div
            key={rowIndex}
            className="grid gap-x-3 gap-y-1.5"
            style={{ gridTemplateColumns: `repeat(${fields.length}, minmax(0, 1fr))` }}
          >
            {fields.map((field) => renderField(field, details, optionsFor, onCommit))}
          </div>
        )
      })}
    </div>
  )
}
