import { useEffect, useMemo, useState } from 'react';
import { useSearchParams } from 'react-router';
import { App, Button, Checkbox, Dropdown, Pagination, Popover, Select, Tooltip } from 'antd';
import { Clock, ImageOff, Webcam } from 'lucide-react';
import { listVisionTasks, type VisionTask } from '@/apis/core-api/task';
import { listVideoCameras, type VideoCamera } from '@/apis/core-api/video';
import {
  batchDeleteVisionEvents,
  deleteVisionEvent,
  eventScreenshotUrl,
  listVisionEvents,
  type VisionEvent,
  type VisionEventSort,
} from '@/apis/core-api/vision-results';
import { ComEmptyState, ComLayout, ProSearch } from '@/components';
import { Filter, OverflowMenuVertical, Renew } from '@/components/lucide-icon/carbon';
import { usePagination, useTranslate } from '@/hooks';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { formatTimestamp } from '@/utils/format';
import ResultDetailModal from './ResultDetailModal';
import shared from './index.module.scss';
import styles from './results.module.scss';

// 事件流轮询间隔(ms),多选模式下暂停避免选择被刷新打断。
const REFRESH_INTERVAL = 10000;
const DEFAULT_PAGE_SIZE = 10;

type ResultFilters = {
  search: string;
  cameraId?: number;
  eventName?: string;
  taskId?: number;
  sort: VisionEventSort;
};

