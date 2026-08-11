import {
  getNotebook,
  getNotebookList,
  getNotebookSnapshotList,
  revertNotebookSnapshot,
  saveAsNotebookSnapshot,
} from '@/apis/core-api/notebook';
import { ButtonPermission } from '@/common-types/button-permission';
import { useTranslate } from '@/hooks';
import { hasPermission } from '@/utils/auth';
import { Close } from '@/components/lucide-icon/carbon';
import { App, Button, Empty, Spin, Tag, Typography } from 'antd';
import { type FC, useCallback, useEffect, useState } from 'react';
import SaveAsModal from './SaveAsModal';
import styles from './VersionDrawer.module.scss';

const formatTimeAgo = (dateStr: string, formatMessage: ReturnType<typeof useTranslate>): string => {
  const now = Date.now();
  const past = new Date(dateStr).getTime();
  const diffMs = now - past;
  if (diffMs < 0) return formatMessage('Notebook.snapshot.justNow', {}, 'just now');
  const minutes = Math.floor(diffMs / 60000);
  if (minutes < 1) return formatMessage('Notebook.snapshot.justNow', {}, 'just now');
  if (minutes < 60) return formatMessage('Notebook.snapshot.minutesAgo', { count: minutes }, `${minutes} min ago`);
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return formatMessage('Notebook.snapshot.hoursAgo', { count: hours }, `${hours} hours ago`);
  const days = Math.floor(hours / 24);
  return formatMessage('Notebook.snapshot.daysAgo', { count: days }, `${days} days ago`);
};

interface SnapshotItem {
  id: number;
  notebookId: number;
  versionName: string;
  description: string;
  snapshotType: string;
  fileSize: number;
  creator: string;
  createdAt: string;
  updatedAt: string;
  isCurrent?: boolean;
}

interface HistoryPanelProps {
  notebookId: string | number;
  selectedSnapshotId?: number | null;
  onClose: () => void;
  onReverted?: () => void | Promise<void>;
}

