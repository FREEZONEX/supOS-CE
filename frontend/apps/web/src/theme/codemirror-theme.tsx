import { createTheme } from '@uiw/codemirror-themes';
// 导入高亮标签
import { tags as t } from '@lezer/highlight';

// 创建并导出CodeMirror主题
export const codemirrorTheme = createTheme({
  // 主题类型：light（浅色主题）
  theme: 'light',
  // 基础设置
  settings: {
    // 编辑器背景色
    background: 'var(--ui-modal-color)',
    fontFamily: "'IBM Plex Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
    // 背景图片
    // backgroundImage: '',
    // 文本前景色
    // foreground: '#75baff',
    // 光标颜色
    caret: 'var(--ui-text-color)',
    // 选中文本背景色
    // selection: '#036dd626',
    // 匹配选中文本的背景色
    // selectionMatch: '#036dd626',
    // 当前行高亮颜色
    lineHighlight: 'transparent',
    // 行号区域背景色
    gutterBackground: 'var(--ui-modal-color)',
    // 行号文本颜色
    gutterForeground: 'var(--ui-text-color)',
    // 行号与内容区分割线（暗色主题由 --ui-codemirror-gutter-border 覆盖）
    gutterBorder: 'var(--ui-codemirror-gutter-border, var(--ui-line-color))',
  },
  // 语法高亮样式配置
  styles: [
    // 注释样式
    { tag: t.comment, color: '#787b8099' },
    // 变量名样式
    { tag: t.variableName, color: '#0080ff' },
    // 字符串和特殊括号样式
    { tag: [t.string, t.special(t.brace)], color: '#5c6166' },
    // 数字样式
    { tag: t.number, color: '#5c6166' },
    // 布尔值样式
    { tag: t.bool, color: '#5c6166' },
    // null值样式
    { tag: t.null, color: '#5c6166' },
    // 关键字样式
    { tag: t.keyword, color: '#5c6166' },
    // 操作符样式
    { tag: t.operator, color: '#5c6166' },
    // 类名样式
    { tag: t.className, color: '#5c6166' },
    // 类型定义样式
    { tag: t.definition(t.typeName), color: '#5c6166' },
    // 类型名样式
    { tag: t.typeName, color: '#5c6166' },
    // 尖括号样式
    { tag: t.angleBracket, color: '#5c6166' },
    // 标签名样式
    { tag: t.tagName, color: '#5c6166' },
    // 属性名样式
    { tag: t.attributeName, color: '#5c6166' },
  ],
});
