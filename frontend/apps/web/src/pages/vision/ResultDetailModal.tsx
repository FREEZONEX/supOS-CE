import { type ReactNode, useEffect, useState } from 'react';
import { App, Button, Collapse, Image, Spin, Tag, Typography } from 'antd';
import ProModal from '@/components/pro-modal';
import { Download } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useTranslate } from '@/hooks';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { formatTimestamp } from '@/utils/format';
import {
  deleteVisionEvent,
  eventScreenshotUrl,
  getVisionEvent,
  type VisionEvent,
  type VisionEventDetail,
} from '@/apis/core-api/vision-results';
import { ImageOff } from 'lucide-react';
import styles from './results.module.scss';

const { Text } = Typography;

type ResultDetailModalProps = {
  event: VisionEvent | null;
  onClose: () => void;
  onViewTask: (event: VisionEvent) => void;
  onDeleted: () => void;
};

// 触发浏览器直接下载当前截图(同源请求自动带 Cookie)。
const downloadScreenshot = (id: number) => {
  const link = document.createElement('a');
  link.href = eventScreenshotUrl(id);
  link.download = `vision-event-${id}.jpg`;
  document.body.appendChild(link);
  link.click();
  link.remove();
};

const ResultDetailModal = ({ event, onClose, onViewTask, onDeleted }: ResultDetailModalProps) => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [detail, setDetail] = useState<VisionEventDetail | null>(null);

  useEffect(() => {
    if (!event) return undefined;
    let cancelled = false;
    getVisionEvent(event.id)
      .then((data) => {
        if (!cancelled) setDetail(data);
      })
      .catch(() => {
        // 详情接口失败时退回列表数据,同时结束 loading。
        if (!cancelled) setDetail({ ...event });
      });
    return () => {
      cancelled = true;
    };
  }, [event]);

  if (!event) return null;
  // 详情接口未返回前先用列表数据渲染,避免弹窗闪空;通过 id 匹配判断是否已加载。
  const loaded = detail?.id === event.id;
  const data: VisionEventDetail = loaded && detail ? detail : event;
  const loading = !loaded;

  const remove = () => {
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          title: formatMessage('Vision.results.deleteTitle'),
          width: 420,
          content: (
            <span className="tier0-delete-confirm-message">
              {formatMessage('Vision.results.deleteConfirm', { name: data.eventName })}
            </span>
          ),
          onOk: async () => {
            await deleteVisionEvent(data.id);
            message.success(formatMessage('common.deleteSuccessfully'));
            onDeleted();
          },
        },
        formatMessage
      )
    );
  };

  const fields: { label: string; value: ReactNode }[] = [
    { label: formatMessage('Vision.results.eventName'), value: data.eventName || '-' },
    { label: formatMessage('Vision.results.time'), value: formatTimestamp(data.createdAt) || '-' },
    {
      label: formatMessage('Vision.results.camera'),
      value: data.cameraName || (data.cameraId ? `#${data.cameraId}` : '-'),
    },
    {
      label: formatMessage('Vision.results.analysisTask'),
      value: data.taskName || (data.taskId ? `#${data.taskId}` : '-'),
    },
    { label: formatMessage('Vision.results.algorithm'), value: data.algorithmName || data.algorithmCode || '-' },
    { label: formatMessage('Vision.results.algorithmVer'), value: data.algorithmVersion || '-' },
    { label: formatMessage('Vision.results.modelVer'), value: data.modelVersion || '-' },
    { label: formatMessage('Vision.results.retention'), value: data.retention || '-' },
  ];

  const techItems: { label: string; value: ReactNode }[] = [
    { label: formatMessage('Vision.results.eventId'), value: data.eventId || '-' },
    { label: formatMessage('Vision.results.evidenceId'), value: data.evidenceId || '-' },
    { label: formatMessage('Vision.results.unsTopic'), value: data.unsTopic || '-' },
  ];

  return (
    <ProModal
      open={Boolean(event)}
      title={formatMessage('Vision.results.detailTitle')}
      width={760}
      fullScreenable={false}
      // 矮视口下内容超高时 body 内部滚动,避免 centered 弹窗头部被顶出视口无法关闭
      styles={{ body: { maxHeight: 'calc(100vh - 140px)', overflowY: 'auto' } }}
      onCancel={onClose}
      footer={null}
      destroyOnHidden
    >
      <Spin spinning={loading}>
        <div className={styles.detailShot}>
          {data.hasScreenshot ? (
            <>
              <Image
                className={styles.detailShotImg}
                src={eventScreenshotUrl(data.id)}
                alt={data.eventName}
                preview={{ mask: formatMessage('Vision.results.expand') }}
              />
              <Button
                className={styles.detailShotDownload}
                size="small"
                icon={<Download {...toolbarIconProps} size={14} />}
                onClick={() => downloadScreenshot(data.id)}
              >
                {formatMessage('Vision.results.downloadShot')}
              </Button>
            </>
          ) : (
            <div className={styles.detailShotEmpty}>
              <ImageOff size={28} strokeWidth={1.5} />
              <span>{formatMessage('Vision.results.noScreenshot')}</span>
            </div>
          )}
        </div>
        <div className={styles.detailGrid}>
          {fields.map((item) => (
            <div key={item.label} className={styles.detailItem}>
              <Text type="secondary" className={styles.detailLabel}>
                {item.label}
              </Text>
              <div className={styles.detailValue}>{item.value}</div>
            </div>
          ))}
        </div>
        <div className={styles.detailBlock}>
          <Text type="secondary" className={styles.detailLabel}>
            {formatMessage('Vision.results.resultContent')}
          </Text>
          <div className={styles.contentTags}>
            {(data.resultContent || []).length === 0 ? (
              <Text type="secondary">--</Text>
            ) : (
              (data.resultContent || []).map((item) => (
                <Tag key={`${item.label}-${item.value}`} color="success">
                  {[item.label, item.value].filter(Boolean).join(' ')}
                </Tag>
              ))
            )}
          </div>
        </div>
        <Collapse
          ghost
          className={styles.techCollapse}
          items={[
            {
              key: 'tech',
              label: formatMessage('Vision.results.techInfo'),
              children: (
                <div className={styles.techGrid}>
                  {techItems.map((item) => (
                    <div key={item.label} className={styles.detailItem}>
                      <Text type="secondary" className={styles.detailLabel}>
                        {item.label}
                      </Text>
                      <div className={styles.detailValue}>{item.value}</div>
                    </div>
                  ))}
                </div>
              ),
            },
          ]}
        />
        <div className={styles.detailFooter}>
          <Button onClick={() => onViewTask(data)}>{formatMessage('Vision.results.viewTask')}</Button>
          <div className={styles.detailFooterRight}>
            <Button danger onClick={remove}>
              {formatMessage('common.delete')}
            </Button>
            <Button type="primary" onClick={onClose}>
              {formatMessage('Vision.results.done')}
            </Button>
          </div>
        </div>
      </Spin>
    </ProModal>
  );
};

export default ResultDetailModal;
