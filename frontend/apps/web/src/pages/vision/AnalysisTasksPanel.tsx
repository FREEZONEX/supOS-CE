import { useEffect, useState } from 'react';
import { App, Dropdown, Popover, Select, Tooltip, Typography, type MenuProps, type TableColumnsType } from 'antd';
import {
  listVisionTasks,
  deleteVisionTask,
  startVisionTask,
  stopVisionTask,
  type VisionTask,
  type VisionTaskStatus,
} from '@/apis/core-api/task';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import { ComEmptyState, ComLayout, ProSearch, ProTable } from '@/components';
import {
  Add,
  Edit,
  Filter,
  Help,
  Information,
  ListChecked,
  OverflowMenuHorizontal,
  Route,
  Run,
  StopOutline,
  TrashCan,
  View,
} from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { usePagination, useTranslate } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import TaskFormModal from './TaskFormModal';
import TaskTestFrameModal from './TaskTestFrameModal';
import TaskRegionModal from './TaskRegionModal';
import TaskDetailModal from './TaskDetailModal';
import TaskPreviewModal from './TaskPreviewModal';
import TaskStatusTag from './TaskStatusTag';
import { TASK_STATUS_META, ACTIVE_STATUSES, spatialModeFor } from './task-meta';
import styles from './index.module.scss';

const { Text } = Typography;

/** Camera 列只展示首个名称,其余用 +N 折叠,避免多标签竖排换行。 */
const MAX_VISIBLE_CAMERAS = 1;
// 存在活动任务时轮询刷新,让状态/统计接近实时。
const POLL_INTERVAL = 5000;

// idle 是「已部署但不在排班窗口」,同样可以停止。
const isStoppable = (status: VisionTaskStatus) => ['running', 'starting', 'recovering', 'idle'].includes(status);

