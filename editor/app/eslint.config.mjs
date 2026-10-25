import levelcraft from '@levelcraft/eslint-config'

export default [
  ...levelcraft({ tsconfigRootDir: import.meta.dirname }),
  {
    // shadcn/ui components are vendored via `npx shadcn add` and regenerated
    // on update, so they aren't held to this project's stricter conventions.
    files: ['src/renderer/src/components/ui/**/*.{ts,tsx}'],
    rules: {
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      'react-refresh/only-export-components': 'off',
      'sonarjs/no-identical-functions': 'off'
    }
  }
]
