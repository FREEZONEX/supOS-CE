import {
  connectCloudSync,
  forceSyncCloudSync,
  testCloudSyncConnect,
  type CloudSyncConfigResp,
  type CloudSyncTestConnectResp,
} from '@/apis/core-api/cloudsync.ts';
import { useTranslate } from '@/hooks';
import { App, Button, Checkbox, Flex, Input, Modal, Typography } from 'antd';
import { useEffect, useState } from 'react';
import styles from './cloudsync.module.scss';
import CloudSyncRootSelector from '@/components/cloud-sync/CloudSyncRootSelector';

const storedTokenMask = '*'.repeat(24);

interface CloudSyncConnectModalProps {
  open: boolean;
  initialConfig?: CloudSyncConfigResp | null;
  onCancel: () => void;
  onSaved: () => void;
}

const CloudSyncConnectModal = ({ open, initialConfig, onCancel, onSaved }: CloudSyncConnectModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const [token, setToken] = useState('');
  const [storeDataAtCenter, setStoreDataAtCenter] = useState(false);
  const [selectedRootNodeIDs, setSelectedRootNodeIDs] = useState<string[]>([]);
  const [validationVisible, setValidationVisible] = useState(false);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<CloudSyncTestConnectResp | null>(null);
  const hasStoredToken = !!initialConfig?.hasToken;

  useEffect(() => {
    if (!open) return;
    setToken(hasStoredToken ? storedTokenMask : '');
    setValidationVisible(false);
    setTestResult(null);
    setStoreDataAtCenter(!!initialConfig?.syncMetricEnabled);
    setSelectedRootNodeIDs(initialConfig?.selectedRootNodeIDs || []);
  }, [hasStoredToken, initialConfig, open]);

  const isUsingStoredToken = hasStoredToken && token === storedTokenMask;
  const rawTokenError = isUsingStoredToken || token.trim() ? '' : formatMessage('uns.cloudSyncTokenRequired');
  const tokenError = validationVisible ? rawTokenError : '';

  const handleTest = async () => {
    setValidationVisible(true);
    if (isUsingStoredToken) {
      message.warning(formatMessage('uns.cloudSyncTokenRequired'));
      return;
    }
    if (!token.trim()) {
      message.warning(formatMessage('uns.cloudSyncTokenRequired'));
      return;
    }
    setTesting(true);
    try {
      const result = await testCloudSyncConnect({ token: token.trim() });
      setTestResult(result);
      message.success(formatMessage('uns.cloudSyncTestSuccess'));
    } finally {
      setTesting(false);
    }
  };

  const handleSave = async () => {
    setValidationVisible(true);
    if (rawTokenError) {
      message.warning(rawTokenError);
      return;
    }
    if (initialConfig?.configured && selectedRootNodeIDs.length === 0) {
      message.warning(formatMessage('uns.cloudSyncSelectTopicRequired'));
      return;
    }
    setSaving(true);
    try {
      await connectCloudSync({
        token: isUsingStoredToken ? '' : token.trim(),
        selectedRootNodeIDs,
        expectedScopeRevision: initialConfig?.scopeRevision ?? 0,
        syncMetricEnabled: storeDataAtCenter,
      });
      message.success(formatMessage(initialConfig?.configured ? 'uns.cloudSyncUpdated' : 'uns.cloudSyncConnected'));
      if (initialConfig?.configured && selectedRootNodeIDs.length > 0) {
        try {
          await forceSyncCloudSync();
        } catch {
          message.error(formatMessage('cluster.forceSyncFailed', {}, 'Force sync failed.'));
        }
      }
      onSaved();
      onCancel();
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      title={formatMessage(initialConfig?.configured ? 'uns.editCloudSync' : 'uns.connectToMasterNode')}
      open={open}
      width={452}
      className={styles['cloudsync-connect-modal']}
      destroyOnHidden
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel}>
          {formatMessage('common.cancel')}
        </Button>,
        <Button key="test" onClick={handleTest} loading={testing} disabled={isUsingStoredToken}>
          {formatMessage('uns.testConnection')}
        </Button>,
        <Button key="save" type="primary" onClick={handleSave} loading={saving}>
          {formatMessage(initialConfig?.configured ? 'common.save' : 'uns.connectToCloud')}
        </Button>,
      ]}
    >
      <Flex vertical gap={14} className={styles['cloudsync-modal']}>
        <div className={styles['cloudsync-required-label']}>
          <span className={styles['cloudsync-required-mark']}>*</span>
          {formatMessage('uns.token')}
        </div>
        <div>
          <Input
            value={token}
            placeholder={formatMessage('uns.cloudSyncTokenPlaceholder')}
            className={styles['cloudsync-token-input']}
            onChange={(event) => {
              setToken(event.target.value);
              setTestResult(null);
            }}
            onFocus={() => {
              if (isUsingStoredToken) {
                setToken('');
              }
            }}
            onBlur={() => {
              if (hasStoredToken && token.trim() === '') {
                setToken(storedTokenMask);
              }
            }}
          />
          {tokenError ? <div className={styles['cloudsync-error']}>{tokenError}</div> : null}
        </div>
        <CloudSyncRootSelector
          value={selectedRootNodeIDs}
          validationVisible={validationVisible}
          required={!!initialConfig?.configured}
          onChange={setSelectedRootNodeIDs}
        />
        <Checkbox checked={storeDataAtCenter} onChange={(event) => setStoreDataAtCenter(event.target.checked)}>
          {formatMessage('uns.cloudSyncStoreAtCenter')}
        </Checkbox>
        {testResult?.ok ? (
          <div className={styles['cloudsync-test-result']}>
            <Typography.Text strong>{formatMessage('uns.cloudSyncTestResultTitle')}</Typography.Text>
            <span>
              {formatMessage('uns.cloudSyncWorkspaceID')}: {testResult.workspaceID || '-'}
            </span>
            <span>
              {formatMessage('uns.mqttClientId')}: {testResult.clientID || '-'}
            </span>
            <span>
              {formatMessage('uns.cloudSyncConnectClientKey')}: {testResult.connectClientKey || '-'}
            </span>
          </div>
        ) : null}
      </Flex>
    </Modal>
  );
};

export default CloudSyncConnectModal;