const AnalysisTasksPanel = () => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<string>();
  const [formOpen, setFormOpen] = useState(false);
  const [testTask, setTestTask] = useState<VisionTask | null>(null);
  const [regionTask, setRegionTask] = useState<VisionTask | null>(null);
  const [editTask, setEditTask] = useState<VisionTask | null>(null);
  const [previewTask, setPreviewTask] = useState<VisionTask | null>(null);
  const [detailTask, setDetailTask] = useState<VisionTask | null>(null);
  const [busyId, setBusyId] = useState<number | null>(null);
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const resourceKeys = (currentUserInfo as (typeof currentUserInfo & { resourceKeys?: string[] }) | undefined)
    ?.resourceKeys;
  const canManage =
    systemInfo?.authEnable === false ||
    currentUserInfo?.superAdmin === true ||
    resourceKeys?.includes('vision.task.manage');
  const { loading, data, pagination, setSearchParams, refreshRequest } = usePagination<VisionTask>({
    fetchApi: listVisionTasks,
  });

  // 有活动任务时按固定间隔轮询,刷新运行状态与统计。
  const hasActive = data.some((task) => ACTIVE_STATUSES.includes(task.status));
  useEffect(() => {
    if (!hasActive) return undefined;
    const timer = window.setInterval(() => refreshRequest(), POLL_INTERVAL);
    return () => window.clearInterval(timer);
  }, [hasActive, refreshRequest]);

  const runTask = async (task: VisionTask) => {
    setBusyId(task.id);
    try {
      await startVisionTask(task.id);
      message.success(formatMessage('Vision.task.startSuccess'));
      refreshRequest();
    } finally {
      setBusyId(null);
    }
  };

  const haltTask = async (task: VisionTask) => {
    setBusyId(task.id);
    try {
      await stopVisionTask(task.id);
      message.success(formatMessage('Vision.task.stopSuccess'));
      refreshRequest();
    } finally {
      setBusyId(null);
    }
  };

  const remove = (task: VisionTask) => {
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          ...createDeleteConfirmOptions({
            title: formatMessage('Vision.task.deleteTitle'),
            name: task.name,
            formatMessage,
          }),
          onOk: async () => {
            await deleteVisionTask(task.id);
            message.success(formatMessage('common.deleteSuccessfully'));
            refreshRequest();
          },
        },
        formatMessage
      )
    );
  };

  // 操作全部收进三点下拉,可执行项按状态给(对齐设计稿的状态-操作矩阵):
  // 运行中 → 详情/预览/停止;启动或停止过渡中 → 只有详情;已停止或出错 → 详情/启动/编辑/(区域)/测试/删除。
  const operationItems = (record: VisionTask): MenuProps['items'] => {
    const transitioning = record.status === 'starting' || record.status === 'stopping';
    const running = isStoppable(record.status) && !transitioning;
    const idle = !isStoppable(record.status);
    const modelReady = !record.algorithmMissing && record.algorithm?.modelStatus === 'available';
    // definition.spatialMode==='none' 的算法没有区域/计数线配置项。
    const regionMode = spatialModeFor(record.algorithm);
    // 启停请求在飞时先禁掉这两项,避免重复下发。
    const busy = busyId === record.id;
    const details = {
      key: 'details',
      icon: <Information size={16} />,
      label: formatMessage('Vision.task.details'),
      onClick: () => setDetailTask(record),
    };
    if (transitioning) return [details];
    if (running) {
      return [
        details,
        {
          key: 'preview',
          icon: <View size={16} />,
          label: formatMessage('Vision.task.previewTitle'),
          onClick: () => setPreviewTask(record),
        },
        {
          key: 'stop',
          icon: <StopOutline size={16} />,
          disabled: busy,
          label: formatMessage('Vision.task.stop'),
          onClick: () => !busy && haltTask(record),
        },
      ];
    }
    return [
      details,
      {
        key: 'start',
        icon: <Run size={16} />,
        disabled: busy || !modelReady,
        label: formatMessage('Vision.task.run'),
        onClick: () => !busy && modelReady && runTask(record),
      },
      {
        key: 'edit',
        icon: <Edit size={16} />,
        label: formatMessage('common.edit'),
        onClick: () => setEditTask(record),
      },
      ...(regionMode !== 'none'
        ? [
            {
              key: 'regions',
              icon: <Route size={16} />,
              label:
                regionMode === 'line'
                  ? formatMessage('Vision.region.configLine')
                  : formatMessage('Vision.region.configZone'),
              onClick: () => setRegionTask(record),
            },
          ]
        : []),
      {
        key: 'test',
        icon: <ListChecked size={16} />,
        disabled: !modelReady,
        label: formatMessage('Vision.task.testFrame'),
        onClick: () => modelReady && setTestTask(record),
      },
      {
        key: 'delete',
        danger: true,
        icon: <TrashCan size={16} />,
        disabled: !idle,
        label: idle ? formatMessage('common.delete') : formatMessage('Vision.task.stopBeforeDelete'),
        onClick: () => idle && remove(record),
      },
    ];
  };

  const columns: TableColumnsType<VisionTask> = [
    {
      title: formatMessage('Vision.task.name'),
      dataIndex: 'name',
      key: 'name',
      width: 223,
      render: (value: string, record) => (
        <div className={styles.algoNameCell}>
          <Text ellipsis={{ tooltip: value }} className={styles.primaryText}>
            {value}
          </Text>
          <Text type="secondary" ellipsis={{ tooltip: record.unsTopic }} className={styles.algoCode}>
            {record.unsTopic}
          </Text>
        </div>
      ),
    },
    {
      title: formatMessage('Vision.task.cameras'),
      dataIndex: 'cameras',
      key: 'cameras',
      width: 160,
      render: (_, record) => {
        const cameras = record.cameras || [];
        if (cameras.length === 0) return '--';
        const shown = cameras.slice(0, MAX_VISIBLE_CAMERAS);
        const rest = cameras.length - shown.length;
        const allNamesTitle = (
          <div className={styles.taskCameraTooltip}>
            {cameras.map((cam) => (
              <div key={cam.id}>{cam.name}</div>
            ))}
          </div>
        );
        return (
          <span className={styles.taskCameraTags}>
            {shown.map((cam) => (
              <Tooltip key={cam.id} title={cam.name}>
                <span className={styles.taskCameraTag}>{cam.name}</span>
              </Tooltip>
            ))}
            {rest > 0 && (
              <Tooltip title={allNamesTitle}>
                <span className={`${styles.taskCameraTag} ${styles.taskCameraMore}`}>+{rest}</span>
              </Tooltip>
            )}
          </span>
        );
      },
    },
    {
      title: formatMessage('Vision.task.algorithm'),
      dataIndex: 'algorithm',
      key: 'algorithm',
      width: 148,
      render: (_, record) =>
        record.algorithmMissing ? (
          <span className={styles.algoMissing}>{formatMessage('Vision.task.algorithmMissing')}</span>
        ) : (
          <div className={styles.algoNameCell}>
            <Text ellipsis={{ tooltip: record.algorithm?.name }}>{record.algorithm?.name}</Text>
            <Text type="secondary" className={styles.algoCode}>
              {record.algorithm ? formatMessage(`Vision.algorithm.type.${record.algorithm.algoType}`) : ''}
            </Text>
          </div>
        ),
    },
    {
      title: formatMessage('Vision.task.status'),
      dataIndex: 'status',
      key: 'status',
      width: 98,
      render: (value: VisionTaskStatus) => <TaskStatusTag status={value} />,
    },
    {
      title: (
        <span className={styles.thWithTip}>
          {formatMessage('Vision.task.frames')}
          <Tooltip title={formatMessage('Vision.task.framesTip')}>
            <Help size={12} className={styles.thTipIcon} />
          </Tooltip>
        </span>
      ),
      dataIndex: 'processedFrames',
      key: 'processedFrames',
      width: 104,
      render: (value: number, record) => (ACTIVE_STATUSES.includes(record.status) ? value.toLocaleString() : '--'),
    },
    {
      title: (
        <span className={styles.thWithTip}>
          {formatMessage('Vision.task.targets')}
          <Tooltip title={formatMessage('Vision.task.targetsTip')}>
            <Help size={12} className={styles.thTipIcon} />
          </Tooltip>
        </span>
      ),
      dataIndex: 'currentTargets',
      key: 'currentTargets',
      width: 83,
      render: (value: number, record) =>
        ACTIVE_STATUSES.includes(record.status) ? (value ?? 0).toLocaleString() : '--',
    },
    {
      title: formatMessage('Vision.task.note'),
      dataIndex: 'note',
      key: 'note',
      width: 180,
      ellipsis: true,
      render: (value: string) => <Text ellipsis={{ tooltip: value || '-' }}>{value || '-'}</Text>,
    },
    {
      title: formatMessage('common.operation'),
      key: 'operation',
      width: 112,
      fixed: 'right',
      render: (_, record) => {
        if (!canManage) return null;
        return (
          <div className={styles.taskOps}>
            <Dropdown trigger={['click']} menu={{ items: operationItems(record) }} placement="bottomRight">
              <button type="button" className={styles.algoCardMore} aria-label={formatMessage('common.operation')}>
                <OverflowMenuHorizontal size={16} />
              </button>
            </Dropdown>
          </div>
        );
      },
    },
  ];

  const filterContent = (
    <Select
      allowClear
      value={status}
      style={{ width: 200 }}
      placeholder={formatMessage('Vision.task.status')}
      onChange={(value) => {
        setStatus(value);
        setSearchParams({ search, status: value });
      }}
      options={(Object.keys(TASK_STATUS_META) as VisionTaskStatus[]).map((value) => ({
        value,
        label: formatMessage(TASK_STATUS_META[value].labelKey),
      }))}
    />
  );

  return (
    <ComLayout loading={loading}>
      <div className={styles.algoPanel}>
        <div className={styles.algoHeader}>
          <div className={styles.algoTitle}>{formatMessage('Vision.task.title')}</div>
          {canManage && (
            <AuthButton
              auth={ButtonPermission['Vision.task.add']}
              type="primary"
              icon={<Add {...toolbarIconProps} />}
              onClick={() => setFormOpen(true)}
            >
              {formatMessage('Vision.task.add')}
            </AuthButton>
          )}
        </div>
        <div className={styles.algoToolbar}>
          <div className={styles.algoToolbarLeft}>
            <Popover content={filterContent} trigger="click" placement="bottomLeft">
              <button
                type="button"
                className={`${styles.iconButton} ${styles.filterButton} ${status ? styles.filterActive : ''}`}
                aria-label={formatMessage('Vision.task.status')}
              >
                <Filter size={16} />
              </button>
            </Popover>
            <ProSearch
              size="sm"
              value={search}
              placeholder={formatMessage('Vision.task.searchPlaceholder')}
              onChange={(event) => setSearch(event.target.value)}
              onSearch={(value) => {
                setSearch(value);
                setSearchParams({ search: value, status });
              }}
              className={styles.search}
            />
          </div>
        </div>
        <div className={styles.algoContent}>
          <ProTable
            className={styles.table}
            rowKey="id"
            columns={columns}
            dataSource={data}
            scroll={{ x: 1108, y: '100%' }}
            locale={{
              emptyText: (
                <div className={styles.tableEmpty}>
                  <ComEmptyState
                    className={styles.panelEmptyInner}
                    variant="inline"
                    title={formatMessage(search || status ? 'Vision.task.noMatchTitle' : 'Vision.task.emptyTitle')}
                    description={formatMessage(search || status ? 'Vision.task.noMatchDesc' : 'Vision.task.emptyDesc')}
                  />
                </div>
              ),
            }}
            pagination={{
              total: pagination.total,
              current: pagination.page,
              pageSize: pagination.pageSize,
              pageSizeOptions: pagination.pageSizes,
              showSizeChanger: true,
              onChange: (page, pageSize) => pagination.onChange({ page, pageSize }),
            }}
          />
        </div>
      </div>
      <TaskFormModal
        open={formOpen || Boolean(editTask)}
        editTask={editTask}
        onCancel={() => {
          setFormOpen(false);
          setEditTask(null);
        }}
        onSaved={() => {
          setFormOpen(false);
          setEditTask(null);
          refreshRequest();
        }}
      />
      <TaskTestFrameModal task={testTask} onClose={() => setTestTask(null)} />
      <TaskPreviewModal task={previewTask} onClose={() => setPreviewTask(null)} />
      <TaskDetailModal
        task={detailTask}
        onClose={() => setDetailTask(null)}
        onEdit={(task) => {
          setDetailTask(null);
          setEditTask(task);
        }}
      />
      <TaskRegionModal
        task={regionTask}
        onClose={() => setRegionTask(null)}
        onSaved={() => {
          setRegionTask(null);
          refreshRequest();
        }}
      />
    </ComLayout>
  );
};

export default AnalysisTasksPanel;
