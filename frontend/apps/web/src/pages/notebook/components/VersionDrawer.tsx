import {
  getNotebook,
  getNotebookList,
  getNotebookSnapshotList,
  revertNotebookSnapshot,
  saveAsNotebookSnapshot,
} from '@/apis/core-api/notebook';
import { ButtonPermission } from '@/common-types/button-permission';
import { useTranslate } from '@/hooks';
import { formatTimestamp } from '@/utils/format';
import { hasPermission } from '@/utils/auth';
import { App, Button, Drawer, Empty, Popconfirm, Spin, Tag, Typography } from 'antd';
import { type FC, useCallback, useEffect, useState } from 'react';
import SaveAsModal from './SaveAsModal';
import styles from './VersionDrawer.module.scss';

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

interface VersionDrawerProps {
  open: boolean;
  notebookId: string | number;
  onClose: () => void;
  onReverted?: () => void | Promise<void>;
}

const VersionDrawer: FC<VersionDrawerProps> = ({ open, notebookId, onClose, onReverted }) => {
  const { message } = App.useApp();
  const formatMessage = useTranslate();
  const canManage = hasPermission(ButtonPermission['Notebook.manage']);
  const [loading, setLoading] = useState(false);
  const [list, setList] = useState<SnapshotItem[]>([]);
  const [revertingId, setRevertingId] = useState<number | null>(null);
  const [saveAsOpen, setSaveAsOpen] = useState(false);
  const [saveAsLoading, setSaveAsLoading] = useState(false);
  const [activeSnapshotId, setActiveSnapshotId] = useState<number | null>(null);
  const [existingNotebookNames, setExistingNotebookNames] = useState<string[]>([]);

  const pageSize = 20;

  const getSnapshotTypeLabel = (snapshotType?: string) => {
    const normalized = String(snapshotType || '').toLowerCase();
    if (normalized === 'auto') return formatMessage('Notebook.snapshot.typeAuto', {}, 'Auto');
    return formatMessage('Notebook.snapshot.typeManual', {}, 'Manual');
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
    if (open) {
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
    } else {
      setList([]);
      setActiveSnapshotId(null);
      setExistingNotebookNames([]);
    }
  }, [open, loadList]);

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
      onClose();
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
      <Drawer
        title={formatMessage('Notebook.snapshot.drawerTitle', {}, 'Version History')}
        open={open}
        onClose={onClose}
        width={380}
        footer={null}
      >
        <Spin spinning={loading}>
          {list.length === 0 && !loading ? (
            <Empty description={formatMessage('Notebook.snapshot.empty', {}, 'No versions yet')} />
          ) : (
            <div className={styles.timeline}>
              {list.map((item) => {
                const selected = activeSnapshotId === item.id;
                return (
                  <div key={item.id} className={styles.timelineItem}>
                    <div className={styles.timelineRail} aria-hidden="true">
                      <span className={`${styles.timelineDot} ${selected ? styles.timelineDotCurrent : ''}`} />
                    </div>

                    <div className={styles.timelineContent}>
                      <div
                        className={`${styles.snapshotCard} ${selected ? styles.snapshotCardSelected : ''}`}
                        onClick={() => setActiveSnapshotId(item.id)}
                        role="button"
                        tabIndex={0}
                        onKeyDown={(event) => {
                          if (event.key === 'Enter' || event.key === ' ') {
                            event.preventDefault();
                            setActiveSnapshotId(item.id);
                          }
                        }}
                      >
                        <div className={styles.snapshotBody}>
                          <div className={styles.snapshotHeader}>
                            <div className={styles.snapshotTitle}>
                              <Typography.Text
                                strong
                                className={styles.snapshotTitleText}
                                ellipsis={{ tooltip: item.versionName }}
                              >
                                {item.versionName}
                              </Typography.Text>
                            </div>
                            <Tag bordered color="purple">
                              {getSnapshotTypeLabel(item.snapshotType)}
                            </Tag>
                          </div>

                          <div className={styles.snapshotMeta}>
                            <Typography.Text
                              type="secondary"
                              className={styles.snapshotMetaName}
                              ellipsis={{ tooltip: item.creator || '-' }}
                            >
                              {item.creator || '-'}
                            </Typography.Text>
                            <Typography.Text type="secondary" className={styles.snapshotMetaTime}>
                              {item.createdAt ? formatTimestamp(item.createdAt) : '-'}
                            </Typography.Text>
                          </div>

                          {item.description ? (
                            <Typography.Paragraph
                              className={styles.snapshotDescription}
                              ellipsis={{ rows: 2, tooltip: item.description }}
                            >
                              {item.description}
                            </Typography.Paragraph>
                          ) : null}

                          {selected && canManage ? (
                            <div className={styles.snapshotActions} onClick={(event) => event.stopPropagation()}>
                              <Button block size="small" onClick={() => handleSaveAsOpen(item.id)}>
                                {formatMessage('Notebook.snapshot.saveAs', {}, 'Save as New')}
                              </Button>
                              <Popconfirm
                                title={formatMessage(
                                  'Notebook.snapshot.revertConfirm',
                                  {},
                                  'Revert to this version? Current content will be replaced.'
                                )}
                                onConfirm={() => void handleRevert(item.id)}
                                okText={formatMessage('common.confirm', {}, 'Confirm')}
                                cancelText={formatMessage('common.cancel', {}, 'Cancel')}
                              >
                                <Button block type="primary" size="small" loading={revertingId === item.id}>
                                  {formatMessage('Notebook.snapshot.revert', {}, 'Revert to This version')}
                                </Button>
                              </Popconfirm>
                            </div>
                          ) : null}
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </Spin>
      </Drawer>
      <SaveAsModal
        open={saveAsOpen}
        confirmLoading={saveAsLoading}
        onCancel={() => setSaveAsOpen(false)}
        onSubmit={handleSaveAsSubmit}
      />
    </>
  );
};

export default VersionDrawer;
