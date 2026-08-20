import { useState } from 'react';
import {
  App,
  Dropdown,
  Pagination,
  Popover,
  Select,
  Tooltip,
  Typography,
  type MenuProps,
  type TableColumnsType,
} from 'antd';
import {
  listVisionAlgorithms,
  deleteVisionAlgorithm,
  type AlgorithmSource,
  type VisionAlgorithm,
} from '@/apis/core-api/algorithm';
import { ButtonPermission } from '@/common-types/button-permission';
import { AuthButton } from '@/components/auth';
import { ComEmptyState, ComLayout, ProSearch, ProTable } from '@/components';
import ViewModeSegmented from '@/components/lucide-icon/ViewModeSegmented';
import { Add, Filter, OverflowMenuHorizontal, TrashCan, Upload, View } from '@/components/lucide-icon/carbon';
import { toolbarIconProps } from '@/components/lucide-icon/icon-props';
import { usePagination, useTranslate, useViewModeStorage, VIEW_MODE_CARD, VIEW_MODE_STORAGE_KEYS } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { mergeDeleteConfirmProps } from '@/utils/delete-confirm-modal';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import AlgorithmCard from './AlgorithmCard';
import AlgorithmDetailModal from './AlgorithmDetailModal';
import AlgorithmFormModal from './AlgorithmFormModal';
import AlgorithmImportModal from './AlgorithmImportModal';
import AlgorithmUpdateModelModal from './AlgorithmUpdateModelModal';
import ModelStatusTag from './ModelStatusTag';
import SourceTag from './SourceTag';
import { fitAlgorithmLabels } from './algorithm-meta';
import styles from './index.module.scss';

const { Text } = Typography;

// 自定义算法上传暂不开放,现阶段只支持内置算法。放开时把这个开关连同判断一起删掉,
// 导入弹窗和后端接口都还在,不需要重新接线。
const ENABLE_CUSTOM_ALGORITHM = false;

