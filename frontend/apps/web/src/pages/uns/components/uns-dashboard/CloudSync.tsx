import { disconnectCloudSync, getCloudSyncConfig, type CloudSyncConfigResp } from '@/apis/core-api/cloudsync.ts';
import useTranslate from '@/hooks/useTranslate.ts';
import HelpTooltip from '@/components/help-tooltip';
import { App, Button, Flex, Popconfirm, Spin, Tag, Typography } from 'antd';
import { CloudApp, ConnectionSignal, ConnectionSignalOff } from '@carbon/icons-react';
import { useCallback, useEffect, useMemo, useState } from 'react';
import CloudSyncConnectModal from './CloudSyncConnectModal.tsx';
import styles from './cloudsync.module.scss';

const statusClassMap: Record<string, string> = {
  connected: 'is-connected',
  connecting: 'is-connecting',
  disconnected: 'is-disconnected',
  error: 'is-error',
};

const CloudSync = () => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [loading, setLoading] = useState(true);
  const [disconnecting, setDisconnecting] = useState(false);
  const [openModal, setOpenModal] = useState(false);
  const [config, setConfig] = useState<CloudSyncConfigResp | null>(null);

  const loadConfig = useCallback(() => {
    setLoading(true);
    getCloudSyncConfig()
      .then((data) => {
        setConfig(data);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const statusText = useMemo(() => {
    switch (config?.connectStatus) {
      case 'connected':
        return formatMessage('uns.cloudSyncStatusConnected');
      case 'connecting':
        return formatMessage('uns.cloudSyncStatusConnecting');
      case 'error':
        return formatMessage('uns.cloudSyncStatusError');
      default:
        return formatMessage('uns.cloudSyncStatusDisconnected');
    }
  }, [config?.connectStatus, formatMessage]);

  const edgeNodeName = config?.edgeNodeName || '-';
  const currentStatus = config?.connectStatus || 'disconnected';
  const isConnected = currentStatus === 'connected';

  const handleDisconnect = async () => {
    setDisconnecting(true);
    try {
      const nextConfig = await disconnectCloudSync();
      setConfig(nextConfig);
      message.success(formatMessage('uns.cloudSyncDisconnected'));
    } finally {
      setDisconnecting(false);
    }
  };

  const renderValue = (value?: string | number | boolean) => {
    if (value === undefined || value === null || value === '') return '-';
    return String(value);
  };

  return (
    <>
      <Spin spinning={loading}>
        <Flex vertical className={styles['cloudsync-card']}>
          <Flex align="center" justify="space-between" className={styles['cloudsync-card-header']}>
            <Flex align="center" gap={8}>
              <CloudApp size={16} />
              <span className={styles['cloudsync-card-title']}>{formatMessage('uns.cloudSync')}</span>
              <HelpTooltip title={formatMessage('uns.cloudSyncDescription')} />
            </Flex>
          </Flex>

          <div className={styles['cloudsync-card-body']}>
            <div className={styles['cloudsync-field']}>
              <div className={styles['cloudsync-label']}>{formatMessage('uns.cloudSyncEdgeNodeName')}</div>
              <Typography.Text className={styles['cloudsync-name']} ellipsis={{ tooltip: edgeNodeName }}>
                {renderValue(edgeNodeName)}
              </Typography.Text>
            </div>

            <div className={styles['cloudsync-field']}>
              <div className={styles['cloudsync-label']}>{formatMessage('uns.cloudSyncEdgeConnectionStatus')}</div>
              <Tag className={`${styles['cloudsync-status-tag']} ${styles[statusClassMap[currentStatus]] || ''}`}>
                <span className={styles['cloudsync-status-dot']} />
                {statusText}
              </Tag>
            </div>

            {config?.lastError ? (
              <div className={styles['cloudsync-error-panel']}>
                <div className={styles['cloudsync-label']}>{formatMessage('common.errorInfo')}</div>
                <Typography.Paragraph
                  className={styles['cloudsync-error-text']}
                  ellipsis={{ rows: 2, tooltip: config.lastError }}
                >
                  {config.lastError}
                </Typography.Paragraph>
              </div>
            ) : null}
          </div>

          <Flex gap={8} className={styles['cloudsync-actions']}>
            {!isConnected ? (
              <Button
                type="primary"
                block
                className={styles['cloudsync-connect-button']}
                icon={<ConnectionSignal size={16} />}
                onClick={() => setOpenModal(true)}
              >
                {formatMessage('uns.connectToCloud')}
              </Button>
            ) : (
              <Popconfirm
                title={formatMessage('uns.disconnectFromCloud')}
                description={formatMessage('uns.cloudSyncDisconnectConfirm')}
                onConfirm={handleDisconnect}
                okButtonProps={{ loading: disconnecting }}
              >
                <Button
                  block
                  className={styles['cloudsync-disconnect-button']}
                  icon={<ConnectionSignalOff size={16} />}
                  loading={disconnecting}
                >
                  {formatMessage('uns.disconnectFromCloud')}
                </Button>
              </Popconfirm>
            )}
          </Flex>
        </Flex>
      </Spin>

      <CloudSyncConnectModal
        open={openModal}
        initialConfig={config}
        onCancel={() => setOpenModal(false)}
        onSaved={loadConfig}
      />
    </>
  );
};

export default CloudSync;
