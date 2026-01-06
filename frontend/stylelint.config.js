/** @type {import('stylelint').Config} */
export default {
  extends: ['stylelint-config-standard', 'stylelint-scss', 'stylelint-prettier/recommended'],
  rules: {
    'at-rule-no-unknown': [
      true,
      {
        ignoreAtRules: ['use', 'forward', 'import', 'include'],
      },
    ],
    'selector-class-pattern': null, // 禁用对类选择器命名的强制要求
    'selector-id-pattern': null, // 禁用对id选择器命名的强制要求
    'no-descending-specificity': null, // 禁用排序要求
    'selector-pseudo-class-no-unknown': [
      true,
      {
        ignorePseudoClasses: ['global'], // 允许使用global的全局伪类
      },
    ],
  },
  overrides: [
    {
      files: ['*.scss', '**/*.scss'],
      customSyntax: 'postcss-scss',
    },
  ],
};