const ResultsPanel = () => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [, setUrlParams] = useSearchParams();
  const [filters, setFilters] = useState<ResultFilters>({ search: '', sort: 'latest' });
  const [searchInput, setSearchInput] = useState('');
  const [cameras, setCameras] = useState<VideoCamera[]>([]);
  const [tasks, setTasks] = useState<VisionTask[]>([]);
  const [eventNames, setEventNames] = useState<string[]>([]);
  const [selectMode, setSelectMode] = useState(false);
  const [selectedIds, setSelectedIds] = useState<number[]>([]);
  const [detailEvent, setDetailEvent] = useState<VisionEvent | null>(null);

  const { loading, data, pagination, setSearchParams, refreshRequest } = usePagination<VisionEvent>({
    fetchApi: listVisionEvents,
    initPageSize: DEFAULT_PAGE_SIZE,
    initPageSizes: [10, 20, 50, 100],
    onSuccessCallback: (result) => setEventNames(result?.eventNames || []),
  });

  useEffect(() => {
    listVideoCameras({ pageNo: 1, pageSize: 100 })
      .then((result) => setCameras(result.data))
      .catch(() => setCameras([]));
    listVisionTasks({ pageNo: 1, pageSize: 100 })
      .then((result) => setTasks(result.data))
      .catch(() => setTasks([]));
  }, []);

  // 多选模式暂停轮询,避免刷新打断选择。
  useEffect(() => {
    if (selectMode) return undefined;
    const timer = window.setInterval(() => refreshRequest(false), REFRESH_INTERVAL);
    return () => window.clearInterval(timer);
  }, [selectMode, refreshRequest]);

  const applyFilters = (next: Partial<ResultFilters>) => {
    const merged = { ...filters, ...next };
    setFilters(merged);
    setSearchParams({
      search: merged.search || undefined,
      cameraId: merged.cameraId,
      eventName: merged.eventName,
      taskId: merged.taskId,
      sort: merged.sort,
    });
  };

  const resetFilters = () => {
    setSearchInput('');
    applyFilters({ search: '', cameraId: undefined, eventName: undefined, taskId: undefined });
  };

  const hasFilter = Boolean(filters.search || filters.cameraId || filters.eventName || filters.taskId);

  const exitSelectMode = () => {
    setSelectMode(false);
    setSelectedIds([]);
  };

  const toggleSelected = (id: number) => {
    setSelectedIds((prev) => (prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id]));
  };

  const allSelected = data.length > 0 && data.every((item) => selectedIds.includes(item.id));

  const toggleSelectAll = () => {
    setSelectedIds(allSelected ? [] : data.map((item) => item.id));
  };

  const goToTaskTab = () => {
    setUrlParams(
      (prev) => {
        const next = new URLSearchParams(prev);
        next.set('tab', 'task');
        return next;
      },
      { replace: true }
    );
  };

  const removeOne = (event: VisionEvent) => {
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          title: formatMessage('Vision.results.deleteTitle'),
          width: 420,
          content: (
            <span className="tier0-delete-confirm-message">
              {formatMessage('Vision.results.deleteConfirm', { name: event.eventName })}
            </span>
          ),
          onOk: async () => {
            await deleteVisionEvent(event.id);
            message.success(formatMessage('common.deleteSuccessfully'));
            setSelectedIds((prev) => prev.filter((id) => id !== event.id));
            refreshRequest();
          },
        },
        formatMessage
      )
    );
  };

  const removeSelected = () => {
    if (selectedIds.length === 0) return;
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          title: formatMessage('Vision.results.deleteTitle'),
          width: 420,
          content: (
            <span className="tier0-delete-confirm-message">
              {formatMessage('Vision.results.deleteBatch', { count: selectedIds.length })}
            </span>
          ),
          onOk: async () => {
            await batchDeleteVisionEvents(selectedIds);
            message.success(formatMessage('common.deleteSuccessfully'));
            exitSelectMode();
            refreshRequest();
          },
        },
        formatMessage
      )
    );
  };

  const cameraOptions = useMemo(() => cameras.map((camera) => ({ value: camera.id, label: camera.name })), [cameras]);

  const eventOptions = useMemo(() => eventNames.map((name) => ({ value: name, label: name })), [eventNames]);

  const taskOptions = useMemo(() => tasks.map((task) => ({ value: task.id, label: task.name })), [tasks]);

  // 漏斗筛选面板:承载不在工具栏常驻的筛选维度,并提供一键清空。
  const filterContent = (
    <div className={styles.filterPanel}>
      <div className={styles.filterField}>
        <span className={styles.filterLabel}>{formatMessage('Vision.results.filterTask')}</span>
        <Select
          allowClear
          showSearch
          optionFilterProp="label"
          value={filters.taskId}
          style={{ width: 220 }}
          placeholder={formatMessage('Vision.results.allTasks')}
          onChange={(value) => applyFilters({ taskId: value })}
          options={taskOptions}
        />
      </div>
      <Button size="small" disabled={!hasFilter} onClick={resetFilters}>
        {formatMessage('Vision.results.resetFilters')}
      </Button>
    </div>
  );

  const renderCard = (event: VisionEvent) => {
    const checked = selectedIds.includes(event.id);
    const cameraLabel = event.cameraName || (event.cameraId ? `#${event.cameraId}` : '-');
    const menuItems = [
      { key: 'viewTask', label: formatMessage('Vision.results.viewTask'), onClick: goToTaskTab },
      { key: 'delete', danger: true, label: formatMessage('common.delete'), onClick: () => removeOne(event) },
    ];
    return (
      <div
        key={event.id}
        className={styles.resultCard}
        role="button"
        tabIndex={0}
        onClick={() => (selectMode ? toggleSelected(event.id) : setDetailEvent(event))}
        onKeyDown={(keyEvent) => {
          if (keyEvent.key === 'Enter' || keyEvent.key === ' ') {
            keyEvent.preventDefault();
            if (selectMode) toggleSelected(event.id);
            else setDetailEvent(event);
          }
        }}
      >
        <div className={styles.thumbWrap}>
          {event.hasScreenshot ? (
            <img className={styles.thumbImg} src={eventScreenshotUrl(event.id)} alt={event.eventName} loading="lazy" />
          ) : (
            <div className={styles.thumbPlaceholder}>
              <ImageOff size={28} strokeWidth={1.5} />
            </div>
          )}
          {selectMode && (
            <Checkbox
              className={styles.cardCheck}
              checked={checked}
              onClick={(clickEvent) => clickEvent.stopPropagation()}
              onChange={() => toggleSelected(event.id)}
            />
          )}
          {!selectMode && (
            <div className={styles.thumbMask}>
              <Button className={styles.viewDetailsBtn}>{formatMessage('Vision.results.viewDetails')}</Button>
            </div>
          )}
          <Dropdown trigger={['click']} menu={{ items: menuItems }}>
            <button
              type="button"
              className={styles.cardMore}
              aria-label={formatMessage('common.operation')}
              onClick={(clickEvent) => clickEvent.stopPropagation()}
            >
              <OverflowMenuVertical size={16} />
            </button>
          </Dropdown>
        </div>
        <div className={styles.cardBody}>
          <div className={styles.cardTitle} title={event.eventName}>
            {event.eventName}
          </div>
          <div className={styles.cardMeta}>
            <span className={styles.cardMetaItem}>
              <Webcam size={12} strokeWidth={1.75} aria-hidden />
              <span title={cameraLabel}>{cameraLabel}</span>
            </span>
            <span className={styles.cardMetaItem}>
              <Clock size={12} strokeWidth={1.75} aria-hidden />
              <span>{formatTimestamp(event.createdAt)}</span>
            </span>
          </div>
        </div>
      </div>
    );
  };

  return (
    <ComLayout loading={loading}>
      <div className={shared.algoPanel}>
        <div className={shared.algoHeader}>
          <div className={shared.algoTitle}>{formatMessage('Vision.results.title')}</div>
          <div className={styles.headerActions}>
            <span className={styles.resultCount}>
              {formatMessage('Vision.results.count', { count: pagination.total })}
            </span>
            {selectMode ? (
              <>
                <Button
                  className={styles.deleteSelectedBtn}
                  disabled={selectedIds.length === 0}
                  onClick={removeSelected}
                >
                  {formatMessage('Vision.results.deleteSelected', { count: selectedIds.length })}
                </Button>
                <Button onClick={toggleSelectAll}>{formatMessage('Vision.results.selectAll')}</Button>
                <Button onClick={exitSelectMode}>{formatMessage('common.cancel')}</Button>
              </>
            ) : (
              <Button onClick={() => setSelectMode(true)}>{formatMessage('Vision.results.select')}</Button>
            )}
            <Tooltip title={formatMessage('common.refresh')}>
              <button
                type="button"
                className={`${shared.iconButton} ${styles.refreshButton}`}
                aria-label={formatMessage('common.refresh')}
                onClick={() => refreshRequest()}
              >
                <Renew size={16} />
              </button>
            </Tooltip>
          </div>
        </div>
        <div className={`${shared.algoToolbar} ${styles.resultsToolbar}`}>
          <div className={shared.algoToolbarLeft}>
            <Popover content={filterContent} trigger="click" placement="bottomLeft">
              <button
                type="button"
                className={`${shared.iconButton} ${shared.filterButton} ${hasFilter ? shared.filterActive : ''}`}
                aria-label={formatMessage('Vision.results.filterTitle')}
              >
                <Filter size={16} />
              </button>
            </Popover>
            <ProSearch
              size="sm"
              value={searchInput}
              placeholder={formatMessage('Vision.results.searchPlaceholder')}
              onChange={(event) => setSearchInput(event.target.value)}
              onSearch={(value) => applyFilters({ search: value })}
              className={styles.resultSearch}
            />
            <Select
              allowClear
              showSearch
              optionFilterProp="label"
              value={filters.cameraId}
              style={{ width: 148 }}
              placeholder={formatMessage('Vision.results.allCameras')}
              onChange={(value) => applyFilters({ cameraId: value })}
              options={cameraOptions}
            />
            <Select
              allowClear
              showSearch
              optionFilterProp="label"
              value={filters.eventName}
              style={{ width: 136 }}
              placeholder={formatMessage('Vision.results.allEvents')}
              onChange={(value) => applyFilters({ eventName: value })}
              options={eventOptions}
            />
          </div>
          <div className={styles.toolbarRight}>
            <span className={styles.sortLabel}>{formatMessage('Vision.results.sortBy')}</span>
            <Select
              value={filters.sort}
              style={{ width: 132 }}
              onChange={(value: VisionEventSort) => applyFilters({ sort: value })}
              options={[
                { value: 'latest', label: formatMessage('Vision.results.sortLatest') },
                { value: 'oldest', label: formatMessage('Vision.results.sortOldest') },
              ]}
            />
          </div>
        </div>
        <div className={`${shared.algoContent} ${styles.resultsContent}`}>
          {data.length === 0 && !loading ? (
            <div className={shared.panelEmpty}>
              <ComEmptyState
                className={shared.panelEmptyInner}
                variant="inline"
                title={formatMessage(hasFilter ? 'Vision.results.noMatchTitle' : 'Vision.results.emptyTitle')}
                description={formatMessage(hasFilter ? 'Vision.results.noMatchDesc' : 'Vision.results.emptyDesc')}
              />
            </div>
          ) : (
            <div className={styles.resultGrid}>{data.map(renderCard)}</div>
          )}
        </div>
        <div className={shared.algoPager}>
          <Pagination
            total={pagination.total}
            current={pagination.page}
            pageSize={pagination.pageSize}
            pageSizeOptions={pagination.pageSizes}
            showSizeChanger
            showQuickJumper
            onChange={(page, pageSize) => pagination.onChange({ page, pageSize })}
          />
        </div>
      </div>
      <ResultDetailModal
        event={detailEvent}
        onClose={() => setDetailEvent(null)}
        onViewTask={() => {
          setDetailEvent(null);
          goToTaskTab();
        }}
        onDeleted={() => {
          setDetailEvent(null);
          refreshRequest();
        }}
      />
    </ComLayout>
  );
};

export default ResultsPanel;