const HistoryPanel: FC<HistoryPanelProps> = ({ notebookId, selectedSnapshotId, onClose, onReverted }) => {
  const { message, modal } = App.useApp();
  const formatMessage = useTranslate();
  const canManage = hasPermission(ButtonPermission['Notebook.manage']);
  const [loading, setLoading] = useState(false);
  const [list, setList] = useState<SnapshotItem[]>([]);
  const [revertingId, setRevertingId] = useState<number | null>(null);
  const [saveAsOpen, setSaveAsOpen] = useState(false);
  const [saveAsLoading, setSaveAsLoading] = useState(false);
  const [activeSnapshotId, setActiveSnapshotId] = useState<number | null>(selectedSnapshotId ?? null);
  const [existingNotebookNames, setExistingNotebookNames] = useState<string[]>([]);

  const pageSize = 50;

  const getSnapshotTypeLabel = (snapshotType?: string) => {
    const normalized = String(snapshotType || '').toLowerCase();
    if (normalized === 'auto') return formatMessage('Notebook.snapshot.typeAuto', {}, 'Auto');
    return formatMessage('Notebook.snapshot.typeManual', {}, 'Manual');
  };

  const getSnapshotTypeColor = (snapshotType?: string) => {
    if (!snapshotType || snapshotType.toLowerCase() === 'manual') return 'purple';
    return 'blue';
  };

  const loadList = useCallback(
    async (page = 1) => {
      if (!notebookId) return;
      setLoading(true);
      try {
        const resp = await getNotebookSnapshotList(notebookId, { pageNo: page, pageSize });
        setList(resp?.list ?? []);
      } catch {
        // silent
      } finally {
        setLoading(false);
      }
    },
    [notebookId]
  );

  useEffect(() => {
    void loadList(1);
    void (async () => {
      try {
        const detail = await getNotebook(notebookId);
        const resp = await getNotebookList({ listType: 'all', folderId: detail?.folderId || 0, search: '' });
        setExistingNotebookNames(
          (resp?.list ?? [])
            .filter((item: { type?: string; name?: string }) => item.type === 'notebook' && item.name)
            .map((item: { name: string }) => item.name)
        );
      } catch {
        setExistingNotebookNames([]);
      }
    })();
  }, [loadList, notebookId]);

  useEffect(() => {
    if (selectedSnapshotId) {
      setActiveSnapshotId(selectedSnapshotId);
    }
  }, [selectedSnapshotId]);

  const handleRevert = async (snapshotId: number) => {
    setRevertingId(snapshotId);
    try {
      await revertNotebookSnapshot(notebookId, snapshotId);
      message.success(formatMessage('Notebook.snapshot.revertSuccess', {}, 'Reverted to this version'));
      void loadList(1);
      await onReverted?.();
      onClose();
    } catch {
      // error handled by interceptor
    } finally {
      setRevertingId(null);
    }
  };

  const handleSaveAsOpen = (snapshotId: number) => {
    setActiveSnapshotId(snapshotId);
    setSaveAsOpen(true);
  };

  const handleSaveAsSubmit = async (newName: string, description: string) => {
    if (!activeSnapshotId) return false;
    const trimmedNewName = newName.trim();
    if (trimmedNewName && existingNotebookNames.includes(trimmedNewName)) {
      message.error(
        formatMessage(
          'Notebook.snapshot.duplicateNameSaveFailed',
          { name: trimmedNewName },
          `The name "${trimmedNewName}" already exists. Choose another name.`
        )
      );
      return false;
    }

    setSaveAsLoading(true);
    try {
      const resp = await saveAsNotebookSnapshot(notebookId, activeSnapshotId, { newName, description });
      if (trimmedNewName && resp?.name && resp.name !== trimmedNewName) {
        message.success(
          formatMessage(
            'Notebook.snapshot.saveAsRenamedSuccess',
            { name: resp.name },
            `Saved as new Notebook. Actual name: ${resp.name}`
          )
        );
      } else {
        message.success(formatMessage('Notebook.snapshot.saveAsSuccess', {}, 'Saved as new Notebook'));
      }
      setSaveAsOpen(false);
      return true;
    } catch {
      message.error(formatMessage('Notebook.snapshot.saveAsFailed', {}, 'Failed to save as new Notebook'));
      return false;
    } finally {
      setSaveAsLoading(false);
    }
  };

  return (
    <>
      <div
        style={{
          height: 'var(--ui-bar-height, 52px)',
          padding: '0 16px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: '1px solid var(--ui-home-border-color)',
          flexShrink: 0,
          fontWeight: 700,
          fontSize: 16,
          color: 'var(--ui-text-color)',
        }}
      >
        <span>{formatMessage('Notebook.snapshot.drawerTitle', {}, 'History')}</span>
        <Close size={20} style={{ cursor: 'pointer' }} onClick={onClose} />
      </div>
      <div style={{ flex: 1, overflow: 'auto', padding: 16 }}>
        <Spin spinning={loading}>
          {list.length === 0 && !loading ? (
            <Empty description={formatMessage('Notebook.snapshot.empty', {}, 'No versions yet')} />
          ) : (
            <div className={styles.timeline}>
              {list.map((item) => {
                const selected = activeSnapshotId === item.id;
                const isCurrent = item.isCurrent === true;
                return (
                  <div key={item.id} className={styles.timelineItem}>
                    <div className={styles.timelineRail} aria-hidden="true">
                      <span className={`${styles.timelineDot} ${selected ? styles.timelineDotCurrent : ''}`} />
                    </div>

                    <div className={styles.timelineContent}>
                      <div
                        className={`${styles.snapshotCard} ${selected ? styles.snapshotCardSelected : ''}`}
                        onClick={() => setActiveSnapshotId(selected ? null : item.id)}
                        role="button"
                        tabIndex={0}
                        onKeyDown={(event) => {
                          if (event.key === 'Enter' || event.key === ' ') {
                            event.preventDefault();
                            setActiveSnapshotId(selected ? null : item.id);
                          }
                        }}
                      >
                        <div className={styles.snapshotBody}>
                          <div className={styles.snapshotHeader}>
                            <div className={styles.snapshotTitle}>
                              <Typography.Text
                                strong
                                style={{ fontSize: 14 }}
                                className={styles.snapshotTitleText}
                                ellipsis={{ tooltip: item.versionName }}
                              >
                                {item.versionName}
                              </Typography.Text>
                            </div>
                            <div style={{ display: 'flex', gap: 4, flexShrink: 0 }}>
                              <Tag bordered={false} color={getSnapshotTypeColor(item.snapshotType)}>
                                {getSnapshotTypeLabel(item.snapshotType)}
                              </Tag>
                              {isCurrent && (
                                <Tag bordered color="lime">
                                  {formatMessage('Notebook.snapshot.current', {}, 'Current')}
                                </Tag>
                              )}
                            </div>
                          </div>

                          {item.description ? (
                            <Typography.Paragraph
                              className={styles.snapshotDescription}
                              ellipsis={{ rows: 2, tooltip: item.description }}
                            >
                              {item.description}
                            </Typography.Paragraph>
                          ) : null}

                          <div className={styles.snapshotMetaRow}>
                            <Typography.Text
                              className={styles.snapshotMetaName}
                              ellipsis={{ tooltip: item.creator || '-' }}
                            >
                              {item.creator || '-'}
                            </Typography.Text>
                            <Typography.Text
                              type="secondary"
                              style={{ fontSize: 12 }}
                              className={styles.snapshotMetaTime}
                            >
                              {item.createdAt ? formatTimeAgo(item.createdAt, formatMessage) : '-'}
                            </Typography.Text>
                          </div>
                        </div>
                      </div>

                      {selected && canManage ? (
                        <div className={styles.snapshotActionsOuter}>
                          <Button block size="small" onClick={() => handleSaveAsOpen(item.id)}>
                            {formatMessage('Notebook.snapshot.saveAs', {}, 'Save as New')}
                          </Button>
                          <Button
                            block
                            type="primary"
                            size="small"
                            loading={revertingId === item.id}
                            onClick={() => {
                              modal.confirm({
                                icon: null,
                                title: formatMessage(
                                  'Notebook.snapshot.rollbackTitle',
                                  { name: item.versionName },
                                  `Roll back to ${item.versionName}?`
                                ),
                                content: formatMessage(
                                  'Notebook.snapshot.rollbackContent',
                                  {},
                                  'This will overwrite any unsaved changes in your workspace and trigger a restart.'
                                ),
                                okText: formatMessage('common.confirm', {}, 'Confirm'),
                                cancelText: formatMessage('common.cancel', {}, 'Cancel'),
                                onOk: () => handleRevert(item.id),
                              });
                            }}
                          >
                            {formatMessage('Notebook.snapshot.revert', {}, 'Revert to This version')}
                          </Button>
                        </div>
                      ) : null}
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </Spin>
      </div>
      <SaveAsModal
        open={saveAsOpen}
        confirmLoading={saveAsLoading}
        onCancel={() => setSaveAsOpen(false)}
        onSubmit={handleSaveAsSubmit}
      />
    </>
  );
};

export default HistoryPanel;
