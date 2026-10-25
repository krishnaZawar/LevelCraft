// The one ESLint baseline every LevelCraft frontend shares.
//
// Apps consume it rather than restating it, so a rule is decided once here and
// takes effect everywhere. An app may still append its own blocks after this
// config for anything genuinely local to it (see editor/app/eslint.config.mjs).

import js from '@eslint/js'
import tseslint from 'typescript-eslint'
import globals from 'globals'
import sonarjs from 'eslint-plugin-sonarjs'
import eslintPluginReact from 'eslint-plugin-react'
import eslintPluginReactHooks from 'eslint-plugin-react-hooks'
import eslintPluginReactRefresh from 'eslint-plugin-react-refresh'
import eslintConfigPrettier from 'eslint-config-prettier'

// Build output and vendored code, which no app should ever lint.
const IGNORES = ['**/node_modules/**', '**/dist/**', '**/out/**', '**/*.d.ts']

// Rules that need no type information. Cheap, and they run on config files too.
const correctnessRules = {
  'no-extra-bind': 'error',
  'no-extra-label': 'error',
  'no-lone-blocks': 'error',
  'no-self-compare': 'error',
  'no-unmodified-loop-condition': 'error',
  'no-unreachable-loop': 'error',
  'no-template-curly-in-string': 'error',
  eqeqeq: ['error', 'always', { null: 'ignore' }],
  'no-console': ['warn', { allow: ['warn', 'error'] }]
}

// Complexity ceilings. These are deliberately generous: they exist to catch a
// function that has quietly become unreviewable, not to force a style.
const complexityRules = {
  complexity: ['warn', 15],
  'max-depth': ['error', 4],
  'max-nested-callbacks': ['error', 4],
  'max-params': ['warn', 5]
}

// Dead code. Unused names prefixed with _ are kept deliberately — a destructured
// rest sibling, or an interface method an implementation does not need.
const deadCodeRules = {
  'no-unused-private-class-members': 'error',
  '@typescript-eslint/no-unused-vars': [
    'error',
    {
      args: 'after-used',
      argsIgnorePattern: '^_',
      varsIgnorePattern: '^_',
      caughtErrors: 'all',
      caughtErrorsIgnorePattern: '^_',
      destructuredArrayIgnorePattern: '^_'
    }
  ]
}

// Async correctness. A floating promise is the one that actually bites: a
// rejected request with no catch surfaces as an unhandled rejection far from
// the call that caused it.
const asyncRules = {
  '@typescript-eslint/no-floating-promises': 'error',
  // React has no way to await an event handler, so an `async onClick` is the
  // normal shape rather than a mistake. Handlers are expected to catch their
  // own failures; the floating-promise rule below still covers bare calls.
  '@typescript-eslint/no-misused-promises': ['error', { checksVoidReturn: { attributes: false } }],
  '@typescript-eslint/await-thenable': 'error',
  '@typescript-eslint/require-await': 'error',
  '@typescript-eslint/return-await': ['error', 'in-try-catch'],
  '@typescript-eslint/promise-function-async': 'off'
}

// Type discipline.
const typeRules = {
  '@typescript-eslint/no-explicit-any': 'error',
  '@typescript-eslint/explicit-function-return-type': [
    'error',
    { allowExpressions: true, allowTypedFunctionExpressions: true }
  ],
  '@typescript-eslint/consistent-type-imports': [
    'error',
    { prefer: 'type-imports', fixStyle: 'inline-type-imports' }
  ],
  '@typescript-eslint/consistent-type-definitions': 'off',
  '@typescript-eslint/no-unnecessary-condition': 'off'
}

