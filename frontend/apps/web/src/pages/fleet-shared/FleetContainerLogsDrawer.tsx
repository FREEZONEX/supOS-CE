import { Button, Drawer, Input, Switch, Tooltip, Typography } from 'antd';
import cx from 'classnames';
import type { DragEndEvent } from '@dnd-kit/core';
import { DndContext, PointerSensor, useSensor } from '@dnd-kit/core';
import { horizontalListSortingStrategy, SortableContext, useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { restrictToHorizontalAxis } from '@dnd-kit/modifiers';
import { Close, Document, Download, Renew, Search, TrashCan } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { useTranslate } from '@/hooks';
import type { FleetContainerLogsResp, FleetContainerSummary } from '@/apis/core-api/fleet';
import type { FleetContainerLogTab } from './useFleetContainerLogs';
import {
  forwardRef,
  useEffect,
  useMemo,
  useRef,
  useState,
  type CSSProperties,
  type Dispatch,
  type HTMLAttributes,
  type PointerEvent as ReactPointerEvent,
  type RefObject,
  type SetStateAction,
  type UIEventHandler,
} from 'react';
import styles from './FleetContainerLogsDrawer.module.scss';

const MIN_DRAWER_HEIGHT = 200;
const MAX_DRAWER_HEIGHT = 720;
const DRAWER_VIEWPORT_GAP = 72;
const DEFAULT_DRAWER_HEIGHT = 320;

const clampDrawerHeight = (height: number) => {
  if (typeof window === 'undefined') return height;
  const maximum = Math.min(MAX_DRAWER_HEIGHT, Math.max(MIN_DRAWER_HEIGHT, window.innerHeight - DRAWER_VIEWPORT_GAP));
  const minimum = Math.min(MIN_DRAWER_HEIGHT, maximum);
  return Math.min(Math.max(height, minimum), maximum);
};

const getDefaultDrawerHeight = () =>
  typeof window === 'undefined' ? DEFAULT_DRAWER_HEIGHT : clampDrawerHeight(DEFAULT_DRAWER_HEIGHT);

const noop = () => undefined;
const noopWithID = () => undefined;
const noopReorder = () => undefined;
const noopSetFilterText: Dispatch<SetStateAction<string>> = () => undefined;

type DraggableTabNodeProps = HTMLAttributes<HTMLDivElement> & {
  'data-node-key': string;
};

const DraggableTabNode = forwardRef<HTMLDivElement, DraggableTabNodeProps>((props, ref) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: props['data-node-key'],
  });

  const style: CSSProperties = {
    ...props.style,
    transform: CSS.Transform.toString(transform && { ...transform, scaleX: 1 }),
    transition,
    display: 'inline-flex',
    maxWidth: '100%',
    zIndex: isDragging ? 2 : undefined,
    opacity: isDragging ? 0.92 : undefined,
  };

  return (
    <div
      ref={(node) => {
        setNodeRef(node);
        if (typeof ref === 'function') {
          ref(node);
        } else if (ref) {
          ref.current = node;
        }
      }}
      style={style}
      {...attributes}
      {...listeners}
    >
      {props.children}
    </div>
  );
});
DraggableTabNode.displayName = 'DraggableTabNode';

type FleetContainerLogsDrawerProps = {
  open: boolean;
  loading: boolean;
  downloading: boolean;
  logs: string;
  container?: FleetContainerSummary;
  live: boolean;
  follow: boolean;
  filterText: string;
  tabs: FleetContainerLogTab[];
  activeContainerID?: string;
  logsRef: RefObject<HTMLPreElement>;
  closeLogs: () => void;
  refreshLogs: () => Promise<void>;
  downloadLogs: () => Promise<FleetContainerLogsResp | undefined>;
  clearLogs: () => void;
  closeLogTab: (id: string) => void;
  reorderTabs?: (fromID: string, toID: string) => void;
  setActiveTab: (id: string) => void;
  setLive: Dispatch<SetStateAction<boolean>>;
  setFollow: Dispatch<SetStateAction<boolean>>;
  setFilterText: Dispatch<SetStateAction<string>>;
  handleScroll: UIEventHandler<HTMLPreElement>;
  onHeightChange?: (height: number) => void;
  onWidthChange?: (width: number) => void;
};

