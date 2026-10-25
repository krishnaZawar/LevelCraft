// Validates component-cards.yaml when it loads.
//
// Every failure is reported rather than thrown one at a time, so a config with
// three mistakes takes one fix round rather than three. A card that silently
// drops a field is worse than one that refuses to load: the field looks absent
// rather than broken.

import {
  INPUT_TYPES,
  MAX_FIELDS_PER_ROW,
  OPTION_SOURCES,
  SUPPORTED_VERSION,
  type CardConfig,
  type Condition,
  type FieldConfig,
  type InputType,
  type RowConfig,
  type SelectOption
} from './types'

export class CardConfigError extends Error {
  constructor(problems: string[]) {
    super(`component-cards.yaml is invalid:\n  - ${problems.join('\n  - ')}`)
    this.name = 'CardConfigError'
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function validateCondition(raw: unknown, where: string, problems: string[]): Condition | undefined {
  if (raw === undefined) return undefined
  if (!isRecord(raw)) {
    problems.push(`${where}: visibleIf must be a mapping`)
    return undefined
  }
  const { field, equals } = raw
  if (typeof field !== 'string' || field === '') {
    problems.push(`${where}: visibleIf.field must be a non-empty string`)
    return undefined
  }
  if (typeof equals !== 'string' && typeof equals !== 'number' && typeof equals !== 'boolean') {
    problems.push(`${where}: visibleIf.equals must be a string, number or boolean`)
    return undefined
  }
  return { field, equals }
}

function validateOptions(
  raw: unknown,
  where: string,
  problems: string[]
): SelectOption[] | undefined {
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
    const label = typeof entry.label === 'string' ? entry.label : entry.value
    options.push({ value: entry.value, label })
  }
  return options
}

function validateField(raw: unknown, where: string, problems: string[]): FieldConfig | undefined {
  if (!isRecord(raw)) {
    problems.push(`${where}: each field must be a mapping`)
    return undefined
  }

  const { field, label, input } = raw
  if (typeof field !== 'string' || field === '') {
    problems.push(`${where}: "field" must name a key on the component`)
    return undefined
  }
  if (!INPUT_TYPES.includes(input as InputType)) {
    problems.push(`${where}: unknown input "${String(input)}" (expected ${INPUT_TYPES.join(', ')})`)
    return undefined
  }

  const config: FieldConfig = {
    field,
    label: typeof label === 'string' ? label : field,
    input: input as InputType,
    visibleIf: validateCondition(raw.visibleIf, where, problems)
  }

  for (const key of ['min', 'max', 'step'] as const) {
    if (raw[key] === undefined) continue
    if (typeof raw[key] !== 'number') {
      problems.push(`${where}: ${key} must be a number`)
      continue
    }
    config[key] = raw[key]
  }

  if (typeof raw.help === 'string') config.help = raw.help

  if (raw.enables !== undefined) {
    if (!isRecord(raw.enables)) {
      problems.push(`${where}: enables must be a mapping of field to value`)
    } else if (config.input !== 'checkbox') {
      problems.push(`${where}: enables only applies to a checkbox`)
    } else {
      const enables: Record<string, string | number | boolean> = {}
      for (const [key, value] of Object.entries(raw.enables)) {
        if (typeof value !== 'string' && typeof value !== 'number' && typeof value !== 'boolean') {
          problems.push(`${where}: enables.${key} must be a string, number or boolean`)
          continue
        }
        enables[key] = value
      }
      config.enables = enables
    }
  }

  // A slider with no track is not a slider, so both bounds are required even
  // though they are optional for a plain number input.
  if (config.input === 'slider' && (config.min === undefined || config.max === undefined)) {
    problems.push(`${where}: a slider needs both "min" and "max"`)
  }
  if (config.min !== undefined && config.max !== undefined && config.min > config.max) {
    problems.push(`${where}: min (${config.min}) is above max (${config.max})`)
  }

  if (raw.optionsFrom !== undefined) {
    if (!OPTION_SOURCES.includes(raw.optionsFrom as (typeof OPTION_SOURCES)[number])) {
      problems.push(`${where}: unknown optionsFrom "${String(raw.optionsFrom)}"`)
    } else {
      config.optionsFrom = raw.optionsFrom as (typeof OPTION_SOURCES)[number]
    }
  }
  config.options = validateOptions(raw.options, where, problems)

  if (config.input === 'select' && !config.options && !config.optionsFrom) {
    problems.push(`${where}: a select needs "options" or "optionsFrom"`)
  }

  return config
}

function validateRow(raw: unknown, where: string, problems: string[]): RowConfig | undefined {
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
    const field = validateField(rawField, `${where}.fields[${index}]`, problems)
    if (field) fields.push(field)
  }

  return { fields, visibleIf: validateCondition(raw.visibleIf, where, problems) }
}

/**
 * Checks a parsed config and returns it in the generator's shape.
 *
 * @throws CardConfigError listing every problem found.
 */
export function validateCardConfig(raw: unknown): CardConfig {
  const problems: string[] = []

  if (!isRecord(raw)) throw new CardConfigError(['the file must contain a mapping'])

  if (raw.version !== SUPPORTED_VERSION) {
    problems.push(`version must be ${SUPPORTED_VERSION}, found ${String(raw.version)}`)
  }
  if (!isRecord(raw.components)) {
    throw new CardConfigError([
      ...problems,
      '"components" must be a mapping of component name to card'
    ])
  }

  const components: Record<string, { description?: string; rows: RowConfig[] }> = {}

  for (const [name, rawCard] of Object.entries(raw.components)) {
    if (!isRecord(rawCard) || !Array.isArray(rawCard.rows)) {
      problems.push(`${name}: needs a "rows" list`)
      continue
    }

    const rows: RowConfig[] = []
    const seen = new Set<string>()

    for (const [index, rawRow] of rawCard.rows.entries()) {
      const row = validateRow(rawRow, `${name}.rows[${index}]`, problems)
      if (!row) continue
      for (const field of row.fields) {
        // Two inputs bound to one key would fight over the same value.
        if (seen.has(field.field)) {
          problems.push(`${name}: "${field.field}" is laid out more than once`)
        }
        seen.add(field.field)
      }
      rows.push(row)
    }

    // A condition on a field the card never shows can never become true, which
    // means the row it guards is dead rather than conditional.
    for (const row of rows) {
      const conditions = [row.visibleIf, ...row.fields.map((field) => field.visibleIf)]
      for (const condition of conditions) {
        if (condition && !seen.has(condition.field)) {
          problems.push(
            `${name}: visibleIf refers to "${condition.field}", which this card never shows`
          )
        }
      }
      for (const field of row.fields) {
        for (const key of Object.keys(field.enables ?? {})) {
          if (!seen.has(key)) {
            problems.push(`${name}: "${field.field}" enables "${key}", which this card never shows`)
          }
        }
      }
    }

    components[name] = {
      description: typeof rawCard.description === 'string' ? rawCard.description : undefined,
      rows
    }
  }

  if (problems.length > 0) throw new CardConfigError(problems)

  return { version: SUPPORTED_VERSION, components }
}
