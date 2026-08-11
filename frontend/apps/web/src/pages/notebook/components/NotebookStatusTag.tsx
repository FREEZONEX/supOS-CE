import { Tag } from 'antd';
import type { FC } from 'react';
import { useTranslate } from '@/hooks';

export type NotebookStatusValue = 'running' | 'stopped' | 'idle' | (string & {});

export const normalizeNotebookStatus = (status?: string): string => {
  const normalized = String(status || 'stopped')
    .trim()
    .toLowerCase();
  if (normalized === 'idle') {
    return 'stopped';
  }
  return normalized || 'stopped';
};

export const getNotebookStatusTagClassName = (status?: string): string => {
  const normalized = normalizeNotebookStatus(status);
  if (normalized === 'running') {
    return 'notebook-status-running';
  }
  if (normalized === 'stopped') {
    return 'notebook-status-stopped';
  }
  return 'notebook-type-tag';
};

export const getNotebookStatusLabel = (
  status: string | undefined,
  formatMessage: ReturnType<typeof useTranslate>
): string => {
  const normalized = normalizeNotebookStatus(status);
  if (normalized === 'running') {
    return formatMessage('common.running', {}, 'Running');
  }
  if (normalized === 'stopped') {
    return formatMessage('common.stopped', {}, 'Stopped');
  }
  return formatMessage('Notebook.statusUnknown', {}, 'Unknown');
};

export const isNotebookRunning = (status?: string): boolean => normalizeNotebookStatus(status) === 'running';

interface NotebookStatusTagProps {
  status?: string;
}

const NotebookStatusTag: FC<NotebookStatusTagProps> = ({ status }) => {
  const formatMessage = useTranslate();

  return (
    <Tag className={getNotebookStatusTagClassName(status)}>{getNotebookStatusLabel(status, formatMessage)}</Tag>
  );
};

export default NotebookStatusTag;