const safeFileName = (value: string) => value.trim().replace(/[^a-zA-Z0-9._-]+/g, '-') || 'container';

export const FleetContainerLogsDrawer = ({
  open,
  loading,
  downloading,
  logs,
  container,
  live,
  follow,
  filterText = '',
  tabs = [],
  activeContainerID,
  logsRef,
  closeLogs,
  refreshLogs,
  downloadLogs,
  clearLogs = noop,
  closeLogTab = noopWithID,
  reorderTabs = noopReorder,
  setActiveTab = noopWithID,
  setLive,
  setFollow,
  setFilterText = noopSetFilterText,
  handleScroll,
  onHeightChange,
  onWidthChange,
}: FleetContainerLogsDrawerProps) => {
  const formatMessage = useTranslate();
  const [drawerHeight, setDrawerHeight] = useState(getDefaultDrawerHeight);
  const [resizing, setResizing] = useState(false);
  const resizeStartRef = useRef<{ y: number; height: number } | undefined>(undefined);
  const tabSensor = useSensor(PointerSensor, { activationConstraint: { distance: 10 } });
  const tabKeys = useMemo(() => tabs.map((tab) => tab.id), [tabs]);
  const filteredLogs = useMemo(() => {
    const keyword = filterText.trim().toLowerCase();
    if (!keyword) return logs;
    return logs
      .split('\n')
      .filter((line) => line.toLowerCase().includes(keyword))
      .join('\n');
  }, [filterText, logs]);

  useEffect(() => {
    const handleViewportResize = () => setDrawerHeight((current) => clampDrawerHeight(current));
    window.addEventListener('resize', handleViewportResize);
    return () => window.removeEventListener('resize', handleViewportResize);
  }, []);

  useEffect(() => {
    onHeightChange?.(open ? drawerHeight : 0);
    onWidthChange?.(0);
  }, [drawerHeight, onHeightChange, onWidthChange, open]);

  useEffect(
    () => () => {
      onHeightChange?.(0);
      onWidthChange?.(0);
    },
    [onHeightChange, onWidthChange]
  );

  const handleResizeStart = (event: ReactPointerEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
    resizeStartRef.current = { y: event.clientY, height: drawerHeight };
    event.currentTarget.setPointerCapture(event.pointerId);
    setResizing(true);
  };

  const handleResizeMove = (event: ReactPointerEvent<HTMLDivElement>) => {
    const start = resizeStartRef.current;
    if (!start) return;
    setDrawerHeight(clampDrawerHeight(start.height + start.y - event.clientY));
  };

  const handleResizeEnd = (event: ReactPointerEvent<HTMLDivElement>) => {
    resizeStartRef.current = undefined;
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
    setResizing(false);
  };

  const handleTabDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return;
    reorderTabs(String(active.id), String(over.id));
  };

  const handleDownloadLogs = async () => {
    if (!container) return;
    const containerName = safeFileName(container?.displayName || container?.serviceName || 'container');
    const result = await downloadLogs();
    if (!result) return;
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const url = URL.createObjectURL(new Blob([result.logs || ''], { type: 'text/plain;charset=utf-8' }));
    const link = document.createElement('a');
    link.href = url;
    link.download = `${containerName}-${timestamp}.log`;
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.setTimeout(() => URL.revokeObjectURL(url), 0);
  };

  return (
    <Drawer
      open={open}
      placement="bottom"
      height={drawerHeight}
      mask={false}
      push={false}
      autoFocus={false}
      rootStyle={{ top: 'auto', height: drawerHeight }}
      title={
        <div className={styles.titleBar}>
          <DndContext sensors={[tabSensor]} modifiers={[restrictToHorizontalAxis]} onDragEnd={handleTabDragEnd}>
            <SortableContext items={tabKeys} strategy={horizontalListSortingStrategy}>
              <div className={styles.tabs}>
                {tabs.map((tab) => {
                  const label = tab.container.displayName || tab.container.serviceName || tab.id.slice(0, 12);
                  const active = tab.id === activeContainerID;
                  return (
                    <DraggableTabNode data-node-key={tab.id} key={tab.id}>
                      <button
                        type="button"
                        className={cx(styles.tabItem, active && styles.tabItemActive)}
                        title={label}
                        onClick={() => setActiveTab(tab.id)}
                      >
                        <span className={styles.tabLabel}>
                          <Document {...toolbarIconProps} />
                          <span>{label}</span>
                        </span>
                        {active ? (
                          <span
                            className={styles.tabClose}
                            role="button"
                            tabIndex={0}
                            aria-label={formatMessage('common.close')}
                            onClick={(event) => {
                              event.stopPropagation();
                              closeLogTab(tab.id);
                            }}
                            onKeyDown={(event) => {
                              if (event.key === 'Enter' || event.key === ' ') {
                                event.preventDefault();
                                event.stopPropagation();
                                closeLogTab(tab.id);
                              }
                            }}
                          >
                            <Close size={12} strokeWidth={2} aria-hidden />
                          </span>
                        ) : null}
                      </button>
                    </DraggableTabNode>
                  );
                })}
              </div>
            </SortableContext>
          </DndContext>
        </div>
      }
      onClose={closeLogs}
      rootClassName={cx(styles.drawer, resizing && styles.resizing)}
    >
      <div
        className={styles.resizeEdge}
        role="separator"
        aria-orientation="horizontal"
        aria-label={formatMessage('common.drag')}
        title={formatMessage('common.drag')}
        onDoubleClick={() => setDrawerHeight(getDefaultDrawerHeight())}
        onPointerDown={handleResizeStart}
        onPointerMove={handleResizeMove}
        onPointerUp={handleResizeEnd}
        onPointerCancel={handleResizeEnd}
      />
      <div className={styles.viewer}>
        <div className={styles.toolbar}>
          <Input
            allowClear
            value={filterText}
            className={styles.search}
            prefix={<Search {...toolbarIconProps} />}
            placeholder={formatMessage('fleet.runtime.logFilterPlaceholder', undefined, formatMessage('common.search'))}
            onChange={(event) => setFilterText(event.target.value)}
          />
          <div className={styles.toolbarActions}>
            <Tooltip title={formatMessage('fleet.runtime.followLatest')}>
              <Button
                size="small"
                className={cx(styles.actionButton, !follow && styles.actionButtonActive)}
                disabled={follow}
                onClick={() => setFollow(true)}
              >
                {formatMessage('fleet.runtime.followLatest')}
              </Button>
            </Tooltip>
            <Tooltip title={formatMessage('common.clear')}>
              <span className={styles.tooltipAnchor}>
                <Button
                  size="small"
                  className={styles.actionButton}
                  icon={<TrashCan {...toolbarIconProps} />}
                  disabled={!logs}
                  onClick={clearLogs}
                >
                  {formatMessage('common.clear')}
                </Button>
              </span>
            </Tooltip>
            <span className={styles.toolbarDivider} aria-hidden />
            <Tooltip title={formatMessage('common.download')}>
              <span className={styles.tooltipAnchor}>
                <Button
                  size="small"
                  className={styles.iconAction}
                  aria-label={formatMessage('common.download')}
                  icon={<Download {...toolbarIconProps} />}
                  disabled={!container}
                  loading={downloading}
                  onClick={() => void handleDownloadLogs()}
                />
              </span>
            </Tooltip>
            <Tooltip title={formatMessage('common.refresh')}>
              <Button
                size="small"
                className={cx(styles.iconAction, loading && styles.loadingButton)}
                aria-label={formatMessage('common.refresh')}
                icon={<Renew {...toolbarIconProps} />}
                loading={loading}
                onClick={() => void refreshLogs()}
              />
            </Tooltip>
            <span className={styles.toolbarDivider} aria-hidden />
            <Tooltip title={formatMessage('fleet.runtime.liveUpdate')}>
              <label className={styles.liveControl}>
                <Switch size="small" checked={live} onChange={setLive} />
                <span>{formatMessage('fleet.runtime.liveUpdate')}</span>
              </label>
            </Tooltip>
          </div>
        </div>
        <pre ref={logsRef} className={styles.terminal} onScroll={handleScroll}>
          {loading
            ? formatMessage('common.loading')
            : filteredLogs || (filterText.trim() ? formatMessage('common.noData') : '-')}
        </pre>
        <footer className={styles.footer}>
          <Typography.Text className={styles.hint}>{formatMessage('fleet.runtime.logsHint')}</Typography.Text>
        </footer>
      </div>
    </Drawer>
  );
};
