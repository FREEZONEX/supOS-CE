import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/api/core/license';

const api = new ApiWrapper(baseUrl);

export type LicensePhase = 'normal' | 'warning' | 'grace' | 'hard_expired' | 'none';

export interface LicenseStatusResp {
  activated: boolean;
  phase: LicensePhase;
  daysLeft: number;
  expireAt?: string;
  licenseKey?: string;
  company?: string;
  lang?: string;
}

export interface LicenseActivateResp {
  success: boolean;
  licenseKey?: string;
  expireAt?: string;
  message?: string;
}

export const queryDeadline = async () => api.get('/deadline', { [CustomAxiosConfigEnum.NoMessage]: true });

// 查询 License 状态信息
export const queryLicenseStatus = async () =>
  api.get('/status', {
    [CustomAxiosConfigEnum.NoMessage]: true,
  }) as Promise<LicenseStatusResp>;

// 更换 License Key（需管理员权限）
export const replaceLicenseKey = async (licenseKey: string) =>
  api.post('/replace', { licenseKey }, { [CustomAxiosConfigEnum.NoMessage]: true });

// 已登录场景下更换离线 Bundle
export const replaceLicenseBundle = async (bundleJson: string): Promise<LicenseActivateResp> => {
  return api.post('/replace', { bundleJson }, { [CustomAxiosConfigEnum.NoMessage]: true });
};

// 登录前检查 License 状态 (无需认证)
export interface LicenseCheckResult {
  status: 'valid' | 'expired' | 'not_activated';
  phase?: LicensePhase;
  daysLeft?: number;
  expireAt?: string;
}

const isExpiredPhase = (phase?: LicensePhase) => phase === 'grace' || phase === 'hard_expired';
const isWarningPhase = (phase?: LicensePhase) => phase === 'warning';

export const checkLicenseBeforeLogin = async (): Promise<LicenseCheckResult> => {
  try {
    const status = await queryLicenseStatus();
    if (!status?.activated || status?.phase === 'none') {
      return { status: 'not_activated', phase: status?.phase, daysLeft: status?.daysLeft, expireAt: status?.expireAt };
    }
    if (isExpiredPhase(status?.phase)) {
      return { status: 'expired', phase: status?.phase, daysLeft: status?.daysLeft, expireAt: status?.expireAt };
    }
    if (isWarningPhase(status?.phase)) {
      return { status: 'valid', phase: status?.phase, daysLeft: status?.daysLeft, expireAt: status?.expireAt };
    }
    return { status: 'valid', phase: status?.phase, daysLeft: status?.daysLeft, expireAt: status?.expireAt };
  } catch (error) {
    console.error('Failed to check license:', error);
    // 不阻断登录：默认视为有效
    return { status: 'valid' };
  }
};

// 激活 License (无需认证)
export const activateLicense = async (licenseKey: string): Promise<LicenseActivateResp> => {
  return api.post('/activate', { licenseKey }, { [CustomAxiosConfigEnum.NoMessage]: true });
};

// 获取设备标识 (无需认证，用于离线激活，返回值已拼接 4 位校验码)
export interface LicenseIdentityResp {
  device_token: string;
}

export const maskDeviceToken = (token?: string) => {
  if (!token || token.length <= 8) return token || '';
  return `${token.slice(0, 4)}****${token.slice(-4)}`;
};

// license 授权的功能模块开关，对应后端 EnterpriseModules；旧版 license 可能缺字段，故全部可选
export interface LicenseModules {
  vision?: boolean;
  unsAgent?: boolean;
  notebook?: boolean;
  anchor?: boolean;
  appBuilder?: boolean;
  launchPad?: boolean;
  audit?: boolean;
  siem?: boolean;
}

export interface LicenseQuotasResp {
  maxClientConnections: number;
  maxUsers: number;
  maxApps: number;
  usedUsers: number;
  usedApps: number;
  modules?: LicenseModules;
}

export const queryQuotas = async () => api.get('/quotas') as Promise<LicenseQuotasResp>;

export const queryIdentity = async (): Promise<LicenseIdentityResp> => {
  return api.get('/identity', { [CustomAxiosConfigEnum.NoMessage]: true });
};

// 离线导入 Bundle JSON (无需认证)
export const importLicenseBundle = async (bundleJson: string): Promise<LicenseActivateResp> => {
  return api.post('/import', { bundleJson }, { [CustomAxiosConfigEnum.NoMessage]: true });
};
