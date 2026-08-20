import { Alert, Card, Descriptions, Space, Tag, Typography } from 'antd';
import dayjs from 'dayjs';
import type { FleetCurrentHA } from '@/apis/core-api/fleet';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

interface FleetHAPanelProps {
  ha?: FleetCurrentHA;
}

const roleColor = (role?: string) => {
  if (role === 'active') return 'success';
  if (role === 'standby') return 'processing';
  if (role === 'faulted') return 'error';
  return 'default';
};

const syncColor = (state?: string) => {
  if (state === 'succeeded') return 'success';
  if (state === 'running') return 'processing';
  if (state === 'error') return 'error';
  return 'default';
};

export const FleetHAPanel = ({ ha }: FleetHAPanelProps) => {
  const formatMessage = useTranslate();
  const status = ha?.status;

  if (!ha?.enabled) return null;
  if (!status) {
    return (
      <Alert
        className={styles.haAlert}
        type="warning"
        showIcon
        message={formatMessage('fleet.ha.statusUnavailable')}
        description={ha.errorCode || formatMessage('fleet.ha.statusUnavailableDescription')}
      />
    );
  }

  const syncState = status.sync.enabled ? status.sync.state : 'disabled';
  const syncError = status.sync.lastError;
  return (
    <Card
      className={styles.haCard}
      title={
        <Space>
          <span>{formatMessage('fleet.ha.title')}</span>
          <Tag color={roleColor(status.role)}>{formatMessage(`fleet.ha.role.${status.role}`)}</Tag>
          <Tag color={status.vipOwned ? 'success' : 'default'}>
            {formatMessage(status.vipOwned ? 'fleet.ha.vipOwned' : 'fleet.ha.vipNotOwned')}
          </Tag>
        </Space>
      }
    >
      {status.lastErrorCode || !status.eligible ? (
        <Alert
          className={styles.haAlert}
          type="warning"
          showIcon
          message={formatMessage('fleet.ha.attention')}
          description={status.lastErrorCode || status.eligibilityReason}
        />
      ) : null}
      {syncError ? (
        <Alert
          className={styles.haAlert}
          type="error"
          showIcon
          message={formatMessage('fleet.ha.syncFailed')}
          description={syncError}
        />
      ) : null}
      <Descriptions className={styles.haDescriptions} size="small" column={{ xs: 1, sm: 2, xl: 4 }}>
        <Descriptions.Item label={formatMessage('fleet.ha.memberID')}>
          <Typography.Text className={styles.monospace} copyable>
            {status.memberId}
          </Typography.Text>
        </Descriptions.Item>
        <Descriptions.Item label={formatMessage('fleet.ha.keepalived')}>
          {status.keepalivedState || '-'}
        </Descriptions.Item>
        <Descriptions.Item label={formatMessage('fleet.ha.businessGroup')}>
          {status.serviceGroup.running} {formatMessage('fleet.ha.running')} / {status.serviceGroup.unhealthy}{' '}
          {formatMessage('fleet.ha.unhealthy')}
        </Descriptions.Item>
        <Descriptions.Item label={formatMessage('fleet.ha.syncState')}>
          <Tag color={syncColor(syncState)}>{formatMessage(`fleet.ha.sync.${syncState}`)}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label={formatMessage('fleet.ha.syncPeers')}>
          <Typography.Text className={styles.monospace}>{status.sync.peers?.join(', ') || '-'}</Typography.Text>
        </Descriptions.Item>
        <Descriptions.Item label={formatMessage('fleet.ha.lastSync')}>
          {status.sync.lastSucceededAt ? dayjs(status.sync.lastSucceededAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
        </Descriptions.Item>
      </Descriptions>
    </Card>
  );
};
