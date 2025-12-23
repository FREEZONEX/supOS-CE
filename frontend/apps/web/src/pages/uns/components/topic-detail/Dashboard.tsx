import { type FC, useEffect, useRef, useState, useCallback } from 'react';
import md5 from 'blueimp-md5';
import { ResizableBox } from 'react-resizable';
import '@/components/resizable-container/index.scss';
import { Result, type TimeRangePickerProps } from 'antd';
import { Flex, DatePicker, Button, Empty } from 'antd';
import { Renew } from '@carbon/icons-react';
import dayjs from 'dayjs';
import type { Dayjs } from 'dayjs';
import { useTranslate } from '@/hooks';
import IframeMask from '@/components/iframe-mask';
import { useBaseStore } from '@/stores/base';

const { RangePicker } = DatePicker;

interface DetailDashboardProps {
  instanceInfo: { [key: string]: any };
}

const DetailDashboard: FC<DetailDashboardProps> = ({ instanceInfo }) => {
  const { dataType, refers, alias } = instanceInfo;

  const formatMessage = useTranslate();
  const hasDashboards = useBaseStore((state) => {
    console.log(
      state.menuGroup?.some((f) => f.url === '/dashboards'),
      state.menuGroup
    );
    return state.menuGroup?.some((f) => f.url === '/dashboards');
  });
  const observer = useRef<MutationObserver | null>(null);
  const newAlias = dataType === 7 ? refers?.[0]?.alias : alias;
  const aliasHash = md5(newAlias).slice(8, 24);
  const iframeName = `${newAlias?.replaceAll('_', '-')}`;

  const [iframeUrl, setIframeUrl] = useState(
    `/grafana/home/d-solo/${aliasHash}/${iframeName}?orgId=1&panelId=1&__feature.dashboardSceneSolo`
  );
  const [dates, setDates] = useState<any>(null);

  useEffect(() => {
    handleDefaultTime();
  }, [instanceInfo]);

  useEffect(() => {
    const timeFrame = dates ? `&from=${dayjs(dates[0]).valueOf()}&to=${dayjs(dates[1]).valueOf()}` : '';
    setIframeUrl(
      `/grafana/home/d-solo/${aliasHash}/${iframeName}?orgId=1&panelId=1&__feature.dashboardSceneSolo${timeFrame}`
    );
  }, [dates]);

  const iframeCallbackRef = useCallback(
    (iframe: HTMLIFrameElement | null) => {
      // ===== 清理阶段（iframe 卸载时）=====
      if (observer.current) {
        observer.current.disconnect();
        observer.current = null;
      }

      // ===== 挂载阶段 =====
      if (!iframe) return;
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
    [iframeUrl]
  ); // 依赖 iframeUrl，确保 URL 变化时重新绑定

  const [isResizing, setIsResizing] = useState(false);

  const rangePresets: TimeRangePickerProps['presets'] = [
    {
      label: <span title={formatMessage('uns.last5minutes')}>{formatMessage('uns.last5minutes')}</span>,
      value: [dayjs().add(-5, 'm'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last30minutes')}>{formatMessage('uns.last30minutes')}</span>,
      value: [dayjs().add(-30, 'm'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last1hour')}>{formatMessage('uns.last1hour')}</span>,
      value: [dayjs().add(-1, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last6hours')}>{formatMessage('uns.last6hours')}</span>,
      value: [dayjs().add(-6, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last24hours')}>{formatMessage('uns.last24hours')}</span>,
      value: [dayjs().add(-24, 'h'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last1week')}>{formatMessage('uns.last1week')}</span>,
      value: [dayjs().add(-1, 'w'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last6weeks')}>{formatMessage('uns.last6weeks')}</span>,
      value: [dayjs().add(-6, 'w'), dayjs()],
    },
    {
      label: <span title={formatMessage('uns.last1year')}>{formatMessage('uns.last1year')}</span>,
      value: [dayjs().add(-1, 'y'), dayjs()],
    },
  ];

  const handleDefaultTime = () => {
    setDates([dataType === 2 ? dayjs().add(-6, 'h') : dayjs().add(-5, 'm'), dayjs()]);
  };

  const onRangeChange = (dates: null | (Dayjs | null)[]) => {
    setDates(dates);
  };
  if (!instanceInfo?.withDashboard) {
    return <Empty />;
  }
  return hasDashboards ? (
    <>
      <Flex gap={10} style={{ marginBottom: '10px' }}>
        <RangePicker
          showTime
          format="YYYY-MM-DD HH:mm:ss"
          value={dates}
          onChange={onRangeChange}
          presets={rangePresets}
        />
        <Button
          color="default"
          variant="filled"
          icon={<Renew />}
          onClick={() => {
            if (dates) {
              setDates([dayjs().add(dates[0] - dates[1], 'ms'), dayjs()]);
            } else {
              handleDefaultTime();
            }
          }}
          style={{
            border: '1px solid #CBD5E1',
            color: 'var(--supos-text-color)',
            backgroundColor: 'var(--supos-uns-button-color)',
          }}
        />
      </Flex>

      <ResizableBox
        className="resizable-container resizable-hover-handles"
        width={900}
        height={300}
        minConstraints={[200, 200]}
        maxConstraints={[1280, 500]}
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
