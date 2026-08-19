import { useTranslate } from '@/hooks';
import type { VisionTaskStatus } from '@/apis/core-api/task';
import { TASK_STATUS_META, type TaskStatusTone } from './task-meta';
import styles from './index.module.scss';

const TONE_CLASS: Record<TaskStatusTone, string> = {
  running: styles.toneRunning,
  starting: styles.toneStarting,
  warning: styles.toneWarning,
  error: styles.toneError,
  idle: styles.toneIdle,
  scheduled: styles.toneScheduled,
};

// 任务状态胶囊：每个状态一套底色/文字色（对齐 Figma 的任务状态图例）。
const TaskStatusTag = ({ status }: { status: VisionTaskStatus }) => {
  const formatMessage = useTranslate();
  const meta = TASK_STATUS_META[status] || TASK_STATUS_META.stopped;
  return <span className={`${styles.taskStatusTag} ${TONE_CLASS[meta.tone]}`}>{formatMessage(meta.labelKey)}</span>;
};

export default TaskStatusTag;
