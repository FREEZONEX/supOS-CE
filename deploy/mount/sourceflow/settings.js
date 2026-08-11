const path = require('path');

const internalToken = process.env.NODERED_INTERNAL_TOKEN || '';
const internalTokenHeader = 'x-tier0-internal-token';
const themeAssetVersion = '20260710-001';

function requireInternalToken(req, res, next) {
  if (!internalToken || req.get(internalTokenHeader) === internalToken) {
    next();
    return;
  }
  res.status(401).send('unauthorized');
}

module.exports = {
  flowFile: 'flows.json',
  flowFilePretty: true,
  uiPort: process.env.PORT || 1880,
  apiMaxLength: '50mb',
  httpServerOptions: { limit: '50mb' },
  httpAdminMiddleware: requireInternalToken,
  httpNodeMiddleware: requireInternalToken,
  lang: process.env.OS_LANG || 'zh-CN',
  diagnostics: { enabled: true, ui: true },
  telemetryEnabled: false,
  telemetry: { enabled: false, updateNotification: false },
  runtimeState: { enabled: false, ui: false },
  logging: { console: { level: 'info', metrics: false, audit: false } },
  contextStorage: {
    default: {
      module: 'localfilesystem',
      config: { dir: '/data/cache/', flushInterval: 10 },
    },
  },
  externalModules: {},
  editorTheme: {
    page: {
      title: 'Flow Editor',
      css: path.join(__dirname, 'themes', `index.${themeAssetVersion}.css`),
      scripts: [path.join(__dirname, 'themes', `index.${themeAssetVersion}.js`)],
    },
    header: { title: 'Flow Editor' },
    tours: false,
    projects: { enabled: false },
    multiplayer: { enabled: false },
  },
  functionExternalModules: true,
  functionTimeout: 0,
  functionGlobalContext: {},
  debugMaxLength: 1000,
  mqttReconnectTime: 15000,
  serialReconnectTime: 15000,
};
