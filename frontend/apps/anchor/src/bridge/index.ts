// 宿主桥接单例（@internal/anchor-bridge 客户端侧）：
//   - 在 iframe 内运行时与主站交换上下文/生命周期消息；
//   - 扫码分享页等无宿主直开场景 hasHost=false，所有发送为 no-op、监听不注册；
//   - 消息到达可能早于各订阅方挂载，最近一次 context/lifecycle 会被缓存供后来者读取。
import { createAnchorClient, type AnchorContextPayload, type AnchorLifecyclePayload } from '@internal/anchor-bridge';

type Listener<T> = (payload: T) => void;

const contextListeners = new Set<Listener<AnchorContextPayload>>();
const lifecycleListeners = new Set<Listener<AnchorLifecyclePayload>>();
let lastContext: AnchorContextPayload | null = null;
let lastLifecycle: AnchorLifecyclePayload = { active: true };

export const anchorClient = createAnchorClient({
  onContext: (payload) => {
    lastContext = payload;
    contextListeners.forEach((listener) => listener(payload));
  },
  onLifecycle: (payload) => {
    lastLifecycle = payload;
    lifecycleListeners.forEach((listener) => listener(payload));
  },
});

export const getHostContext = () => lastContext;
export const getHostLifecycle = () => lastLifecycle;

export function onHostContext(listener: Listener<AnchorContextPayload>) {
  contextListeners.add(listener);
  return () => contextListeners.delete(listener);
}

export function onHostLifecycle(listener: Listener<AnchorLifecyclePayload>) {
  lifecycleListeners.add(listener);
  return () => lifecycleListeners.delete(listener);
}
