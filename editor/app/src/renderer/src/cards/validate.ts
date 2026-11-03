// Validates component-cards.yaml when it loads.
//
// Every failure is collected rather than thrown one at a time, so a config
// with three mistakes takes one fix round rather than three. A card that
// silently drops a field is worse than one that refuses to load: the field
// looks absent rather than broken.
//
// The checks are split one concern per function, each adding to the running
// `problems` list, so no single function carries the whole rule set.

import {
  INPUT_TYPES,
  MAX_FIELDS_PER_ROW,
  OPTION_SOURCES,
  SUPPORTED_VERSION,
  type CardConfig,
  type ComponentCardConfig,
  type Condition,
  type FieldConfig,
  type InputType,
  type OptionSource,
  type RowConfig,
  type SelectOption
} from './types'

export class CardConfigError extends Error {
  constructor(problems: string[]) {
    super(`component-cards.yaml is invalid:\n  - ${problems.join('\n  - ')}`)
    this.name = 'CardConfigError'
  }
}

type Problems = string[]

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function isScalar(value: unknown): value is string | number | boolean {
  return typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean'
}

// String() on an object yields "[object Object]", which tells an author
// nothing about what they actually wrote.
function describe(value: unknown): string {
  return JSON.stringify(value) ?? 'undefined'
}

function readCondition(raw: unknown, where: string, problems: Problems): Condition | undefined {
  if (raw === undefined) return undefined
  if (!isRecord(raw)) {
    problems.push(`${where}: visibleIf must be a mapping`)
    return undefined
  }
  if (typeof raw.field !== 'string' || raw.field === '') {
    problems.push(`${where}: visibleIf.field must be a non-empty string`)
    return undefined
  }
  if (!isScalar(raw.equals)) {
    problems.push(`${where}: visibleIf.equals must be a string, number or boolean`)
    return undefined
  }
  return { field: raw.field, equals: raw.equals }
}

function readOptions(raw: unknown, where: string, problems: Problems): SelectOption[] | undefined {
  if (raw === undefined) return undefined
  if (!Array.isArray(raw) || raw.length === 0) {
    problems.push(`${where}: options must be a non-empty list`)
    return undefined
  }

  const options: SelectOption[] = []
  for (const [index, entry] of raw.entries()) {
    if (!isRecord(entry) || typeof entry.value !== 'string') {
      problems.push(`${where}: options[${index}] needs a string "value"`)
      continue
    }
    options.push({
      value: entry.value,
      label: typeof entry.label === 'string' ? entry.label : entry.value
    })
  }
  return options
}

function readOptionSource(
  raw: unknown,
  where: string,
  problems: Problems
): OptionSource | undefined {
  if (raw === undefined) return undefined
  if (!OPTION_SOURCES.includes(raw as OptionSource)) {
    problems.push(`${where}: unknown optionsFrom ${describe(raw)}`)
    return undefined
  }
  return raw as OptionSource
}

function readEnables(
  raw: unknown,
  input: InputType,
  where: string,
  problems: Problems
): FieldConfig['enables'] {
  if (raw === undefined) return undefined
  if (!isRecord(raw)) {
    problems.push(`${where}: enables must be a mapping of field to value`)
    return undefined
  }
  if (input !== 'checkbox') {
    problems.push(`${where}: enables only applies to a checkbox`)
    return undefined
  }

  const enables: Record<string, string | number | boolean> = {}
  for (const [key, value] of Object.entries(raw)) {
    if (!isScalar(value)) {
      problems.push(`${where}: enables.${key} must be a string, number or boolean`)
      continue
    }
    enables[key] = value
  }
  return enables
}

// min/max/step, and the rules relating them to the widget.
function readBounds(
  raw: Record<string, unknown>,
  config: FieldConfig,
  where: string,
  problems: Problems
): void {
  for (const key of ['min', 'max', 'step'] as const) {
    const value = raw[key]
    if (value === undefined) continue
    if (typeof value !== 'number') {
      problems.push(`${where}: ${key} must be a number`)
      continue
    }
    config[key] = value
  }

  // A slider with no track is not a slider, so both bounds are required even
  // though they are optional on a plain number input.
  if (config.input === 'slider' && (config.min === undefined || config.max === undefined)) {
    problems.push(`${where}: a slider needs both "min" and "max"`)
  }
  if (config.min !== undefined && config.max !== undefined && config.min > config.max) {
    problems.push(`${where}: min (${config.min}) is above max (${config.max})`)
  }
}

