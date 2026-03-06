import { notification } from 'antd';
import { downloadFn } from '@/utils';
import { useTranslate } from '@/hooks/index.ts';

// 快速使用的 Hook 版本
export const useDownloadNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const formatMessage = useTranslate();

  const showDownloadNotification = (props: { data: any; name: string }) => {
    downloadFn(props);
    // 显示通知
    api.success({
      message: formatMessage('uns.exportFileReady'),
      placement: 'bottomRight',
      closable: true,
      description: (
        <div style={{ fontSize: '12px', color: 'rgba(0, 0, 0, 0.45)' }}>
          {formatMessage('uns.downloadAutoMsg')}{' '}
          <a
            href="#"
            onClick={(e) => {
              e.preventDefault();
              downloadFn(props);
            }}
            style={{ color: '#1677ff', textDecoration: 'underline' }}
          >
            {formatMessage('uns.clickHere')}
          </a>{' '}
          {formatMessage('uns.toDownload')}
        </div>
      ),
      style: {
        backgroundColor: '#DEFBE6',
        border: '1px solid #6FDC8C',
        padding: '16px',
      },
    });
  };

  return [showDownloadNotification, contextHolder] as const;
};
