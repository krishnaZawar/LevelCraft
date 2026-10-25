// The shape of config/component-cards.yaml, as the generator understands it.

export const INPUT_TYPES = ['number', 'text', 'checkbox', 'select', 'slider'] as const
export type InputType = (typeof INPUT_TYPES)[number]

// Option lists the config cannot know statically. `inputKeys` is the engine's
// own key mapping, fetched from the backend, so the config never restates it.
export const OPTION_SOURCES = ['inputKeys'] as const
export type OptionSource = (typeof OPTION_SOURCES)[number]

export interface SelectOption {
  value: string
  label: string
}

// One comparison, deliberately not an expression language: it can be checked
// when the config loads and cannot loop.
export interface Condition {
  field: string
  equals: string | number | boolean
}

export interface FieldConfig {
  field: string
  label: string
  input: InputType
  help?: string
  // number and slider
  min?: number
  max?: number
  step?: number
  // select
  options?: SelectOption[]
  optionsFrom?: OptionSource
  visibleIf?: Condition
  // Values committed alongside this field when a checkbox is switched on.
  // Some components validate a flag and its dependent fields together, so
  // enabling the flag alone would be rejected.
  enables?: Record<string, string | number | boolean>
}

export interface RowConfig {
  fields: FieldConfig[]
  visibleIf?: Condition
}

export interface ComponentCardConfig {
  description?: string
  rows: RowConfig[]
}

export interface CardConfig {
  version: number
  components: Record<string, ComponentCardConfig>
}

// The version this build of the generator understands. A config declaring
// anything else is rejected rather than guessed at.
export const SUPPORTED_VERSION = 1

// A row holds at most this many fields: the Attributes panel is too narrow to
// stay legible past it.
export const MAX_FIELDS_PER_ROW = 4
