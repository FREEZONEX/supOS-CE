const DARK_SCOPE = 'html[data-tier0-dark-theme="1"]';

const readThemeToken = (name: string, fallback: string) => {
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  return value || fallback;
};

export const isNotebookDarkThemeActive = () =>
  document.documentElement.classList.contains('dark') || document.documentElement.classList.contains('chartreuseDark');

const buildThemeVariables = (
  bg: string,
  text: string,
  panelBg: string,
  mutedBg: string,
  mutedText: string,
  border: string,
  input: string,
  accentBg: string,
  editorBg: string,
  primary: string,
  primaryForeground: string,
  success: string,
  error: string
) => `
  color-scheme: dark !important;
  --background: ${bg} !important;
  --foreground: ${text} !important;
  --muted: ${mutedBg} !important;
  --muted-foreground: ${mutedText} !important;
  --popover: ${panelBg} !important;
  --popover-foreground: ${text} !important;
  --card: ${panelBg} !important;
  --card-foreground: ${text} !important;
  --border: ${border} !important;
  --input: ${input} !important;
  --secondary: ${mutedBg} !important;
  --secondary-foreground: ${text} !important;
  --accent: ${accentBg} !important;
  --accent-foreground: ${text} !important;
  --primary: ${primary} !important;
  --primary-foreground: ${primaryForeground} !important;
  --ring: ${primary} !important;
  --success: ${success} !important;
  --success-foreground: ${bg} !important;
  --error: ${error} !important;
  --error-foreground: ${text} !important;
  --destructive: ${error} !important;
  --destructive-foreground: ${text} !important;
  --link: ${mutedText} !important;
  --link-visited: ${mutedText} !important;
  --action: ${mutedBg} !important;
  --action-hover: ${border} !important;
  --action-foreground: ${text} !important;
  --stale: rgb(255 255 255 / 10%) !important;
  --base-shadow: rgb(0 0 0 / 35%) !important;
  --base-shadow-darker: rgb(0 0 0 / 55%) !important;
  --cm-background: ${editorBg} !important;
  --cm-comment: ${mutedText} !important;
`;

export const syncNotebookIframeDarkDocument = (doc: Document) => {
  if (doc.documentElement.getAttribute('data-tier0-dark-theme') !== '1') {
    doc.documentElement.setAttribute('data-tier0-dark-theme', '1');
  }
  doc.documentElement.classList.add('dark');
  if (doc.documentElement.style.colorScheme !== 'dark') {
    doc.documentElement.style.colorScheme = 'dark';
  }

  const body = doc.body;
  if (!body) {
    return;
  }

  body.classList.remove('light', 'light-theme');
  body.classList.add('dark', 'dark-theme');
  // Marimo edit mode resolves CodeMirror's theme from its internal config.
  // This supported override makes its own dark theme react immediately instead
  // of leaving light-theme syntax colors underneath our surface overrides.
  if (body.dataset.vscodeThemeKind !== 'vscode-dark') {
    body.dataset.vscodeThemeKind = 'vscode-dark';
  }
  if (body.dataset.theme !== 'dark') {
    body.dataset.theme = 'dark';
  }
  if (body.dataset.mode !== 'dark') {
    body.dataset.mode = 'dark';
  }
  if (body.style.colorScheme !== 'dark') {
    body.style.colorScheme = 'dark';
  }

  doc.getElementById('App')?.classList.add('dark');
  doc.querySelector('.marimo')?.classList.add('dark');
};

export const clearNotebookIframeDarkDocument = (doc: Document) => {
  doc.documentElement.removeAttribute('data-tier0-dark-theme');
  doc.documentElement.classList.remove('dark');
  doc.documentElement.style.removeProperty('color-scheme');

  const body = doc.body;
  if (!body) {
    return;
  }

  body.classList.remove('dark', 'dark-theme');
  body.dataset.vscodeThemeKind = 'vscode-light';
  body.dataset.theme = 'light';
  body.dataset.mode = 'light';
  body.style.removeProperty('color-scheme');

  doc.getElementById('App')?.classList.remove('dark');
  doc.querySelector('.marimo')?.classList.remove('dark');
};

