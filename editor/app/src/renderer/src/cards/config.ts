// Loads the card config once, when the module is first imported.
//
// The YAML ships with the app and is parsed here rather than at build time, so
// the file stays readable as a file instead of becoming a bundler artifact.

import { parse } from 'yaml'
import rawConfig from '../../../../config/component-cards.yaml?raw'
import { validateCardConfig } from './validate'
import type { CardConfig, ComponentCardConfig, FieldConfig, InputType } from './types'

function load(): CardConfig {
  return validateCardConfig(parse(rawConfig))
}

// A broken config is a programming error, not a runtime condition, so it
// surfaces on load rather than as an empty panel later.
export const cardConfig: CardConfig = load()

export function getCardConfig(componentName: string): ComponentCardConfig | undefined {
  return cardConfig.components[componentName]
}

// Picks a widget for a value the config says nothing about, so a component
// added to the engine is editable before anyone lays its card out.
function inferInput(value: unknown): InputType {
  if (typeof value === 'boolean') return 'checkbox'
  if (typeof value === 'number') return 'number'
  return 'text'
}

/**
 * The rows to render for a component.
 *
 * Configured rows come first, in the order the config gives them. Any field
 * the component actually carries that the config does not mention is appended
 * one per row, so adding a field in the engine never makes it invisible — it
 * just makes it unstyled until someone lays it out.
 */
export function resolveRows(
  componentName: string,
  details: Record<string, unknown>
): {
  description?: string
  rows: { fields: FieldConfig[]; visibleIf?: FieldConfig['visibleIf'] }[]
} {
  const configured = getCardConfig(componentName)
  const rows = configured ? [...configured.rows] : []

  const laidOut = new Set(rows.flatMap((row) => row.fields.map((field) => field.field)))
  for (const [field, value] of Object.entries(details)) {
    if (laidOut.has(field)) continue
    rows.push({ fields: [{ field, label: field, input: inferInput(value) }] })
  }

  return { description: configured?.description, rows }
}
