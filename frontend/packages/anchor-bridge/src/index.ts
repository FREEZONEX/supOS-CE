// 主站（apps/web，宿主）与 Anchor 3D 子应用（apps/anchor，iframe 内）的 postMessage 通信契约。
// 设计约束：
//   - 无 React/框架依赖，两侧共用同一份实现，规避 apps -> apps 依赖；
//   - 严格同源：双方都校验 event.origin === window.location.origin，宿主额外校验 event.source 是当前 iframe；
//   - 消息带版本号，收到更高 major 版本时按不识别处理（保持向后兼容的余地）；
//   - context 只携带主题/主色/语言等非敏感展示状态，禁止传 Cookie/Token/API Key。

export const ANCHOR_BRIDGE_VERSION = 1;

const MESSAGE_MARK = 'tier0-anchor-bridge';

export const AnchorMessageType = {
  /** 子 → 宿主：iframe 内应用挂载完成，请求初始上下文 */
  Ready: 'tier0.anchor.ready',
  /** 宿主 → 子：平台上下文（主题/主色/语言），初始及变更时推送 */
  Context: 'tier0.anchor.context',
  /** 宿主 → 子：页签生命周期（active=false 时子应用应暂停渲染与订阅） */
  Lifecycle: 'tier0.anchor.lifecycle',
  /** 子 → 宿主：内部路由变化（仅同步标题与路径展示，不改写宿主路由） */
  Route: 'tier0.anchor.route',
  /** 子 → 宿主：子应用错误上报（展示用） */
  Error: 'tier0.anchor.error',
} as const;

export type AnchorMessageTypeValue = (typeof AnchorMessageType)[keyof typeof AnchorMessageType];

export interface AnchorContextPayload {
  theme: 'light' | 'dark';
  primaryColor: string;
  language: string;
}

export interface AnchorLifecyclePayload {
  active: boolean;
}

export interface AnchorRoutePayload {
  path: string;
  title?: string;
}

export interface AnchorErrorPayload {
  message: string;
}

interface AnchorBridgeMessage<T = unknown> {
  mark: typeof MESSAGE_MARK;
  type: AnchorMessageTypeValue;
  version: number;
  payload?: T;
}

function isBridgeMessage(data: unknown): data is AnchorBridgeMessage {
  if (!data || typeof data !== 'object') return false;
  const record = data as Record<string, unknown>;
  return record.mark === MESSAGE_MARK && typeof record.type === 'string' && typeof record.version === 'number';
}

function makeMessage<T>(type: AnchorMessageTypeValue, payload?: T): AnchorBridgeMessage<T> {
  return { mark: MESSAGE_MARK, type, version: ANCHOR_BRIDGE_VERSION, payload };
}

// ---------- 宿主侧（apps/web） ----------

export interface AnchorHostHandlers {
  onReady?: () => void;
  onRoute?: (payload: AnchorRoutePayload) => void;
  onError?: (payload: AnchorErrorPayload) => void;
}

export interface AnchorHost {
  sendContext: (payload: AnchorContextPayload) => void;
  sendLifecycle: (payload: AnchorLifecyclePayload) => void;
  dispose: () => void;
}

export function createAnchorHost(iframe: HTMLIFrameElement, handlers: AnchorHostHandlers = {}): AnchorHost {
  const origin = window.location.origin;

  const post = (message: AnchorBridgeMessage) => {
    iframe.contentWindow?.postMessage(message, origin);
  };

  const onMessage = (event: MessageEvent) => {
    if (event.origin !== origin) return;
    if (!iframe.contentWindow || event.source !== iframe.contentWindow) return;
    if (!isBridgeMessage(event.data)) return;
    switch (event.data.type) {
      case AnchorMessageType.Ready:
        handlers.onReady?.();
        break;
      case AnchorMessageType.Route:
        handlers.onRoute?.(event.data.payload as AnchorRoutePayload);
        break;
      case AnchorMessageType.Error:
        handlers.onError?.(event.data.payload as AnchorErrorPayload);
        break;
    }
  };

  window.addEventListener('message', onMessage);
  return {
    sendContext: (payload) => post(makeMessage(AnchorMessageType.Context, payload)),
    sendLifecycle: (payload) => post(makeMessage(AnchorMessageType.Lifecycle, payload)),
    dispose: () => window.removeEventListener('message', onMessage),
  };
}

// ---------- 子应用侧（apps/anchor） ----------

export interface AnchorClientHandlers {
  onContext?: (payload: AnchorContextPayload) => void;
  onLifecycle?: (payload: AnchorLifecyclePayload) => void;
}

export interface AnchorClient {
  /** 是否运行在宿主 iframe 内（扫码分享页等无宿主直开场景为 false，所有 send 变为 no-op） */
  readonly hasHost: boolean;
  sendReady: () => void;
  sendRoute: (payload: AnchorRoutePayload) => void;
  sendError: (payload: AnchorErrorPayload) => void;
  dispose: () => void;
}

export function createAnchorClient(handlers: AnchorClientHandlers = {}): AnchorClient {
  const origin = window.location.origin;
  const hasHost = window.parent !== window;

  const post = (message: AnchorBridgeMessage) => {
    if (!hasHost) return;
    window.parent.postMessage(message, origin);
  };

  const onMessage = (event: MessageEvent) => {
    if (event.origin !== origin) return;
    if (event.source !== window.parent) return;
    if (!isBridgeMessage(event.data)) return;
    switch (event.data.type) {
      case AnchorMessageType.Context:
        handlers.onContext?.(event.data.payload as AnchorContextPayload);
        break;
      case AnchorMessageType.Lifecycle:
        handlers.onLifecycle?.(event.data.payload as AnchorLifecyclePayload);
        break;
    }
  };

  if (hasHost) window.addEventListener('message', onMessage);
  return {
    hasHost,
    sendReady: () => post(makeMessage(AnchorMessageType.Ready)),
    sendRoute: (payload) => post(makeMessage(AnchorMessageType.Route, payload)),
    sendError: (payload) => post(makeMessage(AnchorMessageType.Error, payload)),
    dispose: () => {
      if (hasHost) window.removeEventListener('message', onMessage);
    },
  };
}
