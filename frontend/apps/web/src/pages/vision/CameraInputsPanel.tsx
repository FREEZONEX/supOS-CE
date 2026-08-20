import { useState } from 'react';
import { App, Pagination, Popover, Select, Tag, Tooltip, Typography, type TableColumnsType } from 'antd';
import {
  listVideoCameras,
  deleteVideoCamera,
  reconnectVideoCamera,
  type CameraLiveStatus,
  type VideoCamera,
} from '@/apis/core-api/video';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton, AuthWrapper } from '@/components/auth';
import { ComEmptyState, ComLayout, ProSearch, ProTable } from '@/components';
import { Add, Edit, Filter, Renew, TrashCan, View } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { usePagination, useTranslate } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import CameraFormModal from './components/CameraFormModal';
import CameraPreviewModal from './components/CameraPreviewModal';
import styles from './index.module.scss';

const { Text } = Typography;

const LIVE_STATUS_META: Record<CameraLiveStatus, { color: string; labelKey: string }> = {
  online: { color: 'success', labelKey: 'Vision.camera.liveOnline' },
  offline: { color: 'error', labelKey: 'Vision.camera.liveOffline' },
  checking: { color: 'processing', labelKey: 'Vision.camera.liveChecking' },
  unknown: { color: 'default', labelKey: 'Vision.camera.liveUnknown' },
};

