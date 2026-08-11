// 3D 资源挂起信号：宿主页签隐藏（bridge lifecycle active=false）或浏览器窗口进入后台
// （document.visibilitychange）时为 true。恢复条件是两者同时满足：页签激活且窗口可见。
// 无宿主直开（扫码分享页等）时 lifecycle 恒为 active，仅 visibilitychange 生效。
import { useEffect, useState } from 'react';
import { getHostLifecycle, onHostLifecycle } from './index';

function computeSuspended() {
  return !getHostLifecycle().active || document.visibilityState === 'hidden';
}

export function useSuspended(): boolean {
  const [suspended, setSuspended] = useState(computeSuspended);

  useEffect(() => {
    const update = () => setSuspended(computeSuspended());
    const unsubscribe = onHostLifecycle(update);
    document.addEventListener('visibilitychange', update);
    return () => {
      unsubscribe();
      document.removeEventListener('visibilitychange', update);
    };
  }, []);

  return suspended;
}