export const buildNotebookIframeDarkThemeCss = () => {
  const bg = readThemeToken('--ui-bg-color', '#161616');
  const text = readThemeToken('--ui-text-color', '#f4f4f4');
  const panelBg = readThemeToken('--ui-sidebar-bg', '#262626');
  const mutedBg = readThemeToken('--ui-charttop-bg-color', '#393939');
  const mutedText = readThemeToken('--ui-description-card-color', '#a8a8a8');
  const border = readThemeToken('--ui-control-border', '#525252');
  const input = readThemeToken('--ui-input-color', '#393939');
  const accentBg = readThemeToken('--ui-fill-tertiary', '#393939');
  const editorBg = readThemeToken('--ui-modal-color', panelBg);
  const primary = readThemeToken('--ui-theme-color', mutedText);
  const primaryForeground = readThemeToken('--ui-primary-button-text', text);
  const success = readThemeToken('--ui-status-active-text', mutedText);
  const error = readThemeToken('--ui-status-inactive-text', text);
  const themeVars = buildThemeVariables(
    bg,
    text,
    panelBg,
    mutedBg,
    mutedText,
    border,
    input,
    accentBg,
    editorBg,
    primary,
    primaryForeground,
    success,
    error
  );

  return `
    ${DARK_SCOPE},
    ${DARK_SCOPE} :root,
    ${DARK_SCOPE} .marimo {
      ${themeVars}
    }

    ${DARK_SCOPE},
    ${DARK_SCOPE} body,
    ${DARK_SCOPE} body.light,
    ${DARK_SCOPE} body.light-theme,
    ${DARK_SCOPE} body.dark,
    ${DARK_SCOPE} body.dark-theme {
      ${themeVars}
      background: ${bg} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} #root,
    ${DARK_SCOPE} #App,
    ${DARK_SCOPE} #App.dark,
    ${DARK_SCOPE} main,
    ${DARK_SCOPE} .marimo,
    ${DARK_SCOPE} .marimo.dark,
    ${DARK_SCOPE} [data-app],
    ${DARK_SCOPE} [data-testid="chrome-sidebar"],
    ${DARK_SCOPE} nav {
      background: ${bg} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} #app-chrome-sidebar,
    ${DARK_SCOPE} [data-testid="helper"],
    ${DARK_SCOPE} [data-testid="panel"],
    ${DARK_SCOPE} [data-panel],
    ${DARK_SCOPE} [data-panel-body],
    ${DARK_SCOPE} [data-testid*="file-explorer"],
    ${DARK_SCOPE} [data-testid*="files"] {
      background: ${panelBg} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} .bg-background,
    ${DARK_SCOPE} [class*="bg-background"] {
      background-color: ${bg} !important;
    }

    ${DARK_SCOPE} .bg-card,
    ${DARK_SCOPE} .bg-popover,
    ${DARK_SCOPE} [class*="bg-card"],
    ${DARK_SCOPE} [class*="bg-popover"] {
      background-color: ${panelBg} !important;
    }

    ${DARK_SCOPE} .bg-muted,
    ${DARK_SCOPE} .bg-secondary,
    ${DARK_SCOPE} [class*="bg-muted"],
    ${DARK_SCOPE} [class*="bg-secondary"] {
      background-color: ${mutedBg} !important;
    }

    ${DARK_SCOPE} .text-foreground,
    ${DARK_SCOPE} [class*="text-foreground"] {
      color: ${text} !important;
    }

    ${DARK_SCOPE} .text-muted-foreground,
    ${DARK_SCOPE} [class*="text-muted-foreground"] {
      color: ${mutedText} !important;
    }

    ${DARK_SCOPE} .marimo-cell,
    ${DARK_SCOPE} [data-testid*="cell"] {
      background: ${panelBg} !important;
      border-color: ${border} !important;
      color: ${text} !important;
      outline-color: ${border} !important;
      box-shadow: none !important;
    }

    ${DARK_SCOPE} .marimo-cell:hover,
    ${DARK_SCOPE} .marimo-cell:focus-within,
    ${DARK_SCOPE} .marimo-cell:has(.cm-editor.cm-focused),
    ${DARK_SCOPE} [data-testid*="cell"]:hover,
    ${DARK_SCOPE} [data-testid*="cell"]:focus-within {
      border-color: ${border} !important;
      outline: 1px solid ${mutedText} !important;
      outline-offset: -1px;
      box-shadow: none !important;
    }

    ${DARK_SCOPE} .text-primary,
    ${DARK_SCOPE} [class*="text-primary"],
    ${DARK_SCOPE} [class*="text-sage"],
    ${DARK_SCOPE} [class*="text-grass"],
    ${DARK_SCOPE} [class*="text-green"],
    ${DARK_SCOPE} [class*="text-lime"] {
      color: ${mutedText} !important;
    }

    ${DARK_SCOPE} .bg-primary,
    ${DARK_SCOPE} [class*="bg-primary"],
    ${DARK_SCOPE} [class*="bg-sage"],
    ${DARK_SCOPE} [class*="bg-grass"],
    ${DARK_SCOPE} [class*="bg-green"],
    ${DARK_SCOPE} [class*="bg-lime"] {
      background-color: ${mutedBg} !important;
    }

    ${DARK_SCOPE} .border-primary,
    ${DARK_SCOPE} [class*="border-primary"],
    ${DARK_SCOPE} [class*="border-sage"],
    ${DARK_SCOPE} [class*="border-grass"],
    ${DARK_SCOPE} [class*="border-green"],
    ${DARK_SCOPE} [class*="border-lime"],
    ${DARK_SCOPE} [class*="border-emerald"],
    ${DARK_SCOPE} [class*="border-teal"],
    ${DARK_SCOPE} [class*="border-cyan"] {
      border-color: ${border} !important;
    }

    ${DARK_SCOPE} hr,
    ${DARK_SCOPE} [role="separator"],
    ${DARK_SCOPE} [data-orientation="horizontal"],
    ${DARK_SCOPE} [data-orientation="vertical"] {
      background-color: ${border} !important;
      border-color: ${border} !important;
      color: ${border} !important;
    }

    ${DARK_SCOPE} nav button,
    ${DARK_SCOPE} nav [role="button"],
    ${DARK_SCOPE} aside button,
    ${DARK_SCOPE} aside [role="button"],
    ${DARK_SCOPE} [data-testid*="sidebar"] button,
    ${DARK_SCOPE} [data-testid*="sidebar"] [role="button"] {
      color: ${mutedText} !important;
      border-color: ${border} !important;
    }

    ${DARK_SCOPE} nav button[data-state="active"],
    ${DARK_SCOPE} nav [role="button"][data-state="active"],
    ${DARK_SCOPE} aside button[data-state="active"],
    ${DARK_SCOPE} aside [role="button"][data-state="active"],
    ${DARK_SCOPE} nav button[aria-current="page"],
    ${DARK_SCOPE} aside button[aria-current="page"] {
      background: ${mutedBg} !important;
      color: ${text} !important;
      border-color: ${border} !important;
      box-shadow: inset 2px 0 0 ${mutedText} !important;
    }

    ${DARK_SCOPE} nav svg,
    ${DARK_SCOPE} aside svg,
    ${DARK_SCOPE} [data-testid*="sidebar"] svg,
    ${DARK_SCOPE} [data-testid*="files"] svg {
      color: ${mutedText} !important;
      stroke: currentColor !important;
    }

    ${DARK_SCOPE} footer,
    ${DARK_SCOPE} [data-testid*="footer"],
    ${DARK_SCOPE} [data-testid*="status"] {
      background: ${bg} !important;
      border-color: ${border} !important;
      color: ${mutedText} !important;
    }

    ${DARK_SCOPE} .marimo-cell .cm-editor,
    ${DARK_SCOPE} .marimo-cell .cm,
    ${DARK_SCOPE} .marimo-cell .cm-gutters,
    ${DARK_SCOPE} .marimo-cell .cm-scroller,
    ${DARK_SCOPE} .marimo-cell .cm-content,
    ${DARK_SCOPE} .cm-editor,
    ${DARK_SCOPE} .cm,
    ${DARK_SCOPE} .cm-gutters,
    ${DARK_SCOPE} #cell-setup,
    ${DARK_SCOPE} #cell-setup .cm-editor,
    ${DARK_SCOPE} #cell-setup .cm-gutter {
      background-color: ${editorBg} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} .cm-editor .cm-content {
      caret-color: ${text} !important;
    }

    ${DARK_SCOPE} .cm-editor .cm-cursor,
    ${DARK_SCOPE} .cm-editor .cm-dropCursor {
      border-left-color: ${text} !important;
    }

    ${DARK_SCOPE} .cm-editor .cm-selectionBackground,
    ${DARK_SCOPE} .cm-editor.cm-focused .cm-selectionBackground,
    ${DARK_SCOPE} .cm-editor ::selection {
      background-color: ${border} !important;
    }

    ${DARK_SCOPE} .marimo-cell .cm-editor.cm-focused .cm-activeLineGutter,
    ${DARK_SCOPE} .marimo-cell .cm-activeLineGutter {
      background: ${mutedBg} !important;
    }

    ${DARK_SCOPE} .marimo-cell .cm-editor.cm-focused .cm-activeLine:not(.cm-error-line),
    ${DARK_SCOPE} .marimo-cell .cm-activeLine:not(.cm-error-line) {
      background: rgb(255 255 255 / 6%) !important;
    }

    ${DARK_SCOPE} .marimo-cell .cm-gutters.cm-gutters-before,
    ${DARK_SCOPE} .cm-gutters.cm-gutters-before {
      border-right-color: rgb(255 255 255 / 30%) !important;
    }

    ${DARK_SCOPE} .output-area,
    ${DARK_SCOPE} .console-output-area {
      background: ${panelBg} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} [data-sonner-toast],
    ${DARK_SCOPE} .marimo-notification,
    ${DARK_SCOPE} [role="status"],
    ${DARK_SCOPE} [data-radix-toast-viewport] > * {
      background: ${panelBg} !important;
      color: ${text} !important;
      border-color: ${border} !important;
    }

    ${DARK_SCOPE} .border,
    ${DARK_SCOPE} .border-border,
    ${DARK_SCOPE} [class*="border-border"] {
      border-color: ${border} !important;
    }

    ${DARK_SCOPE} button:not(:disabled),
    ${DARK_SCOPE} [role="button"]:not([aria-disabled="true"]) {
      color: ${text};
    }

    ${DARK_SCOPE} button[class*="bg-background"],
    ${DARK_SCOPE} button[class*="bg-card"],
    ${DARK_SCOPE} button[class*="bg-secondary"],
    ${DARK_SCOPE} [role="button"][class*="bg-background"],
    ${DARK_SCOPE} [role="button"][class*="bg-card"],
    ${DARK_SCOPE} [role="button"][class*="bg-secondary"] {
      background-color: ${panelBg} !important;
      border-color: ${border} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} button:hover:not(:disabled),
    ${DARK_SCOPE} [role="button"]:hover:not([aria-disabled="true"]) {
      background-color: ${mutedBg} !important;
    }

    ${DARK_SCOPE} input,
    ${DARK_SCOPE} textarea,
    ${DARK_SCOPE} select {
      background: ${input} !important;
      border-color: ${border} !important;
      color: ${text} !important;
    }

    ${DARK_SCOPE} .resize-handle:hover,
    ${DARK_SCOPE} .resize-handle[data-resize-handle-active],
    ${DARK_SCOPE} .resize-handle-collapsed:hover,
    ${DARK_SCOPE} .resize-handle-collapsed[data-resize-handle-active] {
      background-color: ${mutedText} !important;
    }

    ${DARK_SCOPE} .disconnected-gradient,
    ${DARK_SCOPE} .noise {
      opacity: 0.18 !important;
    }
  `;
};