const CameraInputsPanel = () => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [search, setSearch] = useState('');
  const [status, setStatus] = useState<number>();
  const [transport, setTransport] = useState<string>();
  const [formOpen, setFormOpen] = useState(false);
  const [editingCamera, setEditingCamera] = useState<VideoCamera | null>(null);
  const [previewCamera, setPreviewCamera] = useState<VideoCamera | null>(null);
  const [reconnectingId, setReconnectingId] = useState<number | null>(null);
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const resourceKeys = (currentUserInfo as (typeof currentUserInfo & { resourceKeys?: string[] }) | undefined)
    ?.resourceKeys;
  const canManage =
    systemInfo?.authEnable === false ||
    currentUserInfo?.superAdmin === true ||
    resourceKeys?.includes('vision.camera.manage');
  const { loading, data, pagination, setSearchParams, refreshRequest } = usePagination<VideoCamera>({
    fetchApi: listVideoCameras,
  });

  const openCreate = () => {
    setEditingCamera(null);
    setFormOpen(true);
  };

  const openEdit = (camera: VideoCamera) => {
    setEditingCamera(camera);
    setFormOpen(true);
  };

  const reconnect = async (camera: VideoCamera) => {
    setReconnectingId(camera.id);
    try {
      const updated = await reconnectVideoCamera(camera.id);
      if (updated?.liveStatus === 'online') {
        message.success(formatMessage('Vision.camera.reconnectOnline'));
      } else {
        message.warning(formatMessage('Vision.camera.reconnectOffline'));
      }
      refreshRequest();
    } finally {
      setReconnectingId(null);
    }
  };

  const remove = (camera: VideoCamera) => {
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          ...createDeleteConfirmOptions({
            title: formatMessage('Vision.camera.deleteTitle'),
            name: camera.name,
            formatMessage,
          }),
          onOk: async () => {
            await deleteVideoCamera(camera.id);
            message.success(formatMessage('common.deleteSuccessfully'));
            refreshRequest();
          },
        },
        formatMessage
      )
    );
  };

  const columns: TableColumnsType<VideoCamera> = [
    {
      title: formatMessage('Vision.camera.name'),
      dataIndex: 'name',
      key: 'name',
      width: 180,
      ellipsis: true,
      render: (value: string) => (
        <Text ellipsis={{ tooltip: value }} className={styles.primaryText}>
          {value}
        </Text>
      ),
    },
    {
      title: formatMessage('Vision.camera.code'),
      dataIndex: 'cameraCode',
      key: 'cameraCode',
      width: 160,
      ellipsis: true,
      render: (value: string) => <Text ellipsis={{ tooltip: value }}>{value}</Text>,
    },
    {
      title: formatMessage('Vision.camera.location'),
      dataIndex: 'location',
      key: 'location',
      width: 160,
      ellipsis: true,
      render: (value: string) => <Text ellipsis={{ tooltip: value || '-' }}>{value || '-'}</Text>,
    },
    {
      title: formatMessage('Vision.camera.liveStatus'),
      dataIndex: 'liveStatus',
      key: 'liveStatus',
      width: 108,
      render: (value: CameraLiveStatus) => {
        const meta = LIVE_STATUS_META[value] || LIVE_STATUS_META.unknown;
        return (
          <Tag color={meta.color} bordered={false}>
            {formatMessage(meta.labelKey)}
          </Tag>
        );
      },
    },
    {
      title: formatMessage('Vision.camera.protocol'),
      dataIndex: 'rtpType',
      key: 'rtpType',
      width: 96,
      render: (value: string, record) => (record.transport === 'rtsp' ? value?.toUpperCase() || '-' : '-'),
    },
    {
      title: formatMessage('common.description'),
      dataIndex: 'description',
      key: 'description',
      width: 220,
      ellipsis: true,
      render: (value: string) => <Text ellipsis={{ tooltip: value || '-' }}>{value || '-'}</Text>,
    },
    {
      title: formatMessage('common.operation'),
      key: 'operation',
      width: 154,
      fixed: 'right',
      render: (_, record) => {
        const canPreview = record.liveStatus === 'online';
        const canReconnect = record.liveStatus === 'offline' || record.liveStatus === 'unknown';
        return (
          <span className={`${styles.operationCell} ${styles.cameraOperationCell}`}>
            {canPreview ? (
              <Tooltip title={formatMessage('Vision.camera.preview')}>
                <button
                  type="button"
                  className={styles.iconButton}
                  aria-label={formatMessage('Vision.camera.preview')}
                  onClick={() => setPreviewCamera(record)}
                >
                  <View size={16} />
                </button>
              </Tooltip>
            ) : (
              <Tooltip title={formatMessage('Vision.camera.previewOfflineHint')}>
                <span className={styles.disabledOperation} aria-label={formatMessage('Vision.camera.preview')}>
                  <View size={16} />
                </span>
              </Tooltip>
            )}
            {canManage && canReconnect ? (
              <Tooltip title={formatMessage('Vision.camera.reconnect')}>
                <span className={styles.operationSlot}>
                  <AuthWrapper auth={ButtonPermission['Vision.camera.reconnect']}>
                    <button
                      type="button"
                      className={styles.iconButton}
                      aria-label={formatMessage('Vision.camera.reconnect')}
                      disabled={reconnectingId === record.id}
                      onClick={() => void reconnect(record)}
                    >
                      <Renew size={16} />
                    </button>
                  </AuthWrapper>
                </span>
              </Tooltip>
            ) : (
              <Tooltip title={formatMessage('Vision.camera.reconnectOnlineHint')}>
                <span className={styles.disabledOperation} aria-label={formatMessage('Vision.camera.reconnect')}>
                  <Renew size={16} />
                </span>
              </Tooltip>
            )}
            {canManage && (
              <Tooltip title={formatMessage('common.edit')}>
                <span className={styles.operationSlot}>
                  <AuthWrapper auth={ButtonPermission['Vision.camera.edit']}>
                    <button
                      type="button"
                      className={styles.iconButton}
                      aria-label={formatMessage('common.edit')}
                      onClick={() => openEdit(record)}
                    >
                      <Edit size={16} />
                    </button>
                  </AuthWrapper>
                </span>
              </Tooltip>
            )}
            {canManage && (
              <Tooltip title={formatMessage('common.delete')}>
                <span className={styles.operationSlot}>
                  <AuthWrapper auth={ButtonPermission['Vision.camera.delete']}>
                    <button
                      type="button"
                      className={`${styles.iconButton} ${styles.danger}`}
                      aria-label={formatMessage('common.delete')}
                      onClick={() => remove(record)}
                    >
                      <TrashCan size={16} />
                    </button>
                  </AuthWrapper>
                </span>
              </Tooltip>
            )}
          </span>
        );
      },
    },
  ];

  const hasFilter = Boolean(search || status || transport);

  // 漏斗筛选面板:启用状态 + 接入类型(设计里工具栏只保留漏斗与搜索)。
  const filterContent = (
    <div className={styles.cameraFilterPanel}>
      <Select
        allowClear
        value={status}
        placeholder={formatMessage('Vision.camera.enabled')}
        style={{ width: 200 }}
        onChange={(value) => {
          setStatus(value);
          setSearchParams({ search, status: value, transport });
        }}
        options={[
          { value: 1, label: formatMessage('common.enable') },
          { value: 2, label: formatMessage('common.disable') },
        ]}
      />
      <Select
        allowClear
        value={transport}
        placeholder={formatMessage('Vision.camera.inputType')}
        style={{ width: 200 }}
        onChange={(value) => {
          setTransport(value);
          setSearchParams({ search, status, transport: value });
        }}
        options={[
          { value: 'rtsp', label: 'RTSP' },
          { value: 'rtmp_pull', label: formatMessage('Vision.camera.rtmpPull') },
        ]}
      />
    </div>
  );

  return (
    <ComLayout loading={loading}>
      <div className={styles.algoPanel}>
        <div className={styles.algoHeader}>
          <div className={styles.algoTitle}>{formatMessage('Vision.camera.title')}</div>
          {canManage && (
            <AuthButton
              auth={ButtonPermission['Vision.camera.add']}
              type="primary"
              icon={<Add {...toolbarIconProps} />}
              onClick={openCreate}
            >
              {formatMessage('Vision.camera.add')}
            </AuthButton>
          )}
        </div>
        <div className={styles.algoToolbar}>
          <div className={styles.algoToolbarLeft}>
            <Popover content={filterContent} trigger="click" placement="bottomLeft">
              <button
                type="button"
                className={`${styles.iconButton} ${styles.filterButton} ${hasFilter ? styles.filterActive : ''}`}
                aria-label={formatMessage('Vision.camera.inputType')}
              >
                <Filter size={16} />
              </button>
            </Popover>
            <ProSearch
              size="sm"
              value={search}
              placeholder={formatMessage('Vision.camera.searchPlaceholder')}
              onChange={(event) => setSearch(event.target.value)}
              onSearch={(value) => {
                setSearch(value);
                setSearchParams({ search: value, status, transport });
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
            scroll={{ x: 1078, y: '100%' }}
            locale={{
              emptyText: (
                <div className={styles.tableEmpty}>
                  <ComEmptyState
                    className={styles.panelEmptyInner}
                    variant="inline"
                    title={formatMessage(
                      hasFilter ? 'Vision.camera.noMatchTitle' : 'Vision.camera.emptyTitle'
                    )}
                    description={formatMessage(
                      hasFilter ? 'Vision.camera.noMatchDesc' : 'Vision.camera.emptyDesc'
                    )}
                  />
                </div>
              ),
            }}
            pagination={false}
          />
        </div>
        <div className={styles.algoPager}>
          <Pagination
            total={pagination.total}
            current={pagination.page}
            pageSize={pagination.pageSize}
            pageSizeOptions={pagination.pageSizes}
            showSizeChanger
            onChange={(page, pageSize) => pagination.onChange({ page, pageSize })}
          />
        </div>
      </div>
      <CameraFormModal
        open={formOpen}
        camera={editingCamera}
        onCancel={() => setFormOpen(false)}
        onSaved={() => {
          setFormOpen(false);
          refreshRequest();
        }}
      />
      <CameraPreviewModal open={Boolean(previewCamera)} camera={previewCamera} onClose={() => setPreviewCamera(null)} />
    </ComLayout>
  );
};

export default CameraInputsPanel;
