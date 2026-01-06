import { type FC, useEffect, useRef } from 'react';
import type { PageProps } from '@/common-types';

interface DynamicIframeProps extends PageProps {
  url?: string;
  // 自定义的的url
  iframeRealUrl?: string;
  name?: string;
  code?: string;
}
const DynamicIframe: FC<DynamicIframeProps> = ({ url, name, iframeRealUrl, location, code }) => {
  const { state, search } = location || {};
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const _iframeRealUrl = iframeRealUrl ?? state?.iframeRealUrl;
  const src = url ?? state?.url;
  const iframeSrc = (_iframeRealUrl ? _iframeRealUrl : src) + search;
  const observer = useRef<MutationObserver | null>(null);

  useEffect(() => {
    const iframe = iframeRef.current;
    if (!iframe) return;

    interface CustomHTMLElement extends HTMLElement {
      __handled__?: boolean;
    }

    const handleMutation = (mutationsList: MutationRecord[]) => {
      for (const mutation of mutationsList) {
        if (mutation.type === 'childList') {
          const showOnHoverBtns = Array.from(
            iframe?.contentWindow?.document?.querySelectorAll('.show-on-hover') || []
          ) as CustomHTMLElement[];
          showOnHoverBtns.forEach((btn: CustomHTMLElement) => {
            if (!btn.__handled__) {
              btn.style.display = 'none';
              btn.addEventListener('click', function handleClick(event: Event) {
                event.stopPropagation();
                event.preventDefault();
                btn.removeEventListener('click', handleClick);
              });
              btn.__handled__ = true;
            }
          });
        }
      }
    };

    const startObserving = () => {
      // 创建一个 MutationObserver 实例并定义其回调函数
      observer.current = new MutationObserver(handleMutation);

      // 开始观察 iframe 内的内容变化
      observer.current.observe(iframe.contentWindow?.document.body || document.body, {
        childList: true,
        subtree: true,
      });
    };

    // 监听 iframe 加载完成
    const onLoad = () => {
      if (!iframe) return;
      setTimeout(() => {
        // 获取 iframe 的 document 对象
        const iframeDocument = iframe.contentDocument || iframe.contentWindow?.document;
        if (iframeDocument?.URL?.includes('/hasura/home/')) {
          // 创建一个新的 <style> 元素
          const style: any = iframeDocument?.createElement('style') || {};

          // 设置默认字体
          style.textContent = `
            * {
              font-family: 'IBM Plex Sans', sans-serif !important; /* 设置默认字体 */
            }
          `;

          // 将样式插入到 iframe 的 <head> 中
          iframeDocument?.head.appendChild(style);
        }
        if (code === 'm_150c96b72b9bc52484bfca5ac6c6f88a_dashboard') {
          // 特殊处理日志管理的样式
          startObserving();
        }
      }, 0);
    };
    if (iframe) {
      // 绑定 load 事件
      iframe.addEventListener('load', onLoad);
    }

    // 清理事件监听器
    return () => {
      if (iframe) {
        iframe.removeEventListener('load', onLoad);
      }
      if (observer.current) observer.current?.disconnect?.();
    };
  }, []);
  return (
    <iframe
      ref={iframeRef}
      style={{
        width: '100%',
        height: '100%',
        border: 'none',
      }}
      title={name ?? state.name}
      src={iframeSrc}
    />
  );
};

export default DynamicIframe;
