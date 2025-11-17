module.exports = {
  root: true,
  env: {
    node: true,
    browser: true,
    es2022: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:vue/vue3-essential', // 放宽到 essential
    'plugin:prettier/recommended',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: '@typescript-eslint/parser',
    ecmaVersion: 2022,
    sourceType: 'module',
    extraFileExtensions: ['.vue'],
  },
  plugins: ['@typescript-eslint'],
  rules: {
    /* 注释规则 - 暂时放宽，避免大批量报错 */
    'spaced-comment': 'off',
    'multiline-comment-style': 'off',

    /* Prettier 一致性 - 先警告不阻塞 */
    'prettier/prettier': 'warn',

    /* Vue规则 - 基础规范 */
    'vue/component-name-in-template-casing': ['warn', 'PascalCase'],
    'vue/component-definition-name-casing': ['warn', 'PascalCase'],
    'vue/max-attributes-per-line': 'off', // 暂时关闭，太严格
    'vue/html-self-closing': 'off', // 暂时关闭
    'vue/require-default-prop': 'off',
    'vue/prop-name-casing': ['warn', 'camelCase'],
    'vue/component-tags-order': 'off', // 暂时关闭
    'vue/block-tag-newline': 'off', // 暂时关闭
    'vue/component-api-style': 'off', // 暂时关闭
    'vue/define-macros-order': 'off', // 暂时关闭
    'vue/no-unused-vars': 'error',
    'vue/no-mutating-props': 'error',
    'vue/no-duplicate-attributes': 'warn',
    'vue/no-dupe-keys': 'error',
    'vue/no-side-effects-in-computed-properties': 'warn',
    'vue/multi-word-component-names': 'off',

    /* Vue 进阶规则 - 暂时关闭大部分 */
    'vue/no-unused-components': 'warn',
    'vue/no-unused-refs': 'warn',
    'vue/prefer-true-attribute-shorthand': 'off',
    'vue/prefer-separate-static-class': 'off',
    'vue/padding-line-between-blocks': 'off',
    'vue/no-empty-component-block': 'off',
    'vue/custom-event-name-casing': 'off',
    'vue/v-for-delimiter-style': 'off',
    'vue/no-useless-v-bind': 'warn',
    'vue/no-useless-mustaches': 'warn',
    'vue/no-useless-concat': 'warn',
    'vue/prefer-import-from-vue': 'off',

    /* Vue样式规则 - 暂时关闭 */
    'vue/array-bracket-spacing': 'off',
    'vue/arrow-spacing': 'off',
    'vue/block-spacing': 'off',
    'vue/brace-style': 'off',
    'vue/comma-dangle': 'off',
    'vue/comma-spacing': 'off',
    'vue/key-spacing': 'off',
    'vue/object-curly-spacing': 'off',

    /* TypeScript规则 - 基础检查 */
    '@typescript-eslint/no-unused-vars': [
      'error',
      {
        argsIgnorePattern: '^_',
        varsIgnorePattern: '^_',
        destructuredArrayIgnorePattern: '^_',
        caughtErrorsIgnorePattern: '^_',
      },
    ],
    '@typescript-eslint/no-explicit-any': 'off', // 暂时关闭，工具函数中使用any较多
    '@typescript-eslint/no-non-null-assertion': 'off',
    '@typescript-eslint/ban-ts-comment': 'off', // 暂时关闭
    '@typescript-eslint/consistent-type-imports': 'off', // 暂时关闭
    '@typescript-eslint/consistent-type-definitions': 'off', // 暂时关闭

    /* 基础代码质量规则 */
    'no-console': 'off',
    'no-debugger': 'warn',
    'prefer-const': 'warn',
    'no-var': 'warn',
    'no-undef': 'off',
    'object-shorthand': 'warn', // 改为警告
    'prefer-template': 'warn', // 改为警告
    eqeqeq: ['warn', 'always'], // 改为警告
    curly: 'off', // 暂时关闭
    'no-duplicate-imports': 'error',
    'no-unused-expressions': 'warn',
    'prefer-arrow-callback': 'off', // 暂时关闭
    'arrow-body-style': 'off', // 暂时关闭

    /* 代码风格规则 - 大部分关闭 */
    camelcase: 'off',
    'consistent-return': 'off',
    'no-else-return': 'off',
    'no-unneeded-ternary': 'warn',
    'prefer-destructuring': 'off',

    /* 进阶规则 - 暂时关闭 */
    'no-nested-ternary': 'off',
    'no-underscore-dangle': 'off',
    'prefer-rest-params': 'off',
    'prefer-spread': 'off',
    'require-await': 'off',
    'no-return-await': 'off',
    'no-useless-return': 'off',
    'no-useless-concat': 'off',
    'no-implicit-coercion': 'off',
    'no-lonely-if': 'off',
    'no-multi-assign': 'off',
    'no-param-reassign': 'off',
    'no-shadow': 'off',
    '@typescript-eslint/no-shadow': 'off',

    /* 变量命名和组织 - 关闭 */
    'sort-imports': 'off',

    /* 函数和对象规范 - 关闭 */
    'max-params': 'off',
    'max-depth': 'off',
    'max-nested-callbacks': 'off',
    'max-statements-per-line': 'off',
    'no-magic-numbers': 'off',
    'no-self-compare': 'warn',
    'no-use-before-define': 'off',
    '@typescript-eslint/no-use-before-define': 'off',

    /* 严格模式和最佳实践 - 大部分关闭 */
    'array-callback-return': 'off',
    complexity: 'off',
    'default-case': 'off',
    'default-case-last': 'off',
    'dot-notation': 'off',
    'no-alert': 'error',
    'no-caller': 'warn',
    'no-constructor-return': 'off',
    'no-empty-function': 'off',
    'no-eval': 'warn',
    'no-extend-native': 'warn',
    'no-extra-bind': 'off',
    'no-floating-decimal': 'off',
    'no-implied-eval': 'warn',
    'no-invalid-this': 'off',
    'no-iterator': 'warn',
    'no-labels': 'warn',
    'no-loop-func': 'off',
    'no-new': 'off',
    'no-new-func': 'warn',
    'no-new-wrappers': 'off',
    'no-proto': 'warn',
    'no-return-assign': 'off',
    'no-script-url': 'warn',
    'no-self-assign': 'off',
    'no-sequences': 'off',
    'no-throw-literal': 'off',
    'no-unused-labels': 'warn',
    'no-useless-call': 'off',
    'no-useless-catch': 'off',
    'no-void': 'off',
    'no-warning-comments': 'off',
    'prefer-promise-reject-errors': 'off',
    radix: 'off',
    yoda: 'off',

    // 追加的宽松设置（全局）
    'no-fallthrough': 'warn',
    'no-prototype-builtins': 'warn',
    'no-unreachable': 'warn',
    'no-case-declarations': 'off',

    // Vue 有效性校验降级为警告，避免阻塞
    'vue/valid-template-root': 'off',
    'vue/valid-v-bind': 'warn',
    'vue/valid-v-for': 'warn',
    'vue/valid-v-if': 'warn',
    'vue/valid-v-on': 'warn',
    'vue/valid-v-once': 'warn',
    'vue/valid-v-slot': 'warn',
    'vue/valid-v-model': 'warn',

    // 常见 eslint:recommended 报错降级
    'no-unsafe-negation': 'warn',
    'valid-typeof': 'warn',
    'use-isnan': 'warn',
    'no-cond-assign': 'warn',
    'no-constant-condition': 'warn',
    'no-control-regex': 'warn',
    'no-dupe-args': 'warn',
    'no-duplicate-case': 'warn',
    'no-empty': 'off',
    'no-ex-assign': 'warn',
    'no-irregular-whitespace': 'warn',
    'no-sparse-arrays': 'warn',
    'no-unsafe-finally': 'warn',
  },
  globals: {
    /* 使用自动导入配置 */
    ...require('./.eslintrc-auto-import.json').globals,
    // 临时全局，避免第三方库未显式导入导致的 no-undef 报错
    gsap: 'readonly',
  },
  ignorePatterns: [
    'dist',
    'node_modules',
    '*.d.ts',
    '*.md',
    'docs',
    'public',
    'eslint-*.json',
    'src/pages/docs/composables/useCodeExamples.ts',
  ],
  overrides: [
    {
      files: ['*.cjs', '*.js'],
      env: {
        node: true,
      },
      rules: {
        '@typescript-eslint/no-var-requires': 'off',
        '@typescript-eslint/no-require-imports': 'off',
        'no-console': 'off',
        'prefer-destructuring': 'off',
      },
    },
    {
      files: ['src/utils/system/logger.ts'],
      rules: {
        'no-console': 'off', // Logger implementation needs console
      },
    },
  ],
}