function readField(raw: unknown, where: string, problems: Problems): FieldConfig | undefined {
  if (!isRecord(raw)) {
    problems.push(`${where}: each field must be a mapping`)
    return undefined
  }
  if (typeof raw.field !== 'string' || raw.field === '') {
    problems.push(`${where}: "field" must name a key on the component`)
    return undefined
  }
  if (!INPUT_TYPES.includes(raw.input as InputType)) {
    problems.push(
      `${where}: unknown input ${describe(raw.input)} (expected ${INPUT_TYPES.join(', ')})`
    )
    return undefined
  }

  const input = raw.input as InputType
  const config: FieldConfig = {
    field: raw.field,
    label: typeof raw.label === 'string' ? raw.label : raw.field,
    input,
    help: typeof raw.help === 'string' ? raw.help : undefined,
    visibleIf: readCondition(raw.visibleIf, where, problems),
    optionsFrom: readOptionSource(raw.optionsFrom, where, problems),
    options: readOptions(raw.options, where, problems),
    enables: readEnables(raw.enables, input, where, problems)
  }

  readBounds(raw, config, where, problems)

  if (input === 'select' && !config.options && !config.optionsFrom) {
    problems.push(`${where}: a select needs "options" or "optionsFrom"`)
  }

  return config
}

function readRow(raw: unknown, where: string, problems: Problems): RowConfig | undefined {
  if (!isRecord(raw) || !Array.isArray(raw.fields)) {
    problems.push(`${where}: each row needs a "fields" list`)
    return undefined
  }
  if (raw.fields.length === 0 || raw.fields.length > MAX_FIELDS_PER_ROW) {
    problems.push(
      `${where}: a row holds 1 to ${MAX_FIELDS_PER_ROW} fields, found ${raw.fields.length}`
    )
    return undefined
  }

  const fields: FieldConfig[] = []
  for (const [index, rawField] of raw.fields.entries()) {
    const field = readField(rawField, `${where}.fields[${index}]`, problems)
    if (field) fields.push(field)
  }

  return { fields, visibleIf: readCondition(raw.visibleIf, where, problems) }
}

// Checks that can only be made once every row of a component is known.
function checkCrossReferences(name: string, rows: RowConfig[], problems: Problems): void {
  const allFields = rows.flatMap((row) => row.fields)
  const laidOut = new Set(allFields.map((field) => field.field))

  const seen = new Set<string>()
  for (const field of allFields) {
    // Two inputs bound to one key would fight over the same value.
    if (seen.has(field.field)) {
      problems.push(`${name}: "${field.field}" is laid out more than once`)
    }
    seen.add(field.field)
  }

  // A condition on a field the card never shows can never become true, which
  // makes the row it guards dead rather than conditional.
  const conditions = rows.flatMap((row) => [row.visibleIf, ...row.fields.map((f) => f.visibleIf)])
  for (const condition of conditions) {
    if (condition && !laidOut.has(condition.field)) {
      problems.push(
        `${name}: visibleIf refers to "${condition.field}", which this card never shows`
      )
    }
  }

  for (const field of allFields) {
    for (const key of Object.keys(field.enables ?? {})) {
      if (!laidOut.has(key)) {
        problems.push(`${name}: "${field.field}" enables "${key}", which this card never shows`)
      }
    }
  }
}

function readComponent(
  name: string,
  raw: unknown,
  problems: Problems
): ComponentCardConfig | undefined {
  if (!isRecord(raw) || !Array.isArray(raw.rows)) {
    problems.push(`${name}: needs a "rows" list`)
    return undefined
  }

  const rows: RowConfig[] = []
  for (const [index, rawRow] of raw.rows.entries()) {
    const row = readRow(rawRow, `${name}.rows[${index}]`, problems)
    if (row) rows.push(row)
  }

  checkCrossReferences(name, rows, problems)

  return {
    description: typeof raw.description === 'string' ? raw.description : undefined,
    rows
  }
}

/**
 * Checks a parsed config and returns it in the generator's shape.
 *
 * @throws CardConfigError listing every problem found.
 */
export function validateCardConfig(raw: unknown): CardConfig {
  if (!isRecord(raw)) throw new CardConfigError(['the file must contain a mapping'])

  const problems: Problems = []

  if (raw.version !== SUPPORTED_VERSION) {
    problems.push(`version must be ${SUPPORTED_VERSION}, found ${describe(raw.version)}`)
  }
  if (!isRecord(raw.components)) {
    throw new CardConfigError([
      ...problems,
      '"components" must be a mapping of component name to card'
    ])
  }

  const components: Record<string, ComponentCardConfig> = {}
  for (const [name, rawCard] of Object.entries(raw.components)) {
    const card = readComponent(name, rawCard, problems)
    if (card) components[name] = card
  }

  if (problems.length > 0) throw new CardConfigError(problems)

  return { version: SUPPORTED_VERSION, components }
}
