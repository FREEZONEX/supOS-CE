import { Add, CheckmarkOutline, ChevronDown, Upload } from '@carbon/icons-react';
import { Box, Edit, Launch, MisuseOutline, Renew, TrashCan } from '@/components/lucide-icon/carbon';
import { App, Dropdown, Flex, Spin, Tooltip } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';

import ProTable from '@/components/pro-table';
import { AuthWrapper } from '@/components/auth';
import { ButtonPermission } from '@/common-types/button-permission';
import { useTranslate } from '@/hooks';
import { formatTimestamp } from '@/utils/format';

import { deleteProjectApp, getProjectApps } from '@/apis/core-api';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import type { AppItem } from '@/pages/project/types';
import { useLocationNavigate } from '@/routers';
import { normalizeLocalAppUrl } from '@/utils/app-url';
import { getToken } from '@/utils/auth';
import { createDeleteConfirmOptions } from '@/utils/modal-confirm';
import EditAppModal, { type EditAppModalRef } from '../../EditAppModal';
import ImportAppModal from '../../ImportAppModal';
import ManualAppModal from '../../ManualAppModal';
import ReplaceAppModal, { type ReplaceAppModalRef } from '../../ReplaceAppModal';
import TabLayout from '../../TabLayout';
import styles from './index.module.scss';

const buildAppOpenUrl = (record: AppItem) => {
  const normalizedUrl = normalizeLocalAppUrl(record.url || '', {
    allowPrivateHostRewrite: record.siteType === 'dynamic' && !record.manual && record.appType !== 'manual',
  });
  const token = getToken();
  if (!token || !normalizedUrl) {
    return normalizedUrl;
  }
  const separator = normalizedUrl.includes('?') ? '&' : '?';
  return `${normalizedUrl}${separator}_token=${encodeURIComponent(token)}`;
};

interface AppsProps {
  projectId?: string;
}

export enum AppStatus {
  Active = 'Active',
  Disabled = 'Disabled',
  Deploying = 'Deploying',
  Failed = 'Failed',
  Stopped = 'Stopped',
}

const operationTooltipClassNames = {
  root: styles.operationTooltip,
};

