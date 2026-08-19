import { getNotebookSnapshotList } from '@/apis/core-api/notebook';
import { useTranslate } from '@/hooks';
import ComEmpty from '@/components/com-empty';
import { Button, Spin } from 'antd';
import { type FC, useCallback, useEffect, useState } from 'react';

interface SnapshotItem {
  id: number;
  versionName: string;
}

interface HistoryPopoverProps {
  notebookId: string | number;
  visible?: boolean;
  refreshKey?: number;
  onSelectItem: (snapshotId: number) => void;
  onViewAll: () => void;
}

const HistoryPopover: FC<HistoryPopoverProps> = ({ notebookId, visible, refreshKey, onSelectItem, onViewAll }) => {
  const formatMessage = useTranslate();
  const [loading, setLoading] = useState(false);
  const [list, setList] = useState<SnapshotItem[]>([]);

  const loadList = useCallback(async () => {
    if (!notebookId) return;
    setLoading(true);
    try {
      const resp = await getNotebookSnapshotList(notebookId, { pageNo: 1, pageSize: 50 });
      setList(resp?.list ?? []);
    } catch {
      // silent
    } finally {
      setLoading(false);
    }
  }, [notebookId]);

  useEffect(() => {
    if (visible) {
      void loadList();
    }
  }, [visible, refreshKey, loadList]);

  return (
    <div style={{ width: 240 }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 8,
          fontWeight: 600,
        }}
      >
        <span>{formatMessage('Notebook.snapshot.drawerTitle', {}, 'Version History')}</span>
        <Button type="text" size="small" disabled={list.length === 0} onClick={onViewAll}>
          {formatMessage('Notebook.snapshot.viewAll', {}, 'View All')}
        </Button>
      </div>
      <Spin spinning={loading}>
        {list.length === 0 && !loading ? (
          <ComEmpty description={formatMessage('Notebook.snapshot.empty', {}, 'No versions yet')} />
        ) : (
          <div style={{ display: 'flex', flexDirection: 'column', maxHeight: 300, overflowY: 'auto' }}>
            {list.map((item) => (
              <div
                key={item.id}
                onClick={() => onSelectItem(item.id)}
                style={{
                  padding: '8px 12px',
                  cursor: 'pointer',
                  borderRadius: 4,
                  transition: 'background 0.2s',
                  fontSize: 13,
                  color: 'var(--ui-text-color)',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = 'var(--ui-t-card-hover-bg)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }}
              >
                {item.versionName}
              </div>
            ))}
          </div>
        )}
      </Spin>
    </div>
  );
};

export default HistoryPopover;
