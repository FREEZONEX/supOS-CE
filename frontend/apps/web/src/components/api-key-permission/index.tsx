import type { ReactNode } from 'react';
import './index.scss';

export const apiKeyPermissionValues = ['read', 'write', 'full'] as const;

export type ApiKeyPermissionValue = (typeof apiKeyPermissionValues)[number];

const defaultPermissionLabels: Record<ApiKeyPermissionValue, string> = {
  read: 'Read Only',
  write: 'Write',
  full: 'Full Access',
};

const permissionToneMap: Record<ApiKeyPermissionValue, string> = {
  read: 'read-only',
  write: 'data-writer',
  full: 'full-access',
};

const permissionAliases: Record<string, ApiKeyPermissionValue> = {
  read_only: 'read',
  data_writer: 'write',
  full_access: 'full',
};

const isKnownPermission = (permission?: string): permission is ApiKeyPermissionValue =>
  apiKeyPermissionValues.includes(permission as ApiKeyPermissionValue);

export const normalizeApiKeyPermission = (permission?: string): ApiKeyPermissionValue => {
  if (isKnownPermission(permission)) {
    return permission;
  }
  return permissionAliases[String(permission || '')] || 'read';
};

export const getApiKeyPermissionLabel = (permission?: string, label?: ReactNode) => {
  if (label) {
    return label;
  }
  return defaultPermissionLabels[normalizeApiKeyPermission(permission)];
};

export const getApiKeyPermissionOptions = (formatLabel?: (permission: ApiKeyPermissionValue) => string) =>
  apiKeyPermissionValues.map((permission) => ({
    label: formatLabel?.(permission) || defaultPermissionLabels[permission],
    value: permission,
  }));

export const ApiKeyPermissionBadge = ({ permission, label }: { permission?: string; label?: ReactNode }) => {
  const normalized = normalizeApiKeyPermission(permission);
  const tone = permissionToneMap[normalized];

  return (
    <span className={`api-key-permission-badge api-key-permission-badge--${tone}`}>
      {getApiKeyPermissionLabel(permission, label)}
    </span>
  );
};
