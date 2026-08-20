import {
  connectCloudSync,
  forceSyncCloudSync,
  testCloudSyncConnect,
  type CloudSyncConfigResp,
  type CloudSyncTestConnectResp,
} from '@/apis/core-api/cloudsync.ts';
import HelpTooltip from '@/components/help-tooltip';
import { useTranslate } from '@/hooks';
import { App, Button, Checkbox, Flex, Input, Modal, Typography } from 'antd';
import { useEffect, useState } from 'react';
import classNames from 'classnames';
import styles from '@/pages/uns/components/uns-dashboard/cloudsync.module.scss';
import CloudSyncRootSelector from './CloudSyncRootSelector';

const storedTokenMask = '*'.repeat(24);

const errorText = (error: unknown) => {
  if (error && typeof error === 'object') {
    const record = error as Record<string, unknown>;
    const msg = record.msg || record.message || record.error;
    if (typeof msg === 'string' && msg.trim()) return msg.trim();
  }
  if (error instanceof Error && error.message.trim()) return error.message.trim();
  return '';
};

const displayBaseURL = (value?: string) => {
  const raw = String(value || '').trim();
  if (!raw) return '';
  try {
    const url = new URL(raw);
    url.pathname = '';
    url.search = '';
    url.hash = '';
    return url.toString().replace(/\/$/, '');
  } catch {
    return raw.replace(/\/api\/.*$/, '').replace(/\/$/, '');
  }
};

export type CloudSyncConnectModalMode = 'connect' | 'scope';

interface CloudSyncConnectModalProps {
  open: boolean;
  mode?: CloudSyncConnectModalMode;
  initialConfig?: CloudSyncConfigResp | null;
  titleMessageKey?: string;
  primaryActionMessageKey?: string;
  onCancel: () => void;
  onSaved: (config: CloudSyncConfigResp) => void | Promise<void>;
}

