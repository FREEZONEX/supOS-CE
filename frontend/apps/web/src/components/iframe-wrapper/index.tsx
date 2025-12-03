import { type CSSProperties, type FC, useEffect, useRef } from 'react';

const IframeWrapper: FC<{ src: string; iframeRealUrl?: string; title: string; style?: CSSProperties; id?: string }> = ({
  src,
  iframeRealUrl,
  title,
  style,
  id,
}) => {
  const observer = useRef<MutationObserver | null>(null);

  useEffect(() => {
    const _id = 'm_9e685f1a061cdadf6785ef0d2404e813_dashboard' + '_tab_iframe';
    // 资源监控特殊处理
    if (id + '_iframe' === _id) {
      const iframe = document.getElementById(_id) as HTMLIFrameElement | null;
      if (iframe) {
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

        // 当 iframe 加载完成后开始观察
        iframe.onload = function () {
          const iframeDocument: Document | undefined = iframe.contentDocument || iframe.contentWindow?.document;
          if (!iframeDocument) return;
          startObserving();
        };

        return () => {
          if (observer.current) observer.current.disconnect();
        };
      }
    }
  }, [src]);

  return (
    <div
      style={{
        width: '100%',
        height: '100%',
        ...style,
      }}
      id={id}
    >
      <iframe
        id={id + '_iframe'}
        style={{
          width: '100%',
          height: '100%',
          border: 'none',
        }}
        title={title}
        src={iframeRealUrl ? iframeRealUrl : src}
      ></iframe>
    </div>
  );
};

export default IframeWrapper;