const Apps: FC<AppsProps> = ({ projectId }) => {
  const formatMessage = useTranslate();
  const navigate = useLocationNavigate();
  const { message, modal } = App.useApp();
  const [apps, setApps] = useState<AppItem[]>();
  const [loading, setLoading] = useState(false);
  const appsRequestIdRef = useRef(0);
  const editAppModalRef = useRef<EditAppModalRef>(null);
  const replaceAppModalRef = useRef<ReplaceAppModalRef>(null);
  const importAppModalRef = useRef<any>(null);
  const manualAppModalRef = useRef<any>(null);

  const getApps = useCallback(async (id: string) => {
    const requestId = ++appsRequestIdRef.current;
    setLoading(true);
    try {
      const data = await getProjectApps(id);
      if (requestId === appsRequestIdRef.current) {
        setApps(data);
      }
    } catch {
      if (requestId === appsRequestIdRef.current) {
        setApps([]);
      }
    } finally {
      if (requestId === appsRequestIdRef.current) {
        setLoading(false);
      }
    }
  }, []);

  const refreshRequest = useCallback(() => {
    if (projectId) {
      return getApps(projectId);
    }
  }, [projectId, getApps]);

  const handleDelete = useCallback(
    (record: AppItem) => {
      if (projectId) {
        return deleteProjectApp(projectId, `${record.appId}`).then(() => {
          message.success(formatMessage('common.deleteSuccessfully'));
          getApps(projectId);
        });
      }
    },
    [formatMessage, message, projectId, getApps]
  );

  const openApp = useCallback(
    (record: AppItem) => {
      if (!record.url) {
        return;
      }
      const isManual = record.manual || record.siteType === 'manualExternal' || record.appType === 'manual';
      const appName = record.displayName || record.name;
      if (isManual && record.openInPlatform !== false && projectId) {
        navigate({
          pathname: `/launchpad/${encodeURIComponent(String(projectId))}/${encodeURIComponent(String(record.appId))}`,
          state: {
            tabName: appName,
          },
        });
        return;
      }
      window.open(buildAppOpenUrl(record), '_blank', 'noopener,noreferrer');
    },
    [navigate, projectId]
  );

  useEffect(() => {
    if (!projectId) {
      return;
    }
    const requestTimer = window.setTimeout(() => void getApps(projectId), 0);
    return () => window.clearTimeout(requestTimer);
  }, [projectId, getApps]);

  useActivate(() => {
    if (!projectId) {
      setApps([]);
      return;
    }
    getApps(projectId);
  });

  const columns: any = useMemo(() => {
    return [
      {
        title: () => formatMessage('common.name'),
        dataIndex: 'name',
        key: 'name',
        width: 220,
        render: (value: string, record: AppItem) => (
          <Flex align="center" gap={8} className={styles.nameCell} title={value}>
            <span className={styles.nameIcon}>
              {record.iconUrl ? (
                <img style={{ maxWidth: '100%', maxHeight: '100%' }} src={record.iconUrl} />
              ) : (
                <Box size={14} strokeWidth={1.75} />
              )}
            </span>
            <span className={styles.nameText}>{value}</span>
          </Flex>
        ),
      },
      {
        title: () => formatMessage('common.description'),
        dataIndex: 'description',
        key: 'description',
        width: 260,
        render: (value?: string) => (
          <span className={styles.descriptionCell} title={value || '-'}>
            {value || '-'}
          </span>
        ),
      },
      {
        title: () => formatMessage('common.version'),
        dataIndex: 'version',
        key: 'version',
        width: 100,
      },
      {
        title: () => formatMessage('common.status'),
        dataIndex: 'status',
        key: 'status',
        width: 120,
        render: (status: AppItem['status']) => {
          const statusConfig: Record<string, { label: string; className: string; Icon?: any; spinner?: boolean }> = {
            [AppStatus.Active]: {
              label: formatMessage('common.active'),
              className: styles.statusActive,
              Icon: CheckmarkOutline,
            },
            [AppStatus.Deploying]: {
              label: formatMessage('common.deploying'),
              className: styles.statusDeploying,
              spinner: true,
            },
            [AppStatus.Failed]: {
              label: formatMessage('common.failed'),
              className: styles.statusFailed,
              Icon: MisuseOutline,
            },
            [AppStatus.Stopped]: {
              label: formatMessage('common.stopped'),
              className: styles.statusInactive,
              Icon: MisuseOutline,
            },
            [AppStatus.Disabled]: {
              label: formatMessage('common.stopped'),
              className: styles.statusInactive,
              Icon: MisuseOutline,
            },
          };

          const config = statusConfig[status] || statusConfig[AppStatus.Disabled];
          const { label, className, Icon, spinner } = config;

          return (
            <span className={`${styles.tag} ${className}`}>
              {spinner ? <Spin size="small" /> : Icon ? <Icon size={12} /> : null}
              <span className={styles.tagText}>{label}</span>
            </span>
          );
        },
      },
      {
        title: () => formatMessage('project.accessRoles'),
        dataIndex: 'accessibleRoles',
        key: 'accessibleRoles',
        width: 180,
        render: (roles: AppItem['accessibleRoles']) => (
          <div className={styles.roleTagWrap}>
            {roles?.map((role) => (
              <span key={role} className={`${styles.tag} ${styles.roleTag}`}>
                <span className={styles.tagText} title={role}>
                  {role}
                </span>
              </span>
            ))}
          </div>
        ),
      },
      {
        title: () => formatMessage('project.lastUpdated'),
        dataIndex: 'updatedAt',
        key: 'updatedAt',
        width: 180,
        render: (value: number) => formatTimestamp(value, 'YYYY/MM/DD HH:mm', true),
      },
      {
        title: () => formatMessage('common.operation'),
        key: 'operation',
        fixed: 'right',
        width: 176,
        render: (_: any, record: AppItem) => {
          const isDeploying = record.status === AppStatus.Deploying;
          const isActive = record.status === AppStatus.Active;
          const isFailed = record.status === AppStatus.Failed || record.deployStatus === 'failed';

          return (
            <span className={styles['operation']}>
              {isActive && record.url ? (
                <Tooltip title={formatMessage('project.openApp')} classNames={operationTooltipClassNames}>
                  <button
                    type="button"
                    className={styles.operationButton}
                    aria-label={formatMessage('project.openApp')}
                    onClick={() => openApp(record)}
                  >
                    <Launch size={16} strokeWidth={1.75} />
                  </button>
                </Tooltip>
              ) : isDeploying ? (
                <Tooltip title={formatMessage('project.appDeploying')} classNames={operationTooltipClassNames}>
                  <span className={styles.operationButton}>
                    <Spin size="small" />
                  </span>
                </Tooltip>
              ) : isFailed ? (
                <Tooltip
                  title={record.deployError || formatMessage('project.appDeployFailed')}
                  classNames={operationTooltipClassNames}
                >
                  <span className={`${styles.operationButton} ${styles.operationError}`}>
                    <MisuseOutline size={16} strokeWidth={1.75} />
                  </span>
                </Tooltip>
              ) : null}
              {!isDeploying && (
                <>
                  <AuthWrapper auth={ButtonPermission['Project.edit']}>
                    <Tooltip
                      title={formatMessage('project.appSetting', {}, 'Edit App')}
                      classNames={operationTooltipClassNames}
                    >
                      <button
                        type="button"
                        className={styles.operationButton}
                        aria-label={formatMessage('project.appSetting', {}, 'Edit App')}
                        onClick={() => editAppModalRef.current?.onOpen(record)}
                      >
                        <Edit size={16} strokeWidth={1.75} />
                      </button>
                    </Tooltip>
                  </AuthWrapper>
                  <AuthWrapper auth={ButtonPermission['Project.edit']}>
                    <Tooltip
                      title={formatMessage('project.replace.action', {}, 'Replace App')}
                      classNames={operationTooltipClassNames}
                    >
                      <button
                        type="button"
                        className={styles.operationButton}
                        aria-label={formatMessage('project.replace.action', {}, 'Replace App')}
                        onClick={() => replaceAppModalRef.current?.onOpen(record)}
                      >
                        <Renew size={16} strokeWidth={1.75} />
                      </button>
                    </Tooltip>
                  </AuthWrapper>
                  <Tooltip title={formatMessage('common.delete')} classNames={operationTooltipClassNames}>
                    <button
                      type="button"
                      className={`${styles.operationButton} ${styles.operationButtonDanger}`}
                      aria-label={formatMessage('common.delete')}
                      onClick={() => {
                        modal.confirm({
                          ...createDeleteConfirmOptions({
                            title: formatMessage('project.confirmDeleteAppTitle'),
                            content: formatMessage('project.confirmDeleteAppDesc'),
                            okText: formatMessage('common.delete'),
                            cancelText: formatMessage('common.cancel'),
                          }),
                          onOk: () => handleDelete(record),
                        });
                      }}
                    >
                      <TrashCan size={16} strokeWidth={1.75} />
                    </button>
                  </Tooltip>
                </>
              )}
            </span>
          );
        },
      },
    ];
  }, [formatMessage, handleDelete, modal, openApp]);

  return (
    <TabLayout
      style={{ padding: '0 24px' }}
      toolbarActions={
        <Dropdown.Button
          type="primary"
          trigger={['click']}
          icon={<ChevronDown size={14} />}
          menu={{
            items: [
              {
                key: 'import',
                icon: <Upload size={14} />,
                label: formatMessage('project.importApp'),
              },
              {
                key: 'manual',
                icon: <Add size={14} />,
                label: formatMessage('project.addApp'),
              },
            ],
            onClick: ({ key }) => {
              if (key === 'manual') {
                manualAppModalRef.current?.onOpen?.();
                return;
              }
              importAppModalRef.current?.onOpen?.();
            },
          }}
          onClick={() => importAppModalRef.current?.onOpen?.()}
        >
          <Flex align="center" gap={6}>
            <Upload size={14} />
            {formatMessage('project.importApp')}
          </Flex>
        </Dropdown.Button>
      }
    >
      <ProTable
        loading={loading}
        scroll={{ x: 'max-content', y: 'calc(100vh - 300px)' }}
        columns={columns}
        rowKey="appId"
        dataSource={apps}
        resizeable
        pagination={false}
      />
      <EditAppModal ref={editAppModalRef} refreshRequest={refreshRequest} />
      <ReplaceAppModal ref={replaceAppModalRef} refreshRequest={refreshRequest} />
      <ImportAppModal ref={importAppModalRef} refreshRequest={refreshRequest} existingApps={apps} />
      <ManualAppModal ref={manualAppModalRef} refreshRequest={refreshRequest} />
    </TabLayout>
  );
};

export default Apps;