const AlgorithmsPanel = () => {
  const formatMessage = useTranslate();
  const { modal, message } = App.useApp();
  const [search, setSearch] = useState('');
  const [source, setSource] = useState<string>();
  const [detailAlgorithm, setDetailAlgorithm] = useState<VisionAlgorithm | null>(null);
  const [editAlgorithm, setEditAlgorithm] = useState<VisionAlgorithm | null>(null);
  const [updateModelAlgorithm, setUpdateModelAlgorithm] = useState<VisionAlgorithm | null>(null);
  const [importOpen, setImportOpen] = useState(false);
  const currentUserInfo = useBaseStore((state) => state.currentUserInfo);
  const systemInfo = useBaseStore((state) => state.systemInfo);
  const resourceKeys = (currentUserInfo as (typeof currentUserInfo & { resourceKeys?: string[] }) | undefined)
    ?.resourceKeys;
  const canManage =
    systemInfo?.authEnable === false ||
    currentUserInfo?.superAdmin === true ||
    resourceKeys?.includes('vision.algorithm.manage');
  const { loading, data, pagination, setSearchParams, refreshRequest } = usePagination<VisionAlgorithm>({
    fetchApi: listVisionAlgorithms,
  });
  const [viewMode, setViewMode] = useViewModeStorage(VIEW_MODE_STORAGE_KEYS.visionAlgorithm);
  const isCardView = viewMode === VIEW_MODE_CARD;

  const remove = (algorithm: VisionAlgorithm) => {
    modal.confirm(
      mergeDeleteConfirmProps(
        {
          ...createDeleteConfirmOptions({
            title: formatMessage('Vision.algorithm.deleteTitle'),
            name: algorithm.name,
            formatMessage,
          }),
          onOk: async () => {
            await deleteVisionAlgorithm(algorithm.id);
            message.success(formatMessage('common.deleteSuccessfully'));
            refreshRequest();
          },
        },
        formatMessage
      )
    );
  };

  // 表格和卡片共用同一份操作项:内置算法不允许删除,置灰并说明原因。
  const operationItems = (record: VisionAlgorithm): MenuProps['items'] => [
    {
      key: 'detail',
      icon: <View size={16} />,
      label: formatMessage('Vision.algorithm.details'),
      onClick: () => setDetailAlgorithm(record),
    },
    ...(canManage
      ? [
          {
            key: 'uploadModel',
            icon: <Upload size={16} />,
            label: formatMessage('Vision.algorithm.uploadModel'),
            onClick: () => setUpdateModelAlgorithm(record),
          },
          {
            key: 'delete',
            danger: true,
            disabled: record.source !== 'custom',
            icon: <TrashCan size={16} />,
            label:
              record.source === 'custom' ? (
                formatMessage('common.delete')
              ) : (
                <Tooltip title={formatMessage('Vision.algorithm.builtinNoDelete')}>
                  {formatMessage('common.delete')}
                </Tooltip>
              ),
            onClick: () => remove(record),
          },
        ]
      : []),
  ];

  const columns: TableColumnsType<VisionAlgorithm> = [
    {
      title: formatMessage('Vision.algorithm.name'),
      dataIndex: 'name',
      key: 'name',
      width: 220,
      render: (value: string, record) => (
        <div className={styles.algoNameCell}>
          <Text ellipsis={{ tooltip: value }} className={styles.primaryText}>
            {value}
          </Text>
          <Text type="secondary" ellipsis={{ tooltip: record.algorithmCode }} className={styles.algoCode}>
            {record.algorithmCode}
          </Text>
        </div>
      ),
    },
    {
      title: formatMessage('Vision.algorithm.source'),
      dataIndex: 'source',
      key: 'source',
      width: 110,
      render: (value: AlgorithmSource) => <SourceTag source={value} />,
    },
    {
      title: formatMessage('common.description'),
      dataIndex: 'description',
      key: 'description',
      width: 260,
      ellipsis: true,
      render: (value: string) => <Text ellipsis={{ tooltip: value || '-' }}>{value || '-'}</Text>,
    },
    {
      title: formatMessage('Vision.algorithm.labels'),
      dataIndex: 'labels',
      key: 'labels',
      width: 160,
      render: (value: string[]) => {
        const labels = value || [];
        if (labels.length === 0) return '--';
        const { shown, rest, all } = fitAlgorithmLabels(labels, 148);
        return (
          <span className={styles.algoTableLabelRow}>
            {shown.map((label) => (
              <span key={label} className={styles.algoCardLabelTag}>
                {label}
              </span>
            ))}
            {rest > 0 && (
              <Tooltip title={all.join(', ')}>
                <span className={styles.algoCardLabelTag}>+{rest}</span>
              </Tooltip>
            )}
          </span>
        );
      },
    },
    {
      title: formatMessage('Vision.algorithm.modelStatus'),
      dataIndex: 'modelStatus',
      key: 'modelStatus',
      width: 130,
      render: (_, record) => <ModelStatusTag status={record.modelStatus} />,
    },
    {
      title: formatMessage('common.operation'),
      key: 'operation',
      width: 90,
      fixed: 'right',
      render: (_, record) => (
        <span className={styles.operationCell}>
          <Dropdown trigger={['click']} menu={{ items: operationItems(record) }} placement="bottomRight">
            <button type="button" className={styles.iconButton} aria-label={formatMessage('common.operation')}>
              <OverflowMenuHorizontal size={16} />
            </button>
          </Dropdown>
        </span>
      ),
    },
  ];

  const filterContent = (
    <Select
      allowClear
      value={source}
      style={{ width: 200 }}
      placeholder={formatMessage('Vision.algorithm.source')}
      onChange={(value) => {
        setSource(value);
        setSearchParams({ search, source: value });
      }}
      options={[
        { value: 'builtin', label: formatMessage('Vision.algorithm.sourceBuiltin') },
        { value: 'custom', label: formatMessage('Vision.algorithm.sourceCustom') },
      ]}
    />
  );

  return (
    <ComLayout loading={loading}>
      <div className={styles.algoPanel}>
        <div className={styles.algoHeader}>
          <div className={styles.algoTitle}>{formatMessage('Vision.algorithm.title')}</div>
          {ENABLE_CUSTOM_ALGORITHM && canManage && (
            <AuthButton
              auth={ButtonPermission['Vision.algorithm.add']}
              type="primary"
              icon={<Add {...toolbarIconProps} />}
              onClick={() => setImportOpen(true)}
            >
              {formatMessage('Vision.algorithm.add')}
            </AuthButton>
          )}
        </div>
        <div className={styles.algoToolbar}>
          <div className={styles.algoToolbarLeft}>
            <Popover content={filterContent} trigger="click" placement="bottomLeft">
              <button
                type="button"
                className={`${styles.iconButton} ${styles.filterButton} ${source ? styles.filterActive : ''}`}
                aria-label={formatMessage('Vision.algorithm.source')}
              >
                <Filter size={16} />
              </button>
            </Popover>
            <ProSearch
              size="sm"
              value={search}
              placeholder={formatMessage('Vision.algorithm.searchPlaceholder')}
              onChange={(event) => setSearch(event.target.value)}
              onSearch={(value) => {
                setSearch(value);
                setSearchParams({ search: value, source });
              }}
              className={styles.search}
            />
          </div>
          <ViewModeSegmented
            value={viewMode}
            onChange={setViewMode}
            cardTitle={formatMessage('common.cardMode')}
            listTitle={formatMessage('common.listMode')}
          />
        </div>
        <div className={styles.algoContent}>
          {isCardView ? (
            data.length === 0 && !loading ? (
              <div className={styles.panelEmpty}>
                <ComEmptyState
                  className={styles.panelEmptyInner}
                  variant="inline"
                  title={formatMessage(
                    search || source ? 'Vision.algorithm.noMatchTitle' : 'Vision.algorithm.emptyTitle'
                  )}
                  description={formatMessage(
                    search || source ? 'Vision.algorithm.noMatchDesc' : 'Vision.algorithm.emptyDesc'
                  )}
                />
              </div>
            ) : (
              <div className={styles.algoCardGrid}>
                {data.map((algorithm) => (
                  <AlgorithmCard
                    key={algorithm.id}
                    algorithm={algorithm}
                    canManage={canManage}
                    onDetail={setDetailAlgorithm}
                    onUpdateModel={setUpdateModelAlgorithm}
                    onDelete={remove}
                  />
                ))}
              </div>
            )
          ) : (
            <ProTable
              className={styles.table}
              rowKey="id"
              columns={columns}
              dataSource={data}
              scroll={{ x: 970, y: '100%' }}
              locale={{
                emptyText: (
                  <div className={styles.tableEmpty}>
                    <ComEmptyState
                      className={styles.panelEmptyInner}
                      variant="inline"
                      title={formatMessage(
                        search || source ? 'Vision.algorithm.noMatchTitle' : 'Vision.algorithm.emptyTitle'
                      )}
                      description={formatMessage(
                        search || source ? 'Vision.algorithm.noMatchDesc' : 'Vision.algorithm.emptyDesc'
                      )}
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
          )}
        </div>
        {isCardView && (
          <div className={styles.algoPager}>
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
        )}
      </div>
      <AlgorithmDetailModal
        open={Boolean(detailAlgorithm)}
        algorithm={detailAlgorithm}
        canManage={canManage}
        onEdit={setEditAlgorithm}
        onClose={() => setDetailAlgorithm(null)}
      />
      <AlgorithmImportModal
        open={importOpen}
        onCancel={() => setImportOpen(false)}
        onSaved={() => {
          setImportOpen(false);
          refreshRequest();
        }}
      />
      <AlgorithmFormModal
        open={Boolean(editAlgorithm)}
        algorithm={editAlgorithm}
        onCancel={() => setEditAlgorithm(null)}
        onSaved={() => {
          setEditAlgorithm(null);
          refreshRequest();
        }}
      />
      <AlgorithmUpdateModelModal
        open={Boolean(updateModelAlgorithm)}
        algorithm={updateModelAlgorithm}
        onClose={() => setUpdateModelAlgorithm(null)}
        onSaved={() => {
          setUpdateModelAlgorithm(null);
          refreshRequest();
        }}
      />
    </ComLayout>
  );
};

export default AlgorithmsPanel;
