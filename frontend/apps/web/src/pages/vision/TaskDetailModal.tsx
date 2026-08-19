import type { ReactNode } from 'react';
import { Alert, Button, Tag, Tooltip } from 'antd';
import dayjs from 'dayjs';
import type { VisionTask, VisionTaskStatus } from '@/apis/core-api/task';
import ProModal from '@/components/pro-modal';
import TaskStatusTag from './TaskStatusTag';
import { Help } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import { RETENTION_OPTIONS, samplingToFps, showsPushIntervalFor, spatialModeFor } from './task-meta';
import styles from './components/TaskDetailModal.module.scss';

type TaskDetailModalProps = {
  task: VisionTask | null;
  onClose: () => void;
  onEdit: (task: VisionTask) => void;
};

// 任务详情弹窗(Figma Task Details):双列字段 + 参数/输出标签 + 分相机 UNS topic;
// Error 状态顶部显示错误横幅;Edit 跳编辑 Drawer。
const TaskDetailModal = ({ task, onClose, onEdit }: TaskDetailModalProps) => {
  const formatMessage = useTranslate();
  if (!task) return null;

  const params = task.params;
  const regionMode = spatialModeFor(task.algorithm);
  const regionValues = Object.values(task.regions || {});
  const zoneTotal = regionValues.reduce((sum, r) => sum + (r.zones?.filter((z) => z.enabled !== false).length || 0), 0);
  const hasLine = regionValues.some((r) => r.countingLine);
  const errorText = task.status === 'error' ? task.lastError || task.observedReason : '';
  const retention = RETENTION_OPTIONS.find((o) => o.value === (params?.resultRetention || 'permanent'));

  const labelWithTip = (label: string, tip?: string) => (
    <span className={styles.fieldLabel}>
      {label}
      {tip && (
        <Tooltip title={tip}>
          <Help size={13} className={styles.fieldTip} />
        </Tooltip>
      )}
    </span>
  );

  const field = (label: ReactNode, value: ReactNode, span2 = false) => (
    <div className={span2 ? styles.fieldWide : styles.field}>
      <div className={styles.fieldLabelRow}>{label}</div>
      <div className={styles.fieldValue}>{value}</div>
    </div>
  );

  const areaLabel =
    regionMode === 'line'
      ? labelWithTip(formatMessage('Vision.task.countingLineLabel'), formatMessage('Vision.task.detectionAreaTip'))
      : labelWithTip(formatMessage('Vision.task.detectionArea'), formatMessage('Vision.task.detectionAreaTip'));
  const areaValue =
    regionMode === 'line'
      ? formatMessage(hasLine ? 'Vision.task.customLine' : 'Vision.task.defaultCenterLine')
      : zoneTotal > 0
        ? `${formatMessage('Vision.task.customArea')} (${zoneTotal})`
        : formatMessage('Vision.task.fullFrame');

  const scheduleValue =
    !task.schedule || task.schedule.mode === 'allDay'
      ? formatMessage('Vision.task.runAllDay')
      : formatMessage('Vision.task.customSchedule');

  return (
    <ProModal
      open={Boolean(task)}
      title={formatMessage('Vision.task.detailsTitle')}
      width={640}
      fullScreenable={false}
      styles={{ body: { maxHeight: 'calc(100vh - 140px)', overflowY: 'auto' } }}
      onCancel={onClose}
      destroyOnHidden
      footer={
        <div className={styles.footer}>
          <Button onClick={() => onEdit(task)}>{formatMessage('common.edit')}</Button>
          <Button type="primary" onClick={onClose}>
            {formatMessage('common.done')}
          </Button>
        </div>
      }
    >
      {errorText && <Alert className={styles.errorBanner} type="error" showIcon message={errorText} />}
      <div className={styles.grid}>
        {field(formatMessage('Vision.task.name'), task.name)}
        {field(formatMessage('Vision.task.taskId'), task.id)}
        {field(formatMessage('Vision.task.algorithmLabel'), task.algorithm?.name || '-')}
        {field(formatMessage('Vision.task.taskStatus'), <TaskStatusTag status={task.status as VisionTaskStatus} />)}
        {field(
          formatMessage('Vision.task.camerasLabel'),
          <div className={styles.tagWrap}>
            {(task.cameras || []).map((cam) => (
              <Tag key={cam.id} bordered={false}>
                {cam.name || cam.code}
              </Tag>
            ))}
          </div>,
          true
        )}
        {task.note && field(formatMessage('Vision.task.note'), task.note, true)}
        <div className={styles.divider} />
        {field(
          labelWithTip(formatMessage('Vision.task.processedFrames'), formatMessage('Vision.task.framesTip')),
          (task.processedFrames ?? 0).toLocaleString()
        )}
        {field(
          labelWithTip(formatMessage('Vision.task.detectedTargets'), formatMessage('Vision.task.targetsTip')),
          (task.currentTargets ?? 0).toLocaleString()
        )}
        {regionMode !== 'none' && field(areaLabel, areaValue)}
        {field(formatMessage('Vision.task.runSchedule'), scheduleValue)}
        {field(
          formatMessage('Vision.task.lastProcessed'),
          task.lastHeartbeatAt ? dayjs(task.lastHeartbeatAt).format('YYYY-MM-DD HH:mm:ss') : '-'
        )}
        {showsPushIntervalFor(task.algorithm)
          ? field(
              labelWithTip(formatMessage('Vision.task.pushInterval'), formatMessage('Vision.task.pushIntervalHint')),
              `${((params?.resultPushIntervalMs ?? 0) / 1000).toLocaleString()} s`
            )
          : null}
        {field(
          formatMessage('Vision.task.resultRetention'),
          formatMessage(retention?.labelKey || 'Vision.task.retentionPermanent')
        )}
        {field(
          labelWithTip(formatMessage('Vision.task.recognitionParams'), formatMessage('Vision.task.recogParamsTip')),
          <div className={styles.tagWrap}>
            <Tag>{samplingToFps(params?.samplingIntervalMs)} FPS</Tag>
            <Tag>{params?.samplingIntervalMs ?? 200} ms</Tag>
            <Tag>Conf {params?.confThreshold ?? 0.5}</Tag>
            <Tag>IOU {params?.iouThreshold ?? 0.45}</Tag>
            <Tag>
              {params?.inputWidth ?? 640}×{params?.inputHeight ?? 640}
            </Tag>
          </div>,
          true
        )}
        {field(
          labelWithTip(formatMessage('Vision.task.outputs'), formatMessage('Vision.task.outputsTip')),
          <div className={styles.tagWrap}>
            {(task.outputCategories || []).map((output) => (
              <Tag key={output} color="green">
                {output}
              </Tag>
            ))}
          </div>,
          true
        )}
        <div className={styles.topicCard}>
          <div className={styles.topicTitle}>{formatMessage('Vision.task.unsTopicLabel')}</div>
          <div className={styles.topicList}>
            {(task.cameras || []).map((cam) => (
              <div key={cam.id} className={styles.topicRow}>
                <span className={styles.topicCam}>{cam.name || cam.code}</span>
                <code className={styles.topicCode}>{cam.unsTopic}</code>
              </div>
            ))}
          </div>
        </div>
      </div>
    </ProModal>
  );
};

export default TaskDetailModal;
