import {
  Activity,
  Bell,
  Book,
  Box,
  Boxes,
  CircleHelp,
  Cloud,
  Code2,
  Database,
  FileText,
  GitFork,
  GitMerge,
  Home,
  Landmark,
  LayoutGrid,
  ListTree,
  Play,
  Settings,
  Shield,
  User,
  Workflow,
  Zap,
  type LucideIcon,
} from 'lucide-react';

/** Backend menu icon name strings → Lucide icons (aligned with menu-icon-map.ts). */
const legacyIconMap: Record<string, LucideIcon> = {
  Home: Home,
  Homepage: Home,
  homeNamespace: ListTree,
  Namespace: ListTree,
  UNS: ListTree,
  'menu.tag.uns': ListTree,
  SourceFlow: Workflow,
  homeSourceFlow: Workflow,
  'collection-flow': Workflow,
  EventFlow: Zap,
  homeEventFlow: Zap,
  MqttBroker: GitFork,
  Project: LayoutGrid,
  Launchpad: Play,
  Marimo: Book,
  Notebook: Book,
  SystemCascade: Box,
  Cluster: Box,
  'menu.tag.connections': Workflow,
  'menu.tag.system': Settings,
  'menu.tag.settings': Settings,
  'menu.tag.apps': LayoutGrid,
  'menu.tag.appspace': LayoutGrid,
  'menu.tag.devtools': Code2,
  UserManagement: User,
  Authentication: Shield,
  OpenData: Database,
  MenuConfiguration: Settings,
  Logs: FileText,
  RoutingManagement: GitFork,
  NotificationManagement: Bell,
  PluginManagement: Boxes,
  AppManagement: LayoutGrid,
  Dashboards: LayoutGrid,
  Grafana: LayoutGrid,
  Settings: Settings,
  AboutUs: CircleHelp,
  Localization: Landmark,
  DevTools: Code2,
  Apps: LayoutGrid,
  AppBuilder: LayoutGrid,
  GenApps: LayoutGrid,
  AdvancedUse: CircleHelp,
  Alert: Bell,
  Apm: Activity,
  CICD: Workflow,
  CodeManagement: Code2,
  CollectionGatewayManagement: GitFork,
  homeCollectionGatewayManagement: GitFork,
  Connections: GitFork,
  ContainerManagement: Box,
  DBConnect: Database,
  Dify: Cloud,
  ElasticSearch: Database,
  Emqx: GitFork,
  Gitea: Code2,
  GraphQL: Code2,
  Keycloak: Shield,
  LowCodeTool: Code2,
  Minio: Box,
  NodeRed: Workflow,
  ObjectStorageServer: Box,
  Ollama: Cloud,
  PRIDE: CircleHelp,
  SQLEditor: Code2,
  StreamProcessing: Workflow,
  Swagger: FileText,
  TimescaleDB: Database,
  WebHooks: GitFork,
  GenerativeUI: LayoutGrid,
  McpClient: Code2,
  license: Shield,
  'open-api-docs': FileText,
  dataModeling: Database,
  Simplification: CircleHelp,
  'Appearance Settings': Settings,
  '3d-curve--auto-colon': Boxes,
};

export const resolveLegacyLucideIcon = (iconName?: string): LucideIcon => {
  const raw = String(iconName || '').trim();
  if (!raw) {
    return ListTree;
  }

  const nameWithoutExt = raw.replace(/\.(svg|png|jpg|jpeg|gif|webp|ico)$/i, '');
  if (legacyIconMap[nameWithoutExt]) {
    return legacyIconMap[nameWithoutExt];
  }

  const lower = nameWithoutExt.toLowerCase();
  for (const [key, Icon] of Object.entries(legacyIconMap)) {
    if (key.toLowerCase() === lower) {
      return Icon;
    }
  }

  if (lower.includes('uns') || lower.includes('namespace')) return ListTree;
  if (lower.includes('event') && lower.includes('flow')) return Zap;
  if (lower.includes('source') && lower.includes('flow')) return GitMerge;
  if (lower.includes('collection') && lower.includes('flow')) return GitMerge;
  if (lower.includes('launchpad') || lower === 'launch') return Play;
  if (lower.includes('project')) return LayoutGrid;
  if (lower.includes('home')) return Home;
  if (lower.includes('notebook') || lower.includes('marimo')) return Book;
  if (lower.includes('mqtt')) return GitFork;
  if (lower.includes('cluster')) return Box;
  if (lower.includes('flow')) return Workflow;

  return ListTree;
};
