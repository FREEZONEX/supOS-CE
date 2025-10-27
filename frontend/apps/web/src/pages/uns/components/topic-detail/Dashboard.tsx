import { type FC, useEffect, useRef, useState, useCallback } from 'react';
import { ResizableBox } from 'react-resizable';
import '@/components/resizable-container/index.scss';
import { Result } from 'antd';
import { Empty } from 'antd';
import { useTranslate } from '@/hooks';
import IframeMask from '@/components/iframe-mask';
import { useBaseStore } from '@/stores/base';
import { useSize } from 'ahooks';
import { getDashboardDetail } from '@/apis/inter-api';

interface DetailDashboardProps {
  instanceInfo: { [key: string]: any };
  dashboardInfo: any;
}

const DetailDashboard: FC<DetailDashboardProps> = ({ instanceInfo, dashboardInfo }) => {
  const formatMessage = useTranslate();
  const hasDashboards = useBaseStore((state) => state.menuGroup?.some((f) => f.url === '/dashboards'));
  const observer = useRef<MutationObserver | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const [iframeUrl, setIframeUrl] = useState('');
  useEffect(() => {
    if (instanceInfo && dashboardInfo) {
      if (!dashboardInfo?.type || dashboardInfo?.type === 1) {
        if (dashboardInfo?.id) {
          // grafana
          getDashboardDetail(dashboardInfo?.id).then((res: any) => {
            if (res?.meta?.url) {
              setIframeUrl(`${res?.meta?.url}?kiosk`);
            }
          });
        }
      } else if (dashboardInfo?.type === 2) {
        // fuxa
        setIframeUrl(`/fuxa/home/?id=${dashboardInfo?.id}=lab`);
      }
    }
  }, [instanceInfo, dashboardInfo]);

  const containerSize = useSize(containerRef);

  const iframeCallbackRef = useCallback(
    (iframe: HTMLIFrameElement | null) => {
      // ===== 清理阶段（iframe 卸载时）=====
      if (observer.current) {
        observer.current.disconnect();
        observer.current = null;
      }

      // ===== 挂载阶段 =====
      if (!iframe) return;
      if (dashboardInfo?.type === 2) return;
      const handleMutation = (mutationsList: MutationRecord[]) => {
        const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document;
        if (!iframeDoc) return;

        for (const mutation of mutationsList) {
          if (mutation.type === 'childList') {
            // 遍历新增的节点
            // 使用 querySelectorAll 获取所有匹配的元素，并遍历它们
            iframeDoc.querySelectorAll<HTMLElement>('.show-on-hover').forEach(handleButton);
            // mutation.addedNodes.forEach((node) => {
            //   // 如果新增的是 Element 节点
            //   if (node.nodeType === Node.ELEMENT_NODE) {
            //     const element = node as HTMLElement;

            //     // 检查自身是否是目标按钮
            //     if (element.classList.contains('show-on-hover')) {
            //       handleButton(element);
            //     }

            //     // 检查子树中是否有目标按钮（因为 subtree: true）
            //     const buttons = element.querySelectorAll<HTMLElement>('.show-on-hover');
            //     buttons.forEach(handleButton);
            //   }
            // });
          }
        }
      };

      // 抽离处理逻辑，避免重复
      const handleButton = (btn: HTMLElement) => {
        if (!(btn as any).__handled__) {
          btn.style.display = 'none';
          btn.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
          });
          (btn as any).__handled__ = true;
        }
      };

      // 注入滚动条样式
      const injectStyles = (doc: Document) => {
        const style = doc.createElement('style');
        style.textContent = `
      body::-webkit-scrollbar { width: 8px; height: 8px; background: transparent; }
      body::-webkit-scrollbar-track { margin: 4px 0; border-radius: 8px; }
      body::-webkit-scrollbar-thumb { border-radius: 8px; background: #d3d3d3; cursor: pointer; }
      body::-webkit-scrollbar-thumb:hover { background: #a5a5a5; }
    `;
        doc.head.appendChild(style);
      };

      // onload 处理
      const onLoad = () => {
        const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document;
        if (!iframeDoc) return;

        injectStyles(iframeDoc);

        // 创建并启动 observer
        observer.current = new MutationObserver(handleMutation);
        observer.current.observe(iframeDoc.body, {
          childList: true,
          subtree: true,
        });

        // 立即处理已有元素（防止 missed）
        handleMutation([{ type: 'childList', addedNodes: iframeDoc.body.childNodes } as any]);
      };

      // 绑定 onload
      iframe.onload = onLoad;

      // 注意：如果 iframe 已经加载完成（比如从缓存），可能需要手动触发
      if (iframe.contentDocument?.readyState === 'complete') {
        setTimeout(onLoad, 0); // 确保在下一 tick 执行
      }
    },
    [iframeUrl, dashboardInfo]
  ); // 依赖 iframeUrl，确保 URL 变化时重新绑定

  const [isResizing, setIsResizing] = useState(false);

  if (!instanceInfo?.withDashboard) {
    return <Empty />;
  }
  return hasDashboards ? (
    <>
      <div ref={containerRef}>
        {containerSize?.width && (
          <ResizableBox
            className="resizable-container resizable-hover-handles"
            width={containerSize?.width}
            height={300}
            minConstraints={[200, 200]}
            maxConstraints={[containerSize?.width, 500]}
            axis="both"
            resizeHandles={['se']} // 只允许右下角拖拽
            onResizeStart={() => setIsResizing(true)}
            onResizeStop={() => setIsResizing(false)}
          >
            <>
              <iframe ref={iframeCallbackRef} height="100%" width="100%" id="dashboardIframe" src={iframeUrl} />
              <IframeMask style={{ display: isResizing ? 'block' : 'none' }} />
            </>
          </ResizableBox>
        )}
      </div>
    </>
  ) : (
    <Result
      status="403"
      title={403}
      subTitle={<span style={{ color: 'var(--supos-text-color)' }}>{formatMessage('common.pageNoPermission')}</span>}
    />
  );
};
export default DetailDashboard;
