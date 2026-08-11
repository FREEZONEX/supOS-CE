import type { CloudSyncLogItem } from '@/apis/core-api/cloudsync';

type Translate = (id: string, values?: Record<string, unknown>) => string;

type UserSyncStats = {
  created?: number;
  updated?: number;
  disabled?: number;
  deleted?: number;
  skipped?: number;
  conflicts?: number;
};

const parseUserSyncStats = (details?: string): UserSyncStats | undefined => {
  try {
    return details ? JSON.parse(details).stats : undefined;
  } catch {
    return undefined;
  }
};

export const formatCloudSyncSummary = (record: CloudSyncLogItem, formatMessage: Translate) => {
  const stats = record.syncType === 'user' ? parseUserSyncStats(record.details) : undefined;
  if (stats) {
    return formatMessage('fleet.sync.summary.users', {
      created: Number(stats.created || 0),
      updated: Number(stats.updated || 0),
      disabled: Number(stats.disabled || 0),
      deleted: Number(stats.deleted || 0),
      skipped: Number(stats.skipped || 0),
      conflicts: Number(stats.conflicts || 0),
    });
  }

  switch (record.resultCode) {
    case 'uns.applied':
      return formatMessage('fleet.sync.summary.metadataApplied');
    case 'uns.partial':
      return formatMessage('fleet.sync.summary.metadataPartial');
    case 'uns.duplicate':
      return formatMessage('fleet.sync.summary.duplicateSkipped');
    case 'uns.failed':
      return formatMessage('fleet.sync.summary.unsFailed');
    case 'user.applied':
    case 'user.partial':
      return formatMessage('fleet.sync.summary.usersEmpty');
    case 'user.failed':
      return formatMessage('fleet.sync.summary.usersFailed');
    default:
      break;
  }

  // Compatibility for histories created before resultCode was introduced.
  switch (record.summary) {
    case 'metadata applied':
      return formatMessage('fleet.sync.summary.metadataApplied');
    case 'metadata applied with skippable conflicts':
      return formatMessage('fleet.sync.summary.metadataPartial');
    case 'duplicate metadata snapshot skipped':
      return formatMessage('fleet.sync.summary.duplicateSkipped');
    default:
      return record.summary || '-';
  }
};
