import { type FC, useEffect, useRef } from 'react';
import type { PageProps } from '@/common-types';

interface DynamicIframeProps extends PageProps {
  url?: string;
  // 自定义的的url
  iframeRealUrl?: string;
  name?: string;
  code?: string;
}
const DynamicIframe: FC<DynamicIframeProps> = ({ url, name, iframeRealUrl, location }) => {
  const { state, search } = location || {};
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const _iframeRealUrl = iframeRealUrl ?? state?.iframeRealUrl;
  const src = url ?? state?.url;
  const iframeSrc = (_iframeRealUrl ? _iframeRealUrl : src) + (search ?? '');

  useEffect(() => {
    const iframe = iframeRef.current;
    if (!iframe) return;

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