// Naming. Loose enough to describe what the codebase already does, strict
// enough that a new file cannot drift from it.
const namingRules = {
  '@typescript-eslint/naming-convention': [
    'error',
    // Variables may be camelCase, or UPPER_CASE for module constants, or
    // PascalCase for a component or a class held in a const.
    {
      selector: 'variable',
      format: ['camelCase', 'UPPER_CASE', 'PascalCase'],
      leadingUnderscore: 'allow'
    },
    { selector: 'function', format: ['camelCase', 'PascalCase'] },
    // PascalCase is allowed because a React component passed as a prop is a
    // parameter, and renaming it would break the JSX convention.
    { selector: 'parameter', format: ['camelCase', 'PascalCase'], leadingUnderscore: 'allow' },
    { selector: 'typeLike', format: ['PascalCase'] },
    { selector: 'enumMember', format: ['PascalCase', 'UPPER_CASE'] },
    // Object keys are data as often as they are identifiers — a component name,
    // a wire field, a CSS custom property — so they are left alone.
    { selector: 'objectLiteralProperty', format: null },
    { selector: 'typeProperty', format: null }
  ]
}

// Duplicate code. ESLint core cannot see this; sonarjs can.
const duplicationRules = {
  'sonarjs/no-identical-functions': 'error',
  'sonarjs/no-all-duplicated-branches': 'error',
  'sonarjs/no-identical-conditions': 'error',
  'sonarjs/no-identical-expressions': 'error',
  'sonarjs/no-redundant-boolean': 'error',
  'sonarjs/no-collapsible-if': 'warn',
  'sonarjs/prefer-immediate-return': 'warn'
}

const reactRules = {
  ...eslintPluginReactHooks.configs.recommended.rules,
  ...eslintPluginReactRefresh.configs.vite.rules,
  'react/prop-types': 'off',
  'react/self-closing-comp': 'error',
  'react/jsx-no-useless-fragment': 'error',
  'react/jsx-curly-brace-presence': ['error', { props: 'never', children: 'never' }]
}

/**
 * The shared config, as a flat-config array an app spreads into its own.
 *
 * @param {object} options
 * @param {string} options.tsconfigRootDir  the app root, for type-aware linting
 * @param {boolean} [options.react]         include the React/JSX layer
 * @returns {import('eslint').Linter.Config[]}
 */
export default function levelcraftConfig({ tsconfigRootDir, react = true }) {
  if (!tsconfigRootDir) {
    throw new Error('@levelcraft/eslint-config: tsconfigRootDir is required for type-aware linting')
  }

  return tseslint.config([
    { ignores: IGNORES },

    js.configs.recommended,

    // Type-aware linting is what makes the async and type rules above possible;
    // without a program behind it ESLint cannot tell a promise from a value.
    ...tseslint.configs.recommendedTypeChecked,

    {
      languageOptions: {
        globals: { ...globals.browser, ...globals.node, ...globals.es2024 },
        parserOptions: {
          projectService: true,
          tsconfigRootDir
        }
      },
      plugins: { sonarjs },
      rules: {
        ...correctnessRules,
        ...complexityRules,
        ...deadCodeRules,
        ...asyncRules,
        ...typeRules,
        ...namingRules,
        ...duplicationRules
      }
    },

    // Config and build scripts are plain Node modules outside the app's
    // tsconfig, so type-aware rules have no program to run against.
    {
      files: ['**/*.{js,mjs,cjs}'],
      extends: [tseslint.configs.disableTypeChecked],
      languageOptions: { globals: globals.node },
      rules: {
        '@typescript-eslint/explicit-function-return-type': 'off',
        '@typescript-eslint/naming-convention': 'off'
      }
    },

    ...(react
      ? [
          eslintPluginReact.configs.flat.recommended,
          eslintPluginReact.configs.flat['jsx-runtime'],
          {
            settings: { react: { version: 'detect' } },
            plugins: {
              'react-hooks': eslintPluginReactHooks,
              'react-refresh': eslintPluginReactRefresh
            },
            rules: reactRules
          }
        ]
      : []),

    // Always last: turns off everything Prettier already decides.
    eslintConfigPrettier
  ])
}
