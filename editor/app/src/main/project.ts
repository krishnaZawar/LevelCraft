import { app } from 'electron'
import { randomUUID } from 'crypto'
import { existsSync } from 'fs'
import { mkdir, readdir, readFile, writeFile } from 'fs/promises'
import { join } from 'path'
import {
  PROJECT_SCHEMA_VERSION,
  ProjectApiResult,
  ProjectManifest,
  ProjectOperationError,
  ProjectSummary
} from '../shared/project'

const MANIFEST_FILE = 'project.json'
const SCENE_FILE = 'scene.json'
const RECENT_PROJECTS_FILE = 'recent-projects.json'

// Windows-reserved device names (case-insensitive), invalid as a filename on
// Windows even though they'd work fine on macOS/Linux
const RESERVED_NAMES = new Set([
  'CON',
  'PRN',
  'AUX',
  'NUL',
  'COM1',
  'COM2',
  'COM3',
  'COM4',
  'COM5',
  'COM6',
  'COM7',
  'COM8',
  'COM9',
  'LPT1',
  'LPT2',
  'LPT3',
  'LPT4',
  'LPT5',
  'LPT6',
  'LPT7',
  'LPT8',
  'LPT9'
])

// Unity's convention: ~/Documents/<App> Projects/ on both Windows and macOS,
// not a different location per OS (see docs/client/approach.md §8)
export function getProjectsRoot(): string {
  return join(app.getPath('documents'), 'LevelCraft', 'Projects')
}

async function ensureProjectsRoot(): Promise<string> {
  const root = getProjectsRoot()
  await mkdir(root, { recursive: true })
  return root
}

function validateProjectName(name: string): ProjectOperationError | null {
  const trimmed = name.trim()
  if (trimmed.length === 0) return 'invalid-name'
  if (trimmed.length > 100) return 'invalid-name'
  if (/[/\\:*?"<>|]/.test(trimmed)) return 'invalid-name'
  if (/[. ]$/.test(trimmed)) return 'invalid-name'
  if (RESERVED_NAMES.has(trimmed.toUpperCase())) return 'invalid-name'
  return null
}

function toSummary(manifest: ProjectManifest, projectPath: string): ProjectSummary {
  return {
    id: manifest.id,
    name: manifest.name,
    path: projectPath,
    scenePath: join(projectPath, SCENE_FILE),
    lastOpenedAt: manifest.lastOpenedAt
  }
}

async function readManifest(projectPath: string): Promise<ProjectManifest | null> {
  try {
    const raw = await readFile(join(projectPath, MANIFEST_FILE), 'utf-8')
    return JSON.parse(raw) as ProjectManifest
  } catch {
    return null
  }
}

async function writeManifest(projectPath: string, manifest: ProjectManifest): Promise<void> {
  await writeFile(join(projectPath, MANIFEST_FILE), JSON.stringify(manifest, null, 2), 'utf-8')
}

export async function listProjects(): Promise<ProjectSummary[]> {
  const root = await ensureProjectsRoot()
  const entries = await readdir(root, { withFileTypes: true })

  const summaries: ProjectSummary[] = []
  for (const entry of entries) {
    if (!entry.isDirectory()) continue
    const projectPath = join(root, entry.name)
    const manifest = await readManifest(projectPath)
    if (!manifest) continue
    summaries.push(toSummary(manifest, projectPath))
  }

  summaries.sort((a, b) => b.lastOpenedAt.localeCompare(a.lastOpenedAt))
  return summaries
}

export async function createProject(name: string): Promise<ProjectApiResult> {
  const trimmed = name.trim()
  const nameError = validateProjectName(trimmed)
  if (nameError) {
    return { ok: false, error: nameError, message: 'Project name is invalid.' }
  }

  const root = await ensureProjectsRoot()
  const projectPath = join(root, trimmed)

  if (existsSync(projectPath)) {
    return {
      ok: false,
      error: 'already-exists',
      message: `A project named "${trimmed}" already exists.`
    }
  }

  const now = new Date().toISOString()
  const manifest: ProjectManifest = {
    id: randomUUID(),
    name: trimmed,
    schemaVersion: PROJECT_SCHEMA_VERSION,
    createdAt: now,
    lastOpenedAt: now
  }

  try {
    await mkdir(projectPath, { recursive: false })
    await writeManifest(projectPath, manifest)
    await writeFile(join(projectPath, SCENE_FILE), '{}', 'utf-8')
  } catch (err) {
    return { ok: false, error: 'io-error', message: (err as Error).message }
  }

  return { ok: true, project: toSummary(manifest, projectPath) }
}

export async function openProjectFromPath(projectPath: string): Promise<ProjectApiResult> {
  const manifest = await readManifest(projectPath)
  if (!manifest) {
    return {
      ok: false,
      error: 'invalid-project',
      message: 'That folder does not contain a valid LevelCraft project.json.'
    }
  }

  manifest.lastOpenedAt = new Date().toISOString()
  try {
    await writeManifest(projectPath, manifest)
  } catch (err) {
    return { ok: false, error: 'io-error', message: (err as Error).message }
  }

  return { ok: true, project: toSummary(manifest, projectPath) }
}

// Recent-projects state is app preference state, not user project content,
// so it lives under userData (Application-Support-equivalent) rather than
// inside the Documents/LevelCraft folder (see docs/client/approach.md §8)
function recentProjectsFilePath(): string {
  return join(app.getPath('userData'), RECENT_PROJECTS_FILE)
}

export async function getRecentProjectPaths(): Promise<string[]> {
  try {
    const raw = await readFile(recentProjectsFilePath(), 'utf-8')
    const parsed = JSON.parse(raw)
    return Array.isArray(parsed) ? parsed.filter((p) => typeof p === 'string') : []
  } catch {
    return []
  }
}

export async function rememberRecentProject(projectPath: string): Promise<void> {
  const existing = await getRecentProjectPaths()
  const next = [projectPath, ...existing.filter((p) => p !== projectPath)].slice(0, 10)
  await writeFile(recentProjectsFilePath(), JSON.stringify(next, null, 2), 'utf-8')
}
