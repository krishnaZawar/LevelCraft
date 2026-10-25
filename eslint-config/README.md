# @levelcraft/eslint-config

The shared ESLint and Prettier baseline for every LevelCraft frontend app.

A rule is decided once here and applies everywhere. Apps do not restate it.

## Using it

Each app depends on this package by path, so there is nothing to publish:

```json
"devDependencies": {
  "@levelcraft/eslint-config": "file:../../eslint-config"
}
```

`eslint.config.mjs` in the app:

```js
import levelcraft from '@levelcraft/eslint-config'

export default [...levelcraft({ tsconfigRootDir: import.meta.dirname })]
```

`tsconfigRootDir` is required. The baseline lints with type information, which
means ESLint needs the app's TypeScript program to answer questions like "is
this value a promise". Without it the async and type rules cannot run at all,
so the config throws rather than silently degrading.

Prettier is shared the same way, via `.prettierrc` in the app:

```json
"@levelcraft/eslint-config/prettier"
```

## Adding app-specific rules

Append blocks after the shared config. Later blocks win, so this is where an
app narrows or relaxes a rule for its own files:

```js
export default [
  ...levelcraft({ tsconfigRootDir: import.meta.dirname }),
  {
    files: ['src/renderer/src/components/ui/**/*.{ts,tsx}'],
    rules: { '@typescript-eslint/explicit-function-return-type': 'off' }
  }
]
```

Use this for what is genuinely local to one app — vendored directories, a
generated folder. A rule that should apply to all frontends belongs in
`index.mjs`, not repeated in each app.

`react: false` drops the React/JSX layer, for a future non-React frontend.

## What the baseline covers

| Area | Rules |
| --- | --- |
| Type checking | `typescript-eslint` recommended, type-aware |
| Dead code | `no-unused-vars` (`_` prefix opts out), `no-unused-private-class-members` |
| Async | `no-floating-promises`, `no-misused-promises`, `await-thenable`, `require-await`, `return-await` |
| Types | `no-explicit-any`, `explicit-function-return-type`, `consistent-type-imports` |
| Complexity | `complexity` 15, `max-depth` 4, `max-nested-callbacks` 4, `max-params` 5 |
| Duplication | `sonarjs/no-identical-functions` and friends |
| Correctness | `no-extra-bind`, `eqeqeq`, `no-self-compare`, `no-unmodified-loop-condition`, … |
| Naming | `naming-convention` |
| React | `react`, `react-hooks` (rules of hooks), `react-refresh` |
| Formatting | `eslint-config-prettier` last, so Prettier owns layout |

## Two decisions worth knowing

**`no-misused-promises` does not check JSX attributes.** React cannot await an
event handler, so `onClick={async () => …}` is the normal shape rather than a
mistake. Handlers are expected to catch their own failures.
`no-floating-promises` still covers bare calls, which is where the real risk is.

**Plain `.js`/`.mjs` files skip type-aware rules.** Config and build scripts sit
outside the app's tsconfig, so there is no program to check them against.

## CI

`lint-frontend-apps` in `.github/workflows/ci.yaml` runs format check, typecheck
and `eslint . --max-warnings 0` for each app. Warnings fail CI: they are
advisory while you work, and blocking before merge, so advisory rules cannot
quietly accumulate.
