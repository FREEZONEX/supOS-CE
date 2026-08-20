import { LOGIN_URL } from '@/common-types/constans.ts';

export type LicenseActivationStatus = 'not_activated' | 'expired';
type LicenseGateStatus = 'active' | LicenseActivationStatus;

export const LICENSE_ACTIVATION_PATH = '/license-activation';
const PUBLIC_ACCESS_PATHS = new Set(['/403', '/404', '/share', '/cli-auth', LOGIN_URL]);

export const isLicenseActivationPath = (pathname = window.location.pathname) => {
  void pathname;
  return false;
};
export const isPublicAccessPath = (pathname = window.location.pathname) => PUBLIC_ACCESS_PATHS.has(pathname);
export const buildLicenseActivationUrl = (status: LicenseActivationStatus) => {
  void status;
  return LOGIN_URL;
};
export const replaceWithLicenseActivation = (status: LicenseActivationStatus) => {
  void status;
};
export const resolveLicenseGate = async (): Promise<{
  status: LicenseGateStatus;
  phase?: string;
  daysLeft?: number;
  redirectTo?: string;
}> => ({ status: 'active', phase: 'active' });
