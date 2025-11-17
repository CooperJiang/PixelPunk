module.exports = {
  extends: 'stylelint-config-standard',
  rules: {
    'at-rule-no-unknown': [
      true,
      {
        ignoreAtRules: ['apply', 'tailwind', 'variants', 'responsive', 'screen'],
      },
    ],
  },
}