const CloudSyncConnectModal = ({
  open,
  mode = 'connect',
  initialConfig,
  titleMessageKey,
  primaryActionMessageKey,
  onCancel,
  onSaved,
}: CloudSyncConnectModalProps) => {
  const formatMessage = useTranslate();
  const { message } = App.useApp();
  const isScopeMode = mode === 'scope';
  const resolvedTitleKey =
    titleMessageKey || (isScopeMode ? 'fleet.sync.configureScope' : 'fleet.sync.connectSyncPath');
  const resolvedPrimaryKey = primaryActionMessageKey || (isScopeMode ? 'common.save' : 'fleet.sync.connect');
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
  const rawTokenError =
    isScopeMode || isUsingStoredToken || token.trim() ? '' : formatMessage('uns.cloudSyncTokenRequired');
  const tokenError = validationVisible ? rawTokenError : '';

  const handleTest = async () => {
    setValidationVisible(true);
    if (!token.trim()) {
      message.warning(formatMessage('uns.cloudSyncTokenRequired'));
      return;
    }
    setTesting(true);
    try {
      const result = await testCloudSyncConnect({ token: isUsingStoredToken ? '' : token.trim() });
      setTestResult(result);
      if (result.ok) {
        message.success(formatMessage('uns.cloudSyncTestSuccess'));
      }
    } catch (error) {
      const msg = errorText(error) || formatMessage('uns.cloudSyncTestFailed');
      setTestResult({
        ...(testResult || {}),
        ok: false,
        message: msg,
        error: msg,
      });
    } finally {
      setTesting(false);
    }
  };

  const handleSave = async () => {
    setValidationVisible(true);
    if (!isScopeMode && rawTokenError) {
      message.warning(rawTokenError);
      return;
    }
    if (isScopeMode && !hasStoredToken) {
      message.warning(formatMessage('uns.cloudSyncTokenRequired'));
      return;
    }
    if (isScopeMode && selectedRootNodeIDs.length === 0) {
      message.warning(formatMessage('uns.cloudSyncSelectTopicRequired'));
      return;
    }
    setSaving(true);
    try {
      let config = await connectCloudSync({
        token: isScopeMode || isUsingStoredToken ? '' : token.trim(),
        selectedRootNodeIDs,
        expectedScopeRevision: initialConfig?.scopeRevision ?? 0,
        syncMetricEnabled: storeDataAtCenter,
      });
      message.success(
        formatMessage(isScopeMode || initialConfig?.configured ? 'uns.cloudSyncUpdated' : 'uns.cloudSyncConnected')
      );
      if (selectedRootNodeIDs.length > 0 && (isScopeMode || initialConfig?.configured)) {
        try {
          config = await forceSyncCloudSync();
        } catch {
          message.error(formatMessage('cluster.forceSyncFailed', {}, 'Force sync failed.'));
        }
      }
      await onSaved(config);
      onCancel();
    } catch (error) {
      const msg = errorText(error) || formatMessage('uns.cloudSyncTestFailed');
      if (!isScopeMode) {
        setTestResult({
          ...(testResult || {}),
          ok: false,
          message: msg,
          error: msg,
        });
      } else {
        message.error(msg);
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <Modal
      title={formatMessage(resolvedTitleKey)}
      open={open}
      width={620}
      className={styles['cloudsync-connect-modal']}
      destroyOnHidden
      onCancel={onCancel}
      footer={
        <Flex justify="flex-end" gap={8} className={styles['cloudsync-modal-footer']}>
          <Button onClick={onCancel}>{formatMessage('common.cancel')}</Button>
          {!isScopeMode ? (
            <Button onClick={handleTest} loading={testing}>
              {formatMessage('uns.testConnection')}
            </Button>
          ) : null}
          <Button type="primary" onClick={handleSave} loading={saving}>
            {formatMessage(resolvedPrimaryKey)}
          </Button>
        </Flex>
      }
    >
      <Flex vertical gap={14} className={styles['cloudsync-modal']}>
        {!isScopeMode ? (
          <>
            <div className={styles['cloudsync-required-label']}>
              <span className={styles['cloudsync-required-mark']}>*</span>
              <span>{formatMessage('uns.token')}</span>
              <HelpTooltip title={formatMessage('uns.cloudSyncTokenHint')} />
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
          </>
        ) : null}
        <CloudSyncRootSelector
          value={selectedRootNodeIDs}
          validationVisible={validationVisible}
          required={isScopeMode}
          onChange={setSelectedRootNodeIDs}
        />
        <Checkbox checked={storeDataAtCenter} onChange={(event) => setStoreDataAtCenter(event.target.checked)}>
          {formatMessage('uns.cloudSyncStoreAtCenter')}
        </Checkbox>
        {!isScopeMode && testResult ? (
          <div
            className={classNames(
              styles['cloudsync-test-result'],
              !testResult.ok && styles['cloudsync-test-result-error']
            )}
          >
            <Typography.Text strong>{formatMessage('uns.cloudSyncTestResultTitle')}</Typography.Text>
            <span>
              {formatMessage('common.status')}:{' '}
              {testResult.ok ? formatMessage('mqttAuth.connected') : formatMessage('mqttAuth.disconnected')}
            </span>
            {!testResult.ok ? (
              <span>
                {formatMessage('mqttAuth.lastError')}: {testResult.error || testResult.message || '-'}
              </span>
            ) : null}
            <span>
              {formatMessage('uns.cloudSyncWorkspaceID')}: {testResult.workspaceID || '-'}
            </span>
            <span>
              {formatMessage('uns.mqttClientId')}: {testResult.clientID || '-'}
            </span>
            <span>
              {formatMessage('uns.cloudSyncMqttAuthID')}: {testResult.mqttAuthID || testResult.connectClientKey || '-'}
            </span>
            <span>
              {formatMessage('uns.cloudSyncHttpEndpoint')}: {displayBaseURL(testResult.httpEndpoint) || '-'}
            </span>
            <span>
              {formatMessage('uns.cloudSyncMqttEndpoint')}: {testResult.mqttBrokers || '-'}
            </span>
            <span>
              {formatMessage('mqttAuth.syncMode')}: {testResult.syncMode || '-'}
            </span>
          </div>
        ) : null}
      </Flex>
    </Modal>
  );
};

export default CloudSyncConnectModal;
